package repos

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Marker struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Coordinates [2]float64    `bson:"coordinates"`
	Name        string        `json:"name"`
	Address     string        `json:"address"`
	Website     string        `json:"website"`
}

type MarkerCollection struct {
	Data []Marker `json:"data"`
}

type MarkerResource struct {
	Data Marker `json:"data"`
}

type MarkerRepo struct {
	Coll *mgo.Collection
}

func (r *MarkerRepo) All() (MarkerCollection, error) {
	result := MarkerCollection{[]Marker{}}
	err := r.Coll.Find(nil).All(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

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
	err := r.Coll.UpdateId(marker.Id, marker)
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
