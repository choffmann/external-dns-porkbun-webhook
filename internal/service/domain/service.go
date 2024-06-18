package domain

import (
	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/service"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/service/domain/porkbun"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage"
)

func NewProvider(cfg *config.Config, repos *storage.Repository) *service.Provider {
	return &service.Provider{
		PorkbunProvider: porkbun.NewPorkbunProvider(cfg, repos.Porkbun),
	}
}
