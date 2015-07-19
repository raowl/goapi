package handlers

import (
	"encoding/json"
	"goapi/repos"
	"net/http"
)

/* type AppContext struct {
	Db *mgo.Database
} */

func (c *AppContext) SkillsHandler(w http.ResponseWriter, r *http.Request) {
	repo := repos.SkillCategoryRepo{c.Db.C("skill_categoryd")}
	skills, err := repo.All()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(skills)
}
