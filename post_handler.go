package gohandle

import (
	"encoding/json"
	"net/http"
)

type PostHandler[T any] struct {
	Pattern    string
	HandleFunc func(T) error
}

func (ph *PostHandler[T]) GetPattern() string {
	return ph.Pattern
}

func (ph *PostHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only support POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ph.HandleFunc != nil {
		// Read data from the request body
		var data T
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Call the user-defined handler function
		if err := ph.HandleFunc(data); err != nil {
			http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Set the success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		// TODO: What to put here
		// "msg":    fmt.Sprintf("Expense for %s saved!", exp.Category),
	})
}
