package clients

import (
	"api-gateway/models"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientFactory struct {
	Config      *models.Config
	connections map[string]*grpc.ClientConn
}

func NewGRPCClientFactory(config *models.Config) *GRPCClientFactory {
	return &GRPCClientFactory{
		Config:      config,
		connections: make(map[string]*grpc.ClientConn),
	}
}

func (f *GRPCClientFactory) GetClient(serviceName, address string) (*grpc.ClientConn, error) {
	if conn, exists := f.connections[serviceName]; exists {
		if conn.GetState() == connectivity.Ready {
			return conn, nil
		}
		conn.Close()
		delete(f.connections, serviceName)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to create a new gRPC client for %s at %s: %v", serviceName, address, err)
		return nil, err
	}

	f.connections[serviceName] = conn
	return conn, nil
}
