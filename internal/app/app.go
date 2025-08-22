package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"order-notification/config"
	"order-notification/internal/handler"
	"order-notification/internal/kafka"
	"order-notification/internal/repo/cache"
	"order-notification/internal/repo/pgdb"
	"order-notification/internal/service"
	"order-notification/pkg/httpserver"
	"order-notification/pkg/postgres"
)

// Run Запуск приложения
func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	SetLogrus(cfg.Log.Level)

	logrus.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	ordersRepo := pgdb.NewOrdersRepo(pg)
	cacheRepo := cache.NewCache(cfg.Cache.CacheSize, cfg.Cache.TTL)

	orderService := service.NewOrderService(pg)
	orderService.Repo = ordersRepo
	orderService.Cache = cacheRepo

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewConsumer(cfg.Brokers, cfg.Topic, cfg.GroupID, orderService)
	go func() {
		consumer.Run(ctx)
	}()

	logrus.Info("Starting http server...")
	logrus.Debugf("Server port: %s", cfg.HTTP.Port)

	orderHandler := handler.NewOrderHandler(orderService)
	router := handler.NewRouter(orderHandler)

	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(10*time.Second),
		httpserver.WriteTimeout(10*time.Second),
	)

	logrus.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logrus.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logrus.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	logrus.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		logrus.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
