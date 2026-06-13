package main

import (
	"log"

	"github.com/gauas/account-service/bootstrap"
	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/controller"
	"github.com/gauas/account-service/grpc"
	"github.com/gauas/account-service/http"
	"github.com/gauas/account-service/infra"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/publisher"
	"github.com/gauas/account-service/repository"
	"github.com/gauas/account-service/service"
)

func main() {
	cfg := config.New()

	infraInstance := infra.New(cfg)
	repo := repository.New(infraInstance.DB)
	pub := publisher.New(infraInstance.Queue)
	svc := service.New(repo, pub, cfg, infraInstance)
	ctrl := controller.New(svc, cfg)
	mw := middlewares.New(cfg, infraInstance)

	httpServer := http.Register(ctrl, mw, cfg)
	grpcServer := grpc.Register(cfg.GRPCPort, svc)

	if err := bootstrap.Start(httpServer, grpcServer); err != nil {
		log.Fatal(err)
	}
}
