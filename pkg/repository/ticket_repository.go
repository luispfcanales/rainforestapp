package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/luispfcanales/rainforestapp/pkg/models"
)

const ticketsCollection = "tickets"

// TicketRepository handles database operations for tickets
type TicketRepository struct {
	client *firestore.Client
}

// NewTicketRepository creates a new repository instance
func NewTicketRepository(client *firestore.Client) *TicketRepository {
	return &TicketRepository{
		client: client,
	}
}

// Create saves a new ticket in Firestore
func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	docRef, _, err := r.client.Collection(ticketsCollection).Add(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket: %w", err)
	}

	ticket.ID = docRef.ID
	return ticket, nil
}

// GetByID retrieves a ticket by its ID
func (r *TicketRepository) GetByID(ctx context.Context, id string) (*models.Ticket, error) {
	doc, err := r.client.Collection(ticketsCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	var ticket models.Ticket
	if err := doc.DataTo(&ticket); err != nil {
		return nil, fmt.Errorf("error parsing ticket: %w", err)
	}

	ticket.ID = doc.Ref.ID
	return &ticket, nil
}

// GetAll retrieves all tickets
func (r *TicketRepository) GetAll(ctx context.Context, limit int) ([]*models.Ticket, error) {
	query := r.client.Collection(ticketsCollection).OrderBy("report_date", firestore.Desc)

	if limit > 0 {
		query = query.Limit(limit)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var tickets []*models.Ticket
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating tickets: %w", err)
		}

		var ticket models.Ticket
		if err := doc.DataTo(&ticket); err != nil {
			continue
		}
		ticket.ID = doc.Ref.ID
		tickets = append(tickets, &ticket)
	}

	return tickets, nil
}

// UpdateStatus updates the status of a ticket, and optionally the closed date
func (r *TicketRepository) UpdateStatus(ctx context.Context, id string, status string, closedDate *time.Time) error {
	updates := []firestore.Update{
		{Path: "status", Value: status},
	}

	if closedDate != nil {
		updates = append(updates, firestore.Update{Path: "closed_date", Value: *closedDate})
	}

	_, err := r.client.Collection(ticketsCollection).Doc(id).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("error updating ticket status: %w", err)
	}
	return nil
}

// Update updates an entire ticket
func (r *TicketRepository) Update(ctx context.Context, id string, ticket *models.Ticket) error {
	_, err := r.client.Collection(ticketsCollection).Doc(id).Set(ctx, ticket)
	if err != nil {
		return fmt.Errorf("error updating ticket: %w", err)
	}
	return nil
}

// Delete removes a ticket
func (r *TicketRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(ticketsCollection).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting ticket: %w", err)
	}
	return nil
}
