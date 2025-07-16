# GGtaskAPI 🚀

`GGtaskAPI` is a simple, robust, and containerized RESTful API for managing tasks, built with Go. It provides endpoints to create, retrieve, update, and delete tasks using an in-memory data store.

## ✨ Features

- **CRUD Operations**: Full support for Create, Read, Update, and Delete tasks.
- **In-Memory Storage**: Uses a thread-safe map for fast, non-persistent data storage.
- **RESTful Endpoints**: Clean and predictable API design.
- **Containerized**: Includes a multi-stage `Dockerfile` for lightweight and secure deployments.
- **Tested**: Unit tests for all API endpoints.

## Project Structure

For this project, I've used a **flat structure**, keeping all the core logic in the `main` package. This approach is efficient and easy to navigate for smaller applications and microservices like this one.

### Alternative: Layered Structure

For a larger, production-grade application that is expected to scale, I would adopt a **layered structure**. This pattern separates concerns into distinct packages, improving maintainability and testability.
An example of a layered structure would be:

```
GGtaskAPI/
├── api/          // HTTP handlers
├── model/        // Data structures (e.g., Task struct)
├── store/        // Data storage logic
└── main.go       // Wires everything together
```

This separation makes the project easier to manage as it grows. By choosing the flat structure here, I've selected the most direct and appropriate design for the current scope of the project.

## 🛠️ Prerequisites

- [Go](https://go.dev/doc/install) 1.18+
- [Docker](https://docs.docker.com/get-docker/) (for containerization)

## ⚙️ Running Locally

1.  **Clone the repository (or create the files):**
    ```bash
    # git clone https://github.com/mavyric/GGtaskAPI.git
    # cd GGtaskAPI
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the application:**
    ```bash
    go run main.go
    ```
    The API server will start on `http://localhost:8080`.

## 🐳 Running with Docker

1.  **Build the Docker image:**
    ```bash
    docker build -t ggtask-api .
    ```

2.  **Run the Docker container:**
    ```bash
    docker run -p 8080:8080 ggtask-api
    ```
    The API server will be accessible at `http://localhost:8080`.

## 🧪 Running Tests

To run the unit tests, execute the following command in the project root:

```bash
go test -v
```

## 📜 API Endpoints

All request and response bodies are in JSON format.

#### `Task` Object

```json
{
  "id": "string (uuid)",
  "name": "string",
  "description": "string",
  "status": "integer (0 for incomplete, 1 for completed)"
}
```

---

### **List All Tasks**

-   **Endpoint:** `GET /tasks`
-   **Description:** Retrieves a list of all tasks.
-   **Success Response:** `200 OK`
-   **Example:** `curl http://localhost:8080/tasks`

    ```json
    [
      {
        "id": "f8c3de3d-1fea-4d7c-a8b0-29f63c4c3454",
        "name": "Learn Go",
        "description": "Complete the official Go tour.",
        "status": 1
      }
    ]
    ```

### **Create a New Task**

-   **Endpoint:** `POST /tasks`
-   **Description:** Creates a new task. The `id` is generated automatically.
-   **Success Response:** `201 Created`
-   **Example:** `curl -X POST -H "Content-Type: application/json" -d '{"name": "Build an API", "description": "Use Go and Docker", "status": 0}' http://localhost:8080/tasks`

### **Update an Existing Task**

-   **Endpoint:** `PUT /tasks/{id}`
-   **Description:** Updates the details of a specific task by its ID.
-   **Success Response:** `200 OK`
-   **Error Response:** `404 Not Found` if the task ID does not exist.
-   **Example:** `curl -X PUT -H "Content-Type: application/json" -d '{"name": "Build an API", "description": "Use Go and Docker", "status": 1}' http://localhost:8080/tasks/YOUR_TASK_ID`

### **Delete a Task**

-   **Endpoint:** `DELETE /tasks/{id}`
-   **Description:** Deletes a specific task by its ID.
-   **Success Response:** `204 No Content`
-   **Error Response:** `404 Not Found` if the task ID does not exist.
-   **Example:** `curl -X DELETE http://localhost:8080/tasks/YOUR_TASK_ID`
