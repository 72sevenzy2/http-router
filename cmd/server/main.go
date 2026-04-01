package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Use(router.BasicAuth("user", "hi"))
	r.Use(router.Recoverer())
	r.Use(router.Logger())

	r.Handle(http.MethodGet, "/p", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("responded"))
	})

	go func() {
		r.Handle(http.MethodGet, "/hi", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello user"))
		})
	}()

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
