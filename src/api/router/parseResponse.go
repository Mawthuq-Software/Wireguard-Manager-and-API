package router

import (
	"encoding/json"
	"errors"
	"net/http"
)

func parseResponse(req *http.Request, schema interface{}) error {
	var unmarshalErr *json.UnmarshalTypeError

	headerContentTtype := req.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		return errors.New("content type is not application/json")
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields() //throws error if uneeded JSON is added
	err := decoder.Decode(schema)   //decodes the incoming JSON into the struct
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			return errors.New("Bad Request. Wrong Type provided for field " + unmarshalErr.Field)
		} else {
			return errors.New("Bad Request " + err.Error())
		}
	}
	return nil
}
