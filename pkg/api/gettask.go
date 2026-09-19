package api

import (
	"errors"
	"github.com/yamixdev/go-final-sprint/pkg/db"
	"log"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{
				"error": err.Error(),
			})
			return
		}
		log.Println("db.GetTask error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "ошибка при получении задачи",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}
