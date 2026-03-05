package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luispfcanales/rainforestapp/pkg/config"
	"github.com/luispfcanales/rainforestapp/pkg/database"
	"github.com/luispfcanales/rainforestapp/pkg/handler"
)

func main() {
	// Cargar configuración desde variables de entorno
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Error validando configuración: %v", err)
	}

	// Inicializar handlers una sola vez
	h, err := handler.NewUsuarioHandler(cfg)
	if err != nil {
		log.Fatalf("Error inicializando handler de usuarios: %v", err)
	}

	th, err := handler.NewTicketHandler(cfg)
	if err != nil {
		log.Fatalf("Error inicializando handler de tickets: %v", err)
	}

	// Registrar rutas
	mux := http.NewServeMux()

	// Rutas de usuarios
	mux.HandleFunc("/api/usuarios/listar", h.ListUsuarios)
	mux.HandleFunc("/api/usuarios/registrar", h.CreateUsuario)
	mux.HandleFunc("/api/usuarios/buscar", h.GetUsuario)
	mux.HandleFunc("/api/usuarios/reporte", h.GetUsuarioPDF)

	// Rutas de tickets
	mux.HandleFunc("/api/tickets/create", th.CreateTicket)
	mux.HandleFunc("/api/tickets/list", th.ListTickets)
	mux.HandleFunc("/api/tickets/get", th.GetTicket)
	mux.HandleFunc("/api/tickets/status", th.UpdateStatus)
	mux.HandleFunc("/api/tickets/update", th.UpdateTicket)
	mux.HandleFunc("/api/tickets/delete", th.DeleteTicket)

	// Puerto configurable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Apagando servidor...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Error en shutdown: %v", err)
		}
		database.Close()
	}()

	log.Println("=== Servidor de desarrollo ===")
	log.Printf("Rutas registradas:")
	log.Printf("  GET    /api/usuarios/listar")
	log.Printf("  POST   /api/usuarios/registrar")
	log.Printf("  GET    /api/usuarios/buscar?dni=...")
	log.Printf("  GET    /api/usuarios/reporte?dni=...")
	log.Printf("  POST   /api/tickets/create")
	log.Printf("  GET    /api/tickets/list")
	log.Printf("  GET    /api/tickets/get?id=...")
	log.Printf("  PATCH  /api/tickets/status")
	log.Printf("  PUT    /api/tickets/update")
	log.Printf("  DELETE /api/tickets/delete?id=...")
	log.Printf("Escuchando en http://localhost:%s", port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error en servidor: %v", err)
	}

	log.Println("Servidor detenido")
}
