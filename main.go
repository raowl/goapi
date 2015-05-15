package main

import (
	"coworkingApi/handlers" //controllers
	"coworkingApi/repos" //models
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
	"net/http"
)

//TODO get this like tns-restful-json-api/v8 from a routes file with only the array option and routes
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
	router.Get("/markers", commonMiddleware.ThenFunc(appC.MarkersHandler))
	router.Post("/markers", commonMiddleware.Append(handlers.ContentTypeHandler, handlers.BodyHandler(repos.MarkerResource{})).ThenFunc(appC.CreateMarkerHandler))
	http.ListenAndServe(":8080", router)
}
