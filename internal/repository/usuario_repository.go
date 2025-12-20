package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/luispfcanales/rainforestapp/internal/models"
)

const usuariosCollection = "usuarios"

// UsuarioRepository maneja las operaciones de base de datos para usuarios
type UsuarioRepository struct {
	client *firestore.Client
}

// NewUsuarioRepository crea una nueva instancia del repositorio
func NewUsuarioRepository(client *firestore.Client) *UsuarioRepository {
	return &UsuarioRepository{
		client: client,
	}
}

// Create guarda un nuevo usuario en Firestore
func (r *UsuarioRepository) Create(ctx context.Context, usuario *models.Usuario) (*models.Usuario, error) {
	docRef, _, err := r.client.Collection(usuariosCollection).Add(ctx, usuario)
	if err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	usuario.ID = docRef.ID
	return usuario, nil
}

// GetByID obtiene un usuario por su ID
func (r *UsuarioRepository) GetByID(ctx context.Context, id string) (*models.Usuario, error) {
	doc, err := r.client.Collection(usuariosCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuario: %w", err)
	}

	var usuario models.Usuario
	if err := doc.DataTo(&usuario); err != nil {
		return nil, fmt.Errorf("error parseando usuario: %w", err)
	}

	usuario.ID = doc.Ref.ID
	return &usuario, nil
}

// GetAll obtiene todos los usuarios
func (r *UsuarioRepository) GetAll(ctx context.Context, limit int) ([]*models.Usuario, error) {
	query := r.client.Collection(usuariosCollection).OrderBy("created_at", firestore.Desc)

	if limit > 0 {
		query = query.Limit(limit)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var usuarios []*models.Usuario
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterando usuarios: %w", err)
		}

		var usuario models.Usuario
		if err := doc.DataTo(&usuario); err != nil {
			continue
		}
		usuario.ID = doc.Ref.ID
		usuarios = append(usuarios, &usuario)
	}

	return usuarios, nil
}

// Update actualiza un usuario existente
func (r *UsuarioRepository) Update(ctx context.Context, id string, usuario *models.Usuario) error {
	_, err := r.client.Collection(usuariosCollection).Doc(id).Set(ctx, usuario, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}
	return nil
}

// Delete elimina un usuario
func (r *UsuarioRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(usuariosCollection).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %w", err)
	}
	return nil
}
