package services

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/reservation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mycache"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mydatabase"
)

// Reservation implements the reservation service
type Reservation struct {
	name string
	port int
	reservation.ReservationServiceServer

	mu sync.Mutex
	// reservationStore map[string][]byte
	reservationPopularity map[string]int

	reservationCacheClient    mycache.CacheServiceClient       // Add reservation grpc cache client for communicating with reservation cache server
	reservationDatabaseClient mydatabase.DatabaseServiceClient // Add reservation grpc storage client for communicating with reservation storage server

	CACHE_FLAG bool
}

// NewReservation returns a new server
func NewReservation(name string, reservationPort int, reservationCacheAddr string, reservationDatabaseAddr string) *Reservation {
	return &Reservation{
		name: name,
		port: reservationPort,
		// reservationStore: make(map[string][]byte),
		reservationPopularity:     make(map[string]int),
		reservationCacheClient:    mycache.NewCacheServiceClient(dial(reservationCacheAddr)),          // Initialize and establish cxn using specified address
		reservationDatabaseClient: mydatabase.NewDatabaseServiceClient(dial(reservationDatabaseAddr)), // Initialize and establish cxn using specified address
		CACHE_FLAG:                true,
	}
}

// Run starts the Reservation gRPC server and listens for incoming requests.
// It returns an error if the server fails to start or encounters an error.
func (s *Reservation) Run() error {
	// Create a new gRPC server instance.
	srv := grpc.NewServer()

	// Register the Reservation server implementation with the gRPC server.
	reservation.RegisterReservationServiceServer(srv, s)

	// Create a TCP listener that listens for incoming requests on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// (Optional) Log a message indicating that the server is running and listening on the specified port.
	log.Printf("reservation server running at port: %d", s.port)

	// Start serving incoming requests using the registered implementation.
	return srv.Serve(lis)
}

func (s *Reservation) GetReservation(ctx context.Context, req *reservation.GetReservationRequest) (*reservation.GetReservationResponse, error) {
	username := req.GetUserName()
	reservationResponse := &reservation.GetReservationResponse{}

	var err error

	getCacheItemMsg := &mycache.GetItemRequest{
		Key: username,
	}

	var getItemResponse *mycache.GetItemResponse
	var errGetItem error

	// Check the cache for the reservation
	getItemResponse, errGetItem = s.reservationCacheClient.GetItem(ctx, getCacheItemMsg)
	if errGetItem == nil {
		rawValue := getItemResponse.Item.Value
		err = proto.Unmarshal(rawValue, reservationResponse)
		if err != nil {
			err = status.Errorf(codes.Internal, "Failed to deserialize data")
		} else {
			err = status.Errorf(codes.OK, "Found value with Key: %s", username)
		}
	} else {
		// Make a call to the storage layer due to a cache-miss
		getRecordMsg := &mydatabase.GetRecordRequest{
			Key: username,
		}

		if getRecordResponse, errGetRecord := s.reservationDatabaseClient.GetRecord(ctx, getRecordMsg); errGetRecord == nil {
			rawValue := getRecordResponse.Record.Value
			err = proto.Unmarshal(rawValue, reservationResponse)
			if err != nil {
				err = status.Errorf(codes.Internal, "Failed to deserialize data")
			} else {
				// If we found the item in the storage layer, then update the cache to hold this element
				if s.CACHE_FLAG {
					updateCache(ctx, s.reservationCacheClient, username, rawValue)
				}

				err = status.Errorf(codes.OK, "Found value with Key: %s", username)
			}
		} else {
			return reservationResponse, status.Errorf(codes.NotFound, "Item with Key: %s does not exist", username)
		}
	}
	return reservationResponse, err
}

func (s *Reservation) MakeReservation(ctx context.Context, req *reservation.MakeReservationRequest) (*reservation.MakeReservationResponse, error) {
	username := req.GetUserName()
	restaurant_name := req.GetRestaurantName()
	time := req.GetTime()

	// Create a new GetReservationResponse object with the reservation to save.
	msg := &reservation.GetReservationResponse{
		UserName:       username,
		RestaurantName: restaurant_name,
		Time:           time,
	}

	// Initialize an empty response object.
	reservationResponse := &reservation.MakeReservationResponse{Status: false}

	// Serialize the Go object into protobuf format.
	data, err := proto.Marshal(msg)
	if err != nil {
		// If serialization fails, return an internal error.
		return reservationResponse, status.Errorf(codes.Internal, "Failed to serialize data")
	}

	// Use locking when updating the local data structure to prevent concurrency issues
	s.mu.Lock()
	count := 0
	if val, exists := s.reservationPopularity[restaurant_name]; exists {
		count = val
	}
	s.reservationPopularity[restaurant_name] = count + 1
	s.mu.Unlock()

	reservationResponse.Status, err = updateDB(ctx, s.reservationCacheClient, s.reservationDatabaseClient, username, data, s.CACHE_FLAG)

	// Return the response object and any error.
	return reservationResponse, err
}

func (s *Reservation) MostPopular(ctx context.Context, req *reservation.MostPopularRequest) (*reservation.MostPopularResponse, error) {
	topK := req.GetTopK()

	pq := make(PriorityQueue, 0)
	counter := 0
	s.mu.Lock()
	for k, v := range s.reservationPopularity {
		heap.Push(&pq, &Item{value: k, priority: v, index: counter})

		counter = counter + 1
	}
	s.mu.Unlock()

	var numElements int32 = 0
	topKRestaurants := make([]string, 0)

	for numElements < topK {
		if pq.Len() > 0 {
			topKRestaurants = append(topKRestaurants, heap.Pop(&pq).(*Item).value)
		} else {
			return &reservation.MostPopularResponse{}, status.Errorf(codes.Internal, "Less than K restaurants in WELP")
		}
		numElements = numElements + 1
	}

	mostPopularResponse := &reservation.MostPopularResponse{TopKRestaurants: topKRestaurants}
	return mostPopularResponse, nil
}

// Priority Queue implementation used from Go Documentation

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
