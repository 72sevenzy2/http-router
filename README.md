# An http router built from scratch, built ontop of the go stdlib (net/http).
Built for my own educational reasons and use.

# Features / context

- This router is built ontop the golang stdlib for managing http routes (http.handlerFunc).
- works likewise some of the bigger http routers such as chi/mux, you handle routes with this router using the "r.Handle(...)" func which you then pass in the necessary parameters (eg: path, path method, middleware, and the the actual handler needed to run the logic for that endpoint.)
- also supports route-specific middleware, by passing your desired middleware to the "r.Handle()" func as the 4th parameters, its important that i note you can add more than one middleware, but if you dont need route-specific middlewares, you can apply global middlewares, an example usage will be shown below.
- also includes middlewares, such as basicAuth, recoverer mw, logger mw, a bearerAuth mw, and a timeout middleware, example usages for all will be shown below.

# Example usage:

```
 package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Handle(http.MethodGet, "/resp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("responded"))
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
> That of course is a example with no middlewares attached yet.

Example usage with all the middlewares:

```

package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Use(router.Recoverer()) // recoverer middleware always goes first, prevents server crashes when a bug has occured.
	r.Use(router.Logger()) // standard logging middleware (to view request details.)
	r.Use(router.BearerAuth("secretKey")) // can be any token (which has to be a string),
	r.Use(router.BasicAuth("user1", "password1234")) // parameters username and password need to be included when using.
	r.Use(router.Timeout(5)) // can be any time (its in seconds) depending on how long you want the time limit on every request.
	// the Timeout mw is used for prevent slow requests by setting a timeout in which the request should last.

	r.Handle(http.MethodGet, "/resp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("responded"))
	})

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
> But when using an authentication middleware, make sure to choose 1, either BearerAuth or BasicAuth.


Example usage with route-specific middleware:

```
package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Handle(http.MethodGet, "/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}, router.Recoverer(), router.Logger()) // you can do route-specific middleware(s) like this (can be 1 or many).

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

```
