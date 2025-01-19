package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"api-gateway/clients"
	"api-gateway/models"

	"github.com/codepumper/protos/auth"
)

type GRPCHandler struct {
	clientFactory *clients.GRPCClientFactory
	routes        map[string]models.RouteConfig
}

func NewGRPCHandler(clientFactory *clients.GRPCClientFactory, routes map[string]models.RouteConfig) *GRPCHandler {
	return &GRPCHandler{
		clientFactory: clientFactory,
		routes:        routes,
	}
}

func (h *GRPCHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	route, exists := h.routes[r.URL.Path]
	if !exists {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	serviceConfig, ok := h.clientFactory.Config.Services[route.Service]
	if !ok {
		http.Error(w, "Service configuration not found", http.StatusInternalServerError)
		return
	}

	conn, err := h.clientFactory.GetClient(route.Service, serviceConfig.Address)
	if err != nil {
		http.Error(w, "Failed to connect to gRPC service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	var response interface{}
	switch route.GRPCService {
	case "AuthService":
		client := auth.NewAuthServiceClient(conn)
		switch route.Method {
		case "Login":
			req := &auth.LoginRequest{
				Email:    payload[route.RequestMap["email"]].(string),
				Password: payload[route.RequestMap["password"]].(string),
			}
			response, err = client.Login(context.Background(), req)
		case "Register":
			req := &auth.RegisterRequest{
				Email:    payload[route.RequestMap["email"]].(string),
				Password: payload[route.RequestMap["password"]].(string),
			}
			response, err = client.Register(context.Background(), req)
		default:
			http.Error(w, "Unknown method", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unknown gRPC service", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
