package main

import (
	"log"
	"sync"
	"time"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/service/domain"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage/api"

	"sigs.k8s.io/external-dns/provider"
	externalDNSApi "sigs.k8s.io/external-dns/provider/webhook/api"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error getting config: %v", err)
	}

	repo := api.NewRepository(cfg)
	prov := domain.NewProvider(cfg, repo)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(prov provider.Provider) {
		defer wg.Done()
		externalDNSApi.StartHTTPApi(prov, nil, 1*time.Second, 1*time.Second, "5000")
	}(prov.PorkbunProvider)

	wg.Wait()
}
