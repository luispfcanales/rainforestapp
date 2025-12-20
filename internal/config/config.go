package config

import (
	"fmt"
	"os"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Firebase FirebaseConfig
}

// FirebaseConfig configuración de Firebase
type FirebaseConfig struct {
	ProjectID   string
	Credentials string
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
	firebaseCredentials := os.Getenv("FIREBASE_CREDENTIALS")

	if firebaseCredentials == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS no está configurado")
	}

	return &Config{
		Firebase: FirebaseConfig{
			ProjectID:   firebaseProjectID,
			Credentials: firebaseCredentials,
		},
	}, nil
}

// Validate valida que la configuración sea correcta
func (c *Config) Validate() error {
	if c.Firebase.Credentials == "" {
		return fmt.Errorf("firebase credentials son requeridas")
	}
	return nil
}
