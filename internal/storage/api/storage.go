package api

import (
	"net/http"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage/api/porkbun"
)

func NewRepository(cfg *config.Config) (*storage.Repository, error) {
	porkbunRepo, err := porkbun.NewPorkbunRepository(cfg, http.DefaultClient)
	if err != nil {
		return nil, err
	}
	return &storage.Repository{
		Porkbun: porkbunRepo,
	}, nil
}
