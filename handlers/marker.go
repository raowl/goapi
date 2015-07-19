package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"goapi/repos"
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

func (c *AppContext) MarkerWithUsersHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := repos.MarkerRepo{c.Db.C("markers")}
	marker, err := repo.Find(params.ByName("id"))

	//fmt.Printf("%+v\n", marker)

	repoUsers := repos.UserRepo{c.Db.C("users")}

	oids := make([]bson.ObjectId, len(marker.Data.CheckIns))
	for i := range marker.Data.CheckIns {
		oids[i] = marker.Data.CheckIns[i].CheckUser
	}

	fmt.Println("OIDS")
	fmt.Printf("%+v\n", oids)
	//users, err := repoUsers.Coll.Find(bson.ObjectId("55aae2e98ae0044f89000003"))
	users, err := repoUsers.GetByIds(oids)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", users)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(users)
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
	fmt.Printf("userId")
	userId := bson.ObjectIdHex(context.Get(r, "userid").(string))
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))
	body.Data.CheckIns = []repos.CheckIn{{userId}}
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
