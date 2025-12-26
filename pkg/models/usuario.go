package models

import (
	"fmt"
	"strings"
	"time"
)

// Usuario representa un usuario en el sistema
type Usuario struct {
	ID        string    `json:"id,omitempty" firestore:"-"`
	Nombre    string    `json:"nombre" firestore:"nombre"`
	Apellido  string    `json:"apellido" firestore:"apellido"`
	Email     string    `json:"email,omitempty" firestore:"email,omitempty"`
	Dni       string    `json:"dni,omitempty" firestore:"dni,omitempty"`
	Telefono  string    `json:"telefono,omitempty" firestore:"telefono,omitempty"`
	Foto      string    `json:"foto,omitempty" firestore:"foto,omitempty"`
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" firestore:"updated_at,omitempty"`
}

// CreateUsuarioRequest DTO para crear usuario
type CreateUsuarioRequest struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email,omitempty"`
	Dni      string `json:"dni"`
	Telefono string `json:"telefono,omitempty"`
	Foto     string `json:"foto,omitempty"`
}

// Validate valida los datos del usuario
func (u *CreateUsuarioRequest) Validate() error {
	if strings.TrimSpace(u.Nombre) == "" {
		return fmt.Errorf("el nombre es requerido")
	}
	if strings.TrimSpace(u.Apellido) == "" {
		return fmt.Errorf("el apellido es requerido")
	}
	if len(u.Nombre) < 2 {
		return fmt.Errorf("el nombre debe tener al menos 2 caracteres")
	}
	if len(u.Apellido) < 2 {
		return fmt.Errorf("el apellido debe tener al menos 2 caracteres")
	}
	return nil
}

// ToUsuario convierte el request a un modelo Usuario
func (u *CreateUsuarioRequest) ToUsuario() *Usuario {
	return &Usuario{
		Nombre:    strings.TrimSpace(u.Nombre),
		Apellido:  strings.TrimSpace(u.Apellido),
		Email:     strings.TrimSpace(u.Email),
		Dni:       strings.TrimSpace(u.Dni),
		Telefono:  strings.TrimSpace(u.Telefono),
		Foto:      strings.TrimSpace(u.Foto),
		CreatedAt: time.Now(),
	}
}
