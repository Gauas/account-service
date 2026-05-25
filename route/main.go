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
	api := r.Server.Group("/v1/account", r.Middleware.Device())

	public := api.Group("")
	{
		public.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
		})

		public.POST("/register", r.Controller.Authentication.Register)
		public.POST("/login", r.Controller.Authentication.Login)
		public.POST("/oathh2", r.Controller.Authentication.OAuth2)
	}

	private := api.Group("", r.Middleware.Auth())
	{
		private.GET("", r.Controller.GetProfile)

		private.POST("/logout", r.Controller.Authentication.Logout)
	}

	//profile := api.Group("/profile")
	//profile.Use(auth)
	//profile.PUT("", r.controller.UpdateProfile)
	//profile.PATCH("/avatar", r.controller.UpdateAvatar)
	//
	//mfa := api.Group("/mfa")
	//mfa.Use(auth)
	//mfa.GET("/totp/qr", r.controller.GenerateTOTPQR)
	//mfa.POST("/totp/enable", r.controller.EnableTOTP)
	//mfa.POST("/totp/verify", r.controller.VerifyTOTP)
	//api.GET("/verify-email/:token", r.controller.VerifyEmail)
	//api.POST("/send-verification/:user_id", r.controller.SendVerificationEmail)

	//sso := api.Group("/oauth2")
	//sso.POST("/google", r.controller.)
}
