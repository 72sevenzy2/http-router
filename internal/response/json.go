package response 

import (
	"net/http"
	"encoding/json"
)

type JsonOptions struct { // this struct will be modified via the ConfigOpts func
	w      http.ResponseWriter
	Data   interface{}
	Status int
}

type ConfigOpts func(*JsonOptions) // any func which returns this type ONLY will use a pointer to the JsonOptions struct like used here

// status param func
func WithStatus(status int) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Status = status
	}
}

// param func for extracting data
func ExtractReq(data interface{}) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Data = data
	}
}

// data param func
func WithData(data interface{}) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Data = data
	}
}

// the response format we will be using, will be making another struct for so
type Response struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

func JSON(w http.ResponseWriter, opts ...ConfigOpts) {
	// assigning the default values if i were to assign no params when calling the JSON func
	options := &JsonOptions{ // these options will be replaced if there were opts included when calling this func with the data in those opts
		w:      w,
		Status: http.StatusOK,
		Data:   nil,
	}

	// initialising each opt to the appropriate param func
	for _, opt := range opts {
		opt(options) // each opt is a func that takes a pointer to the JsonOptions struct
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(options.Status)

	response := &Response{
		Data:   options.Data,
		Status: options.Status,
	} // initialising the response

		err := json.NewEncoder(options.w).Encode(response) // handling errors while encoding it aswell
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
}

// error json response helper
func ERROR(w http.ResponseWriter, status int) {

	JSON(w, WithStatus(status), WithData(map[string]string{
		"error": http.StatusText(status),
	}))
}