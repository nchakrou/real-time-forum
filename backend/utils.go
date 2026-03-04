package backend

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError sends a JSON response with a specific status code and error message.
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": "error",
	})
}
