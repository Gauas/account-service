package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gauas/account-service/controller"
	middleware "github.com/gauas/account-service/middlewares"
)

type Router struct {
	server     *echo.Echo
	controller *controller.Controller
	middleware *middleware.Middleware
}

func New(server *echo.Echo, ctrl *controller.Controller, mw *middleware.Middleware) *Router {
	return &Router{server: server, controller: ctrl, middleware: mw}
}

func (r *Router) RegisterRoutes() {
	api := r.server.Group("/v1/account")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})

	basic := api.Group("/basic")
	basic.POST("/register", r.controller.Register)
	basic.POST("/login", r.controller.Login)

	api.GET("/verify-email/:token", r.controller.VerifyEmail)
	api.POST("/send-verification/:user_id", r.controller.SendVerificationEmail)

	sso := api.Group("/sso")
	sso.POST("/google", r.controller.LoginWithGoogle)

	auth := r.middleware.Auth()

	api.POST("/logout", r.controller.Logout, auth)

	profile := api.Group("/profile")
	profile.Use(auth)
	profile.GET("", r.controller.GetProfile)
	profile.PUT("", r.controller.UpdateProfile)
	profile.PATCH("/avatar", r.controller.UpdateAvatar)

	mfa := api.Group("/mfa")
	mfa.Use(auth)
	mfa.GET("/totp/qr", r.controller.GenerateTOTPQR)
	mfa.POST("/totp/enable", r.controller.EnableTOTP)
	mfa.POST("/totp/verify", r.controller.VerifyTOTP)
}
