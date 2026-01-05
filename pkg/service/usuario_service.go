package service

import (
	"context"
	"fmt"
	"time"

	"github.com/luispfcanales/rainforestapp/pkg/models"
	"github.com/luispfcanales/rainforestapp/pkg/repository"
)

// UsuarioService maneja la lógica de negocio de usuarios
type UsuarioService struct {
	repo *repository.UsuarioRepository
}

// NewUsuarioService crea una nueva instancia del servicio
func NewUsuarioService(repo *repository.UsuarioRepository) *UsuarioService {
	return &UsuarioService{
		repo: repo,
	}
}

// CreateUsuario crea un nuevo usuario
func (s *UsuarioService) CreateUsuario(ctx context.Context, req *models.CreateUsuarioRequest) (*models.Usuario, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Convertir a modelo
	usuario := req.ToUsuario()

	// Guardar en base de datos
	createdUsuario, err := s.repo.Create(ctx, usuario)
	if err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	return createdUsuario, nil
}

// GetUsuarioByDNI obtiene un usuario por DNI
func (s *UsuarioService) GetUsuarioByDNI(ctx context.Context, dni string) (*models.Usuario, error) {
	if dni == "" {
		return nil, fmt.Errorf("ID es requerido")
	}

	usuario, err := s.repo.GetByDNI(ctx, dni)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	return usuario, nil
}

// ListUsuarios lista todos los usuarios
func (s *UsuarioService) ListUsuarios(ctx context.Context, limit int) ([]*models.Usuario, error) {
	usuarios, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("error listando usuarios: %w", err)
	}

	return usuarios, nil
}

// UpdateUsuario actualiza un usuario existente
func (s *UsuarioService) UpdateUsuario(ctx context.Context, id string, req *models.CreateUsuarioRequest) (*models.Usuario, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	// Verificar que el usuario existe
	existingUsuario, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Actualizar campos
	// Preservar ID y timestamp de creación
	currentID := existingUsuario.ID
	createdAt := existingUsuario.CreatedAt

	// Convertir request a modelo con los nuevos datos
	newUsuario := req.ToUsuario()

	// Actualizar el usuario existente con los nuevos datos
	*existingUsuario = *newUsuario

	// Restaurar ID y metadatos
	existingUsuario.ID = currentID
	existingUsuario.CreatedAt = createdAt
	existingUsuario.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, id, existingUsuario); err != nil {
		return nil, fmt.Errorf("error actualizando usuario: %w", err)
	}

	return existingUsuario, nil
}

// DeleteUsuario elimina un usuario
func (s *UsuarioService) DeleteUsuario(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("ID es requerido")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("error eliminando usuario: %w", err)
	}

	return nil
}
