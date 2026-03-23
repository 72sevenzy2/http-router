package main

import (
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

func main() {
	fmt.Println("hi")
}
