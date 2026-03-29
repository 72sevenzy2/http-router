package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/72sevenzy2/http-router/internal/response/helpers"
)

type Router struct { // initializing the router struct to hold all the routes
	routes      map[string]map[string]http.HandlerFunc
	Middlewares []Middleware // storing our middlewares here (type is our Middleware function type)
}

type Middleware func(http.HandlerFunc) http.HandlerFunc // the middleware type (takes in the current handler and returns a new one)

// middleware (logger)
func Logger() Middleware { // returns the middleware type (which takes in a handler and returns a new one)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Request has started with method: ", r.Method)

			hf(w, r)

			fmt.Println("Request has ended")
		}
	}
}

// auth middleware (no real authorization yet but will be adding this for a test case)

func Auth(AuthKey string) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authLab := r.Header.Get("Authorization") // grabbing the token

			token := strings.TrimPrefix(authLab, "Bearer ") // removing the "bearer " part of the token to then compare it to the authkey

			if token == authLab || token != AuthKey { // check if the authkey is matching
				helpers.FAILED(w, http.StatusForbidden, "Invalid Token") // if not then throw a failed json response
				return
			}

			hf(w, r)
		}
	}
}

// Use func to use the middewares (also appending it to the Middlewares type in router struct
func (r *Router) Use(s Middleware) {
	r.Middlewares = append(r.Middlewares, s)
}

// func to apply the middlewares
func (r *Router) applyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}

func NewRouter() *Router {
	// contructing the router upon the func being called
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// adding routes, and assigning the method of the route aswell as the url to the handler which then is executed in the ServeHTTP func
func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc) // assign the path to the method type (GET, POST, PUT etc)
	}

	r.routes[path][method] = handler // assign both url and method to the handler (handler type is http.handlerFunc)
	// we're basically taking the path which will be something like "/hi": and the method name, or its type we can call it
	// for example:       "/hi":
	//                       "GET": "and some handler here, (in this case, it will be the http handlerfunc we used)"
}

// core routing logic for my router
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method, ok := s.routes[r.URL.Path] // search for if the url path exists

	if !ok {
		http.NotFound(w, r)
		return
	}

	handler, ok := method[r.Method] // assign the handler to the method type

	if !ok {
		helpers.FAILED(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	wra := s.applyMiddlewares(handler)
	wra(w, r)

}
