package api

import (
	"net/http"

	"github.com/yamixdev/go-final-sprint/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]any{
			"error": "не указан идентификатор",
		})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, map[string]any{})
}
