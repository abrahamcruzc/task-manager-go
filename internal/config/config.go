package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config contiene la configuración de la aplicación mapeada desde variables de entorno
// Campos corresponden a las variables de entorno con la nomenclatura DB_*
type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`     // Host de la base de datos
	DBPort     string `mapstructure:"DB_PORT"`     // Puerto de la base de datos
	DBUser     string `mapstructure:"DB_USER"`     // Usuario de la base de datos
	DBPassword string `mapstructure:"DB_PASSWORD"` // Contraseña del usuario
	DBName     string `mapstructure:"DB_NAME"`     // Nombre de la base de datos
	SSLMode    string `mapstructure:"SSL_MODE"`    // Modo SSL para la conexión
}

// LoadConfig carga la configuración desde variables de entorno y .env
// Flujo de ejecución:
// 1. Intenta cargar archivo .env
// 2. Asigna variables de entorno del sistema
// 3. Valida campos requeridos
func (c *Config) LoadConfig() error {
	// 1. Carga archivo .env (si existe)
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env. Usando variables de entorno del sistema.")
	}

	// 2. Asigna valores directamente desde las variables de entorno
	c.DBHost = os.Getenv("DB_HOST")
	c.DBPort = os.Getenv("DB_PORT")
	c.DBUser = os.Getenv("DB_USER")
	c.DBPassword = os.Getenv("DB_PASSWORD")
	c.DBName = os.Getenv("DB_NAME")
	c.SSLMode = os.Getenv("SSL_MODE")
	

	// 3. Valida campos obligatorios
	if c.DBUser == "" || c.DBPassword == "" || c.DBName == "" {
		return fmt.Errorf("configuración incompleta: DB_USER, DB_PASSWORD y DB_NAME son requeridos")
	}

	log.Printf("Configuración cargada: Host=%s, Port=%s, DB=%s", c.DBHost, c.DBPort, c.DBName)
	return nil
}

// InitDb establece la conexión con PostgreSQL usando GORM
// Retorna:
// - Instancia de GORM DB para operaciones de base de datos
// - error detallado si falla la conexión
func (c *Config) InitDb() (*gorm.DB, error) {
	// Construye el DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.SSLMode,
	)

	// Establece la conexión con PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("fallo en conexión a PostgreSQL: %v\nDSN usado: %s", err, dsn)
	}

	// Configuración adicional 
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)    // Conexiones inactivas máximas
	sqlDB.SetMaxOpenConns(100)   // Conexiones abiertas máximas
	sqlDB.SetConnMaxLifetime(30) // Tiempo máximo de vida de conexión (minutos)

	log.Println("Conexión a PostgreSQL establecida exitosamente")
	return db, nil
}
