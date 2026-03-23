package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Router struct {
	routes map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]http.HandlerFunc),
	}
}

func (r *Router) Handle(path string, handler http.HandlerFunc) {
	r.routes[path] = handler
}

func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := s.routes[r.URL.Path]

	if !ok {
		http.NotFound(w, r)
		return
	}

	handler(w, r)
}

// for json responding so the server can read the incoming info

func JSON(w http.ResponseWriter, status int, data string) {
	resp := map[string]interface{}{
		"message": data,
		"status":  status,
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	r := NewRouter()

	r.Handle("/hi", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, "test")
	})
	fmt.Println("api running")
	http.ListenAndServe(":8080", r)
}
