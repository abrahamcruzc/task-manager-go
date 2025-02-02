package repository

import (
	"errors"
	"fmt"

	"github.com/abrahamcruzc/task-manager-go/internal/models"
	"gorm.io/gorm"
)

// TaskRepository define la interfaz para las operaciones CRUD de tareas
// Contrato que garantiza la implementación de los métodos esenciales
type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetTasks() ([]models.Task, error)
	GetTaskByID(id uint) (*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id uint) error
}

// repository implementación concreta de TaskRepository
// Encapsula la conexión a la base de datos usando GORM
type repository struct {
	db *gorm.DB
}

// NewTaskRepository factory para crear instancias del repositorio
// Recibe: conexión a la base de datos (*gorm.DB)
// Retorna: implementación de TaskRepository lista para usar
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &repository{db: db}
}

// CreateTask crea una nueva tarea en la base de datos
// Recibe: puntero a modelo Task
// Retorna: error de GORM si falla la operación
func (r *repository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

// GetTasks obtiene todas las tareas almacenadas
// Retorna: slice de tareas y error de GORM si ocurre
func (r *repository) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	if result := r.db.Find(&tasks); result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

// GetTaskByID busca una tarea por su ID
// Recibe: ID de la tarea (uint)
// Retorna: 
//   - Tarea encontrada o nil
//   - error detallado (incluye ErrRecordNotFound si no existe el registro)
func (r *repository) GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	
	// Manejo específico para registro no encontrado
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("task with ID %d not found: %w", id, result.Error)
	}
	
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

// UpdateTask actualiza una tarea existente usando actualización parcial
// Recibe: puntero a modelo Task con los campos a actualizar
// Retorna: error de GORM si falla la operación
// Nota: Usa Updates con mapa para evitar sobrescritura de campos no modificados
func (r *repository) UpdateTask(task *models.Task) error {
	result := r.db.Model(task).Updates(map[string]interface{}{
		"name":        task.Name,
		"description": task.Description,
		"status":      task.Status,
	})
	return result.Error
}

// DeleteTask elimina una tarea por su ID
// Recibe: ID de la tarea (uint)
// Retorna: 
//   - error de GORM si falla la operación
//   - error personalizado si el ID no existe
// Valida que se afectó al menos 1 registro con RowsAffected
func (r *repository) DeleteTask(id uint) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}
	return nil
}