package api

import (
	"github.com/yamixdev/go-final-sprint/pkg/db"
	"log"
	"net/http"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		log.Println("db.Tasks error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "ошибка при получении списка задач",
		})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
