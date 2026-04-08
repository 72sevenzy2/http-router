package helpers

import (
	"github.com/72sevenzy2/http-router/internal/response"
	"net/http"
)

// quick method for responses without needing to throw in any paramters (except the data needed and the http.responseWriter) as options for if you were to actually use the JSON(...) func.
func Ok(w http.ResponseWriter, data any) {
	response.JSON(w, response.WithData(data), response.WithStatus(http.StatusOK))
}

// a method for failed aswell. (though status and msg is both needed.)
func Failed(w http.ResponseWriter, status int, msg string) {
	response.JSON(w, response.WithStatus(status), response.WithError(msg))
}
