package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	SSLMode    string `mapstructure:"SSL_MODE"`
}

// Carga la configuración desde el archivo .env
func (c *Config) LoadConfig() error {
	// Carga el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró el archivo .env. Usando variables de entorno del sistema.")
	}

	// Lee las variables de entorno del sistema
	viper.AutomaticEnv()

	// Asignar valores por defecto
	viper.SetDefault("DB_HOST", "0.0.0.0")
	viper.SetDefault("DB_PORT", "80")
	viper.SetDefault("SSL_MODE", "disable")

	// Cargar la configuración en la estructura LoadConfig
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("Error al cargar la configuración %v", err)
	}

	return nil
}

func (c *Config) InitDb() (*gorm.DB, error) {
	// Crea el DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode,
	)

	// Conecta a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error al conectarse a la base de datos %v", err)
	}

	log.Println("Conexión a la base de datos establecida")
	return db, nil
}
