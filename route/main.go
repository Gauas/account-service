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
	health := r.Server.Group("/v1/account")
	{
		health.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
		})
	}

	api := r.Server.Group("/v1/account", r.Middleware.Device())

	public := api.Group("")
	{
		public.POST("/register", r.Controller.Authentication.Register)
		public.POST("/login", r.Controller.Authentication.Login)
		public.POST("/oauth2", r.Controller.Authentication.OAuth2)
	}

	private := api.Group("", r.Middleware.Auth())
	{
		private.POST("/logout", r.Controller.Authentication.Logout)
	}

	info := private.Group("/info")
	{
		info.GET("", r.Controller.Info.GetInfo)
		info.PUT("", r.Controller.Info.UpdateInfo)
		info.PATCH("/avatar", r.Controller.Info.UpdateAvatar)
	}

	mfa := private.Group("/mfa")
	{
		mfa.GET("/totp/qr", r.Controller.MFA.GenerateTOTP)
		mfa.POST("/totp/enable", r.Controller.MFA.EnableTOTP)
		mfa.POST("/totp/verify", r.Controller.MFA.VerifyTOTP)
	}
}
