package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/luispfcanales/rainforestapp/pkg/config"
	"github.com/luispfcanales/rainforestapp/pkg/database"
	"github.com/luispfcanales/rainforestapp/pkg/models"
	"github.com/luispfcanales/rainforestapp/pkg/repository"
	"github.com/luispfcanales/rainforestapp/pkg/response"
	"github.com/luispfcanales/rainforestapp/pkg/service"
)

// TicketHandler handles HTTP requests for tickets
type TicketHandler struct {
	service *service.TicketService
}

// NewTicketHandler creates a new handler instance
func NewTicketHandler(cfg *config.Config) (*TicketHandler, error) {
	ctx := context.Background()

	firestoreClient, err := database.GetFirestoreClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.NewTicketRepository(firestoreClient)
	svc := service.NewTicketService(repo)

	return &TicketHandler{
		service: svc,
	}, nil
}

// CreateTicket handles ticket creation
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req models.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid data: "+err.Error())
		return
	}
	defer r.Body.Close()

	ticket, err := h.service.CreateTicket(ctx, &req)
	if err != nil {
		log.Printf("Error creating ticket: %v", err)
		response.InternalServerError(w, "Error creating ticket: "+err.Error())
		return
	}

	response.Created(w, "Ticket created successfully", ticket)
}

// GetTicket handles retrieving a ticket by ID
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.BadRequest(w, "ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	ticket, err := h.service.GetTicketByID(ctx, id)
	if err != nil {
		log.Printf("Error getting ticket: %v", err)
		response.NotFound(w, "Ticket not found")
		return
	}

	response.Success(w, "Ticket found", ticket)
}

// ListTickets handles listing all tickets
func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	tickets, err := h.service.ListTickets(ctx, 200)
	if err != nil {
		log.Printf("Error listing tickets: %v", err)
		response.InternalServerError(w, "Error listing tickets")
		return
	}

	response.Success(w, "Tickets retrieved successfully", tickets)
}

// UpdateStatus handles status update (drag & drop)
func (h *TicketHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PATCH" {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req models.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid data: "+err.Error())
		return
	}
	defer r.Body.Close()

	if err := h.service.UpdateStatus(ctx, &req); err != nil {
		log.Printf("Error updating status: %v", err)
		response.InternalServerError(w, "Error updating status: "+err.Error())
		return
	}

	response.Success(w, "Status updated successfully", nil)
}

// UpdateTicket handles full ticket updates (PUT)
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PUT" {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req models.UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid data: "+err.Error())
		return
	}
	defer r.Body.Close()

	updatedTicket, err := h.service.UpdateTicket(ctx, &req)
	if err != nil {
		log.Printf("Error updating ticket: %v", err)
		response.InternalServerError(w, "Error updating ticket: "+err.Error())
		return
	}

	response.Success(w, "Ticket updated successfully", updatedTicket)
}

// DeleteTicket handles ticket deletion
func (h *TicketHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "DELETE" {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response.BadRequest(w, "ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.DeleteTicket(ctx, id); err != nil {
		log.Printf("Error deleting ticket: %v", err)
		response.InternalServerError(w, "Error deleting ticket: "+err.Error())
		return
	}

	response.Success(w, "Ticket deleted successfully", nil)
}
