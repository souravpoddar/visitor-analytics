package main

import (
	"net/http"

	"github.com/souravpoddar/visitor-analytics/internal/analytics"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	analyticsService := analytics.NewAnalyticsService()
	analyticsService.WireUpAnalytics(router)
	http.ListenAndServe(":8080", router)
}
