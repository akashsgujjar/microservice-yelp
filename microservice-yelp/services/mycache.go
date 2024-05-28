package services

import (
	"context"
	"fmt"
	"log"
	"net"

	apps "gitlab.cs.washington.edu/syslab/cse453-welp/applications"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mycache"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// MyCache represents a gRPC service for interacting with a cache.
type MyCache struct {
	name string
	port int
	mycache.CacheServiceServer
	// app *apps.FIFOCacheApp
	app *apps.LRUCacheApp
	// app *apps.MRUCacheApp
}

// NewMyCache creates a new instance of MyCache.
// serverName: The name of the cache server.
// cachePort: The port on which the server should listen.
// capacity: The maximum capacity of the cache.
func NewMyCache(serverName string, cachePort int, capacity int) *MyCache {
	return &MyCache{
		name: serverName,
		port: cachePort,
		// app:  apps.NewFIFOCacheApp(capacity),
		app: apps.NewLRUCacheApp(capacity),
		// app: apps.NewCacheApp(capacity),
	}
}

// Run starts the MyCache gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *MyCache) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Cache server implementation with the gRPC server.
	mycache.RegisterCacheServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("cache server <%s> running at port: %d", s.name, s.port)
	return srv.Serve(lis)
}

// GetItem retrieves an item from the cache.
func (s *MyCache) GetItem(ctx context.Context, req *mycache.GetItemRequest) (*mycache.GetItemResponse, error) {
	cacheItem, err := s.app.Get(req.Key)
	getItemResponse := &mycache.GetItemResponse{
		Item: cacheItem,
	}
	return getItemResponse, err
}

// SetItem sets an item in the cache.
func (s *MyCache) SetItem(ctx context.Context, req *mycache.SetItemRequest) (*mycache.SetItemResponse, error) {
	err := s.app.Set(req.Item)
	setItemResponse := &mycache.SetItemResponse{
		Success: err == nil,
	}
	return setItemResponse, err
}

// DeleteItem deletes an item from the cache.
func (s *MyCache) DeleteItem(ctx context.Context, req *mycache.DeleteItemRequest) (*mycache.DeleteItemResponse, error) {
	err := s.app.Delete(req.Key)
	deleteItemResponse := &mycache.DeleteItemResponse{
		Success: err == nil,
	}
	return deleteItemResponse, err

}
