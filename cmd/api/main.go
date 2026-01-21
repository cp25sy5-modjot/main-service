package main

import (
	"log"

	"github.com/cp25sy5-modjot/main-service/internal/database"
	server "github.com/cp25sy5-modjot/main-service/internal/httpapi"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	
)

func main() {
	conf := config.LoadConfig()
	db := database.NewPostgresDatabase(conf)

	// grpcConn, err := grpc.Dial(conf.AIService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to gRPC server: %v", err)
	// }
	// defer grpcConn.Close()
	grpcConn, err := grpc.NewClient(
		conf.AIService.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	defer func() {
		if err := grpcConn.Close(); err != nil {
			log.Printf("failed to close grpc connection: %v", err)
		}
	}()

	aiClient := pb.NewAiWrapperServiceClient(grpcConn)
	server.NewFiberServer(conf, db, aiClient).Start()
}
