package database

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"github.com/luispfcanales/rainforestapp/internal/config"
)

var (
	firestoreClient *firestore.Client
	once            sync.Once
	initError       error
)

// GetFirestoreClient obtiene o crea una instancia del cliente de Firestore
func GetFirestoreClient(ctx context.Context, cfg *config.Config) (*firestore.Client, error) {
	once.Do(func() {
		firestoreClient, initError = initializeFirestore(ctx, cfg)
	})

	return firestoreClient, initError
}

// initializeFirestore inicializa el cliente de Firestore
func initializeFirestore(ctx context.Context, cfg *config.Config) (*firestore.Client, error) {
	// Configurar Firebase
	firebaseConfig := &firebase.Config{
		ProjectID: cfg.Firebase.ProjectID,
	}

	opt := option.WithCredentialsJSON([]byte(cfg.Firebase.Credentials))

	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		return nil, fmt.Errorf("error inicializando Firebase App: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo cliente Firestore: %w", err)
	}

	return client, nil
}

// Close cierra la conexión de Firestore
func Close() error {
	if firestoreClient != nil {
		return firestoreClient.Close()
	}
	return nil
}
