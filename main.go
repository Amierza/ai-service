package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/Amierza/ai-service/config/database"
	"github.com/Amierza/ai-service/jwt"
	"github.com/Amierza/ai-service/logger"
	"github.com/Amierza/ai-service/middleware"
	pb "github.com/Amierza/ai-service/proto"
	"github.com/Amierza/ai-service/repository"
	"github.com/Amierza/ai-service/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// setup potgres connection
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Zap logger
	zapLogger, err := logger.New(true) // true = dev, false = prod
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer zapLogger.Sync() // flush buffer

	// baca API Key dari environment
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("missing OPENAI_API_KEY environment variable")
	}

	var (
		// JWT
		jwt = jwt.NewJWT()

		// Summary Task With LLM GPT
		summaryRepo    = repository.NewSummaryRepository(db)
		summaryService = service.NewSummaryService(summaryRepo, zapLogger, jwt, openaiKey)
		_              = summaryService
		// summaryHandler = handler.NewSummaryHandler(summaryService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes.Summary(server, summaryHandler, jwt)

	server.Static("/uploads", "./uploads")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	time.Local, _ = time.LoadLocation("Asia/Jakarta")

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	// setup gRPC client ke AI Service
	grpcPortTarget := os.Getenv("GRPC_PORT")
	if grpcPortTarget == "" {
		grpcPortTarget = "50051" // default fallback
	}

	go func() {
		lis, err := net.Listen("tcp", grpcPortTarget)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterSummaryServiceServer(grpcServer, service.NewGRPCSummaryServer(summaryService, zapLogger))

		zapLogger.Info("🚀 gRPC server running on port 50051...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
