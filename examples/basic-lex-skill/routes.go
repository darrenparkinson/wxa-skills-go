package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) routes() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/ping", http.HandlerFunc(ping))
	mainRouter.Handle("/metrics", promhttp.Handler())
	mainRouter.HandleFunc("/", app.handleSkills).Methods(http.MethodPost)
	mainRouter.HandleFunc("/", app.handleHealthCheck).Methods(http.MethodGet)
	return app.metrics(app.recoverPanic(app.logRequest(secureHeaders(mainRouter))))
}
