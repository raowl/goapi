package repos

import (
	"code.google.com/p/go.crypto/bcrypt"
	"errors"
	"fmt"
	"goapi/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

/* type UserSkills struct {
	//Created   time.Time     `bson:"created" json:"created"`
	//CheckUser bson.ObjectId `bson:"user" json:"user"`
	SkillId string `bson:"skill" json:"skill"`
} */
type User struct {
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username string        `bson:"username,omitempty" json:"username,omitempty"`
	Email    string        `bson:"email,omitempty" json:"email,omitempty"`
	Verified bool          `bson:"verified,omitempty" json:"verified,omitempty"`
	Password string        `json:"password,omitempty"` //only used for parsing incoming json
	Hash     string        `bson:"hash,omitempty"`
	Salt     string        `bson:"salt,omitempty"`
	Created  time.Time     `bson:"created,omitempty" json:"created,omitempty"`
	Updated  time.Time     `bson:"updated,omitempty" json:"updated,omitempty"`
	UrlToken string        `bson:"urltoken,omitempty" json:"urltoken,omitempty"`
	About    string        `bson:"about,omitempty" json:"about,omitempty"`
	Image    string        `bson:"image,omitempty" json:"image,omitempty"`
	// TODO: make bson object objects id better...
	Skills    []bson.ObjectId `bson:"skills,omitempty" json:"skills,omitempty"`
	Following []bson.ObjectId `bson:"following,omitempty" json:"following,omitempty"`
	//Skills string `bson:"skills" json:"skills"`
	//	UserMarker []Marker      `bson:"markers" json:"markers"`
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

func (r *UserRepo) GetByIds(ids []bson.ObjectId) (UserCollection, error) {
	result := UserCollection{[]User{}}
	err := r.Coll.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) Authenticate(user User) (UserResource, error) {
	// dry use next function
	result := UserResource{}
	err := r.Coll.Find(bson.M{"username": user.Username}).One(&result.Data)

	if err != nil {
		return UserResource{}, err
	}

	print(result.Data.Password)
	print(user.Password)

	if passwordMatch(user.Password, result) {
		return result, nil
	}

	return UserResource{}, errors.New("Passwd dont match")
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

/* func (r *UserRepo) Update(user *User) error {
	err := r.Coll.UpdateId(user.Id, bson.M{"$set": user})
	if err != nil {
		return err
	}

	return nil
} */

func (r *UserRepo) Update(user *User) error {
	fmt.Printf("Entered to Update")
	err := r.Coll.UpdateId(user.Id, bson.M{"$addToSet": bson.M{"following": bson.M{"$each": user.Following}}})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Unfollow(userId bson.ObjectId, unfollowId bson.ObjectId) error {
	fmt.Printf("Entered to Update")
	err := r.Coll.UpdateId(userId, bson.M{"$unset": bson.M{"following": unfollowId}})
	if err != nil {
		return err
	}

	return nil
}

/* func (r *UserRepo) Update(user *User) error {
	err := r.Coll.UpdateId(user.Id, user)
	if err != nil {
		return err
	}

	return nil
} */

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
