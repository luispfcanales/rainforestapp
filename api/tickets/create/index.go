package handler

import (
	"net/http"

	"github.com/luispfcanales/rainforestapp/pkg/config"
	ticketHandler "github.com/luispfcanales/rainforestapp/pkg/handler"
	"github.com/luispfcanales/rainforestapp/pkg/response"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Load()
	if err != nil {
		response.InternalServerError(w, "Configuration error")
		return
	}

	h, err := ticketHandler.NewTicketHandler(cfg)
	if err != nil {
		response.InternalServerError(w, "Error initializing handler")
		return
	}

	h.CreateTicket(w, r)
}
