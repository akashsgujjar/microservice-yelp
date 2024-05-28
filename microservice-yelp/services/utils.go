package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mycache"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mydatabase"
)

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	start := time.Now()
	defer func() {
		in, _ := json.Marshal(req)
		out, _ := json.Marshal(reply)
		inStr, outStr := string(in), string(out)
		duration := int64(time.Since(start).Microseconds())

		delimiter := ";"
		errStr := fmt.Sprintf("%v", err)
		if err == nil {
			errStr = "<nil>"
		}
		logMessage := fmt.Sprintf("grpc%s%s%s%s%s%s%s%s%s%d", delimiter, method, delimiter, inStr, delimiter, outStr, delimiter, errStr, delimiter, duration)
		log.Println(logMessage)

	}()

	return invoker(ctx, method, req, reply, cc, opts...)
}

func logMsg(handler, inStr, outStr, errStr string, duration int64) {
	delimiter := ";"

	// Log the entire time it takes to execute
	logMessage := fmt.Sprintf("%s%s%s%s%s%s%s%s%d", handler, delimiter, inStr, delimiter, outStr, delimiter, errStr, delimiter, duration)
	log.Println(logMessage)
}

// dial creates a new gRPC client connection to the specified address and returns a client connection object.
func dial(addr string) *grpc.ClientConn {
	// Define gRPC dial options for the client connection.
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor),
	}

	// Create a new gRPC client connection to the specified address using the dial options.
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		// If there was an error creating the client connection, panic with an error message.
		panic(fmt.Sprintf("ERROR: dial error: %v", err))
	}

	// Return the created client connection object.
	return conn
}

// Retrieves UUID associated with a particular restaurant and user pair. Useful for mapping reviews and/or reservations to data
func GetQueryUUID(restaurantName, userName string) (string, error) {
	combinedString := restaurantName + userName
	uuidNamespace := uuid.NewSHA1(uuid.Nil, []byte(combinedString))
	uuidBytes := uuidNamespace[:]
	uuidObj, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return "", err
	}
	return uuidObj.String(), nil
}

/******************************************************
 * TODO: place your utilities / helper functions here *
 ******************************************************/

func updateCache(ctx context.Context, cacheClient mycache.CacheServiceClient, key string, val []byte) {
	cacheItem := &mycache.CacheItem{
		Key:   key,
		Value: val,
	}
	setItemRequest := &mycache.SetItemRequest{
		Item: cacheItem,
	}

	cacheClient.SetItem(ctx, setItemRequest)
}

func updateDB(ctx context.Context, cacheClient mycache.CacheServiceClient, dbClient mydatabase.DatabaseServiceClient, key string, val []byte, cacheFlag bool) (bool, error) {
	databaseRecord := &mydatabase.DatabaseRecord{
		Key:   key,
		Value: val,
	}
	setRecordMsg := &mydatabase.SetRecordRequest{
		Record: databaseRecord,
	}

	var err error
	statusVal := false
	if setRecordResponse, errSet := dbClient.SetRecord(ctx, setRecordMsg); errSet == nil {
		if !setRecordResponse.Success {
			err = status.Errorf(codes.Internal, "Failed to update data storage")
		} else {
			statusVal = true
			err = status.Errorf(codes.OK, "Updated data storage with key: ", key)
			if cacheFlag {
				updateCache(ctx, cacheClient, key, val)
			}
		}
	} else {
		err = status.Errorf(codes.Internal, "Error in updating data storage")
	}

	return statusVal, err
}
