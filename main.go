package main

import (
	"github.com/gmakloren-lab/go_final_project/pkg/api"
	"github.com/gmakloren-lab/go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
)

// main — точка входа в приложение.
// Инициализирует базу данных, API и запускает HTTP-сервер.
func main() {
	// Работа с базой данных
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	api.Init()
	// Сервер
	webDir := "./web"
	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540" // значение по умолчанию
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on :" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
