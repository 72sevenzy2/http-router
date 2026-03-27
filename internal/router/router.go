package router

import (
	"net/http"
)

type Router struct { // initializing the router struct to hold all the routes
	routes map[string]map[string]http.HandlerFunc
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
		r.routes[path] = make(map[string]http.HandlerFunc)
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
		ERROR(w, http.StatusMethodNotAllowed)
		return
	}

	handler(w, r) // if all pass, then execute
}