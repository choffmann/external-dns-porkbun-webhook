package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/entities"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/logger"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/server/http"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/service/domain"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage/api"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error getting config: %v", err)
	}

	lgg, err := logger.GetLogger(cfg.LogFormat)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	slog.SetDefault(lgg)
	_ = slog.SetLogLoggerLevel(cfg.LogLevel.ToSLog())

	slog.Info("Starting external-dns porkbun webhook")

	repo := api.NewRepository(cfg)
	prov := domain.NewProvider(cfg, repo)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	healthStatus := entities.NewHealthStatus()
	heathServer := http.NewHealthServer(healthStatus, cfg.HealthPort, cfg.Health)
	webhook := http.NewWebhookServer(prov.PorkbunProvider, cfg.WebhookPort, cfg.Webhook)

	startedChan := make(chan bool)
	defer close(startedChan)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		heathServer.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		webhook.Run(ctx, startedChan)
	}()

	slog.Debug("Waiting for servers to start")

	started := <-startedChan
	if !started {
		slog.Error("Server could not be started")
		cancel()
    wg.Wait()
    return
	}

	healthStatus.SetReady(true)
	healthStatus.SetHealth(true)

	slog.Info("Server started!")

	wg.Wait()
}
