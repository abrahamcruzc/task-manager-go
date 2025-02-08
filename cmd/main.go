package main

import (
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
	"github.com/abrahamcruzc/task-manager-go/internal/config"
	"github.com/abrahamcruzc/task-manager-go/internal/models"
	"github.com/abrahamcruzc/task-manager-go/internal/routes"
)

// main es el punto de entrada de la aplicación.
// Este programa se encarga de cargar la configuración, inicializar la conexión a la base de datos,
// ejecutar las migraciones necesarias, configurar las rutas de la API y arrancar el servidor HTTP.
func main() {
	// 1. Cargar la configuración
	// Se crea una instancia de la estructura de configuración y se cargan los valores (por ejemplo, desde .env o variables de entorno).
	cfg := &config.Config{}
	if err := cfg.LoadConfig(); err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	// 2. Inicializar la base de datos
	// Se utiliza la configuración cargada para establecer una conexión con la base de datos a través de GORM.
	var db *gorm.DB
	var err error
	maxRetries := 10               // Número máximo de intentos
	retryInterval := 5 * time.Second // Intervalo entre intentos

	for i := 1; i <= maxRetries; i++ {
		db, err = cfg.InitDb()
		if err == nil {
			log.Println("¡La base de datos está disponible!")
			break // Salimos del bucle cuando la conexión es exitosa
		}
		log.Printf("Intento %d/%d: La base de datos aún no está disponible: %v", i, maxRetries, err)
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
	}

	// 3. Migrar modelos
	// Se ejecuta la migración automática de la estructura del modelo Task, creando o actualizando la tabla correspondiente en la base de datos.
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}

	// 4. Configurar el router de la API
	// Se llama a la función SetupRoutes, pasando la conexión a la base de datos,
	// para que se instancien el repositorio, los handlers y se configuren las rutas y middlewares.
	handler := routes.SetupRoutes(db)

	// 5. Arrancar el servidor HTTP
	// Se define la dirección y el puerto en los que se ejecutará el servidor.
	port := ":8080" // Puerto en el que se escucharán las solicitudes HTTP.
	addr := "0.0.0.0" // Dirección en la que el servidor estará disponible (escucha en todas las interfaces).
	log.Printf("Listening at: %s", addr+port)
	
	// Se arranca el servidor utilizando ListenAndServe, pasando la dirección completa y el router configurado.
	if err := http.ListenAndServe(addr+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}