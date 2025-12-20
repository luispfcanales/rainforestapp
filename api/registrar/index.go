package handler

import (
	"encoding/json"
	"net/http"
)

// Handler principal para Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS (¡IMPORTANTE para frontend!)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Manejar preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo aceptar POST
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Procesar la solicitud
	// response := procesarRegistro(r)

	// Enviar respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"hola": "how to change"})
}
