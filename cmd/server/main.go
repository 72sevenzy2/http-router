package main

import (
	"encoding/json"
	"fmt"
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
		ERROR(w, http.StatusMethodNotAllowed)
		return
	}

	handler(w, r) // if all pass, then execute
}

// for json responding so the server can read the incoming info

// configuration structs for json helpers below:

// using function options to have optional parameters for main JSON helper func (this pattern is also known as the functional params pattern)
// this pattern is also used by bigger http routers than mine like the chi or mux router.
type JsonOptions struct { // this struct will be pointed from the ConfigOpts func
	w      http.ResponseWriter
	Data   interface{}
	Status int
}

type ConfigOpts func(*JsonOptions) // any func which returns this type ONLY will use a pointer to the JsonOptions struct like used here

// status param func
func WithStatus(status int) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Status = status
	}
}

// data param func
func WithData(data interface{}) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Data = data
	}
}

// the response format we will be using, will be making another struct for so
type Response struct {
	Data   interface{}
	Status int
}

func JSON(w http.ResponseWriter, opts ...ConfigOpts) {
	// assigning the default values if i were to assign no params when calling the JSON func
	options := &JsonOptions{ // these options will be replaced if there were opts included when calling this func with the data in those opts
		w:      w,
		Status: http.StatusOK,
		Data:   nil,
	}

	// initialising each opt to the appropriate param func
	for _, opt := range opts {
		opt(options) // each opt is a func that takes a pointer to the JsonOptions struct
	}

	w.WriteHeader(options.Status)
	w.Header().Set("Content-Type", "application/json")

	response := &Response{
		Data:   options.Data,
		Status: options.Status,
	}

	if options.Data != nil {
		json.NewEncoder(options.w).Encode(response)
	}
}

// error json response helper
func ERROR(w http.ResponseWriter, status int) {

	JSON(w, WithStatus(status), WithData(map[string]string{
		"error": "bad request",
	}))
}

func main() {
	r := NewRouter()

	r.Handle("GET", "/hi", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, WithStatus(http.StatusOK), WithData(map[string]string{
			"message": "sup",
		}))
	})
	fmt.Println("api running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
