package api

import (
	"net/http"
	"time"

	"github.com/yamixdev/go-final-sprint/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]any{
			"error": "не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		next, nextErr := NextDate(time.Now(), task.Date, task.Repeat)
		if nextErr != nil {
			writeJSON(w, map[string]any{
				"error": nextErr.Error(),
			})
			return
		}

		err = db.UpdateDate(next, id)
	}

	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, map[string]any{})
}
