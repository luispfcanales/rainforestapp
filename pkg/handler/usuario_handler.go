package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/luispfcanales/rainforestapp/pkg/config"
	"github.com/luispfcanales/rainforestapp/pkg/database"
	"github.com/luispfcanales/rainforestapp/pkg/models"
	"github.com/luispfcanales/rainforestapp/pkg/pdf"
	"github.com/luispfcanales/rainforestapp/pkg/repository"
	"github.com/luispfcanales/rainforestapp/pkg/response"
	"github.com/luispfcanales/rainforestapp/pkg/service"
)

// UsuarioHandler maneja las peticiones HTTP para usuarios
type UsuarioHandler struct {
	service *service.UsuarioService
	pdfGen  *pdf.PDFGenerator
}

// NewUsuarioHandler crea una nueva instancia del handler
func NewUsuarioHandler(cfg *config.Config) (*UsuarioHandler, error) {
	ctx := context.Background()

	// Inicializar cliente de Firestore
	firestoreClient, err := database.GetFirestoreClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Crear repositorio y servicio
	repo := repository.NewUsuarioRepository(firestoreClient)
	svc := service.NewUsuarioService(repo)

	pdfGen := pdf.NewPDFGenerator()

	return &UsuarioHandler{
		service: svc,
		pdfGen:  pdfGen,
	}, nil
}

// setupCORS configura los headers CORS
func setupCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// CreateUsuario maneja la creación de usuarios
func (h *UsuarioHandler) CreateUsuario(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	// Manejar preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Validar método
	if r.Method != "POST" {
		response.Error(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Parsear request
	var req models.CreateUsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Datos inválidos: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Crear usuario
	usuario, err := h.service.CreateUsuario(ctx, &req)
	if err != nil {
		log.Printf("Error creando usuario: %v", err)
		response.InternalServerError(w, "Error al crear usuario: "+err.Error())
		return
	}

	// Respuesta exitosa
	response.Created(w, "Usuario registrado exitosamente", usuario)
}

// GetUsuario maneja la obtención de un usuario por DNI
func (h *UsuarioHandler) GetUsuario(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		response.Error(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Obtener ID de query params
	dni := r.URL.Query().Get("dni")
	if dni == "" {
		response.BadRequest(w, "DNI es requerido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	usuario, err := h.service.GetUsuarioByDNI(ctx, dni)
	if err != nil {
		log.Printf("Error obteniendo usuario: %v", err)
		response.NotFound(w, "Usuario no encontrado")
		return
	}

	response.Success(w, "Usuario encontrado", usuario)
}

// ListUsuarios maneja la lista de todos los usuarios
func (h *UsuarioHandler) ListUsuarios(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		response.Error(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	usuarios, err := h.service.ListUsuarios(ctx, 100) // Límite de 100
	if err != nil {
		log.Printf("Error listando usuarios: %v", err)
		response.InternalServerError(w, "Error al listar usuarios")
		return
	}

	response.Success(w, "Usuarios obtenidos exitosamente", usuarios)
}

// GetUsuarioPDF maneja la obtención de un usuario específico en PDF
func (h *UsuarioHandler) GetUsuarioPDF(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		response.Error(w, http.StatusMethodNotAllowed, "Método no permitido")
		return
	}

	// Obtener ID del usuario desde query params
	usuarioDNI := r.URL.Query().Get("dni")
	if usuarioDNI == "" {
		response.Error(w, http.StatusBadRequest, "DNI de usuario requerido")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Obtener usuario por ID
	usuario, err := h.service.GetUsuarioByDNI(ctx, usuarioDNI)
	if err != nil {
		log.Printf("Error obteniendo usuario: %v", err)
		response.Error(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	// Generar PDF individual
	pdfBytes, err := h.pdfGen.GenerateUsuarioPDF(ctx, usuario)
	if err != nil {
		log.Printf("Error generando PDF: %v", err)
		response.InternalServerError(w, "Error generando PDF")
		return
	}

	// Configurar headers para descarga
	w.Header().Set("Content-Type", "application/pdf")
	filename := fmt.Sprintf("usuario_%s_%s_%s.pdf", usuario.Nombres, usuario.ApellidoPaterno, time.Now().Format("20060102"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	// Escribir PDF en la respuesta
	if _, err := w.Write(pdfBytes); err != nil {
		log.Printf("Error escribiendo PDF: %v", err)
	}
}
