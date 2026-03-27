package main

import (
	"fmt"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/response"
	"github.com/72sevenzy2/http-router/internal/router"
)

func main() {
	r := router.NewRouter()

	r.Handle(http.MethodGet, "/hi", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, response.WithStatus(http.StatusOK), response.WithData(map[string]string{
			"message": "hello",
		}))
	})

	fmt.Println("server running on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
