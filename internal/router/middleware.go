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

// bearer auth middleware (this includes having a bearer token which will then be compared to the authkey )

func BearerAuth(AuthKey string) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authLab := r.Header.Get("Authorization") // grabbing the token

			token := strings.TrimPrefix(authLab, "Bearer ") // removing the "bearer " part of the token to then compare it to the authkey

			if token == authLab || token != AuthKey { // check if the authkey is matching
				helpers.Failed(w, http.StatusForbidden, "Invalid Token") // if not then throw a failed json response
				return
			}

			hf(w, r)
		}
	}
}

// basic auth middleware (this auth includes having a user and password inorder to access the endpoint)

func BasicAuth(user, password string) Middleware { // implements the middleware type which returns a handler
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			authUser, authPassword, ok := r.BasicAuth()

			if !ok || authUser != user || authPassword != password {
				helpers.Failed(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}

			hf(w, r)
		}
	}
}

// recoverer middleware (for preventing server crashes)

func Recoverer() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("caught: ", err)
				}
			}()

			hf(w, r)
		}
	}

}


// func to apply the middlewares
func (r *Router) ApplyMiddlewares(h http.HandlerFunc) http.HandlerFunc {
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		h = r.Middlewares[i](h)
	}

	return h
}


// Use func to use the middewares (also appending it to the Middlewares type in router struct
func (r *Router) Use(s Middleware) {
	r.Middlewares = append(r.Middlewares, s)
}