package porkbun

import (
	"context"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
)

type PorkbunProvider struct {
	repo storage.PorkbunRepository
}

func NewPorkbunProvider(cfg *config.Config, repo storage.PorkbunRepository) *PorkbunProvider {
	return &PorkbunProvider{repo}
}

func (p *PorkbunProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	return nil, nil
}

func (p *PorkbunProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	return nil
}

func (p *PorkbunProvider) AdjustEndpoints(endpoints []*endpoint.Endpoint) ([]*endpoint.Endpoint, error) {
	return nil, nil
}

func (p *PorkbunProvider) GetDomainFilter() endpoint.DomainFilter {
	return endpoint.NewDomainFilter([]string{})
}
