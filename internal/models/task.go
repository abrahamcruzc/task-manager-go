package models

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
)

// Status define los posibles estados de una tarea
type Status string

const (
	ToDo       Status = "To do"       // Tarea pendiente
	InProgress Status = "In progress" // Tarea en progreso
	Completed  Status = "Completed"   // Tarea finalizada
)

// Scan implementa la interfaz Scanner para leer el valor desde la base de datos
func (s *Status) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*s = Status(v)
	case string:
		*s = Status(v)
	default:
		return fmt.Errorf("tipo %T no compatible para Status", value)
	}
	return s.IsValid()
}

// Value implementa la interfaz Valuer para escribir el valor en la base de datos
func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

// IsValid verifica si el estado actual es válido
func (s Status) IsValid() error {
	switch s {
	case ToDo, InProgress, Completed:
		return nil
	default:
		return fmt.Errorf("estado inválido: %s", s)
	}
}

// ValidValues retorna la lista de valores permitidos como strings
func (Status) ValidValues() []string {
	return []string{
		string(ToDo),
		string(InProgress),
		string(Completed),
	}
}

// Task representa una tarea en el sistema
type Task struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null;size:100"`             // Nombre único de la tarea
	Description string `gorm:"size:255;not null"`                         // Descripción detallada
	Status      Status `gorm:"type:varchar(20);default:'To do';not null"` // Estado actual
}

// BeforeSave hook de GORM que valida el estado antes de guardar
func (t *Task) BeforeSave(tx *gorm.DB) error {
	if err := t.Status.IsValid(); err != nil {
		return fmt.Errorf("error validando status: %w", err)
	}
	return nil
}
