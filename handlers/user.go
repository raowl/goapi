package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"goapi/repos"
	"net/http"
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

//POST: /api/v1/user/login/ handler
func (c *AppContext) AuthUserHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*repos.UserResource)
	repo := repos.UserRepo{c.Db.C("users")}
	status, user_resource, err := repo.Authenticate(body.Data)
	if err != nil {
		panic(err)
	}
	data := map[string]string{"token": ""}

	if status {
		// Create JWT token
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["userid"] = user_resource.Data.Id
		// Expire in 5 mins
		token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
		tokenString, err := token.SignedString([]byte("SecretKey12345"))
		if err != nil {
			panic(err)
		}

		data = map[string]string{
			"token": tokenString,
		}
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(data)
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
