package main

import (
	"log"

	"github.com/abrahamcruzc/task-manager-go/internal/config"
	"github.com/abrahamcruzc/task-manager-go/internal/models"
)

func main() {
	// 1. Cargar la configuración
	cfg := &config.Config{}
	
	if err := cfg.LoadConfig(); err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// 2. Inicializar la base de datos
	db, err := cfg.InitDb()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// 3. Migrar modelos 
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf("Error al migrar modelos: %v", err)
	}

}
