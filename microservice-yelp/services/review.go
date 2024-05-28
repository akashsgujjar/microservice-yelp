package services

import (
	"context"
	"fmt"
	"log"
	"net"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/review"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mycache"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mydatabase"
)

// Review implements the review service
type Review struct {
	name string
	port int
	review.ReviewServiceServer

	// reviewTracker map[string][]byte
	reviewCacheClient    mycache.CacheServiceClient       // Add review grpc cache client for communicating with review cache server
	reviewDatabaseClient mydatabase.DatabaseServiceClient // Add review grpc storage client for communicating with review storage server

	CACHE_FLAG bool
}

// NewReview returns a new server
func NewReview(name string, reviewPort int, reviewCacheAddr string, reviewDatabaseAddr string) *Review {
	return &Review{
		name: name,
		port: reviewPort,

		reviewCacheClient:    mycache.NewCacheServiceClient(dial(reviewCacheAddr)),          // Initialize and establish cxn using specified address
		reviewDatabaseClient: mydatabase.NewDatabaseServiceClient(dial(reviewDatabaseAddr)), // Initialize and establish cxn using specified address
		CACHE_FLAG:           true,
	}
}

// Run starts the Review gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *Review) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Review server implementation with the gRPC server.
	review.RegisterReviewServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("review server running at port: %d", s.port)

	// Start serving incoming requests using the registered implementation.
	return srv.Serve(lis)
}

func (s *Review) GetReview(ctx context.Context, req *review.GetReviewRequest) (*review.GetReviewResponse, error) {
	username := req.GetUserName()
	restaurant_name := req.GetRestaurantName()

	reviewResponse := &review.GetReviewResponse{}
	searchResponse := &review.SearchReviewsResponse{}

	var err error

	getCacheItemMsg := &mycache.GetItemRequest{
		Key: restaurant_name,
	}

	var getItemResponse *mycache.GetItemResponse
	var errGetItem error

	// Check the cache for the review
	getItemResponse, errGetItem = s.reviewCacheClient.GetItem(ctx, getCacheItemMsg)
	if errGetItem == nil {
		rawValue := getItemResponse.Item.Value
		err = proto.Unmarshal(rawValue, searchResponse)
		if err != nil {
			err = status.Errorf(codes.Internal, "Failed to deserialize data")
		} else {
			err = status.Errorf(codes.OK, "Found value with Key: %s", restaurant_name)
		}
	} else {
		// Make a call to the storage layer due to a cache-miss
		getRecordMsg := &mydatabase.GetRecordRequest{
			Key: restaurant_name,
		}

		if getRecordResponse, errGetRecord := s.reviewDatabaseClient.GetRecord(ctx, getRecordMsg); errGetRecord == nil {
			rawValue := getRecordResponse.Record.Value
			err = proto.Unmarshal(rawValue, searchResponse)
			if err != nil {
				err = status.Errorf(codes.Internal, "Failed to deserialize data")
			} else {
				// If we found the item in the storage layer, then update the cache to hold this element
				if s.CACHE_FLAG {
					updateCache(ctx, s.reviewCacheClient, restaurant_name, rawValue)
				}
				err = status.Errorf(codes.OK, "Found value with Key: %s", restaurant_name)
			}
		} else {
			return reviewResponse, status.Errorf(codes.NotFound, "Item with Key: %s does not exist", restaurant_name)
		}
	}
	reviewResponse = searchResponse.ReviewsMap[username]
	return reviewResponse, err
}

func (s *Review) SearchReviews(ctx context.Context, req *review.SearchReviewsRequest) (*review.SearchReviewsResponse, error) {
	restaurant_name := req.GetRestaurantName()
	searchResponse := &review.SearchReviewsResponse{}

	var err error

	getCacheItemMsg := &mycache.GetItemRequest{
		Key: restaurant_name,
	}

	var getItemResponse *mycache.GetItemResponse
	var errGetItem error

	// Check the cache for the restaurant
	getItemResponse, errGetItem = s.reviewCacheClient.GetItem(ctx, getCacheItemMsg)
	if errGetItem == nil {
		rawValue := getItemResponse.Item.Value
		err = proto.Unmarshal(rawValue, searchResponse)
		if err != nil {
			err = status.Errorf(codes.Internal, "Failed to deserialize data")
		} else {
			err = status.Errorf(codes.OK, "Found value with Key: %s", restaurant_name)
		}
	} else {
		// Make a call to the storage layer due to a cache-miss
		getRecordMsg := &mydatabase.GetRecordRequest{
			Key: restaurant_name,
		}

		if getRecordResponse, errGetRecord := s.reviewDatabaseClient.GetRecord(ctx, getRecordMsg); errGetRecord == nil {
			rawValue := getRecordResponse.Record.Value
			err = proto.Unmarshal(rawValue, searchResponse)
			if err != nil {
				err = status.Errorf(codes.Internal, "Failed to deserialize data")
			} else {
				// If we found the item in the storage layer, then update the cache to hold this element
				if s.CACHE_FLAG {
					updateCache(ctx, s.reviewCacheClient, restaurant_name, rawValue)
				}
				err = status.Errorf(codes.OK, "Found value with Key: %s", restaurant_name)
			}
		} else {
			return searchResponse, status.Errorf(codes.NotFound, "Item with Key: %s does not exist", restaurant_name)
		}
	}
	return searchResponse, err
}

func (s *Review) PostReview(ctx context.Context, req *review.PostReviewRequest) (*review.PostReviewResponse, error) {
	username := req.GetUserName()
	restaurant_name := req.GetRestaurantName()
	userReview := req.GetReview()
	rating := req.GetRating()

	reviewResponse := &review.PostReviewResponse{Status: false}

	var err error

	// Create a new GetReviewResponse object with the details to save.
	msg := &review.GetReviewResponse{
		UserName:       username,
		RestaurantName: restaurant_name,
		Review:         userReview,
		Rating:         rating,
	}

	getCacheItemMsg := &mycache.GetItemRequest{
		Key: restaurant_name,
	}

	searchResponse := &review.SearchReviewsResponse{}
	username_reviews := make(map[string]*review.GetReviewResponse)

	var getItemResponse *mycache.GetItemResponse
	var errGetItem error

	// Check the cache to see if restaurant already exists
	getItemResponse, errGetItem = s.reviewCacheClient.GetItem(ctx, getCacheItemMsg)
	if errGetItem == nil {
		rawValue := getItemResponse.Item.Value
		err = proto.Unmarshal(rawValue, searchResponse)
		if err != nil {
			err = status.Errorf(codes.Internal, "Failed to deserialize data")
			return reviewResponse, err
		} else {
			err = status.Errorf(codes.OK, "Found value with Key: %s", restaurant_name)
			// Store the map so we can add the new review
			username_reviews = searchResponse.ReviewsMap
		}
	} else {
		getRecordMsg := &mydatabase.GetRecordRequest{
			Key: restaurant_name,
		}
		if getRecordResponse, errGetItem := s.reviewDatabaseClient.GetRecord(ctx, getRecordMsg); errGetItem == nil {
			rawValue := getRecordResponse.Record.Value
			err = proto.Unmarshal(rawValue, searchResponse)
			if err != nil {
				err = status.Errorf(codes.Internal, "Failed to deserialize data")
				return reviewResponse, err
			} else {
				err = status.Errorf(codes.OK, "Found value with Key: %s", restaurant_name)
				// update the map to add the new review
				username_reviews = searchResponse.ReviewsMap
			}
		}
	}

	username_reviews[username] = msg
	searchResponse.ReviewsMap = username_reviews

	data, err := proto.Marshal(searchResponse)
	if err != nil {
		// If serialization fails, return an internal error.
		return reviewResponse, status.Errorf(codes.Internal, "Failed to serialize data")
	}

	reviewResponse.Status, err = updateDB(ctx, s.reviewCacheClient, s.reviewDatabaseClient, restaurant_name, data, s.CACHE_FLAG)
	return reviewResponse, err
}
