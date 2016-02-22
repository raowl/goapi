package repos

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	//"time"
)

type Geolocation struct {
	Type        string     `bson:"type" json:"type"`
	Coordinates [2]float64 `bson:"coordinates" json:"coordinates"`
}

type CheckIn struct {
	//Created   time.Time     `bson:"created" json:"created"`
	//CheckUser bson.ObjectId `bson:"user" json:"user"`
	CheckUser bson.ObjectId `bson:"user" json:"user"`
}

type Marker struct {
	Id          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Geolocation Geolocation   `json:"geolocation,omitempty"`
	Name        string        `json:"name,omitempty"`
	Address     string        `json:"address,omitempty"`
	Website     string        `json:"website,omitempty"`
	Kind        int           `json:"kind,omitempty"`
	Else        string        `json:"else,omitempty"`
	Author      string        `json:"author,omitempty"`
	CheckIns    []CheckIn     `json:"checkins,omitempty"`
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

func (r *MarkerRepo) All(lat string, lng string, km string) (MarkerCollection, error) {

	latf, err := strconv.ParseFloat(lat, 64)
	lngf, err := strconv.ParseFloat(lng, 64)
	//kmf, err := strconv.ParseFloat(km, 64)
	kmf, err := strconv.ParseInt(km, 10, 64)
	result := MarkerCollection{[]Marker{}}
	fmt.Println("gobend")
	fmt.Println(latf)
	fmt.Println(lngf)
	fmt.Println(kmf)
	//err = r.Coll.Find(bson.M{"geolocation.coordinates": bson.M{"$near": []float64{latf, lngf}, "$maxDistance": kmf / 111.12}}).All(&result.Data)
	//err = r.Coll.Find(bson.M{"geolocation.coordinates": bson.M{"$near": bson.M{"$geometry": {"type":"Point", "coordinates": []float64{latf, lngf}}},"$maxDistance": kmf * 100000}}).All(&result.Data)
	err = r.Coll.Find(bson.M{"geolocation.coordinates": bson.M{"$near": bson.M{"$geometry": bson.M{"type":"Point", "coordinates": []float64{latf, lngf}},"$maxDistance": kmf * 1000}}}).All(&result.Data)
	if err != nil {
		return result, err
	}

	fmt.Println(result)

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
