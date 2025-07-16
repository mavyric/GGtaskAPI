package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// setupRouter initializes the router and handlers for testing.
func setupRouter() (*mux.Router, *Handlers) {
	store := NewTaskStore()
	h := &Handlers{store: store}
	router := mux.NewRouter()
	router.HandleFunc("/tasks", h.getTasksHandler).Methods("GET")
	router.HandleFunc("/tasks", h.createTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/{id}", h.updateTaskHandler).Methods("PUT")
	router.HandleFunc("/tasks/{id}", h.deleteTaskHandler).Methods("DELETE")
	return router, h
}

func TestGetTasksHandler(t *testing.T) {
	router, h := setupRouter()

	// Pre-populate store with a task
	task := Task{ID: "1", Name: "Test Task", Description: "A test task", Status: 0}
	h.store.tasks["1"] = task

	req, _ := http.NewRequest("GET", "/tasks", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var tasks []Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Could not parse response body: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Name != "Test Task" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestCreateTaskHandler(t *testing.T) {
	router, _ := setupRouter()
	
	taskPayload := []byte(`{"name": "New Task", "description": "A new test task", "status": 0}`)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskPayload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdTask Task
	json.Unmarshal(rr.Body.Bytes(), &createdTask)
	if createdTask.Name != "New Task" || createdTask.ID == "" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	router, h := setupRouter()

	// Pre-populate store with a task
	taskID := "1"
	h.store.tasks[taskID] = Task{ID: taskID, Name: "Old Name", Description: "Old Desc", Status: 0}

	updatePayload := []byte(`{"name": "Updated Name", "description": "Updated Desc", "status": 1}`)
	req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if h.store.tasks[taskID].Name != "Updated Name" || h.store.tasks[taskID].Status != 1 {
		t.Errorf("task was not updated correctly in the store")
	}

	// Test update non-existent task
	req, _ = http.NewRequest("PUT", "/tasks/nonexistent", bytes.NewBuffer(updatePayload))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent task: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	router, h := setupRouter()
	
	// Pre-populate store with a task
	taskID := "1"
	h.store.tasks[taskID] = Task{ID: taskID, Name: "To Be Deleted", Description: "", Status: 0}
	
	req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if _, ok := h.store.tasks[taskID]; ok {
		t.Errorf("task was not deleted from the store")
	}

	// Test delete non-existent task
	req, _ = http.NewRequest("DELETE", "/tasks/nonexistent", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent task: got %v want %v", status, http.StatusNotFound)
	}
}
