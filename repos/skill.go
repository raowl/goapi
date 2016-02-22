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

func (r *SkillRepo) All(lang string) (SkillCollection, error) {
	result := SkillCollection{[]SkillTemp{}}
	var err error
	// err := r.Coll.Find(nil, bson.M{"$slice": []int{1,1}}).All(&result.Data)
	if lang == "es" {
		err = r.Coll.Find(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result.Data)
	} else {
		err = r.Coll.Find(bson.M{"name": bson.M{"$slice": []int{0, 1}}}).All(&result.Data)
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SkillRepo) GetByCategoryId(id bson.ObjectId, lang string) ([]Skill, error) {
	result := []SkillTemp{}
	var err error
	print("LANGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG======================================")
	print(lang)
	if lang == "es" {
		err = r.Coll.Find(bson.M{"category": id}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result)
	} else {
		err = r.Coll.Find(bson.M{"category": id}).Select(bson.M{"name": bson.M{"$slice": []int{0, 1}}}).All(&result)
	}
	if err != nil {
		log.Fatal("error")
	}

	skillsallnew := make([]Skill, len(result))
	for i := range result {
		skillsallnew[i] = Skill{Id: result[i].Id, Name: result[i].Name[0]}
	}

	return skillsallnew, nil
}
func (r *SkillRepo) GetByIds(ids []bson.ObjectId, lang string) (SkillCollection, error) {
	result := SkillCollection{[]SkillTemp{}}
	var err error
	fmt.Println("SKILL IDS....")
	fmt.Printf("%+v\n", ids)
	/* oids := make([]bson.ObjectId, len(ids))
	for i := range ids {
		oids[i] = bson.ObjectId(ids[i])
	} */
	if lang == "es" {
		err = r.Coll.Find(bson.M{"_id": bson.M{"$in": ids}}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result.Data)
	} else {
		err = r.Coll.Find(bson.M{"_id": bson.M{"$in": ids}}).Select(bson.M{"name": bson.M{"$slice": []int{0, 1}}}).All(&result.Data)
	}

	fmt.Println("REEEEEEEEEEEEEEEEEEEEEEEESSSSSULT")
	fmt.Printf("%v\n", result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SkillCategoryRepo) GetById(id bson.ObjectId, lang string) (SkillCategoryResource, error) {
	result := SkillCategoryResource{}
	var err error
	// err := r.Coll.FindId(id).One(&result.Data)
	if lang == "es" {
	err = r.Coll.FindId(id).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).One(&result.Data)
	} else {
	err = r.Coll.FindId(id).Select(bson.M{"name": bson.M{"$slice": []int{0, 1}}}).One(&result.Data)
	}
	//err := r.Coll.Find(bson.M{"id": id}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).one(&result)
	fmt.Println("ACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Printf("%+v\n", result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SkillCategoryRepo) All(lang string) (SkillCategoryCollection, error) {
	result := SkillTempCategoryCollection{[]SkillTempCategory{}}
	var err error
	if lang == "es" {
		err = r.Coll.Find(bson.M{}).Select(bson.M{"name": bson.M{"$slice": []int{1, 1}}}).All(&result.Data)
	} else {
		err = r.Coll.Find(bson.M{}).Select(bson.M{"name": bson.M{"$slice": []int{0, 1}}}).All(&result.Data)
        }

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
