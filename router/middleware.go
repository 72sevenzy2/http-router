package router

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/72sevenzy2/json-parser/helpers"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc // the middleware type (takes in the current handler and returns a new one)

func Logger() Middleware { // returns the middleware type (which takes in a handler and returns a new one)
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now() // setting the current time (before the request has ended)
			fmt.Printf("Request has started with method: %s, in time: %s\n", r.Method, start)

			endTime := time.Since(start) // after the request has ended, in which we will print below
			fmt.Println("Request has ended:\n ", endTime)

			var buf bytes.Buffer                           // a buffer which will hold the r.Body
			lim := io.LimitReader(r.Body, 1024)            // limit size to 1 kilobyte of data to prevent large copis which can be time consuming
			r.Body = io.NopCloser(io.TeeReader(lim, &buf)) // using io.NopCloser as io.TeeReader does not implement io.ReadCloser.
			// io.TeeReader allows the current handler to read the request body data, whilst also allowing copying.
			hf(w, r) // calling the next function to continue to the next handler
			// by calling hf() before printing, we give time to the io Readers above to read the request body data.

			fmt.Println("Request body (1 kilobyte of body data):")
			fmt.Println(buf.String())

			fmt.Println("Request headers:", r.Header)
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
				return                                                   // exit the request
			}

			hf(w, r) // continue to next handler
		}
	}
}

// basic auth middleware (this auth includes having a user and password inorder to access the endpoint)

func BasicAuth(user, password string) Middleware { // implements the middleware type which returns a handler
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			authUser, authPassword, ok := r.BasicAuth() // extracting the user and password and if it exists (ok) from the r.BasicAuth() func, which is a built in method in go to do so, instead of manually parsing it ourselves.

			if !ok || authUser != user || authPassword != password { // run the necessary logic
				helpers.Failed(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}

			hf(w, r) // continue to next handler
		}
	}
}

// timeout middleware

func Timeout(seconds int) Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(seconds)*time.Second) // initialising timeout (in seconds)

			defer cancel() // cancelling at the end of the func (current handler)

			hf(w, r.WithContext(ctx)) // ServeHTTP(w, and "r" with the context 'ctx')
		}
	}
}

// recoverer middleware (for preventing server crashes)

func Recoverer() Middleware {
	return func(hf http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() { // catches any crashses and recovers the request, while printing the err in return.
				if err := recover(); err != nil {
					fmt.Println("caught: ", err)
				}
			}()

			hf(w, r) // next handler
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
