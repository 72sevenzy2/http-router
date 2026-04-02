package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
	"github.com/72sevenzy2/http-router/internal/test-handler"
)

func main() {
	r := router.NewRouter()


	r.Handle(http.MethodGet, "/p", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("responded"))
	})

	r.Handle(http.MethodPost, "/user", handler.HiHandler(), router.Logger())

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
