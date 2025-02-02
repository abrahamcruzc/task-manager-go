package models

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
)

// Status define los posibles estados de una tarea
// Implementa Scanner/Valuer para integración con la base de datos
type Status string

const (
	ToDo       Status = "To do"       // Estado inicial de la tarea (no iniciada)
	InProgress Status = "In progress" // Tarea actualmente en desarrollo
	Completed  Status = "Completed"   // Tarea finalizada exitosamente
)

// Scan implementa la interfaz Scanner para convertir valores de la base de datos
// Recibe: valor de la base de datos ([]byte o string)
// Retorna: error si el tipo no es compatible o el estado es inválido
// Nota: Se ejecuta automáticamente al leer de la base de datos
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

// Value implementa la interfaz Valuer para convertir valores a formato compatible con la base de datos
// Retorna: representación string del estado lista para almacenar
func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

// IsValid verifica si el valor actual es un estado permitido
// Retorna: error descriptivo si el estado no está en la lista blanca
// Uso: Validación manual o en hooks de GORM
func (s Status) IsValid() error {
	switch s {
	case ToDo, InProgress, Completed:
		return nil
	default:
		return fmt.Errorf("estado inválido: %s", s)
	}
}

// ValidValues retorna los valores permitidos como strings
// Propósito: Validación en capas superiores (ej: API, formularios)
func (Status) ValidValues() []string {
	return []string{
		string(ToDo),
		string(InProgress),
		string(Completed),
	}
}

// Task representa una entidad de tarea en el sistema
// Campos GORM: ID, CreatedAt, UpdatedAt, DeletedAt (embedded)
type Task struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null;size:100"` // Nombre único con máximo 100 caracteres
	Description string `gorm:"size:255;not null"`             // Descripción con máximo 255 caracteres
	Status      Status `gorm:"type:varchar(20);default:'To do';not null"` // Estado con valor por defecto
}

// BeforeSave hook de ciclo de vida de GORM para validación automática
// Se ejecuta antes de cualquier operación Create/Update
// Retorna: error si la validación falla, abortando la operación
func (t *Task) BeforeSave(tx *gorm.DB) error {
	if err := t.Status.IsValid(); err != nil {
		return fmt.Errorf("validación fallida al guardar tarea: %w", err)
	}
	return nil
}