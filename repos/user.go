package repos

import (
	"code.google.com/p/go.crypto/bcrypt"
	"goapi/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username string        `bson:"username" json:"username"`
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
	Coll *mgo.Collection
}

func (r *UserRepo) All() (UserCollection, error) {
	result := UserCollection{[]User{}}
	err := r.Coll.Find(nil).All(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) Authenticate(user User) (bool, UserResource, error) {
	// dry use next function
	result := UserResource{}
	err := r.Coll.Find(bson.M{"username": user.Username}).One(&result.Data)

	if err != nil {
		return false, UserResource{}, err
	}

	print(result.Data.Password)
	print(user.Password)

	if passwordMatch(user.Password, result) {
		return true, result, nil
	}

	return false, UserResource{}, nil
}

func (r *UserRepo) Find(id string) (UserResource, error) {
	result := UserResource{}
	err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) Create(user *User) error {
	id := bson.NewObjectId()
	hash, salt := utils.GenerateHashAndSalt(user.Password)
	//if err != nil {
	//	return err
	//}

	user.Hash = hash
	user.Salt = salt

	_, err := r.Coll.UpsertId(id, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Update(user *User) error {
	err := r.Coll.UpdateId(user.Id, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Delete(id string) error {
	err := r.Coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	return nil
}

func passwordMatch(guess string, suser UserResource) bool {

	salted_guess := utils.Combine(suser.Data.Salt, guess)

	return bcrypt.CompareHashAndPassword([]byte(suser.Data.Hash), []byte(salted_guess)) == nil
}
