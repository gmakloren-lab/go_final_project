package api

import (
	"errors"
	"time"

	"github.com/gmakloren-lab/go_final_project/pkg/db"
)

// checkDate — валидирует дату задачи.
// Проверяет формат, подставляет текущую дату,
// и корректирует её с учётом repeat (если нужно).
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("дата представлена в неверном формате")
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}
