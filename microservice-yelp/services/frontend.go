package services

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/detail"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/reservation"
	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/review"
)

// Frontend implements a service that acts as an interface to interact with different microservices.
type Frontend struct {
	port          int
	detailClient1 detail.DetailServiceClient
	detailClient2 detail.DetailServiceClient
	detailClient3 detail.DetailServiceClient

	reviewClient1 review.ReviewServiceClient
	reviewClient2 review.ReviewServiceClient
	reviewClient3 review.ReviewServiceClient

	mu sync.Mutex

	LOAD_BALANCING_ALG           string
	sequentialKeyToReplicaDetail map[string]int
	sequentialKeyToReplicaReview map[string]int
	sequentialNextReplicaDetail  int
	sequentialNextReplicaReview  int

	reservationClient reservation.ReservationServiceClient
	User              string
}

// NewFrontend creates a new Frontend instance with the specified configuration.
func NewFrontend(port int, detailaddr1 string, detailaddr2 string, detailaddr3 string, reviewaddr1 string, reviewaddr2 string, reviewaddr3 string, reservationaddr string) *Frontend {
	f := &Frontend{
		port:          port,
		detailClient1: detail.NewDetailServiceClient(dial(detailaddr1)),
		detailClient2: detail.NewDetailServiceClient(dial(detailaddr2)),
		detailClient3: detail.NewDetailServiceClient(dial(detailaddr3)),

		reviewClient1: review.NewReviewServiceClient(dial(reviewaddr1)),
		reviewClient2: review.NewReviewServiceClient(dial(reviewaddr2)),
		reviewClient3: review.NewReviewServiceClient(dial(reviewaddr3)),

		LOAD_BALANCING_ALG: "hash",

		sequentialKeyToReplicaDetail: make(map[string]int),
		sequentialKeyToReplicaReview: make(map[string]int),

		sequentialNextReplicaDetail: 1,
		sequentialNextReplicaReview: 1,

		reservationClient: reservation.NewReservationServiceClient(dial(reservationaddr)),
		User:              "None",
	}
	return f
}

// Run starts the Frontend server and listens for incoming requests on the specified port.
func (s *Frontend) Run() error {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/get-detail", s.getDetailHandler)
	http.HandleFunc("/post-detail", s.postDetailHandler)
	http.HandleFunc("/get-review", s.getReviewHandler)
	http.HandleFunc("/post-review", s.postReviewHandler)
	http.HandleFunc("/search-reviews", s.searchReviewsHandler)
	http.HandleFunc("/get-reservation", s.getReservationHandler)
	http.HandleFunc("/make-reservation", s.makeReservationHandler)
	http.HandleFunc("/most-popular", s.mostPopularHandler)

	log.Printf("frontend server running at port: %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

// getDetailHandler handles requests for retrieving restaurant details.
func (s *Frontend) getDetailHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ctx := r.Context()

	restaurant_name := r.URL.Query().Get("restaurant_name")
	if restaurant_name == "" {
		http.Error(w, "Malformed request to `/get-detail` endpoint!", http.StatusBadRequest)
		return
	}

	replicaNum := 1
	if s.LOAD_BALANCING_ALG == "hash" {
		hashCode := int(hash(restaurant_name))
		replicaNum = hashCode%3 + 1
	} else if s.LOAD_BALANCING_ALG == "none" {
		replicaNum = 1
	} else {
		replicaNum = determineReplicaNext(s, "detail", restaurant_name)
	}

	req := &detail.GetDetailRequest{RestaurantName: restaurant_name}

	var reply *detail.GetDetailResponse
	var err error

	if replicaNum == 1 {
		reply, err = s.detailClient1.GetDetail(ctx, req)
	} else if replicaNum == 2 {
		reply, err = s.detailClient2.GetDetail(ctx, req)
	} else if replicaNum == 3 {
		reply, err = s.detailClient3.GetDetail(ctx, req)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"

	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}

	logMsg("frontend.getDetailHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// postDetailHandler handles requests for posting restaurant details.
func (s *Frontend) postDetailHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ctx := r.Context()
	restaurant_name := r.URL.Query().Get("restaurant_name")
	location := r.URL.Query().Get("location")
	style := r.URL.Query().Get("style")
	capacity, errCap := strconv.Atoi(r.URL.Query().Get("capacity"))

	if restaurant_name == "" || location == "" || style == "" || errCap != nil {
		http.Error(w, "Malformed request to `/post-detail` endpoint!", http.StatusBadRequest)
		return
	}

	replicaNum := 1
	if s.LOAD_BALANCING_ALG == "hash" {
		hashCode := int(hash(restaurant_name))
		replicaNum = hashCode%3 + 1
	} else if s.LOAD_BALANCING_ALG == "none" {
		replicaNum = 1
	} else {
		replicaNum = determineReplicaNext(s, "detail", restaurant_name)
	}

	req := &detail.PostDetailRequest{
		RestaurantName: restaurant_name,
		Location:       location,
		Style:          style,
		Capacity:       int32(capacity),
	}

	var reply *detail.PostDetailResponse
	var err error

	if replicaNum == 1 {
		reply, err = s.detailClient1.PostDetail(ctx, req)
	} else if replicaNum == 2 {
		reply, err = s.detailClient2.PostDetail(ctx, req)
	} else if replicaNum == 3 {
		reply, err = s.detailClient3.PostDetail(ctx, req)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}

	logMsg("frontend.postDetailHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// getReviewHandler handles requests for retrieving reviews.
func (s *Frontend) getReviewHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ctx := r.Context()
	restaurant_name := r.URL.Query().Get("restaurant_name")
	user_name := r.URL.Query().Get("user_name")

	if restaurant_name == "" || user_name == "" {
		http.Error(w, "Malformed request to `/get-review` endpoint!", http.StatusBadRequest)
		return
	}

	replicaNum := 1
	if s.LOAD_BALANCING_ALG == "hash" {
		hashCode := int(hash(restaurant_name))
		replicaNum = hashCode%3 + 1
	} else if s.LOAD_BALANCING_ALG == "none" {
		replicaNum = 1
	} else {
		replicaNum = determineReplicaNext(s, "review", restaurant_name)
	}

	req := &review.GetReviewRequest{
		RestaurantName: restaurant_name,
		UserName:       user_name,
	}

	var reply *review.GetReviewResponse
	var err error

	if replicaNum == 1 {
		reply, err = s.reviewClient1.GetReview(ctx, req)
	} else if replicaNum == 2 {
		reply, err = s.reviewClient2.GetReview(ctx, req)
	} else if replicaNum == 3 {
		reply, err = s.reviewClient3.GetReview(ctx, req)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}
	logMsg("frontend.getReviewHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// postReviewHandler handles requests for posting reviews.
func (s *Frontend) postReviewHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ctx := r.Context()
	user_name := r.URL.Query().Get("user_name")
	restaurant_name := r.URL.Query().Get("restaurant_name")
	restaurant_review := r.URL.Query().Get("review")
	restaurant_rating, errRating := strconv.Atoi(r.URL.Query().Get("rating"))

	if restaurant_name == "" || user_name == "" || restaurant_review == "" || errRating != nil {
		http.Error(w, "Malformed request to `/post-review` endpoint!", http.StatusBadRequest)
		return
	}

	replicaNum := 1
	if s.LOAD_BALANCING_ALG == "hash" {
		hashCode := int(hash(restaurant_name))
		replicaNum = hashCode%3 + 1
	} else if s.LOAD_BALANCING_ALG == "none" {
		replicaNum = 1
	} else {
		replicaNum = determineReplicaNext(s, "review", restaurant_name)
	}

	req := &review.PostReviewRequest{
		UserName:       user_name,
		RestaurantName: restaurant_name,
		Review:         restaurant_review,
		Rating:         int32(restaurant_rating),
	}

	var reply *review.PostReviewResponse
	var err error

	if replicaNum == 1 {
		reply, err = s.reviewClient1.PostReview(ctx, req)
	} else if replicaNum == 2 {
		reply, err = s.reviewClient2.PostReview(ctx, req)
	} else if replicaNum == 3 {
		reply, err = s.reviewClient3.PostReview(ctx, req)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}
	logMsg("frontend.postReviewHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// searchReviewHandler handles requests for searching reviews.
func (s *Frontend) searchReviewsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	restaurant_name := r.URL.Query().Get("restaurant_name")
	if restaurant_name == "" {
		http.Error(w, "Malformed request to `/search-reviews` endpoint!", http.StatusBadRequest)
		return
	}

	replicaNum := 1
	if s.LOAD_BALANCING_ALG == "hash" {
		hashCode := int(hash(restaurant_name))
		replicaNum = hashCode%3 + 1
	} else if s.LOAD_BALANCING_ALG == "none" {
		replicaNum = 1
	} else {
		replicaNum = determineReplicaNext(s, "review", restaurant_name)
	}

	req := &review.SearchReviewsRequest{
		RestaurantName: restaurant_name,
	}

	var reply *review.SearchReviewsResponse
	var err error

	if replicaNum == 1 {
		reply, err = s.reviewClient1.SearchReviews(ctx, req)
	} else if replicaNum == 2 {
		reply, err = s.reviewClient2.SearchReviews(ctx, req)
	} else if replicaNum == 3 {
		reply, err = s.reviewClient3.SearchReviews(ctx, req)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}
	logMsg("frontend.searchReviewsHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// getReservationHandler handles requests for retrieving reservations.
func (s *Frontend) getReservationHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ctx := r.Context()
	user_name := r.URL.Query().Get("user_name")

	if user_name == "" {
		http.Error(w, "Malformed request to `/get-reservation` endpoint!", http.StatusBadRequest)
		return
	}

	req := &reservation.GetReservationRequest{UserName: user_name}
	reply, err := s.reservationClient.GetReservation(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}
	logMsg("frontend.getReservationHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// makeReservationHandler handles requests for making reservations.
func (s *Frontend) makeReservationHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	user_name := r.URL.Query().Get("user_name")
	restaurant_name := r.URL.Query().Get("restaurant_name")
	year, year_err := strconv.Atoi(r.URL.Query().Get("year"))
	month, month_err := strconv.Atoi(r.URL.Query().Get("month"))
	day, day_err := strconv.Atoi(r.URL.Query().Get("day"))

	if restaurant_name == "" || user_name == "" || year_err != nil || month_err != nil || day_err != nil {
		http.Error(w, "Malformed request to `/make-reservation` endpoint!", http.StatusBadRequest)
		return
	}

	req := &reservation.MakeReservationRequest{
		UserName:       user_name,
		RestaurantName: restaurant_name,
		Time:           &reservation.Date{Year: int32(year), Month: int32(month), Day: int32(day)},
	}
	reply, err := s.reservationClient.MakeReservation(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}
	logMsg("frontend.makeReservationHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

// mostPopularHandler handles requests for retrieving most popular restaurants.
func (s *Frontend) mostPopularHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := r.Context()
	topk, err := strconv.Atoi(r.URL.Query().Get("topk"))

	if err != nil {
		http.Error(w, "Malformed request to `/most-popular` endpoint!", http.StatusBadRequest)
		return
	}

	req := &reservation.MostPopularRequest{
		TopK: int32(topk),
	}
	reply, err := s.reservationClient.MostPopular(ctx, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate the duration in microseconds
	duration := int64(time.Since(start).Microseconds())
	in, _ := json.Marshal(req)
	// out, _ := json.Marshal(reply)
	inStr, outStr := string(in), "{}"
	errStr := fmt.Sprintf("%v", err)
	if err == nil {
		errStr = "<nil>"
	}
	logMsg("frontend.mostPopularHandler", inStr, outStr, errStr, duration)

	err = json.NewEncoder(w).Encode(reply)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func determineReplicaNext(s *Frontend, service string, restaurant_name string) int {
	s.mu.Lock()

	if service == "detail" {
		currReplicaNum, exists := s.sequentialKeyToReplicaDetail[restaurant_name]

		if exists {
			// if yes: send to that replica
			s.mu.Unlock()
			return currReplicaNum
		} else {
			// if no: send to sequentialNextReplica and set sequentialNextReplica
			replicaNum := s.sequentialNextReplicaDetail
			s.sequentialNextReplicaDetail++
			if s.sequentialNextReplicaDetail > 3 {
				s.sequentialNextReplicaDetail = 1
			}

			s.sequentialKeyToReplicaDetail[restaurant_name] = replicaNum
			s.mu.Unlock()
			return replicaNum
		}
	} else if service == "review" {
		currReplicaNum, exists := s.sequentialKeyToReplicaReview[restaurant_name]

		if exists {
			// if yes: send to that replica
			s.mu.Unlock()
			return currReplicaNum
		} else {
			// if no: send to sequentialNextReplica and set sequentialNextReplica
			replicaNum := s.sequentialNextReplicaReview
			s.sequentialNextReplicaReview++
			if s.sequentialNextReplicaReview > 3 {
				s.sequentialNextReplicaReview = 1
			}

			s.sequentialKeyToReplicaReview[restaurant_name] = replicaNum
			s.mu.Unlock()
			return replicaNum
		}
	}

	s.mu.Unlock()
	return 0
}
