package api

import (
	"net/http"

	"github.com/gmakloren-lab/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler — возвращает список задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
