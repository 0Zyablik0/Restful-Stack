package main

import (
	"github.com/gorilla/mux"
)

func (app *app) router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", app.MainHandler).Methods("GET")
	router.HandleFunc("/stack/top", app.StackTopHandler).Methods("GET")
	router.HandleFunc("/stack/pop", app.StackPopHandler).Methods("DELETE")
	router.HandleFunc("/stack/push", app.StackPushHandler).Methods("POST")
	router.HandleFunc("/stack/push/range", app.StackPushRangeHandler).Methods("POST")

	if app.config.enableLogging {
		router.Use(app.loggingMiddleware)
	}

	return router
}
