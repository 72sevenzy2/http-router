package handler

import (
	"encoding/json"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/model"
	"github.com/72sevenzy2/http-router/internal/response"
)

func hiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i entity.Entity

		err := json.NewDecoder(r.Body).Decode(&i)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest)
		}

		response.JSON(w, response.WithStatus(http.StatusOK), response.WithData(map[string]interface{}{
			"User": i.User,
			"Id":   i.Id,
		}))
	}
}
