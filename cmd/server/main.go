package main

import (
	"log"
	"money-go/internal/database"
	"money-go/internal/server"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := database.NewConnection()

	log.Println("Database connection established")

	log.Printf("Server listening on %s", "50051")

	srv := server.NewGRPCServer(db)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	defer db.Close()
}
