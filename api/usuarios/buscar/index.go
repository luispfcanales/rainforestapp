package handler

import (
	"net/http"

	"github.com/luispfcanales/rainforestapp/pkg/config"
	usuarioHandler "github.com/luispfcanales/rainforestapp/pkg/handler"
	"github.com/luispfcanales/rainforestapp/pkg/response"
)

// Handler es la función principal para Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		response.InternalServerError(w, "Error de configuración")
		return
	}

	h, err := usuarioHandler.NewUsuarioHandler(cfg)
	if err != nil {
		response.InternalServerError(w, "Error inicializando handler")
		return
	}

	// Delegar al handler para obtener detalle de usuario
	h.GetUsuario(w, r)
}
