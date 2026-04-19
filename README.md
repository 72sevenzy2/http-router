# An http router built from scratch, built ontop of the go stdlib (net/http).

# Features / context

- This router is built ontop the golang stdlib for managing http routes (http.handlerFunc).
- works likewise some of the bigger http routers such as chi/mux, you handle routes with this router using the "r.Handle(...)" func which you then pass in the necessary parameters (eg: path, path method, middleware, and the the actual handler needed to run the logic for that endpoint.)
- also supports route-specific middleware, by passing your desired middleware to the "r.Handle()" func as the 4th parameters, its important that i note you can add more than one middleware, but if you dont need route-specific middlewares, you can apply global middlewares, an example usage will be shown below.
- middlewares. I added about 5 middlewares in this project, the first being a simple logging middleware, second and third being authentication, (second is an BearerToken auth, which uses a token to verify a users credentials), (third being basicAuth auth, which involves a username and password to verify credentials.), and fourth being a recoverer middleware (if used make sure you use it FIRST before any other middleware), it prevents server crashes when a bug is occured. And lastly the timeout middleware, which the purpose of it is to cut off slow requests, for example when making an endpoint and a user goes to it and for some reason its slow or takes too much time to fetch the necessary data for the page to render, the timeout middleware will cancel the request, though it can also cancel due to other factors aswell (not just the particular website being slow itself). The timeout middleware takes in one paramters, which is the number of seconds the request should last before being cancelled.

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
> That of course is a simple way to use it though;

Example usage with all the middleware:

```

package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Use(router.Recoverer()) // recoverer middleware always goes first
	r.Use(router.Logger())
	r.Use(router.BearerAuth("secretKey")) // can be any token (which has to be a string)
	r.Use(router.BasicAuth("user1", "password1234")) // parameters username and password need to be included when using.
	r.Use(router.Timeout(5)) // can be any time (its in seconds) depending on how long you want the time limit on every request.

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
