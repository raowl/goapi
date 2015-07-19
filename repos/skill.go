package repos

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"time"
)

type SkillCategory struct {
	Id     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name,omitempty"`
	Skills []Skill       `json:"skills,omitempty"`
}

type Skill struct {
	Id   bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name,omitempty"`
}

type SkillCategoryCollection struct {
	Data []SkillCategory `json:"data"`
}

type SkillCategoryResource struct {
	Data SkillCategory `json:"data"`
}

type SkillCategoryRepo struct {
	Coll *mgo.Collection
}

func (r *SkillCategoryRepo) All() (SkillCategoryCollection, error) {
	result := SkillCategoryCollection{[]SkillCategory{}}
	err := r.Coll.Find(nil).All(&result.Data)
	if err != nil {
		return result, err
	}

	fmt.Println(result)

	return result, nil
}

/*
func (r *MarkerRepo) Find(id string) (MarkerResource, error) {
	result := MarkerResource{}
	err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *MarkerRepo) Create(marker *Marker) error {
	id := bson.NewObjectId()
	_, err := r.Coll.UpsertId(id, marker)
	if err != nil {
		return err
	}

	marker.Id = id

	return nil
}

func (r *MarkerRepo) Update(marker *Marker) error {
	// no need, Update Id below does updating over that id
	//currentMarker := MarkerResource{}
	//err := r.Coll.FindId(marker.Id).One(&currentMarker.Data)
	//if err != nil {
	//	return err
	//}

	// TODO: update all fields merge data:
	//http://stackoverflow.com/questions/18926303/iterate-through-a-struct-in-go
	//currentMarker.Data.CheckIns = marker.CheckIns

	//err := r.Coll.UpdateId(marker.Id, marker)
	//	err := r.Coll.UpdateId(marker.Id, bson.M{"$push": bson.M{"checkins": bson.M{"$each": marker.CheckIns}}}, true)
	err := r.Coll.UpdateId(marker.Id, bson.M{"$addToSet": bson.M{"checkins": bson.M{"$each": marker.CheckIns}}})
	//err := r.Coll.Update(bson.M{"id": marker.Id}, bson.M{"$push": bson.M{"checkins": bson.M{"$each": marker.CheckIns}}})
	//err = r.Coll.UpdateId(marker.Id, currentMarker.Data)
	if err != nil {
		return err
	}

	return nil
}

func (r *MarkerRepo) Delete(id string) error {
	err := r.Coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	return nil
}

*/
