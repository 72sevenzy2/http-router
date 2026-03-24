package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}

	r.routes[path][method] = handler
}

// core routing logic for my router
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method, ok := s.routes[r.URL.Path] // search for if the url path exists

	if !ok {
		http.NotFound(w, r)
		return
	}

	handler, ok := method[r.Method] // and check if the method is valid

	if !ok {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handler(w, r) // if all pass, then execute
}

// for json responding so the server can read the incoming info

func JSON(w http.ResponseWriter, status int, data string) {
	// creating a temp response format for now
	resp := map[string]any{ 
		"message": data,
		"status":  status,
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	r := NewRouter()

	r.Handle("GET", "/hi", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, "working api")
	})
	fmt.Println("api running")
	http.ListenAndServe(":8080", r)
}
