package handler

import (
	"encoding/json"
	"net/http"

	"github.com/72sevenzy2/http-router/internal/model"
	"github.com/72sevenzy2/http-router/internal/response"
)

func hiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i entity.Entity // initilising the entity

		err := json.NewDecoder(r.Body).Decode(&i) // decoding the body to get the data we want
		if err != nil {                           // if there is no data which we needed in the body, throw an json error msg
			response.ERROR(w, http.StatusBadRequest)
			return
		}
		// respond with json returning the users user and the users id
		response.JSON(w, response.WithStatus(http.StatusOK), response.WithData(map[string]interface{}{
			"User": i.User,
			"Id":   i.Id,
		}))
	}
}
