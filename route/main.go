package route

import (
	"net/http"

	"github.com/gauas/account-service/controller"
	middleware "github.com/gauas/account-service/middlewares"
	"github.com/labstack/echo/v4"
)

type Router struct {
	Server     *echo.Echo
	Controller *controller.Controller
	Middleware *middleware.Middleware
}

func New(server *echo.Echo, ctrl *controller.Controller, mw *middleware.Middleware) *Router {
	return &Router{Server: server, Controller: ctrl, Middleware: mw}
}

func (r *Router) RegisterRoutes() {
	api := r.Server.Group("/v1/account")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})

	api.POST("/register", r.Controller.Authentication.Register)
	api.POST("/login", r.Controller.Authentication.Login)

	//api.GET("/verify-email/:token", r.controller.VerifyEmail)
	//api.POST("/send-verification/:user_id", r.controller.SendVerificationEmail)

	//sso := api.Group("/oauth2")
	//sso.POST("/google", r.controller.)

	authed := api.Group("", r.Middleware.Auth())
	{
		authed.POST("/logout", r.Controller.Authentication.Logout)
	}

	//profile := api.Group("/profile")
	//profile.Use(auth)
	//profile.GET("", r.controller.GetProfile)
	//profile.PUT("", r.controller.UpdateProfile)
	//profile.PATCH("/avatar", r.controller.UpdateAvatar)
	//
	//mfa := api.Group("/mfa")
	//mfa.Use(auth)
	//mfa.GET("/totp/qr", r.controller.GenerateTOTPQR)
	//mfa.POST("/totp/enable", r.controller.EnableTOTP)
	//mfa.POST("/totp/verify", r.controller.VerifyTOTP)
}
