package handlers

import (
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type appContext struct {
	db *mgo.Database
}

//POST: /api/v1/auth/ handler
func (c *appContext) createUserHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*UserResource)
	repo := UserRepo{c.db.C("users")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)

}

//POST: /api/v1/user/login/ handler
func (s *Server) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	//s.logger.SetPrefix("LoginDev: ")

	//get email password from request
	//email, password := getUserPassword(r)

	//get collection
	//session := s.db.GetSession()
	//defer session.Close()
	//c := session.DB("").C(dName)

	//developer := Developer{}

	//find and login developer

	body := context.Get(r, "body").(*UserResource)
	repo := UserRepo{c.db.C("users")}
	err := repo.Authenticate(&body.Data)
	if err != nil {
		panic(err)



	s.logger.Printf("find developer: %s %s", email, password)
	if email != "" && password != "" {
		if findErr := c.Find(bson.M{"email": email}).One(&developer); findErr != nil {
			s.notFound(r, w, findErr, email+" user not found")
			return
		}
		//match password and set session
		if matchPassword(password, developer.Hash, developer.Salt) {
			access_token, err := s.genAccessToken(developer.Email)
			if err != nil {
				s.internalError(r, w, err, email+" generate access token")
			}
			//respond with developer profile
			response := LoginResponse{ObjectId: developer.UrlToken, AccessToken: access_token, Status: "ok"}
			s.serveJson(w, &response)

		} else {

			s.notFound(r, w, nil, email+" password match failed")
		}

	} else {
		s.notFound(r, w, nil, "email empty")
	}

}

//POST:/api/v1/developers/logout/ handler
func (s *Server) LogoutDev(w http.ResponseWriter, r *http.Request) {
	s.logger.SetPrefix("LogoutDev: ")

	//parse body
	logoutRequest := &LogoutRequest{}

	if err := s.readJson(logoutRequest, r, w); err != nil {
		s.badRequest(r, w, err, "malformed logout request")
		return
	}

	//logout
	if logout_err := s.logout(logoutRequest.Email); logout_err != nil {
		s.internalError(r, w, logout_err, logoutRequest.Email+" : could not logout")
	}

	//response
	response := LogoutResponse{Status: "ok"}
	s.serveJson(w, &response)

}
func (c *appContext) markersHandler(w http.ResponseWriter, r *http.Request) {
	repo := MarkerRepo{c.db.C("markers")}
	markers, err := repo.All()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(markers)
}

func (c *appContext) markerHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := MarkerRepo{c.db.C("markers")}
	marker, err := repo.Find(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(marker)
}

func (c *appContext) createMarkerHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*MarkerResource)
	repo := MarkerRepo{c.db.C("markers")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (c *appContext) updateMarkerHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*MarkerResource)
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
	repo := MarkerRepo{c.db.C("markers")}
	err := repo.Update(&body.Data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *appContext) deleteMarkerHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := MarkerRepo{c.db.C("markers")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}
