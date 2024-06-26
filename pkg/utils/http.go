package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadBody(w http.ResponseWriter, r *http.Request, strct interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return err
	}

	defer r.Body.Close()

	if err := json.Unmarshal(body, strct); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return err
	}

	return nil
}
