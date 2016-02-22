package repos

import (
	"errors"
	"fmt"
	"github.com/raowl/goapi/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

/* type UserSkills struct {
	//Created   time.Time     `bson:"created" json:"created"`
	//CheckUser bson.ObjectId `bson:"user" json:"user"`
	SkillId string `bson:"skill" json:"skill"`
} */

type UserF struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string        `bson:"username,omitempty" json:"username,omitempty"`
	FirstName  string        `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName   string        `bson:"lastname,omitempty" json:"lastname,omitempty"`
	FacebookId string        `bson:"facebookid,omitempty" json:"facebookid,omitempty"`
	Image      string        `bson:"image,omitempty" json:"image,omitempty"`
}

type User struct {
	Id           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string        `bson:"username,omitempty" json:"username,omitempty"`
	FirstName    string        `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName     string        `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Email        string        `bson:"email,omitempty" json:"email,omitempty"`
	Verified     bool          `bson:"verified,omitempty" json:"verified,omitempty"`
	Facebook     bool          `bson:"facebook,omitempty" json:"facebook,omitempty"`
	ShowEmail    bool          `bson:"showEmail,omitempty" json:"showEmail,omitempty"`
	Password     string        `json:"password,omitempty"` //only used for parsing incoming json
	FacebookId   string        `bson:"facebookid,omitempty" json:"facebookid,omitempty"`
	Hash         string        `bson:"hash,omitempty"`
	Salt         string        `bson:"salt,omitempty"`
	Created      time.Time     `bson:"created,omitempty" json:"created,omitempty"`
	Updated      time.Time     `bson:"updated,omitempty" json:"updated,omitempty"`
	UrlToken     string        `bson:"urltoken,omitempty" json:"urltoken,omitempty"`
	About        string        `bson:"about,omitempty" json:"about,omitempty"`
	FollowInfo   []UserF       `bson:"followinfo,omitempty" json:"followinfo,omitempty"`
	FollowedInfo []UserF       `bson:"followedinfo,omitempty" json:"followedinfo,omitempty"`
	Image        string        `bson:"image,omitempty" json:"image,omitempty"`
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

func (r *UserRepo) GetFollowers(id bson.ObjectId) ([]UserF, error) {
	result := []UserF{}
	err := r.Coll.Find(bson.M{"following": id}).All(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) GetFByIds(ids []bson.ObjectId) ([]UserF, error) {
	result := []UserF{}
	err := r.Coll.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&result)
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

	if result.Data.Facebook {
		return result, nil
	} else {
		if passwordMatch(user.Password, result) {
			return result, nil
		}
	}

	return UserResource{}, errors.New("Passwd dont match")
}

func (r *UserRepo) UserAlreadyExists(username string) (bool, error) {
	count, err := r.Coll.Find(bson.M{"username": username}).Count()
	if err != nil {
		return false, err
	}

	fmt.Println(username)
	fmt.Println(count)
	if count > 0 {
		fmt.Println("user already exists")
		return true, nil
	} else {
		fmt.Println("user dont exists")
		return false, nil
	}
}
func (r *UserRepo) Find(id string) (UserResource, error) {
	result := UserResource{}
	err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	if err != nil {
		fmt.Println("inside repo error----------")
		return result, err
	}

	fmt.Println("inside repo ussssssssssssssssssssssssssssssser")
	fmt.Printf("%+v\n", result)

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

	user.Username = strings.ToLower(user.Username)

	fmt.Printf("%+v\n", user)

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
	fmt.Printf("USER DATA")
	fmt.Printf("%+v\n", user)
	err := r.Coll.UpdateId(user.Id, bson.M{"$addToSet": bson.M{"following": bson.M{"$each": user.Following}}})
	err = r.Coll.UpdateId(user.Id, bson.M{"$set": bson.M{"skills": user.Skills, "email": user.Email, "image": user.Image, "about": user.About}})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Unfollow(userId bson.ObjectId, unfollowId bson.ObjectId) error {
	fmt.Printf("Entered to Unfollow..........................")
	fmt.Println(userId)
	fmt.Println(unfollowId)
	fmt.Printf("%+v\n", unfollowId)
	err := r.Coll.UpdateId(userId, bson.M{"$pull": bson.M{"following": unfollowId}})
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
