package api

import (
	"errors"
	"github.com/yamixdev/go-final-sprint/pkg/db"
	"log"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		next, nextErr := NextDate(time.Now(), task.Date, task.Repeat)
		if nextErr != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": nextErr.Error(),
			})
			return
		}

		err = db.UpdateDate(next, id)
	}

	if err != nil {
		log.Println("doneTask error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "ошибка при обновлении задачи",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
