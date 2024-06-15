package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/server/http"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/service/domain"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error getting config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	service := domain.NewService()
	httpServer := http.NewServer(cfg, service)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		httpServer.Run(ctx)
	}()

	wg.Wait()
}
