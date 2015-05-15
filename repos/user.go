package repos

import (
	"coworkingApi/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Lng      string        `json:"lng"`
	ParentId string        `bson:"parentId" json:"parentId"` //unused. change string to bson.ObjectId
	Name     string        `bson:"name" json:"name"`
	Email    string        `bson:"email" json:"email"`
	Verified bool          `bson:"verified" json:"verified"`
	Password string        `json:"password"` //only used for parsing incoming json
	Hash     string        `bson:"hash"`
	Salt     string        `bson:"salt"`
	Created  time.Time     `bson:"created" json:"created"`
	Updated  time.Time     `bson:"updated" json:"updated"`
	UrlToken string        `bson:"urltoken" json:"urltoken"`
}

type UserCollection struct {
	Data []User `json:"data"`
}

type UserResource struct {
	Data User `json:"data"`
}

type UserRepo struct {
	coll *mgo.Collection
}

func (r *UserRepo) All() (UserCollection, error) {
	result := UserCollection{[]User{}}
	err := r.coll.Find(nil).All(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) Find(id string) (UserResource, error) {
	result := UserResource{}
	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) Create(user *User) error {
	id := bson.NewObjectId()
	hash, salt, err := utils.GenerateHashAndSalt(user.Password)
	if err != nil {
		return err
	}

	user.Hash = hash
	user.Salt = salt

	_, err := r.coll.UpsertId(id, user)
	if err != nil {
		return err
	}

	user.Id = id

	return nil
}

func (r *UserRepo) Update(user *User) error {
	err := r.coll.UpdateId(user.Id, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Delete(id string) error {
	err := r.coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	return nil
}
