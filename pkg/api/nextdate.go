package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102" // формат даты YYYYMMDD

// nextDayHandler — HTTP-обработчик для запроса следующей даты выполнения задачи.
// Читает параметры now, date и repeat из формы запроса,
// вычисляет следующую дату с помощью NextDate и возвращает её в ответе.
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now() // если now не указан, берём текущую дату
	} else {
		now, err = time.Parse(layout, nowStr)
		if err != nil {
			http.Error(w, "invalid now", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

// afterNow — проверяет, что date находится после now (игнорируя время).
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()

	dt1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	dt2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return dt1.After(dt2)
}

// lastDayOfMonth — возвращает номер последнего дня месяца для даты t.
func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// NextDate — вычисляет следующую дату выполнения задачи
// на основе текущей даты now, даты старта dstart и правила repeat.
// Поддерживаются правила повторения:
// d N — каждые N дней
// w D1,D2 — по дням недели
// m D1,D2,... M1,M2,... — по дням месяца и месяцам
// y — каждый год
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	case "d":
		// повторение через N дней
		if len(parts) != 2 {
			return "", errors.New("invalid d format")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid day interval")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		// повторение каждый год
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		// повторение по дням недели
		if len(parts) != 2 {
			return "", errors.New("invalid w format")
		}

		var days [8]bool
		for _, s := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(s)
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("invalid weekday")
			}
			days[d] = true
		}

		for {
			date = date.AddDate(0, 0, 1)

			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}

			if days[wd] && afterNow(date, now) {
				break
			}
		}

	case "m":
		// повторение по дням месяца и месяцам
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid m format")
		}

		var dayMask [34]bool
		var monthMask [13]bool

		for _, s := range strings.Split(parts[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < -2 || n == 0 || n > 31 {
				return "", errors.New("invalid day")
			}
			dayMask[n+2] = true
		}

		if len(parts) == 3 {
			for _, s := range strings.Split(parts[2], ",") {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("invalid month")
				}
				monthMask[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthMask[i] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			if !monthMask[int(date.Month())] {
				continue
			}

			day := date.Day()
			last := lastDayOfMonth(date)

			ok := false

			if dayMask[day+2] {
				ok = true
			}

			if dayMask[1] && day == last {
				ok = true
			}

			if dayMask[0] && day == last-1 {
				ok = true
			}

			if ok && afterNow(date, now) {
				break
			}
		}

	default:
		return "", errors.New("unknown rule")
	}

	return date.Format(layout), nil
}
