package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJSONSuccess(w http.ResponseWriter, i interface{}, status int) error {
	w.WriteHeader(status)

	// if the data was passed in as an empty string then don't write a response
	if i == "" {
		return nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("Unable to marshal to json: %v", err)
	}

	w.Write(b)

	return nil
}

func writeJSONFail(w http.ResponseWriter, code int, s string) {
	w.WriteHeader(code)
	w.Write([]byte(s))
}
