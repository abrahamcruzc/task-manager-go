package routes

import (
	"net/http"

	"github.com/abrahamcruzc/task-manager-go/internal/handlers"
	"github.com/abrahamcruzc/task-manager-go/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

// SetupRoutes configura el router principal de la API y sus dependencias.
// 
// Parámetros:
//   - db *gorm.DB: Conexión a la base de datos inyectada desde la capa de configuración.
//
// Retorno:
//   - http.Handler: Router configurado con todas las rutas y middlewares.
//
// Flujo de trabajo:
//   1. Crea un nuevo router Chi
//   2. Aplica middlewares globales
//   3. Inicializa capas de repositorio y handlers
//   4. Define las rutas agrupadas
//   5. Retorna el router listo para usar
//
func SetupRoutes(db *gorm.DB) http.Handler {
	// Inicializa el router Chi con sus configuraciones básicas
	r := chi.NewRouter()
	
	// Middlewares globales aplicados a todas las rutas en orden de ejecución:
	// 1. Logger: Registra detalles de cada solicitud (método, ruta, duración)
	// 2. Recoverer: Maneja pánicos y retorna error HTTP 500
	r.Use(
		middleware.Logger,    // Formato de log: [INFO] GET /tasks 200 12.34ms
		middleware.Recoverer, // Previene caídas de la aplicación
	)
	
	// Inicialización de dependencias (patrón de inyección de dependencias)
	// Capa de acceso a datos -> Capa de manejo de requests
	taskRepo := repository.NewTaskRepository(db)      // Repositorio con operaciones DB
	taskHandler := handlers.NewTaskHandler(taskRepo)  // Handler con lógica HTTP

	// Grupo de rutas para operaciones CRUD de tareas
	// Todas las rutas comienzan con /tasks
	r.Route("/tasks", func(r chi.Router) {
		// GET /tasks - Obtener todas las tareas
		r.Get("/", taskHandler.GetTasksHandler)
		
		// POST /tasks - Crear nueva tarea
		r.Post("/", taskHandler.CreateTaskHandler)
		
		// GET /tasks/{id} - Obtener tarea por ID
		r.Get("/{id}", taskHandler.GetTaskByIDHandler)
		
		// PUT /tasks/{id} - Actualizar tarea existente
		r.Put("/{id}", taskHandler.UpdateTaskHandler)
		
		// DELETE /tasks/{id} - Eliminar tarea
		r.Delete("/{id}", taskHandler.DeleteTaskHandler)
	})

	return r
}