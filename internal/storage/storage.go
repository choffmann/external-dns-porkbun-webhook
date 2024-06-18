package storage

import (
	"context"

	"github.com/choffmann/external-dns-porkbun-webhook/internal/entities/porkbun"
)

type PorkbunRepository interface {
	GetZones(ctx context.Context) ([]*porkbun.Zone, error)
	GetRecords(ctx context.Context, domain string) ([]*porkbun.Record, error)
	CreateRecord(ctx context.Context, domain string, record *porkbun.Record) (*porkbun.Record, error)
	UpdateRecord(ctx context.Context, domain string, record *porkbun.Record) (*porkbun.Record, error)
	DeleteRecord(ctx context.Context, domain string, record *porkbun.Record) error
}

type Repository struct {
	Porkbun PorkbunRepository
}
