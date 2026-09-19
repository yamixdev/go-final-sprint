package api

import (
	"encoding/json"
	"errors"
	"github.com/yamixdev/go-final-sprint/pkg/db"
	"log"
	"net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "не указан заголовок задачи",
		})
		return
	}

	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "не указан идентификатор",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{
				"error": err.Error(),
			})
			return
		}
		log.Println("db.UpdateTask error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "ошибка при обновлении задачи",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
