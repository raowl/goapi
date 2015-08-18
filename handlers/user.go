package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"goapi/config"
	"goapi/repos"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	//"path/filepath"
	"fmt"
	"time"
)

//type AppContext struct {
//	db *mgo.Database
//}

//POST: /api/v1/auth/ handler
func (c *AppContext) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*repos.UserResource)
	repo := repos.UserRepo{c.Db.C("users")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)

}

func (c *AppContext) UserHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := repos.UserRepo{c.Db.C("users")}

	var err error
	var user repos.UserResource

	if params.ByName("id") != "undefined" {
		fmt.Println("entro por aca")
		user, err = repo.Find(params.ByName("id"))
	} else {
		userId := context.Get(r, "userid").(string)
		user, err = repo.Find(userId)
	}

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(user)
}

func (c *AppContext) UserWithSkillsHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	userRepo := repos.UserRepo{c.Db.C("users")}
	user, err := userRepo.Find(params.ByName("id"))

	//fmt.Printf("%+v\n", user)

	repoSkills := repos.SkillRepo{c.Db.C("skill")}
	repoCategorySkills := repos.SkillCategoryRepo{c.Db.C("skill_category")}

	//fmt.Printf("%+v\n", user.Data.Skills)
	skills, err := repoSkills.GetByIds(user.Data.Skills)
	//fmt.Printf("%+v\n", skills)

	if err != nil {
		panic(err)
	}

	type SkillCompleteLocal struct {
		Id           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
		Name         string        `json:"name,omitempty"`
		Category     bson.ObjectId `json:"category_id,omitempty" bson:"category_id,omitempty"`
		CategoryName string        `json:"category_name,omitempty" bson:"category_name,omitempty"`
	}

	//other way is an $in in every category also in this loop, but gettting by id prob be quicker, check...
	CatSkillInfo := make([]SkillCompleteLocal, len(skills.Data))
	for i := range skills.Data {
		category, err := repoCategorySkills.GetById(skills.Data[i].Category)
		if err != nil {
			panic(err)
		}
		CatSkillInfo[i] = SkillCompleteLocal{skills.Data[i].Id, skills.Data[i].Name, skills.Data[i].Category, category.Data.Name}
	}

	type AllInfo struct {
		repos.UserResource
		CatSkillInfo []SkillCompleteLocal
	}

	AllInfoI := AllInfo{user, CatSkillInfo}
	//allInfoI := allInfo{userInfo: user}

	//fmt.Printf("%+v\n", CatSkillInfo)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(AllInfoI)
}

//POST: /api/v1/user/login/ handler
func (c *AppContext) AuthUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		privateKey []byte
	)
	body := context.Get(r, "body").(*repos.UserResource)
	repo := repos.UserRepo{c.Db.C("users")}
	user_resource, err := repo.Authenticate(body.Data)
	if err != nil {
		panic(err)
	}
	//data := map[string]interface{"token": ""}
	//var data interface{}

	// Create JWT token
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["userid"] = user_resource.Data.Id
	// Expire in 5 mins
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	//token.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	//cprivateKey, _ := filepath.Abs(config.Info.PrivateKey)
	print("private key")
	print(config.Info.PrivateKey)
	privateKey, err = ioutil.ReadFile(config.Info.PrivateKey)
	if err != nil {
		panic(err)
	}
	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		panic(err)
	}

	fmt.Println("aca")
	//fmt.Println(user_resource.Data.Skills)
	data := map[string]interface{}{
		"token":     tokenString,
		"id":        user_resource.Data.Id,
		"skills":    user_resource.Data.Skills,
		"following": user_resource.Data.Following,
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(data)
}

func (c *AppContext) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("enter to updateuser handler")
	//params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*repos.UserResource)
	fmt.Printf("userId")
	userId := context.Get(r, "userid").(string)
	body.Data.Id = bson.ObjectIdHex(userId)
	//body.Data.Skills = []bson.ObjectId{bson.ObjectIdHex(params.ByName("id"))}
	//body.Data.Following = followingIdArr
	//body.Data.Skills = []repos.Skill{{userId}}
	repo := repos.UserRepo{c.Db.C("users")}
	err := repo.Update(&body.Data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *AppContext) UserUnfollowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("enter to updateuser handler")
	params := context.Get(r, "params").(httprouter.Params)
	fmt.Printf("userId")
	//userId := context.Get(r, "userid").(string)
	userId := bson.ObjectIdHex(context.Get(r, "userid").(string))
	unfollowId := bson.ObjectIdHex(params.ByName("id"))
	repo := repos.UserRepo{c.Db.C("users")}
	err := repo.Unfollow(userId, unfollowId)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

//POST:/api/v1/developers/logout/ handler
//func (s *Server) LogoutDev(w http.ResponseWriter, r *http.Request) {
//	s.logger.SetPrefix("LogoutDev: ")
//
//	//parse body
//	logoutRequest := &LogoutRequest{}
//
//	if err := s.readJson(logoutRequest, r, w); err != nil {
//		s.badRequest(r, w, err, "malformed logout request")
//		return
//	}
//
//	//logout
//	if logout_err := s.logout(logoutRequest.Email); logout_err != nil {
//		s.internalError(r, w, logout_err, logoutRequest.Email+" : could not logout")
//	}
//
//	//response
//	response := LogoutResponse{Status: "ok"}
//	s.serveJson(w, &response)
//
//}
//func (c *appContext) markersHandler(w http.ResponseWriter, r *http.Request) {
//	repo := MarkerRepo{c.db.C("markers")}
//	markers, err := repo.All()
//	if err != nil {
//		panic(err)
//	}
//
//	w.Header().Set("Content-Type", "application/vnd.api+json")
//	json.NewEncoder(w).Encode(markers)
//}
//
//func (c *appContext) markerHandler(w http.ResponseWriter, r *http.Request) {
//	params := context.Get(r, "params").(httprouter.Params)
//	repo := MarkerRepo{c.db.C("markers")}
//	marker, err := repo.Find(params.ByName("id"))
//	if err != nil {
//		panic(err)
//	}
//
//	w.Header().Set("Content-Type", "application/vnd.api+json")
//	json.NewEncoder(w).Encode(marker)
//}
//
//func (c *appContext) createMarkerHandler(w http.ResponseWriter, r *http.Request) {
//	body := context.Get(r, "body").(*MarkerResource)
//	repo := MarkerRepo{c.db.C("markers")}
//	err := repo.Create(&body.Data)
//	if err != nil {
//		panic(err)
//	}
//
//	w.Header().Set("Content-Type", "application/vnd.api+json")
//	w.WriteHeader(201)
//	json.NewEncoder(w).Encode(body)
//}
//
//func (c *appContext) updateMarkerHandler(w http.ResponseWriter, r *http.Request) {
//	params := context.Get(r, "params").(httprouter.Params)
//	body := context.Get(r, "body").(*MarkerResource)
//	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
//	repo := MarkerRepo{c.db.C("markers")}
//	err := repo.Update(&body.Data)
//	if err != nil {
//		panic(err)
//	}
//
//	w.WriteHeader(204)
//	w.Write([]byte("\n"))
//}
//
//func (c *appContext) deleteMarkerHandler(w http.ResponseWriter, r *http.Request) {
//	params := context.Get(r, "params").(httprouter.Params)
//	repo := MarkerRepo{c.db.C("markers")}
//	err := repo.Delete(params.ByName("id"))
//	if err != nil {
//		panic(err)
//	}
//
//	w.WriteHeader(204)
//	w.Write([]byte("\n"))
//}
