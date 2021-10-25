package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() != "/metrics" {
			app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_requests_received_total",
		Help: "The total number of processed requests",
	})
	totalResponsesSent := promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_responses_sent_total",
		Help: "The total number of processed responses",
	})
	totalProcessingTimeMicroseconds := promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_processing_time_microseconds_total",
		Help: "The total amount of processing time in microseconds",
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalRequestsReceived.Inc()
		next.ServeHTTP(w, r)
		totalResponsesSent.Inc()
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(float64(duration))
	})
}
