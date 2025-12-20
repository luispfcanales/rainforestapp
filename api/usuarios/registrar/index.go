package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/luispfcanales/rainforestapp/internal/config"
	usuarioHandler "github.com/luispfcanales/rainforestapp/internal/handler"
	"github.com/luispfcanales/rainforestapp/pkg/response"
)

var (
	handlerInstance *usuarioHandler.UsuarioHandler
	once            sync.Once
	initError       error
)

// getHandler obtiene o crea la instancia del handler
func getHandler() (*usuarioHandler.UsuarioHandler, error) {
	once.Do(func() {
		// Cargar configuración
		cfg, err := config.Load()
		if err != nil {
			initError = err
			log.Printf("Error cargando configuración: %v", err)
			return
		}

		// Validar configuración
		if err := cfg.Validate(); err != nil {
			initError = err
			log.Printf("Error validando configuración: %v", err)
			return
		}

		// Crear handler
		handlerInstance, initError = usuarioHandler.NewUsuarioHandler(cfg)
		if initError != nil {
			log.Printf("Error inicializando handler: %v", initError)
		}
	})

	return handlerInstance, initError
}

// Handler es el punto de entrada de Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Obtener handler inicializado
	h, err := getHandler()
	if err != nil {
		log.Printf("Error obteniendo handler: %v", err)
		response.InternalServerError(w, "Error de configuración del servidor")
		return
	}

	// Delegar al handler de usuario
	h.CreateUsuario(w, r)
}
