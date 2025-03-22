package main

import (
	"fmt"
	"github.com/vestamart/loms/internal/app/loms"
	"github.com/vestamart/loms/internal/config"
	"github.com/vestamart/loms/internal/delivery"
	"github.com/vestamart/loms/internal/mw"
	"github.com/vestamart/loms/internal/repository"
	desc "github.com/vestamart/loms/pkg/api/loms/v1"
	"google.golang.org/grpc"
	"log"
	"net"
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
		mw.Panic,
	))

	//reflection.Register(grpcServer)

	ordersRepo := repository.NewInMemoryOrderRepository(100)
	stocksRepo, err := repository.NewInMemoryStocksRepositoryFromFile()
	if err != nil {
		panic(err)
	}
	service := loms.NewService(ordersRepo, stocksRepo)

	controller := delivery.NewServer(*service)

	desc.RegisterLomsServer(grpcServer, controller)
	log.Print("Server running on port: " + cfg.LOMSServer.Port)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
