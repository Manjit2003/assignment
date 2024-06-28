package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type HTTPReponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ReadBody(w http.ResponseWriter, r *http.Request, strct interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendResponse(w, http.StatusBadRequest, HTTPReponse{
			Error:   true,
			Message: "failed to read request body",
		})
		return err
	}

	defer r.Body.Close()

	if err := json.Unmarshal(body, strct); err != nil {
		SendResponse(w, http.StatusBadRequest, HTTPReponse{
			Error:   true,
			Message: "unable to decode json body",
		})
		return err
	}

	return nil
}

func SendResponse(w http.ResponseWriter, status int, res HTTPReponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
