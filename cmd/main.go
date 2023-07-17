package main

import (
	"log"
	"net"
	"net/http"
	"service2/pb"
	"service2/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	AdminService := &service.AdminDashboard{}

	pb.RegisterAdminDashboardServer(grpcServer, AdminService)

	listener, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Admin dashboard service is running...")
	go grpcServer.Serve(listener)

	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start health check server: %v", err)
	}
}
