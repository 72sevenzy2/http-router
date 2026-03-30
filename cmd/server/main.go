package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/router"
	"github.com/72sevenzy2/http-router/internal/test-handler"
)

func main() {
	r := router.NewRouter()

	r.Use(router.Logger(), )
	r.Use(router.Auth("siam"))

	r.Handle(http.MethodPost, "/p", handler.HiHandler())

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err);
	}
}
