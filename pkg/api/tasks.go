package api

import (
	"net/http"

	"github.com/yamixdev/go-final-sprint/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}