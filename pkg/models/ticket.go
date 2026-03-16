package models

import (
	"fmt"
	"strings"
	"time"
)

// Ticket represents a support/management ticket
type Ticket struct {
	ID           string     `json:"id,omitempty" firestore:"-"`
	TicketNumber string     `json:"ticket_number" firestore:"ticket_number"`
	Subject      string     `json:"subject" firestore:"subject"`
	RequestArea  string     `json:"request_area" firestore:"request_area"`
	Priority     string     `json:"priority" firestore:"priority"`
	Status       string     `json:"status" firestore:"status"`
	Tags         []string   `json:"tags,omitempty" firestore:"tags,omitempty"`
	Link         string     `json:"link" firestore:"link"`
	ReportDate   time.Time  `json:"report_date" firestore:"report_date"`
	ClosedDate   *time.Time `json:"closed_date,omitempty" firestore:"closed_date,omitempty"`
}

// CreateTicketRequest DTO for creating a ticket
type CreateTicketRequest struct {
	TicketNumber string   `json:"ticket_number"`
	Subject      string   `json:"subject"`
	RequestArea  string   `json:"request_area"`
	Priority     string   `json:"priority"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags,omitempty"`
	Link         string   `json:"link"`
}

// UpdateStatusRequest DTO for updating ticket status (drag & drop)
type UpdateStatusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Valid values
var (
	validPriorities = map[string]bool{
		"critico":           true,
		"Alta (Seguridad)":  true,
		"Alta (Desarrollo)": true,
		"Media (Operativo)": true,
		"Baja":              true,
	}

	validStatuses = map[string]bool{
		"nuevo":                 true,
		"en_proceso":            true,
		"en_revision":           true,
		"en_cotizacion":         true,
		"atendido_parcialmente": true,
		"cerrado":               true,
	}

	validTags = map[string]bool{
		"2do nivel":              true,
		"3er nivel":              true,
		"requiere reunión Teams": true,
	}
)

// Validate validates the ticket data
func (t *CreateTicketRequest) Validate() error {
	if strings.TrimSpace(t.TicketNumber) == "" {
		return fmt.Errorf("ticket number is required")
	}
	if strings.TrimSpace(t.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if len(t.Subject) > 400 {
		return fmt.Errorf("subject must not exceed 400 characters")
	}
	if strings.TrimSpace(t.RequestArea) == "" {
		return fmt.Errorf("request area is required")
	}
	if strings.TrimSpace(t.Link) == "" {
		return fmt.Errorf("link is required")
	}

	// Validate priority
	priority := strings.TrimSpace(t.Priority)
	if priority == "" {
		t.Priority = "Media (Operativo)"
	} else if !validPriorities[priority] {
		return fmt.Errorf("invalid priority: %s", priority)
	}

	// Validate status
	status := strings.TrimSpace(t.Status)
	if status == "" {
		t.Status = "nuevo"
	} else if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Validate tags
	for _, tag := range t.Tags {
		if !validTags[tag] {
			return fmt.Errorf("invalid tag: %s", tag)
		}
	}

	return nil
}

// FormatTicketNumber formats the ticket number as "#HT" + 6 digits
func FormatTicketNumber(number string) string {
	padded := number
	for len(padded) < 6 {
		padded = "0" + padded
	}
	return "#HT" + padded
}

// ToTicket converts the request to a Ticket model
func (t *CreateTicketRequest) ToTicket() *Ticket {
	return &Ticket{
		TicketNumber: FormatTicketNumber(strings.TrimSpace(t.TicketNumber)),
		Subject:      strings.TrimSpace(t.Subject),
		RequestArea:  strings.TrimSpace(t.RequestArea),
		Priority:     strings.TrimSpace(t.Priority),
		Status:       strings.TrimSpace(t.Status),
		Tags:         t.Tags,
		Link:         strings.TrimSpace(t.Link),
		ReportDate:   time.Now(),
	}
}

// ValidateStatus validates the status update request
func (u *UpdateStatusRequest) ValidateStatus() error {
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("ticket ID is required")
	}
	if !validStatuses[strings.TrimSpace(u.Status)] {
		return fmt.Errorf("invalid status: %s", u.Status)
	}
	return nil
}

// UpdateTicketRequest DTO for updating a ticket completely (PUT)
type UpdateTicketRequest struct {
	ID           string   `json:"id"`
	TicketNumber string   `json:"ticket_number"`
	Subject      string   `json:"subject"`
	RequestArea  string   `json:"request_area"`
	Priority     string   `json:"priority"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags,omitempty"`
	Link         string   `json:"link"`
	ReportDate   string   `json:"report_date,omitempty"` // RFC3339 format
	ClosedDate   string   `json:"closed_date,omitempty"` // RFC3339 format
}

// ValidateUpdate validates the update request
func (t *UpdateTicketRequest) ValidateUpdate() error {
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("ticket ID is required")
	}
	if strings.TrimSpace(t.TicketNumber) == "" {
		return fmt.Errorf("ticket number is required")
	}
	if strings.TrimSpace(t.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if len(t.Subject) > 400 {
		return fmt.Errorf("subject must not exceed 400 characters")
	}
	if strings.TrimSpace(t.RequestArea) == "" {
		return fmt.Errorf("request area is required")
	}
	if strings.TrimSpace(t.Link) == "" {
		return fmt.Errorf("link is required")
	}

	// Validate priority
	priority := strings.TrimSpace(t.Priority)
	if priority == "" {
		t.Priority = "Media (Operativo)"
	} else if !validPriorities[priority] {
		return fmt.Errorf("invalid priority: %s", priority)
	}

	// Validate status
	status := strings.TrimSpace(t.Status)
	if status == "" {
		t.Status = "nuevo"
	} else if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Validate tags
	for _, tag := range t.Tags {
		if !validTags[tag] {
			return fmt.Errorf("invalid tag: %s", tag)
		}
	}

	// Validate date if present
	if strings.TrimSpace(t.ReportDate) != "" {
		_, err := time.Parse(time.RFC3339, t.ReportDate)
		if err != nil {
			return fmt.Errorf("invalid report_date format, must be RFC3339 (e.g. 2026-03-04T08:38:35Z)")
		}
	}

	if strings.TrimSpace(t.ClosedDate) != "" {
		_, err := time.Parse(time.RFC3339, t.ClosedDate)
		if err != nil {
			return fmt.Errorf("invalid closed_date format, must be RFC3339 (e.g. 2026-03-04T08:38:35Z)")
		}
	}

	return nil
}
