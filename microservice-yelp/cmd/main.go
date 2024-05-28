package main

import (
	"flag"
	"log"
	"os"

	services "gitlab.cs.washington.edu/syslab/cse453-welp/services"
)

type server interface {
	Run() error
}

func main() {
	// Define the flags to specify port numbers and addresses
	var (
		frontendPort     = flag.Int("frontend", 8080, "frontend server port")
		reservationsPort = flag.Int("reservationport", 8083, "reservation service port")

		// ports for each replica
		detailPort1 = flag.Int("detailport1", 8081, "detail-1 service port")
		reviewPort1 = flag.Int("reviewport1", 8082, "review-1 service port")

		detailPort2 = flag.Int("detailport2", 8084, "detail-2 service port")
		detailPort3 = flag.Int("detailport3", 8085, "detail-3 service port")

		reviewPort2 = flag.Int("reviewport2", 8086, "review-2 service port")
		reviewPort3 = flag.Int("reviewport3", 8087, "review-3 service port")

		// addresses for each replica
		detailAddr1 = flag.String("detailaddr1", "detail-1:8081", "detail-1 service address")
		detailAddr2 = flag.String("detailaddr2", "detail-2:8084", "detail-2 service address")
		detailAddr3 = flag.String("detailaddr3", "detail-3:8085", "detail-3 service address")

		reviewAddr1 = flag.String("reviewaddr1", "review-1:8082", "review-1 service addr")
		reviewAddr2 = flag.String("reviewaddr2", "review-2:8086", "review-2 service addr")
		reviewAddr3 = flag.String("reviewaddr3", "review-3:8087", "review-3 service addr")

		reservationAddr = flag.String("reservationaddr", "reservation:8083", "reservation service addr")

		// caches for each replica
		cachePort1           = flag.Int("cacheport1", 11211, "port used by all caches-1")
		detailCacheAddr1     = flag.String("detail_mycache_addr1", "mycache-detail-1:11211", "detail-1 mycache address")
		reviewCacheAddr1     = flag.String("review_mycache_addr1", "mycache-review-1:11211", "review-1 mycache address")
		reservationCacheAddr = flag.String("reservation_mycache_addr", "mycache-reservation:11211", "reservation mycache address")

		cachePort2       = flag.Int("cacheport2", 11212, "port used by all caches-2")
		detailCacheAddr2 = flag.String("detail_mycache_addr2", "mycache-detail-2:11212", "detail-2 mycache address")
		reviewCacheAddr2 = flag.String("review_mycache_addr2", "mycache-review-2:11212", "review-2 mycache address")

		cachePort3       = flag.Int("cacheport3", 11213, "port used by all caches-3")
		detailCacheAddr3 = flag.String("detail_mycache_addr3", "mycache-detail-3:11213", "detail-3 mycache address")
		reviewCacheAddr3 = flag.String("review_mycache_addr3", "mycache-review-3:11213", "review-3 mycache address")

		detailCacheCapacity      = flag.Int("detail_mycache_capacity", 100, "maximum number of K-V entries allowed in the detail cache service")
		reviewCacheCapacity      = flag.Int("review_mycache_capacity", 100, "maximum number of K-V entries allowed in the review cache service")
		reservationCacheCapacity = flag.Int("reservation_mycache_capacity", 100, "maximum number of K-V entries allowed in the reservation cache service")

		// database for each replica
		databasePort1           = flag.Int("databaseport1", 27017, "port used by all databases-1")
		storageDeviceType       = flag.String("storage_device_type", "ssd", "specifies emulated storage device type, e.g. option `ssd` or `disk`")
		detailDatabaseAddr1     = flag.String("detail_mydatabase_addr1", "mydatabase-detail-1:27017", "details-1 mydatabase address")
		reviewDatabaseAddr1     = flag.String("review_mydatabase_addr1", "mydatabase-review-1:27017", "review-1 mydatabase address")
		reservationDatabaseAddr = flag.String("reservation_mydatabase_addr", "mydatabase-reservation:27017", "reservation mydatabase address")

		databasePort2       = flag.Int("databaseport2", 27018, "port used by all databases-2")
		detailDatabaseAddr2 = flag.String("detail_mydatabase_addr2", "mydatabase-detail-2:27018", "details-2 mydatabase address")
		reviewDatabaseAddr2 = flag.String("review_mydatabase_addr2", "mydatabase-review-2:27018", "review-2 mydatabase address")

		databasePort3       = flag.Int("databaseport3", 27019, "port used by all databases-3")
		detailDatabaseAddr3 = flag.String("detail_mydatabase_addr3", "mydatabase-detail-3:27019", "details-3 mydatabase address")
		reviewDatabaseAddr3 = flag.String("review_mydatabase_addr3", "mydatabase-review-3:27019", "review-3 mydatabase address")
	)

	// Parse the flags
	flag.Parse()

	var srv server
	var cmd = os.Args[1]

	// Switch statement to create the correct service based on the command
	switch cmd {
	case "frontend":
		// Create a new frontend service with the specified ports and addresses
		srv = services.NewFrontend(
			*frontendPort,
			*detailAddr1,
			*detailAddr2,
			*detailAddr3,

			*reviewAddr1,
			*reviewAddr2,
			*reviewAddr3,

			*reservationAddr,
		)
	case "detail-1":
		switch {
		case len(os.Args) < 3:
			// Create a new detail service with the specified port
			srv = services.NewDetail(
				"detail-1",
				*detailPort1,
				*detailCacheAddr1,
				*detailDatabaseAddr1,
			)
		case os.Args[2] == "cache-1":
			srv = services.NewMyCache(
				"detail-1-cache",
				*cachePort1,
				*detailCacheCapacity,
			)
		case os.Args[2] == "database-1":
			srv = services.NewMyDatabase(
				"detail-1-database",
				*databasePort1,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for detail service: %s", os.Args[2])
		}
	case "detail-2":
		switch {
		case len(os.Args) < 3:
			// Create a new detail service with the specified port
			srv = services.NewDetail(
				"detail-2",
				*detailPort2,
				*detailCacheAddr2,
				*detailDatabaseAddr2,
			)
		case os.Args[2] == "cache-2":
			srv = services.NewMyCache(
				"detail-2-cache",
				*cachePort2,
				*detailCacheCapacity,
			)
		case os.Args[2] == "database-2":
			srv = services.NewMyDatabase(
				"detail-2-database",
				*databasePort2,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for detail service: %s", os.Args[2])
		}
	case "detail-3":
		switch {
		case len(os.Args) < 3:
			// Create a new detail service with the specified port
			srv = services.NewDetail(
				"detail-3",
				*detailPort3,
				*detailCacheAddr3,
				*detailDatabaseAddr3,
			)
		case os.Args[2] == "cache-3":
			srv = services.NewMyCache(
				"detail-3-cache",
				*cachePort3,
				*detailCacheCapacity,
			)
		case os.Args[2] == "database-3":
			srv = services.NewMyDatabase(
				"detail-3-database",
				*databasePort3,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for detail service: %s", os.Args[2])
		}
	case "reservation":
		switch {
		case len(os.Args) < 3:
			// Create a new reservation service with the specified port
			srv = services.NewReservation(
				"reservation",
				*reservationsPort,
				*reservationCacheAddr,
				*reservationDatabaseAddr,
			)
		case os.Args[2] == "cache":
			srv = services.NewMyCache(
				"reservation-cache",
				*cachePort1,
				*reservationCacheCapacity,
			)
		case os.Args[2] == "database":
			srv = services.NewMyDatabase(
				"reservation-database",
				*databasePort1,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for reservation service: %s", os.Args[2])
		}
	case "review-1":
		switch {
		case len(os.Args) < 3:
			// Create a new review service with the specified port
			srv = services.NewReview(
				"review-1",
				*reviewPort1,
				*reviewCacheAddr1,
				*reviewDatabaseAddr1,
			)
		case os.Args[2] == "cache-1":
			srv = services.NewMyCache(
				"review-1-cache",
				*cachePort1,
				*reviewCacheCapacity,
			)
		case os.Args[2] == "database-1":
			srv = services.NewMyDatabase(
				"review-1-database",
				*databasePort1,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for review service: %s", os.Args[2])
		}
	case "review-2":
		switch {
		case len(os.Args) < 3:
			// Create a new review service with the specified port
			srv = services.NewReview(
				"review-2",
				*reviewPort2,
				*reviewCacheAddr2,
				*reviewDatabaseAddr2,
			)
		case os.Args[2] == "cache-2":
			srv = services.NewMyCache(
				"review-2-cache",
				*cachePort2,
				*reviewCacheCapacity,
			)
		case os.Args[2] == "database-2":
			srv = services.NewMyDatabase(
				"review-2-database",
				*databasePort2,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for review service: %s", os.Args[2])
		}
	case "review-3":
		switch {
		case len(os.Args) < 3:
			// Create a new review service with the specified port
			srv = services.NewReview(
				"review-3",
				*reviewPort3,
				*reviewCacheAddr3,
				*reviewDatabaseAddr3,
			)
		case os.Args[2] == "cache-3":
			srv = services.NewMyCache(
				"review-3-cache",
				*cachePort3,
				*reviewCacheCapacity,
			)
		case os.Args[2] == "database-3":
			srv = services.NewMyDatabase(
				"review-3-database",
				*databasePort3,
				*storageDeviceType,
			)
		default:
			log.Fatalf("unknown subcmd for review service: %s", os.Args[2])
		}
	default:
		// If an unknown command is provided, log an error and exit
		log.Fatalf("unknown cmd: %s", cmd)
	}

	// Start the server and log any errors that occur
	if err := srv.Run(); err != nil {
		log.Fatalf("run %s error: %v", cmd, err)
	}
}
