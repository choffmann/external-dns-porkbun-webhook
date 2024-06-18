package porkbun

import (
	"github.com/choffmann/external-dns-porkbun-webhook/config"
)

type PorkbunRepository struct {
}

func NewPorkbunRepository(cfg *config.Config) *PorkbunRepository {
	return &PorkbunRepository{}
}
