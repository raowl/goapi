package handlers

import (
	"encoding/json"
	"fmt"
	"goapi/repos"
	"net/http"
)

/* type AppContext struct {
	Db *mgo.Database
} */

func (c *AppContext) SkillsHandler(w http.ResponseWriter, r *http.Request) {
	repo := repos.SkillCategoryRepo{c.Db.C("skill_category")}
	skillcat, err := repo.All()
	repo2 := repos.SkillRepo{c.Db.C("skill")}
	//skillcat, err := repo.All()

	for i := range skillcat.Data {
		skills, _ := repo2.GetByCategoryId(skillcat.Data[i].Id)
		fmt.Print(skills)
		skillcat.Data[i].Skill = skills
	}

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(skillcat)
}

/* func (c *AppContext) UserWithSkillsHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := repos.userRepo{c.Db.C("users")}
	user, err := repo.Find(params.ByName("id"))

	//fmt.Printf("%+v\n", marker)

	repoSkills := repos.SkillRepo{c.Db.C("skills")}

	//oids := make([]bson.ObjectId, len(marker.Data.CheckIns))
	for i := range user.Data.Skills {
		oids[i] = marker.Data.CheckIns[i].CheckUser
	}

	fmt.Println("OIDS")
	fmt.Printf("%+v\n", oids)
	//users, err := repoUsers.Coll.Find(bson.ObjectId("55aae2e98ae0044f89000003"))
	users, err := repoSkills.GetByIds(oids)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", users)
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(users)
} */
