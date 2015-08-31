package repos

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	//"time"
)

type SkillTempCategory struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name [1]string     `bson:"name,omitempty" json:"name,omitempty"`
}
type SkillCategory struct {
	Id    bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name  string        `bson:"name,omitempty" json:"name,omitempty"`
	Skill []Skill       `bson:"skills,omitempty" json:"skills,omitempty"`
}

type SkillTemp struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     [1]string     `bson:"name,omitempty" json:"name,omitempty"`
	Category bson.ObjectId `json:"category,omitempty" bson:"category,omitempty"`
}

type Skill struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name,omitempty"`
	Category bson.ObjectId `json:"category,omitempty" bson:"category,omitempty"`
}

type SkillCategoryCollection struct {
	Data []SkillCategory `json:"data"`
}

type SkillTempCategoryCollection struct {
	Data []SkillTempCategory `json:"data"`
}

type SkillCategoryResource struct {
	Data SkillTempCategory `json:"data"`
}

type SkillCategoryRepo struct {
	Coll *mgo.Collection
}

type SkillRepo struct {
	Coll *mgo.Collection
}

type SkillCollection struct {
	Data []SkillTemp `json:"data"`
}

func (r *SkillRepo) All() (SkillCollection, error) {
	result := SkillCollection{[]SkillTemp{}}
	// err := r.Coll.Find(nil, bson.M{"$slice": []int{1,1}}).All(&result.Data)
	err := r.Coll.Find(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SkillRepo) GetByCategoryId(id bson.ObjectId) ([]Skill, error) {
	result := []SkillTemp{}
	err := r.Coll.Find(bson.M{"category": id}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result)
	if err != nil {
		log.Fatal("error")
	}

	skillsallnew := make([]Skill, len(result))
	for i := range result {
		skillsallnew[i] = Skill{Id: result[i].Id, Name: result[i].Name[0]}
	}

	return skillsallnew, nil
}
func (r *SkillRepo) GetByIds(ids []bson.ObjectId) (SkillCollection, error) {
	result := SkillCollection{[]SkillTemp{}}
	fmt.Println("SKILL IDS....")
	fmt.Printf("%+v\n", ids)
	/* oids := make([]bson.ObjectId, len(ids))
	for i := range ids {
		oids[i] = bson.ObjectId(ids[i])
	} */
	err := r.Coll.Find(bson.M{"_id": bson.M{"$in": ids}}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result.Data)

	fmt.Println("REEEEEEEEEEEEEEEEEEEEEEEESSSSSULT")
	fmt.Printf("%v\n", result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SkillCategoryRepo) GetById(id bson.ObjectId) (SkillCategoryResource, error) {
	result := SkillCategoryResource{}
	err := r.Coll.FindId(id).One(&result.Data)
	fmt.Println("ACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Printf("%+v\n", result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SkillCategoryRepo) All() (SkillCategoryCollection, error) {
	result := SkillTempCategoryCollection{[]SkillTempCategory{}}
	err := r.Coll.Find(bson.M{}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result.Data)
	if err != nil {
		log.Fatal("Error")
	}

	skillsallnew := make([]SkillCategory, len(result.Data))
	for i := range result.Data {
		skillsallnew[i] = SkillCategory{Id: result.Data[i].Id, Name: result.Data[i].Name[0]}
	}
	skillsallnewcol := SkillCategoryCollection{skillsallnew}
	fmt.Println(result)

	//return result, nil
	return skillsallnewcol, nil
}
