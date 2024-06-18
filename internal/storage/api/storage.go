package api

import (
	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage/api/porkbun"
)

func NewRepository(cfg *config.Config) *storage.Repository {
	porkbunRepo := porkbun.NewPorkbunRepository(cfg)
	return &storage.Repository{
		Porkbun: porkbunRepo,
	}
}
