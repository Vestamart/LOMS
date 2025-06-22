package main

import (
	"context"
	"fmt"
	"github.com/vestamart/loms/internal/app/loms"
	"github.com/vestamart/loms/internal/config"
	"github.com/vestamart/loms/internal/delivery"
	"github.com/vestamart/loms/internal/mw"
	"github.com/vestamart/loms/internal/repository/postgres"
	desc "github.com/vestamart/loms/pkg/api/loms/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("App started")
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.LOMSServer.Port))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		mw.Panic,
		mw.Logger,
	))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	dbConn, err := mw.ConnectWithRetry(context.Background(), dsn, 10, 5*time.Second)
	if err != nil {
		log.Fatal("Failed to connect to database: " + err.Error())
	}
	defer dbConn.Close(context.Background())

	orderRepoPostgres := postgres.NewOrderRepositoryPostgres(dbConn)
	//ordersRepo := repository.NewInMemoryOrderRepository(100)
	//stocksRepo, err := repository.NewInMemoryStocksRepositoryFromFile()
	stocksRepoPostgres := postgres.NewStocksRepositoryPostgres(dbConn)
	//if err != nil {
	//	panic(err)
	//}
	service := loms.NewService(orderRepoPostgres, stocksRepoPostgres)

	controller := delivery.NewServer(*service)

	desc.RegisterLomsServer(grpcServer, controller)

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on port: %s", cfg.LOMSServer.Port)
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received")

	// Graceful shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	log.Println("Server gracefully stopped")
}
