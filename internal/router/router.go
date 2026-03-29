package router

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/response/helpers"
)

type Router struct { // initializing the router struct to hold all the routes
	routes map[string]map[string]http.HandlerFunc
}

type Middleware func(http.HandlerFunc) http.HandlerFunc // the middleware type (takes in the current handler and returns a new one)

// middleware (logger)
func Logger(next http.HandlerFunc) http.HandlerFunc { // takes in a our original handler of the request being made, and returns a new handler.
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request has started with method: ", r.Method)

		next(w, r) // runs the original handler

		fmt.Println("Request has ended")
	}
}

// auth middleware (no real authorization yet but will be adding this for a test case)

func Auth(next http.HandlerFunc, AuthKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != AuthKey {
			helpers.FAILED(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
	}
		next(w, r)
	}
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

	core := Auth(Logger(handler), "sia")
	core(w, r)
}
