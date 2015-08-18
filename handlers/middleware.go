package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"goapi/config"
	"goapi/utils"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

var (
	PublicKey []byte
)

func init() {
	PublicKey, _ = ioutil.ReadFile(config.Info.PublicKey)
}

func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				utils.WriteError(w, utils.ErrInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

//func (c *appContext) authHandler(next http.Handler) http.Handler {
func AuthHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// authToken := r.Header.Get("Authorization")
		// print(authToken)
		/*
		   We could alternatively make getUser a method of appContext to be cleaner. Or wrap *sql.DB
		   in a custom struct and add getUser as a method of this custom struct so we could call it
		   simply with c.db.getUser(token).
		*/
		//user, err := getUser(c.db, authToken)
		token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
			return PublicKey, nil
		})

		if err != nil {
			utils.WriteError(w, utils.ErrInternalServer)
			return
		}
		if token.Valid == false {
			fmt.Printf("token malo")
			//YAY!
			utils.WriteError(w, utils.ErrInternalServer)
			return
		}

		fmt.Printf("token valido\n")
		//fmt.Printf("%+v\n", token)

		/*	fmt.Printf("user id")
			fmt.Printf("\n%s\n", token.Claims["userid"])
		*/
		/* if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		} */

		context.Set(r, "userid", token.Claims["userid"].(string))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func AcceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.api+json" {
			utils.WriteError(w, utils.ErrNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ContentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/vnd.api+json" {
			utils.WriteError(w, utils.ErrUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func BodyHandler(v interface{}) func(http.Handler) http.Handler {
	fmt.Printf("Inside BodyHandler\n")
	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)

			if err != nil {
				utils.WriteError(w, utils.ErrBadRequest)
				return
			}

			if next != nil {
				context.Set(r, "body", val)
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}

	return m
}
