package router

import (
	"net/http"

	controller "qiscus-test/controllers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(webhookController *controller.AgentController) http.Handler {
	r := chi.NewRouter()

	r.Post("/webhook", webhookController.Webhook)

	return r
}
