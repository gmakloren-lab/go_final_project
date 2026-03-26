Проект представляет собой веб-приложение для планирования задач.
Приложение позволяет:
- добавлять задачи;
- получать список ближайших задач;
- редактировать задачи;
- удалять задачи;
- отмечать задачи выполненными;
- работать с повторяющимися задачами;
- выполнять поиск задач.
Бэкенд написан на Go, данные хранятся в SQLite, фронтенд подключён из готовой директории `web`.

В проекте выполнены следующие задания повышенной трудности:
- поиск задач по строке поиска;
- поиск задач по дате;
- аутентификация по паролю;
- сборка и запуск приложения через Docker.

  Проект поддерживает следующие переменные окружения:
- `TODO_PORT` — порт запуска приложения;
- `TODO_DBFILE` — путь к файлу базы данных SQLite;
- `TODO_PASSWORD` — пароль для аутентификации. Если не задан, аутентификация отключена.
Запуск без аутентификации:
     export TODO_PORT=7540
     export TODO_DBFILE=scheduler.db
     go run main.go
  С аутентификацией:
     export TODO_PORT=7540
     export TODO_DBFILE=scheduler.db
     export TODO_PASSWORD=12345
     go run main.go
  Адрес в браузере: http://localhost:7540

Запуск всех тестов используется команда:
   go test ./tests
Запуск отдельных тестов:
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
  Для проверки поиска нужно включить соответствующую настройку поиска в tests/settings.go.
  Для проверки аутентификации нужно:
  запустить сервер с установленной переменной TODO_PASSWORD;
  выполнить запрос к /api/signin;
  получить токен из ответа;
  вставить этот токен в поле Token в файле tests/settings.go.

Сборка и запуск проекта на Docker:
сборка : docker build -t go-final-project .
запуск без аутентификации :
docker run --name scheduler-app -p 7540:7540 -v /path/to/project/data:/data go-final-project
запуск с аутентификацией :
docker run --name scheduler-app -p 7540:7540 -v /path/to/project/data:/data go-final-project
Доступно по адресу:
http://localhost:7540