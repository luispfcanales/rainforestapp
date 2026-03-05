package service

import (
	"context"
	"fmt"
	"time"

	"github.com/luispfcanales/rainforestapp/pkg/models"
	"github.com/luispfcanales/rainforestapp/pkg/repository"
)

// TicketService handles business logic for tickets
type TicketService struct {
	repo *repository.TicketRepository
}

// NewTicketService creates a new service instance
func NewTicketService(repo *repository.TicketRepository) *TicketService {
	return &TicketService{
		repo: repo,
	}
}

// CreateTicket creates a new ticket
func (s *TicketService) CreateTicket(ctx context.Context, req *models.CreateTicketRequest) (*models.Ticket, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	ticket := req.ToTicket()

	createdTicket, err := s.repo.Create(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket: %w", err)
	}

	return createdTicket, nil
}

// GetTicketByID retrieves a ticket by ID
func (s *TicketService) GetTicketByID(ctx context.Context, id string) (*models.Ticket, error) {
	if id == "" {
		return nil, fmt.Errorf("ID is required")
	}

	ticket, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	return ticket, nil
}

// ListTickets lists all tickets
func (s *TicketService) ListTickets(ctx context.Context, limit int) ([]*models.Ticket, error) {
	tickets, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("error listing tickets: %w", err)
	}

	return tickets, nil
}

// UpdateStatus updates only the status of a ticket (for drag & drop)
func (s *TicketService) UpdateStatus(ctx context.Context, req *models.UpdateStatusRequest) error {
	if err := req.ValidateStatus(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	_, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	var closedDate *time.Time
	if req.Status == "cerrado" {
		now := time.Now()
		closedDate = &now
	}

	return s.repo.UpdateStatus(ctx, req.ID, req.Status, closedDate)
}

// UpdateTicket updates an entire ticket
func (s *TicketService) UpdateTicket(ctx context.Context, req *models.UpdateTicketRequest) (*models.Ticket, error) {
	if err := req.ValidateUpdate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existingTicket, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Update fields
	existingTicket.TicketNumber = models.FormatTicketNumber(req.TicketNumber)
	existingTicket.Subject = req.Subject
	existingTicket.RequestArea = req.RequestArea
	existingTicket.Priority = req.Priority
	existingTicket.Status = req.Status
	existingTicket.Tags = req.Tags
	existingTicket.Link = req.Link

	if req.ReportDate != "" {
		t, _ := time.Parse(time.RFC3339, req.ReportDate)
		existingTicket.ReportDate = t
	}

	// Handle ClosedDate
	if req.ClosedDate != "" {
		t, _ := time.Parse(time.RFC3339, req.ClosedDate)
		existingTicket.ClosedDate = &t
	} else if req.Status == "cerrado" && existingTicket.ClosedDate == nil {
		// Auto-set if changing to cerrado and no explicit date is provided
		now := time.Now()
		existingTicket.ClosedDate = &now
	} else if req.Status != "cerrado" {
		// Clear closed date if status is no longer cerrado
		existingTicket.ClosedDate = nil
	}

	if err := s.repo.Update(ctx, req.ID, existingTicket); err != nil {
		return nil, fmt.Errorf("error updating ticket: %w", err)
	}

	return existingTicket, nil
}

// DeleteTicket deletes a ticket
func (s *TicketService) DeleteTicket(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("ID is required")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting ticket: %w", err)
	}

	return nil
}
