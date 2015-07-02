package main

import (
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"goapi/handlers" //controllers
	"goapi/repos"    //models
	"gopkg.in/mgo.v2"
	"net/http"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := handlers.AppContext{session.DB("test")}
	commonMiddleware := alice.New(context.ClearHandler, handlers.LoggingHandler, handlers.RecoverHandler, handlers.AcceptHandler)
	router := NewRouter()
	router.Get("/markers/:id", commonMiddleware.ThenFunc(appC.MarkerHandler))
	router.Put("/markers/:id", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.BodyHandler(repos.MarkerResource{})).ThenFunc(appC.UpdateMarkerHandler))
	router.Delete("/markers/:id", commonMiddleware.ThenFunc(appC.DeleteMarkerHandler))
	router.Get("/markers", commonMiddleware.Append(handlers.AuthHandler).ThenFunc(appC.MarkersHandler))
	router.Post("/markers", commonMiddleware.Append(handlers.AuthHandler, handlers.ContentTypeHandler, handlers.BodyHandler(repos.MarkerResource{})).ThenFunc(appC.CreateMarkerHandler))
	router.Post("/api/v1/user/auth", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.BodyHandler(repos.UserResource{})).ThenFunc(appC.AuthUserHandler))
	router.Post("/api/v1/user", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.BodyHandler(repos.UserResource{})).ThenFunc(appC.CreateUserHandler))
	http.ListenAndServe(":8080", router)
}
