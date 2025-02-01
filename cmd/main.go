package main

import (
	"log"

	"github.com/abrahamcruzc/task-manager-go/internal/config" 
)

func main() {
	var conf *config.Config
	
	// Carga la configuración de la db
	if err := conf.LoadConfig(); err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}
	
	// Inicializa la conexión a la base de datos
	db, err := conf.InitDb()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	
	log.Println("Base de datos iniciada correctamente")
	log.Print(db)
	
	
	
}
