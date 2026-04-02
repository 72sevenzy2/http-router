package router

import (
	"net/http"

	"github.com/72sevenzy2/http-router/internal/response/helpers"
)

type Router struct { // initializing the router struct to hold all the routes
	Routes      map[string]map[string]http.HandlerFunc
	Middlewares []Middleware // storing our middlewares here (type is our Middleware function type)
}

func NewRouter() *Router {
	// contructing the router upon the func being called
	return &Router{
		Routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// adding routes, and assigning the method of the route aswell as the url to the handler which then is executed in the ServeHTTP func
func (r *Router) Handle(method, path string, handler http.HandlerFunc, mws ...Middleware) {
	if r.Routes[path] == nil {
		r.Routes[path] = make(map[string]http.HandlerFunc) // assign the path to the method type (GET, POST, PUT etc)
	}

	for i := len(mws) -1; i >= 0; i-- {
		handler = mws[i](handler)
	} 

	r.Routes[path][method] = handler // assign method to the handler (handler type is http.handlerFunc)
	// we're basically taking the path which will be something like "/hi": and the method name, or its type we can call it
	// for example:       "/hi":
	//                       "GET": "and some handler here, (in this case, it will be the http handlerfunc we used)"
}


// core routing logic for my router
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method, ok := s.Routes[r.URL.Path] // search for if the url path exists

	if !ok {
		http.NotFound(w, r)
		return
	}

	handler, ok := method[r.Method] // assign the handler to the method type

	if !ok {
		helpers.Failed(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	finalHandler := s.ApplyMiddlewares(handler)
	finalHandler(w, r)
}

