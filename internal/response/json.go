package response

import (
	"encoding/json"
	"net/http"
)

type JsonOptions struct { // this struct will be modified via the ConfigOpts func
	w      http.ResponseWriter
	Data   interface{}
	Status int
	Error  string
}

type ConfigOpts func(*JsonOptions) // any func which returns this type ONLY will use a pointer to the JsonOptions struct like used here

// status param func
func WithStatus(status int) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Status = status
	}
}

// with error param func
func WithError(msg string) ConfigOpts {
	return func(jo *JsonOptions) {
		jo.Error = msg
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
	Error  string      `json:"error"`
}

func JSON(w http.ResponseWriter, opts ...ConfigOpts) {
	// assigning the default values if i were to assign no params when calling the JSON func
	options := &JsonOptions{ // these options will be replaced if there were opts included when calling this func with the data in those opts
		w:      w,
		Status: http.StatusOK,
		Data:   nil,
		Error:  "",
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
		Error:  options.Error,
	} // initialising the response

	err := json.NewEncoder(options.w).Encode(response) // handling errors while encoding it aswell
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
