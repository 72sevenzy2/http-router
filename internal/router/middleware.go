package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/72sevenzy2/http-router/internal/response/helpers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc // the middleware type (takes in the current handler and returns a new one)

func Logger() Middleware { // returns the middleware type (which takes in a handler and returns a new one)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			fmt.Printf("Request has started with method: %s, in time: %s\n", r.Method, start)

			hf(w, r)

			endTime := time.Since(start)
			fmt.Println("Request has ended:\n ", endTime)
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
func (r *Router )Use(s Middleware) {
	r.Middlewares = append(r.Middlewares, s)
}


// func to apply the middlewares
func (r *Router) ApplyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}
