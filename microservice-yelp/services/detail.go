package services

import (
	"context"
	"fmt"
	"log"
	"net"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/detail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mycache"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mydatabase"
)

// Detail implements the detail service
type Detail struct {
	name string
	port int
	detail.DetailServiceServer

	// mu sync.Mutex
	// dataStore map[string][]byte

	detailCacheClient    mycache.CacheServiceClient       // Add detail grpc cache client for communicating with detail cache server
	detailDatabaseClient mydatabase.DatabaseServiceClient // Add detail grpc storage client for communicating with detail storage server

	CACHE_FLAG bool
}

// NewDetail returns a new server
func NewDetail(name string, detailPort int, detailCacheAddr string, detailDatabaseAddr string) *Detail {
	return &Detail{
		name: name,
		port: detailPort,
		// dataStore: make(map[string][]byte),
		detailCacheClient:    mycache.NewCacheServiceClient(dial(detailCacheAddr)),          // Initialize and establish cxn using specified address
		detailDatabaseClient: mydatabase.NewDatabaseServiceClient(dial(detailDatabaseAddr)), // Initialize and establish cxn using specified address
		CACHE_FLAG:           true,
	}
}

// Run starts the Detail gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *Detail) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Detail server implementation with the gRPC server.
	detail.RegisterDetailServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("detail server running at port: %d", s.port)

	// Start serving incoming requests using the registered implementation.
	return srv.Serve(lis)
}

// GetDetail retrieves the details of a restaurant.
// It first checks if the data is cached in mycache.
// If not, it retrieves the data from mydb and stores it in mycache for future use.
// It returns an error if the requested restaurant does not exist.
func (s *Detail) GetDetail(ctx context.Context, req *detail.GetDetailRequest) (*detail.GetDetailResponse, error) {
	// Get the name of the requested restaurant
	restaurantName := req.GetRestaurantName()
	detailResponse := &detail.GetDetailResponse{}

	var err error

	getCacheItemMsg := &mycache.GetItemRequest{
		Key: restaurantName,
	}

	// Check if the details are stored in the cache
	var getItemResponse *mycache.GetItemResponse
	var errGetItem error

	getItemResponse, errGetItem = s.detailCacheClient.GetItem(ctx, getCacheItemMsg)
	if errGetItem == nil {
		rawValue := getItemResponse.Item.Value
		err = proto.Unmarshal(rawValue, detailResponse)
		if err != nil {
			err = status.Errorf(codes.Internal, "Failed to deserialize data")
		} else {
			err = status.Errorf(codes.OK, "Found value with Key: %s", restaurantName)
		}
	} else {
		// Make a call to the storage layer due to a cache-miss
		getRecordMsg := &mydatabase.GetRecordRequest{
			Key: restaurantName,
		}

		if getRecordResponse, errGetRecord := s.detailDatabaseClient.GetRecord(ctx, getRecordMsg); errGetRecord == nil {
			rawValue := getRecordResponse.Record.Value
			err = proto.Unmarshal(rawValue, detailResponse)
			if err != nil {
				err = status.Errorf(codes.Internal, "Failed to deserialize data")
			} else {
				// If we found the item in the storage layer, then update the cache to hold this element
				if s.CACHE_FLAG {
					updateCache(ctx, s.detailCacheClient, restaurantName, rawValue)
				}

				err = status.Errorf(codes.OK, "Found value with Key: %s", restaurantName)
			}
		} else {
			return detailResponse, status.Errorf(codes.NotFound, "Item with Key: %s does not exist", restaurantName)
		}
	}
	return detailResponse, err
}

// PostDetail adds or updates the details of a restaurant in the in-memory dataStore.
func (s *Detail) PostDetail(ctx context.Context, req *detail.PostDetailRequest) (*detail.PostDetailResponse, error) {
	// Extract fields from the request.
	restaurantName := req.GetRestaurantName()
	location := req.GetLocation()
	capacity := req.GetCapacity()
	style := req.GetStyle()

	// Create a new GetDetailResponse object with the details to save.
	msg := &detail.GetDetailResponse{
		RestaurantName: restaurantName,
		Location:       location,
		Style:          style,
		Capacity:       capacity,
	}

	// Initialize an empty response object.
	detailResponse := &detail.PostDetailResponse{Status: false}

	// Serialize the Go object into protobuf format.
	data, err := proto.Marshal(msg)
	if err != nil {
		// If serialization fails, return an internal error.
		return detailResponse, status.Errorf(codes.Internal, "Failed to serialize data")
	}

	// add detail to DB (updates cache as well)
	detailResponse.Status, err = updateDB(ctx, s.detailCacheClient, s.detailDatabaseClient, restaurantName, data, s.CACHE_FLAG)

	// Return the response object and any error.
	return detailResponse, err
}
