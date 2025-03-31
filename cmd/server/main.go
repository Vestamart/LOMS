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
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		mw.Panic,
		mw.Logger,
	))

	dsn := "postgres://root:root@postgres:5432/loms_db?sslmode=disable"
	dbConn, err := mw.ConnectWithRetry(context.Background(), dsn, 10, 5*time.Second)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

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
	log.Print("Server running on port: " + cfg.LOMSServer.Port)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
