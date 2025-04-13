package main

import (
	"binancetrading/config"
	"binancetrading/internal/application/handler"
	"binancetrading/internal/application/service"
	"binancetrading/internal/grpc/server"
	"binancetrading/internal/infrastructure/exchange"
	"binancetrading/internal/infrastructure/storage"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	store, err := storage.NewGORMStorage(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer store.Close()

	exchangeClient := exchange.NewBinanceClient(cfg.Binance.RetryDelay, cfg.Binance.MaxRetries)
	candleService := service.NewCandlestickService(exchangeClient, store)

	if err := candleService.StartAggregation(cfg.Symbols); err != nil {
		log.Fatalf("Failed to start aggregation: %v", err)
	}

	// Fiber HTTP server
	app := fiber.New()
	apiHandler := handler.NewAPIHandler(store)
	apiHandler.SetupRoutes(app)
	go func() {
		if err := app.Listen(":" + cfg.APIPort); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// gRPC server
	go func() {
		if err := server.StartGRPCServer(candleService); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	candleService.Wait()
}
