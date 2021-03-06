package main

import (
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"github.com/raowl/goapi/handlers" //controllers
	"github.com/raowl/goapi/repos"    //models
	"gopkg.in/mgo.v2"
	"net/http"
)

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := handlers.AppContext{session.DB("test")}
	commonMiddleware := alice.New(context.ClearHandler, handlers.LoggingHandler, handlers.RecoverHandler, handlers.AcceptHandler)
	router := NewRouter()
	router.Get("/markers/:id", commonMiddleware.ThenFunc(appC.MarkerHandler))
	router.Get("/markers/:id/users", commonMiddleware.ThenFunc(appC.MarkerWithUsersHandler))
	router.Put("/markers/:id", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.AuthHandler, handlers.BodyHandler(repos.MarkerResource{})).ThenFunc(appC.UpdateMarkerHandler))
	router.Delete("/markers/:id", commonMiddleware.ThenFunc(appC.DeleteMarkerHandler))
	router.Get("/markersnear/:lat/:lng/:km", commonMiddleware.ThenFunc(appC.MarkersHandler))
	router.Post("/markers", commonMiddleware.Append(handlers.AuthHandler, handlers.ContentTypeHandler, handlers.BodyHandler(repos.MarkerResource{})).ThenFunc(appC.CreateMarkerHandler))
	router.Get("/skills", commonMiddleware.ThenFunc(appC.SkillsHandler))
	router.Post("/api/v1/user/auth", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.BodyHandler(repos.UserResource{})).ThenFunc(appC.AuthUserHandler))
	router.Post("/api/v1/user", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.BodyHandler(repos.UserResource{})).ThenFunc(appC.CreateUserHandler))
	router.Put("/api/v1/user", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.AuthHandler, handlers.BodyHandler(repos.UserResource{})).ThenFunc(appC.UpdateUserHandler))
	router.Get("/api/v1/user/:id", commonMiddleware.Append(handlers.AuthHandler).ThenFunc(appC.UserHandler))
	//router.Get("/api/v1/user/username/:username", commonMiddleware.Append(handlers.AuthHandler).ThenFunc(appC.UserHandler))
	router.Get("/api/v1/user/:id/unfollow", commonMiddleware.Append(handlers.AuthHandler).ThenFunc(appC.UserUnfollowHandler))
	router.Get("/api/v1/user/:id/skills", commonMiddleware.ThenFunc(appC.UserWithSkillsHandler))
	http.ListenAndServe(":8080", router)
}
