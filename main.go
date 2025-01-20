package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
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

	// Обработчик для получения всех задач
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, "Ошибка при получении задач", http.StatusInternalServerError)
		}
	})

	// Обработчик для отправки новой задачи
	r.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Невозможно декодировать задачу", http.StatusBadRequest)
			return
		}
		tasks[task.ID] = task
		w.WriteHeader(http.StatusCreated)
	})

	// Обработчик для получения задачи по ID
	r.Get("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		task, exists := tasks[id]
		if !exists {
			http.Error(w, "Задача не найдена", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, "Ошибка при получении задачи", http.StatusInternalServerError)
		}
	})

	// Обработчик для удаления задачи по ID
	r.Delete("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		_, exists := tasks[id]
		if !exists {
			http.Error(w, "Задача не найдена", http.StatusBadRequest)
			return
		}
		delete(tasks, id)
		w.WriteHeader(http.StatusOK)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
	}
}
