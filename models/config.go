package models

type Config struct {
	Routes   map[string]RouteConfig   `json:"routes"`
	Services map[string]ServiceConfig `json:"services"`
}

type RouteConfig struct {
	Service     string            `json:"service"`
	Method      string            `json:"method"`
	GRPCService string            `json:"grpc_service"`
	RequestMap  map[string]string `json:"request_map"`
}

type ServiceConfig struct {
	Address string `json:"address"`
}
