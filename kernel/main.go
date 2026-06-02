package kernel

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gauas/account-service/config"
	"github.com/gauas/account-service/controller"
	"github.com/gauas/account-service/middlewares"
	"github.com/gauas/account-service/packages/httpresp"
	"github.com/gauas/account-service/route"
	"github.com/labstack/echo/v4"
)

type Kernel struct {
	controller *controller.Controller
	middleware *middlewares.Middleware
	config     config.Config
}

func New(ctrl *controller.Controller, mw *middlewares.Middleware, cfg *config.Config) *Kernel {
	return &Kernel{controller: ctrl, middleware: mw, config: *cfg}
}

func (k *Kernel) Start() {
	server := echo.New()
	server.HideBanner = true
	server.HTTPErrorHandler = func(err error, c echo.Context) {
		var e *httpresp.Error
		if errors.As(err, &e) {
			_ = c.JSON(e.Code, httpresp.Response{Status: e.Code, Error: e.Message})
			return
		}

		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			code := httpErr.Code
			_ = c.JSON(code, httpresp.Response{Status: code, Error: fmt.Sprintf("%v", httpErr.Message)})
			return
		}

		_ = c.JSON(http.StatusInternalServerError, httpresp.Response{Status: http.StatusInternalServerError, Error: "internal server error"})
	}

	k.middleware.RegisterGlobal(server)

	route.New(server, k.controller, k.middleware).RegisterRoutes() // <----

	addr := fmt.Sprintf(":%s", k.config.Port)
	log.Printf("account-service listening on %s", addr)

	if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
