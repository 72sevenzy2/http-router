package helpers

import (
	"github.com/72sevenzy2/http-router/internal/response"
	"net/http"
)

func Ok(w http.ResponseWriter, data any) {
	response.JSON(w, response.WithData(data), response.WithStatus(http.StatusOK))
}

func Failed(w http.ResponseWriter, status int, msg string) {
	response.JSON(w, response.WithStatus(status), response.WithError(msg))
}
