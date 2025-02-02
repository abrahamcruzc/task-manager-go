package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/abrahamcruzc/task-manager-go/internal/models"
    "github.com/abrahamcruzc/task-manager-go/internal/repository"
    "github.com/go-chi/chi/v5"
)

// TaskHandler define la interfaz para manejar operaciones CRUD relacionadas con tareas.
type TaskHandler interface {
    CreateTaskHandler(w http.ResponseWriter, r *http.Request)  // Maneja la creación de una nueva tarea.
    GetTasksHandler(w http.ResponseWriter, r *http.Request)    // Maneja la obtención de todas las tareas.
    GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) // Maneja la obtención de una tarea por su ID.
    UpdateTaskHandler(w http.ResponseWriter, r *http.Request)  // Maneja la actualización de una tarea existente.
    DeleteTaskHandler(w http.ResponseWriter, r *http.Request)  // Maneja la eliminación de una tarea.
}

// taskHandler implementa la interfaz TaskHandler y contiene una referencia al repositorio de tareas.
type taskHandler struct {
    repo repository.TaskRepository // Repositorio para interactuar con los datos de las tareas.
}

// NewTaskHandler crea una nueva instancia de taskHandler e inyecta el repositorio de tareas.
func NewTaskHandler(repo repository.TaskRepository) TaskHandler {
    return &taskHandler{repo: repo}
}

// CreateTaskHandler maneja la creación de una nueva tarea.
// Método HTTP: POST
// Ruta: /tasks
func (h *taskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close() // Cerrar el cuerpo de la solicitud al finalizar.

    var task models.Task
    // Decodificar el cuerpo de la solicitud en una estructura Task.
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validar que el nombre de la tarea no esté vacío.
    if task.Name == "" {
        http.Error(w, "Task name is required", http.StatusBadRequest)
        return
    }

    // Validar que el estado de la tarea sea válido.
    if err := task.Status.IsValid(); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Crear la tarea en el repositorio.
    if err := h.repo.CreateTask(&task); err != nil {
        log.Printf("Error creating task: %v", err) // Registrar el error para depuración.
        http.Error(w, "Error creating task", http.StatusInternalServerError)
        return
    }

    // Responder con un código de estado 201 Created y devolver la tarea creada.
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(task); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
    }
}

// GetTasksHandler maneja la obtención de todas las tareas.
// Método HTTP: GET
// Ruta: /tasks
func (h *taskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
    // Obtener todas las tareas del repositorio.
    tasks, err := h.repo.GetTasks()
    if err != nil {
        log.Printf("Error retrieving tasks: %v", err) // Registrar el error para depuración.
        http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
        return
    }

    // Responder con un código de estado 200 OK y devolver las tareas como JSON.
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(tasks); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
    }
}

// GetTaskByIDHandler maneja la obtención de una tarea por su ID.
// Método HTTP: GET
// Ruta: /tasks/{id}
func (h *taskHandler) GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Extraer el ID de la URL.
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    // Obtener la tarea por su ID desde el repositorio.
    task, err := h.repo.GetTaskByID(uint(id))
    if err != nil {
        log.Printf("Error retrieving task: %v", err) // Registrar el error para depuración.
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    // Responder con un código de estado 200 OK y devolver la tarea como JSON.
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(task); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
    }
}

// UpdateTaskHandler maneja la actualización de una tarea existente.
// Método HTTP: PUT
// Ruta: /tasks/{id}
func (h *taskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close() // Cerrar el cuerpo de la solicitud al finalizar.

    // Extraer el ID de la URL.
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var task models.Task
    // Decodificar el cuerpo de la solicitud en una estructura Task.
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validar que el nombre de la tarea no esté vacío.
    if task.Name == "" {
        http.Error(w, "Name is required", http.StatusBadRequest)
        return
    }

    // Validar que el estado de la tarea sea válido si se proporciona.
    if task.Status != "" {
        if err := task.Status.IsValid(); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }

    // Asignar el ID extraído de la URL a la tarea.
    task.ID = uint(id)

    // Actualizar la tarea en el repositorio.
    if err := h.repo.UpdateTask(&task); err != nil {
        log.Printf("Error updating task: %v", err) // Registrar el error para depuración.
        http.Error(w, "Error updating task", http.StatusInternalServerError)
        return
    }

    // Responder con un código de estado 200 OK y devolver la tarea actualizada como JSON.
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(task); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
    }
}

// DeleteTaskHandler maneja la eliminación de una tarea.
// Método HTTP: DELETE
// Ruta: /tasks/{id}
func (h *taskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    // Extraer el ID de la URL.
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    // Eliminar la tarea del repositorio.
    if err := h.repo.DeleteTask(uint(id)); err != nil {
        log.Printf("Error deleting task: %v", err) // Registrar el error para depuración.
        http.Error(w, "Error deleting task", http.StatusInternalServerError)
        return
    }

    // Responder con un código de estado 204 No Content.
    w.WriteHeader(http.StatusNoContent)
}