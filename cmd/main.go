package main

import (
	"log"

	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/controller"
	"github.com/gauas/account-service/infra"
	"github.com/gauas/account-service/kernel"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/repository"
	"github.com/gauas/account-service/service"
)

func main() {
	Config := config.New()

	infraInstance := infra.New(Config)

	repositoryInstance := repository.New(infraInstance.DB)

	serviceInstance := service.New(repositoryInstance, Config, infraInstance)

	controllerInstance := controller.New(serviceInstance, Config)

	middlewareInstance := middlewares.New(Config, infraInstance)

	kernel.New(controllerInstance, middlewareInstance, Config).Start()

	log.Println("account-service started")
}
