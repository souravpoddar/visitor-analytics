package analytics

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type AnalyticsService struct {
	VisitorStore VisitorStore
}
type VisitorStore struct {
	store map[string]map[string]bool
	mu    sync.Mutex
}

func NewVisitorStore() *VisitorStore {
	return &VisitorStore{
		store: make(map[string]map[string]bool),
	}
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		VisitorStore: *NewVisitorStore(),
	}
}

func (a *AnalyticsService) WireUpAnalytics(router *mux.Router) {
	router.HandleFunc("/track", a.VisitorStore.trackVisitor).Methods("GET")
	router.HandleFunc("/analytics", a.VisitorStore.getAnalytics).Methods("GET")
}
func (v *VisitorStore) RecordVisitor(url, visitorID string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, ok := v.store[url]; !ok {
		v.store[url] = make(map[string]bool)
	}
	v.store[url][visitorID] = true
}

func (v *VisitorStore) trackVisitor(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	visitorID := r.URL.Query().Get("visitorID")

	if url == "" || visitorID == "" {
		http.Error(w, "url and visitorID parameters are required", http.StatusBadRequest)
		return
	}

	v.RecordVisitor(url, visitorID)
}

func (v *VisitorStore) getAnalytics(w http.ResponseWriter, r *http.Request) {

	visitorMap := v.GetUniqueVisitors()
	for key, value := range visitorMap {
		kvw := bytes.NewBufferString(key + ":" + strconv.Itoa(value) + "\n")
		if _, err := kvw.WriteTo(w); err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func (v *VisitorStore) GetUniqueVisitors() map[string]int {
	v.mu.Lock()
	defer v.mu.Unlock()
	returnMap := make(map[string]int)
	for k := range v.store {
		returnMap[k] = len(v.store[k])
	}

	return returnMap
}
