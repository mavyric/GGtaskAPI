package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Task represents a to-do item.
type Task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"` // 0: incomplete, 1: completed
}

// TaskStore is an in-memory store for tasks.
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]Task),
	}
}

type Handlers struct {
	store *TaskStore
}

func main() {
	store := NewTaskStore()
	h := &Handlers{store: store}

	r := mux.NewRouter()
	r.HandleFunc("/tasks", h.getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", h.createTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", h.updateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", h.deleteTaskHandler).Methods("DELETE")

	log.Println("Starting API server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Handler methods

func (h *Handlers) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	tasks := make([]Task, 0, len(h.store.tasks))
	for _, task := range h.store.tasks {
		tasks = append(tasks, task)
	}
	respondJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if task.Name == "" || (task.Status != 0 && task.Status != 1) {
		respondError(w, http.StatusBadRequest, "Name is required and status must be 0 or 1")
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	task.ID = uuid.New().String()
	h.store.tasks[task.ID] = task
	respondJSON(w, http.StatusCreated, task)
}

func (h *Handlers) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	task, exists := h.store.tasks[id]
	if !exists {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}

	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if updated.Name == "" || (updated.Status != 0 && updated.Status != 1) {
		respondError(w, http.StatusBadRequest, "Name is required and status must be 0 or 1")
		return
	}
	updated.ID = task.ID
	h.store.tasks[id] = updated
	respondJSON(w, http.StatusOK, updated)
}

func (h *Handlers) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, exists := h.store.tasks[id]; !exists {
		respondError(w, http.StatusNotFound, "Task not found")
		return
	}
	delete(h.store.tasks, id)
	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
