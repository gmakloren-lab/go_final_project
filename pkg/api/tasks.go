package api

import (
	"net/http"

	"github.com/gmakloren-lab/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler — обрабатывает GET /api/tasks
// Возвращает список ближайших задач (с лимитом),
// отсортированных по дате.
// Поддерживает поиск (если реализован).
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
