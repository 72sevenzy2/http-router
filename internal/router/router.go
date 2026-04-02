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
		Routes: make(map[string]map[string]http.HandlerFunc), // initialising the map of map)
	} // which is just: "PATH": "...": "METHOD": ... (method can be either get, post, put, etc)
}

// adding routes, and assigning the method of the route aswell as the url to the handler which then is executed in the ServeHTTP func
func (r *Router) Handle(method, path string, handler http.HandlerFunc, mws ...Middleware) {
	if r.Routes[path] == nil { // checking if route endpoint itself doesnt exist before creating the path and assigning it to another map (map[string]http.HandlerFunc) which will be for the method map
		r.Routes[path] = make(map[string]http.HandlerFunc) // assign the path to the method type (GET, POST, PUT etc)
	}

	for i := len(mws) -1; i >= 0; i-- { // applying the middleware in reverse order
		handler = mws[i](handler) // assign the handler to the middleware (the middleware returns a new handler after performing its programmed task.)
	}  // assigning it to the index of [i] which is each middleware that gets passed in as the mws parameter, so it can be 1 or many.

	r.Routes[path][method] = handler // assign method to the handler (handler type is http.handlerFunc)
	// we're basically taking the path which will be something like "/hi": and the method name, or its type we can call it
	// for example:       "/hi":
	//                       "GET": "and some handler here, (in this case, it will be the http handlerfunc we used)"
}


// core routing logic for my router
func (s *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method, ok := s.Routes[r.URL.Path] // search for if the url path exists

	if !ok {
		http.NotFound(w, r) // check if path exists on the request being sent.
		return
	}

	handler, ok := method[r.Method] // assign the handler to the method type

	if !ok { // and check if the method of the path is valid.
		helpers.Failed(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	finalHandler := s.ApplyMiddlewares(handler)
	finalHandler(w, r)
}

