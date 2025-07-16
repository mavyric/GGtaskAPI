package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Task represents the model for a task.
// The struct tags `json:"..."` are used to control how the struct is encoded to/decoded from JSON.
type Task struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"` // 0: incomplete, 1: completed
}

// TaskStore is an in-memory database for tasks.
// It uses a sync.RWMutex to handle concurrent read/write operations safely.
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

// NewTaskStore creates and returns a new TaskStore.
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]Task),
	}
}

// Handlers hold dependencies for the HTTP handlers, like the task store.
type Handlers struct {
	store *TaskStore
}

// getTasksHandler retrieves all tasks from the store.
func (h *Handlers) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	h.store.mu.RLock() // Lock for reading
	defer h.store.mu.RUnlock()

	// Convert map to a slice for JSON array response
	tasks := make([]Task, 0, len(h.store.tasks))
	for _, task := range h.store.tasks {
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// createTaskHandler creates a new task.
func (h *Handlers) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if task.Name == "" || (task.Status != 0 && task.Status != 1) {
		http.Error(w, `{"error": "name is required and status must be 0 or 1"}`, http.StatusBadRequest)
		return
	}

	h.store.mu.Lock() // Lock for writing
	defer h.store.mu.Unlock()

	task.ID = uuid.New().String()
	h.store.tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// updateTaskHandler updates an existing task.
func (h *Handlers) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.store.mu.Lock() // Lock for writing
	defer h.store.mu.Unlock()

	if _, ok := h.store.tasks[id]; !ok {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Basic validation
	if updatedTask.Name == "" || (updatedTask.Status != 0 && updatedTask.Status != 1) {
		http.Error(w, `{"error": "name is required and status must be 0 or 1"}`, http.StatusBadRequest)
		return
	}


	updatedTask.ID = id // Keep the original ID
	h.store.tasks[id] = updatedTask

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// deleteTaskHandler deletes a task.
func (h *Handlers) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.store.mu.Lock() // Lock for writing
	defer h.store.mu.Unlock()

	if _, ok := h.store.tasks[id]; !ok {
		http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
		return
	}

	delete(h.store.tasks, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	store := NewTaskStore()
	h := &Handlers{store: store}

	r := mux.NewRouter()

	// Define API endpoints
	r.HandleFunc("/tasks", h.getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", h.createTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", h.updateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", h.deleteTaskHandler).Methods("DELETE")

	log.Println("Starting API server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
