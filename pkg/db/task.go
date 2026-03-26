package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask — добавляет новую задачу в базу данных.
// Возвращает id созданной записи.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Tasks — возвращает список задач с лимитом,
// отсортированных по дате.
func Tasks(limit int, search string) ([]*Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if search == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	} else {
		t, parseErr := time.Parse("02.01.2006", search)
		if parseErr == nil {
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, t.Format("20060102"), limit)
		} else {
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			pattern := "%" + search + "%"
			rows, err = DB.Query(query, pattern, pattern, limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []*Task{}

	for rows.Next() {
		var t Task
		var id int64

		err = rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		t.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask — получает одну задачу по id.
// Возвращает ошибку, если не найдена.
func GetTask(id string) (*Task, error) {
	var t Task
	var dbID int64

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача не найдена")
	}

	t.ID = strconv.FormatInt(dbID, 10)
	return &t, nil
}

// UpdateTask — обновляет все поля задачи по id.
// Возвращает ошибку, если задача не найдена.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask — удаляет задачу по id.
// Возвращает ошибку, если задача не найдена.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate — обновляет только дату задачи.
// Используется для повторяющихся задач.
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
