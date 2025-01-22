package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task описывает задачу.
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postman",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func main() {
	r := chi.NewRouter()

	// Обработчик для получения всех задач.
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"Ошибка при получении задач"}`, http.StatusInternalServerError)
			return
		}
	})

	// Обработчик для отправки новой задачи.
	r.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error":"Невозможно декодировать задачу"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, не существует ли уже задача с таким ID.
		if _, exists := tasks[task.ID]; exists {
			http.Error(w, `{"error":"Задача с таким ID уже существует"}`, http.StatusConflict)
			return
		}

		tasks[task.ID] = task
		w.WriteHeader(http.StatusCreated)
		// Можно вернуть созданную задачу в ответе (опционально):
		// json.NewEncoder(w).Encode(task)
	})

	// Обработчик для получения задачи по ID.
	r.Get("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := chi.URLParam(r, "id")

		task, exists := tasks[id]
		if !exists {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, `{"error":"Ошибка при получении задачи"}`, http.StatusInternalServerError)
		}
	})

	// Обработчик для удаления задачи по ID.
	r.Delete("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := chi.URLParam(r, "id")

		if _, exists := tasks[id]; !exists {
			http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
			return
		}

		delete(tasks, id)
		w.WriteHeader(http.StatusOK)

	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s\n", err.Error())
	}
}
