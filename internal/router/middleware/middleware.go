package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/72sevenzy2/http-router/internal/response/helpers"
	"github.com/72sevenzy2/http-router/internal/router"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc // the middleware type (takes in the current handler and returns a new one)

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
func (r router.Router) Use(s Middleware) {
	r.Middlewares = append(r.Middlewares, s)
}

// func to apply the middlewares
func (r *Router) applyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}
