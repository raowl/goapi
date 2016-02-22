package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/raowl/goapi/config"
	"github.com/raowl/goapi/repos"
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
	var err error
	var userregistered bool
	type Validationerror struct {
		Error_msg string `bson:"error_msg,omitempty" json:"error_msg,omitempty"`
		Error_code string `bson:"error_code,omitempty" json:"error_code,omitempty"`
	}
	type ValidationerrorResource struct {
		Data Validationerror `json:"data"`
	}

	body := context.Get(r, "body").(*repos.UserResource)
	repo := repos.UserRepo{c.Db.C("users")}
	userregistered, err = repo.UserAlreadyExists(body.Data.Username)
	if err != nil {
		panic("err")
	}
	if userregistered {
		
	     fmt.Println("el usuario ya existe");
             /* panic("error") */
	     w.Header().Set("Content-Type", "application/vnd.api+json")
	     w.WriteHeader(400)
	     // e := ValidationerrorResource{Validationerror{"duplicate username", "125"}}
	     e := Validationerror{"duplicate username", "125"}
	     json.NewEncoder(w).Encode(e)
	} else {
	     	fmt.Println("el usuario no existe");
		err = repo.Create(&body.Data)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(body)
	}
}

func (c *AppContext) UserHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := repos.UserRepo{c.Db.C("users")}

	var err error
	var user repos.UserResource

	//fmt.Printf("userId")
	userId := context.Get(r, "userid").(string)
	

	if params.ByName("id") != "undefined" {
		fmt.Println("entro por aca")
		user, err = repo.Find(params.ByName("id"))
	} /*  else {
		userId := context.Get(r, "userid").(string)
		user, err = repo.Find(userId)
	} */

	/* if params.ByName("username") != "undefined" {
		fmt.Println("entro por aca")
		user, err = repo.Find(params.ByName("username"))
	} /*  else { */
	if err != nil {
		panic(err)
	}

	following, err := repo.GetFByIds(user.Data.Following)
	user.Data.FollowInfo = following
	fmt.Printf("Currently following00000000000000000000000000000000000000000000000000000000000000000000000000...\n")
	fmt.Printf("%+v\n", following)
	followed, err := repo.GetFollowers(bson.ObjectIdHex(userId))
	user.Data.FollowedInfo = followed
	fmt.Printf("USER===00000000000000000000000000000000000000000000000000000000000000000000000000...\n")
	fmt.Printf("%+v\n", user)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(user)
}

func (c *AppContext) UserWithSkillsHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	lang := r.FormValue("lang")
	print("language")
	print(lang)
	userRepo := repos.UserRepo{c.Db.C("users")}
	categorySkillsRepo := repos.SkillCategoryRepo{c.Db.C("skill_category")}
	skillRepo := repos.SkillRepo{c.Db.C("skill")}

	user, err := userRepo.Find(params.ByName("id"))

	fmt.Println("***********************************************************")
	fmt.Printf("%+v\n", user.Data.Skills)
	fmt.Println("***********************************************************")
	skills, err := skillRepo.GetByIds(user.Data.Skills, lang)
	following, err := userRepo.GetByIds(user.Data.Following)

	fmt.Printf("Currently following...\n")
	fmt.Printf("%+v\n", following)

	fmt.Printf("Skills...\n")
	fmt.Printf("%+v\n", skills)

	if err != nil {
		panic(err)
	}

	type SkillCompleteLocal struct {
		Id           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
		Name         string        `json:"skill_name,omitempty" bson:"skill_name,omitempty"`
		Category     bson.ObjectId `json:"category,omitempty" bson:"category,omitempty"`
		CategoryName string        `json:"category_name,omitempty" bson:"category_name,omitempty"`
	}

	//other way is an $in in every category also in this loop, but gettting by id prob be quicker, check...
	CatSkillInfo := make([]SkillCompleteLocal, len(skills.Data))
	for i := range skills.Data {
		category, err := categorySkillsRepo.GetById(skills.Data[i].Category, lang)
		print("category......")
		fmt.Printf("%+v\n", skills)
		if err != nil {
			print(err)
		}
		CatSkillInfo[i] = SkillCompleteLocal{skills.Data[i].Id, skills.Data[i].Name[0], skills.Data[i].Category, category.Data.Name[0]}
	}

	type AllInfo struct {
		repos.UserResource
		CatSkillInfo []SkillCompleteLocal
		Follow       repos.UserCollection
	}

	AllInfoI := AllInfo{user, CatSkillInfo, following}
	//AllInfoI := ""

	fmt.Printf("ALL INFO\n")
	fmt.Printf("%+v\n", AllInfoI)
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
		"facebookid": user_resource.Data.FacebookId,
		"firstname": user_resource.Data.FirstName,
		"lastname": user_resource.Data.LastName,
		"following": user_resource.Data.Following,
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(data)
}

func (c *AppContext) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("enter to updateuser handler")
	//params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*repos.UserResource)
	//fmt.Printf("userId")
	userId := context.Get(r, "userid").(string)
	body.Data.Id = bson.ObjectIdHex(userId)
	//fmt.Printf("BODY DATA")
	//fmt.Printf("%+v\n", body.Data)
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
