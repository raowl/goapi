package handlers

import (
	"coworkingApi/repos"
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type AppContext struct {
	Db *mgo.Database
}

func (c *AppContext) MarkersHandler(w http.ResponseWriter, r *http.Request) {
	repo := repos.MarkerRepo{c.Db.C("markers")}
	markers, err := repo.All()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(markers)
}

func (c *AppContext) MarkerHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := repos.MarkerRepo{c.Db.C("markers")}
	marker, err := repo.Find(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(marker)
}

func (c *AppContext) CreateMarkerHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*repos.MarkerResource)
	repo := repos.MarkerRepo{c.Db.C("markers")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (c *AppContext) UpdateMarkerHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*repos.MarkerResource)
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
	repo := repos.MarkerRepo{c.Db.C("markers")}
	err := repo.Update(&body.Data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (c *AppContext) DeleteMarkerHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := repos.MarkerRepo{c.Db.C("markers")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}
