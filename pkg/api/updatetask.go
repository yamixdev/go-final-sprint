package api

import (
	"encoding/json"
	"net/http"

	"github.com/yamixdev/go-final-sprint/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]any{
			"error": "не указан заголовок задачи",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, map[string]any{})
}
