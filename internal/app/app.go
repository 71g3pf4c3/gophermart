// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/71g3pf4c3/gophermart/config"
	"github.com/71g3pf4c3/gophermart/internal/controller/restapi"
	"github.com/71g3pf4c3/gophermart/internal/repo/postgres"
	"github.com/71g3pf4c3/gophermart/internal/usecase/accrual"
	"github.com/71g3pf4c3/gophermart/internal/usecase/balance"
	ucorder "github.com/71g3pf4c3/gophermart/internal/usecase/order"
	"github.com/71g3pf4c3/gophermart/internal/usecase/user"
	"github.com/71g3pf4c3/gophermart/pkg/httpserver"
	"github.com/71g3pf4c3/gophermart/pkg/jwt"
	"github.com/71g3pf4c3/gophermart/pkg/logger"
	pgpkg "github.com/71g3pf4c3/gophermart/pkg/postgres"
)

type servers struct {
	http *httpserver.Server
}

func initServers(cfg *config.Config, jwtManager *jwt.Manager, l logger.Interface, userService *user.UseCase, orderService *ucorder.UseCase, balanceService *balance.UseCase) servers {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, jwtManager, l, userService, orderService, balanceService)

	return servers{http: httpServer}
}

func (s *servers) startServers() {
	s.http.Start()
}

func (s *servers) waitForShutdown(cancel context.CancelFunc, l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	cancel()
	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.LogLevel)

	pg, err := pgpkg.New(cfg.DatabaseURI, pgpkg.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	if err = runMigrations(cfg.DatabaseURI); err != nil {
		l.Fatal(fmt.Errorf("app - Run - migrate: %w", err))
	}

	jwtManager := jwt.New(cfg.JWTSecret, cfg.JWT.TokenExpiry)
	userRepo := postgres.NewUserRepo(pg)
	orderRepo := postgres.NewOrderRepository(pg)
	balanceRepo := postgres.NewBalanceRepository(pg)

	userService := user.New(userRepo, jwtManager)
	orderService := ucorder.New(orderRepo)
	balanceService := balance.New(balanceRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accrualProcessor := accrual.New(cfg.AccrualSystemAddress, cfg.Accrual.PollInterval, l, orderRepo)
	go accrualProcessor.Run(ctx)

	s := initServers(cfg, jwtManager, l, userService, orderService, balanceService)
	s.startServers()
	s.waitForShutdown(cancel, l)
}
