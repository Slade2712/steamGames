package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJSONResponse(response http.ResponseWriter, status int, payload any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	encodeError := json.NewEncoder(response).Encode(payload)

	if encodeError != nil {
		fmt.Fprintln(response, "Internal error : ", encodeError)
	}
}
