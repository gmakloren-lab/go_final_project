package api

import (
	"encoding/json"
	"strconv"

	"github.com/gmakloren-lab/go_final_project/pkg/db"
	"net/http"
)

// addTaskHandler — обрабатывает POST /api/task.
// Принимает JSON с задачей, валидирует данные,
// проверяет дату и добавляет задачу в базу.
// Возвращает id созданной задачи или ошибку.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}
