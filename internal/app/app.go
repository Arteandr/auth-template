package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"mzhn/auth/internal/config"
	"mzhn/auth/internal/handlers"
	"mzhn/auth/internal/services/authservice"

	mw "mzhn/auth/internal/middleware"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

type App struct {
	app *echo.Echo
	cfg *config.Config

	as *authservice.AuthService
}

func newApp(cfg *config.Config, as *authservice.AuthService) *App {
	return &App{
		app: echo.New(),
		cfg: cfg,
		as:  as,
	}
}

func (a *App) initApp() {
	a.app.Use(emw.Logger())
	// a.app.Use(emw.Recover())
	a.app.Use(emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowCredentials: true,
	}))

	tokguard := mw.Token
	authguard := mw.RequireAuth(a.as, a.cfg)

	a.app.POST("/register", handlers.Register(a.as))
	a.app.POST("/login", handlers.Login(a.as))
	a.app.POST("/refresh", handlers.Refresh(a.as), tokguard())
	a.app.GET("/profile", handlers.Profile(a.as), tokguard(), authguard())
	a.app.POST("/logout", handlers.Logout(a.as), tokguard(), authguard())
}

func (a *App) Run() {

	a.initApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		port := a.cfg.App.Port
		addr := fmt.Sprintf(":%d", port)
		slog.Info("running server", slog.String("addr", addr))
		a.app.Start(addr)
	}()

	sig := <-sigChan
	slog.Info(fmt.Sprintf("Signal %v received, stopping server...\n", sig))
	a.app.Shutdown(context.Background())
}
