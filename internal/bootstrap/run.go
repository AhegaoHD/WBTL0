package bootstrap

import (
	"context"
	"fmt"
	"github.com/AhegaoHD/WBTL0/internal/repository/pgrepo"
	"os"
	"os/signal"
	"syscall"

	"github.com/AhegaoHD/WBTL0/config"
	"github.com/AhegaoHD/WBTL0/internal/controller/http"
	"github.com/AhegaoHD/WBTL0/internal/repository/cache"
	"github.com/AhegaoHD/WBTL0/internal/repository/queue_subscribe"
	"github.com/AhegaoHD/WBTL0/internal/service"
	"github.com/AhegaoHD/WBTL0/pkg/httpserver"
	"github.com/AhegaoHD/WBTL0/pkg/logger"
	"github.com/AhegaoHD/WBTL0/pkg/nats"
	"github.com/AhegaoHD/WBTL0/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New("")

	pg, err := postgres.New(postgres.GetConnString(&cfg.Db), postgres.MaxPoolSize(cfg.Db.MaxPoolSize))
	if err != nil {
		l.Fatal("APP - START - POSTGRES INI PROBLEM: %v", err)
	}
	defer pg.Close()

	err = pg.Pool.Ping(context.Background())
	if err != nil {
		l.Fatal("APP - START - POSTGRES INI PROBLEM: %v", err)
		return
	}

	sc, err := nats.New("test-cluster", "client-1", &cfg.Nats.URL)
	if err != nil {
		l.Fatal("APP - START - NATS INI PROBLEM: %v", err)
	}

	pgRepo := pgrepo.NewOrderRepo(pg)
	cacheRepo := cache.NewOrderCache(pgRepo)
	natsRepo := queue_subscribe.NewNatsRepository(sc.Conn)
	natsService := service.NewNatsService(l, natsRepo, cacheRepo, pgRepo)

	_, err = natsService.StartListening("subject", "queueGroup")
	if err != nil {
		l.Fatal("APP - START - NATS INI PROBLEM: %v", err)
	}

	orderService := service.NewOrderService(pgRepo)
	httpController := http.NewOrderController(l, orderService)
	if err != nil {
		l.Fatal("APP - START - CONTROLLER INIT: %v", err)
	}

	httpServer := httpserver.New(httpController,
		httpserver.Port(cfg.HttpServer.Addr),
		httpserver.ReadTimeout(cfg.HttpServer.ReadTimeout),
		httpserver.WriteTimeout(cfg.HttpServer.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.HttpServer.ShutdownTimeout),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("RUNNING APP:%v VERSION:%v", cfg.App.Name, cfg.App.Version)

	select {
	case s := <-interrupt:
		l.Info("APP - RUN - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("APP - RUN - HTTPSERVER.NOTIFY: %v", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("APP - RUN - HTPPSERVER.SHUTDOWN: %v", err))
	}
}
