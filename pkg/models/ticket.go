package models

import (
	"fmt"
	"strings"
	"time"
)

// Observation represents a custom note or comment on a ticket with its creation date
type Observation struct {
	Text string `json:"text" firestore:"text"`
	Date string `json:"date" firestore:"date"`
}

// Ticket represents a support/management ticket
type Ticket struct {
	ID               string        `json:"id,omitempty" firestore:"-"`
	TicketNumber     string        `json:"ticket_number" firestore:"ticket_number"`
	Subject          string        `json:"subject" firestore:"subject"`
	RequestArea      string        `json:"request_area" firestore:"request_area"`
	Priority         string        `json:"priority" firestore:"priority"`
	Status           string        `json:"status" firestore:"status"`
	Tags             []string      `json:"tags,omitempty" firestore:"tags,omitempty"`
	Link             string        `json:"link" firestore:"link"`
	ReportDate       time.Time     `json:"report_date" firestore:"report_date"`
	ClosedDate       *time.Time    `json:"closed_date,omitempty" firestore:"closed_date,omitempty"`
	ThemResponseDate *time.Time    `json:"them_response_date,omitempty" firestore:"them_response_date,omitempty"`
	ThemResponse     string        `json:"them_response" firestore:"them_response"`
	MyResponseDate   *time.Time    `json:"my_response_date,omitempty" firestore:"my_response_date,omitempty"`
	MyResponse       string        `json:"my_response" firestore:"my_response"`
	RelatedTicketID  string        `json:"related_ticket_id" firestore:"related_ticket_id"`
	IsPaid           bool          `json:"is_paid" firestore:"is_paid"`
	QuotePDFData     string        `json:"quote_pdf_data" firestore:"quote_pdf_data"`
	Type             string        `json:"type" firestore:"type"`
	AssignedTo       string        `json:"assigned_to" firestore:"assigned_to"`
	Observations     []Observation `json:"observations,omitempty" firestore:"observations,omitempty"`
}

// CreateTicketRequest DTO for creating a ticket
type CreateTicketRequest struct {
	TicketNumber     string   `json:"ticket_number"`
	Subject          string   `json:"subject"`
	RequestArea      string   `json:"request_area"`
	Priority         string   `json:"priority"`
	Status           string   `json:"status"`
	Tags             []string `json:"tags,omitempty"`
	Link             string   `json:"link"`
	ReportDate       string   `json:"report_date,omitempty"` // RFC3339 or blank
	ClosedDate       string   `json:"closed_date,omitempty"` // RFC3339 or blank
	ThemResponseDate string   `json:"them_response_date,omitempty"`
	ThemResponse     string   `json:"them_response,omitempty"`
	MyResponseDate   string   `json:"my_response_date,omitempty"`
	MyResponse       string   `json:"my_response,omitempty"`
	RelatedTicketID  string   `json:"related_ticket_id,omitempty"`
	IsPaid           bool     `json:"is_paid,omitempty"`
	QuotePDFData     string   `json:"quote_pdf_data,omitempty"`
	Type             string        `json:"type,omitempty"`
	AssignedTo       string        `json:"assigned_to,omitempty"`
	Observations     []Observation `json:"observations,omitempty"`
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
		"cotización":             true,
		"cotizacion":             true,
		"atendido parcialmente":  true,
		"en revisión por RFE":    true,
		"en revision por RFE":    true,
	}

	// statusAliases maps Spanish display values (from Excel) to internal status keys
	statusAliases = map[string]string{
		"en proceso":             "en_proceso",
		"en revisión":            "en_revision",
		"en revision":            "en_revision",
		"en cotización":          "en_cotizacion",
		"en cotizacion":          "en_cotizacion",
		"atendido parcialmente":  "atendido_parcialmente",
		"resuelto":               "cerrado",
		"cerrado":                "cerrado",
		"nuevo":                  "nuevo",
		"2do nivel":              "en_revision",
		"3er nivel":              "en_revision",
		"en 2do nivel":           "en_revision",
		"en 3er nivel":           "en_revision",
	}

	// priorityAliases maps common Spanish display values to internal priority keys
	priorityAliases = map[string]string{
		"critico":            "critico",
		"crítico":            "critico",
		"alta (seguridad)":   "Alta (Seguridad)",
		"alta (desarrollo)":  "Alta (Desarrollo)",
		"alta":               "Alta (Seguridad)",
		"media (operativo)": "Media (Operativo)",
		"media":              "Media (Operativo)",
		"baja":               "Baja",
	}
)

// NormalizeStatus converts a free-text status (e.g. from Excel) to a valid internal key.
// Returns the normalized value and true if found, otherwise the original and false.
func NormalizeStatus(raw string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(raw))
	if validStatuses[key] {
		return key, true
	}
	if v, ok := statusAliases[key]; ok {
		return v, true
	}
	return raw, false
}

// NormalizePriority converts a free-text priority to a valid internal key.
func NormalizePriority(raw string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(raw))
	if validPriorities[raw] {
		return raw, true
	}
	if v, ok := priorityAliases[key]; ok {
		return v, true
	}
	return raw, false
}

func appendIfMissing(slice []string, val string) []string {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return slice
		}
	}
	return append(slice, val)
}

// Validate validates the ticket data
func (t *CreateTicketRequest) Validate() error {
	if strings.TrimSpace(t.TicketNumber) == "" {
		return fmt.Errorf("ticket number is required")
	}
	if strings.TrimSpace(t.Subject) == "" {
		return fmt.Errorf("subject is required")
	}
	if len(t.Subject) > 400 {
		t.Subject = t.Subject[:400]
	}
	if strings.TrimSpace(t.RequestArea) == "" {
		return fmt.Errorf("request area is required")
	}
	// link is optional — it may not be present in imported rows

	// Normalize + validate priority
	priority := strings.TrimSpace(t.Priority)
	if priority == "" {
		t.Priority = "Media (Operativo)"
	} else {
		norm, ok := NormalizePriority(priority)
		if !ok {
			t.Priority = "Media (Operativo)" // fallback
		} else {
			t.Priority = norm
		}
	}

	// Normalize + validate status
	status := strings.TrimSpace(t.Status)
	if status == "" {
		t.Status = "nuevo"
	} else {
		norm, ok := NormalizeStatus(status)
		if !ok {
			t.Status = "nuevo" // fallback
		} else {
			if norm == "en_cotizacion" {
				t.Status = "en_proceso"
				t.Tags = appendIfMissing(t.Tags, "cotización")
			} else if norm == "atendido_parcialmente" {
				t.Status = "en_proceso"
				t.Tags = appendIfMissing(t.Tags, "atendido parcialmente")
			} else if norm == "en_revision" {
				t.Status = "en_proceso"
				t.Tags = appendIfMissing(t.Tags, "en revisión por RFE")
			} else {
				t.Status = norm
			}
		}
	}

	// Validate tags
	for _, tag := range t.Tags {
		if !validTags[tag] {
			return fmt.Errorf("invalid tag: %s", tag)
		}
	}

	// Validate date formats if present
	if strings.TrimSpace(t.ClosedDate) != "" {
		_, err := time.Parse(time.RFC3339, t.ClosedDate)
		if err != nil {
			return fmt.Errorf("invalid closed_date format, must be RFC3339")
		}
	}
	if strings.TrimSpace(t.ThemResponseDate) != "" {
		_, err := time.Parse(time.RFC3339, t.ThemResponseDate)
		if err != nil {
			return fmt.Errorf("invalid them_response_date format, must be RFC3339")
		}
	}
	if strings.TrimSpace(t.MyResponseDate) != "" {
		_, err := time.Parse(time.RFC3339, t.MyResponseDate)
		if err != nil {
			return fmt.Errorf("invalid my_response_date format, must be RFC3339")
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
	reportDate := time.Now()
	if rd := strings.TrimSpace(t.ReportDate); rd != "" {
		if parsed, err := time.Parse(time.RFC3339, rd); err == nil {
			reportDate = parsed
		} else if parsed, err := time.Parse("2006-01-02", rd[:10]); err == nil {
			reportDate = parsed
		}
	}

	ticket := &Ticket{
		TicketNumber:    strings.TrimSpace(t.TicketNumber),
		Subject:         strings.TrimSpace(t.Subject),
		RequestArea:     strings.TrimSpace(t.RequestArea),
		Priority:        strings.TrimSpace(t.Priority),
		Status:          strings.TrimSpace(t.Status),
		Tags:            t.Tags,
		Link:            strings.TrimSpace(t.Link),
		ReportDate:      reportDate,
		ThemResponse:    t.ThemResponse,
		MyResponse:      t.MyResponse,
		RelatedTicketID: t.RelatedTicketID,
		IsPaid:          t.IsPaid,
		QuotePDFData:    t.QuotePDFData,
		Type:            strings.TrimSpace(t.Type),
		AssignedTo:      strings.TrimSpace(t.AssignedTo),
		Observations:    t.Observations,
	}

	if strings.TrimSpace(t.ClosedDate) != "" {
		dt, _ := time.Parse(time.RFC3339, t.ClosedDate)
		ticket.ClosedDate = &dt
	} else if strings.TrimSpace(t.Status) == "cerrado" {
		now := time.Now()
		ticket.ClosedDate = &now
	}

	if strings.TrimSpace(t.ThemResponseDate) != "" {
		dt, _ := time.Parse(time.RFC3339, t.ThemResponseDate)
		ticket.ThemResponseDate = &dt
	}
	if strings.TrimSpace(t.MyResponseDate) != "" {
		dt, _ := time.Parse(time.RFC3339, t.MyResponseDate)
		ticket.MyResponseDate = &dt
	}

	return ticket
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
	ID               string   `json:"id"`
	TicketNumber     string   `json:"ticket_number"`
	Subject          string   `json:"subject"`
	RequestArea      string   `json:"request_area"`
	Priority         string   `json:"priority"`
	Status           string   `json:"status"`
	Tags             []string `json:"tags,omitempty"`
	Link             string   `json:"link"`
	ReportDate       string   `json:"report_date,omitempty"` // RFC3339 format
	ClosedDate       string   `json:"closed_date,omitempty"` // RFC3339 format
	ThemResponseDate string   `json:"them_response_date,omitempty"`
	ThemResponse     string   `json:"them_response,omitempty"`
	MyResponseDate   string   `json:"my_response_date,omitempty"`
	MyResponse       string   `json:"my_response,omitempty"`
	RelatedTicketID  string   `json:"related_ticket_id,omitempty"`
	IsPaid           bool     `json:"is_paid,omitempty"`
	QuotePDFData     string   `json:"quote_pdf_data,omitempty"`
	Type             string        `json:"type,omitempty"`
	AssignedTo       string        `json:"assigned_to,omitempty"`
	Observations     []Observation `json:"observations,omitempty"`
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
	} else {
		norm, ok := NormalizeStatus(status)
		if !ok {
			return fmt.Errorf("invalid status: %s", status)
		}
		if norm == "en_cotizacion" {
			t.Status = "en_proceso"
			t.Tags = appendIfMissing(t.Tags, "cotización")
		} else if norm == "atendido_parcialmente" {
			t.Status = "en_proceso"
			t.Tags = appendIfMissing(t.Tags, "atendido parcialmente")
		} else if norm == "en_revision" {
			t.Status = "en_proceso"
			t.Tags = appendIfMissing(t.Tags, "en revisión por RFE")
		} else {
			t.Status = norm
		}
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

	if strings.TrimSpace(t.ThemResponseDate) != "" {
		_, err := time.Parse(time.RFC3339, t.ThemResponseDate)
		if err != nil {
			return fmt.Errorf("invalid them_response_date format, must be RFC3339 (e.g. 2026-03-04T08:38:35Z)")
		}
	}

	if strings.TrimSpace(t.MyResponseDate) != "" {
		_, err := time.Parse(time.RFC3339, t.MyResponseDate)
		if err != nil {
			return fmt.Errorf("invalid my_response_date format, must be RFC3339 (e.g. 2026-03-04T08:38:35Z)")
		}
	}

	return nil
}
