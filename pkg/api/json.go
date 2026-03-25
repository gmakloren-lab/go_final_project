package api

import (
	"encoding/json"
	"net/http"
)

// writeJSON — отправляет ответ в формате JSON.
// Устанавливает Content-Type и кодирует данные.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}
