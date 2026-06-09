package http

import (
	"net/http"
	"subServ/internal/domain"
)

type SubscriptionHandler struct {
	service domain.SubscriptionService
	log     domain.Logger
}

func NewSubscriptionHandler(service domain.SubscriptionService, log domain.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		log:     log,
	}
}

func (h *SubscriptionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/subscriptions", h.Create)
	mux.HandleFunc("GET /api/v1/subscriptions", h.List)
	mux.HandleFunc("GET /api/v1/subscriptions/total-cost", h.TotalCost)
	mux.HandleFunc("GET /api/v1/subscriptions/{id}", h.GetByID)
	mux.HandleFunc("PUT /api/v1/subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/subscriptions/{id}", h.Delete)

	mux.HandleFunc("GET /swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
							<html>
							<head>
								<title>Subscription API</title>
								<meta charset="utf-8"/>
								<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
							</head>
							<body>
							<div id="swagger-ui"></div>
							<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
							<script>
								SwaggerUIBundle({
									url: "/swagger/spec",
									dom_id: '#swagger-ui',
									presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
									layout: "BaseLayout"
								})
							</script>
							</body>
							</html>`))
	})

	mux.HandleFunc("GET /swagger/spec", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, "api/openapi/swagger.yaml")
	})
}
