package porkbun

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/entities/porkbun"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/storage"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/utils"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

type PorkbunProvider struct {
	provider.BaseProvider
	repo         storage.PorkbunRepository
	domainFilter endpoint.DomainFilter
}

func getDomainFilter(cfg config.ExternalDnsDomainConfig) endpoint.DomainFilter {
	var domainFilter endpoint.DomainFilter

	if cfg.RegexDomainFilter != "" {
		domainFilter = endpoint.NewRegexDomainFilter(
			regexp.MustCompile(cfg.RegexDomainFilter),
			regexp.MustCompile(cfg.RegexDomainExclusion),
		)
	} else {
		domainFilter = endpoint.NewDomainFilterWithExclusions(cfg.DomainFilter, cfg.ExcludeDomains)
	}

	return domainFilter
}

func NewPorkbunProvider(cfg *config.Config, repo storage.PorkbunRepository) *PorkbunProvider {
	return &PorkbunProvider{repo: repo, domainFilter: getDomainFilter(cfg.DomainConfig)}
}

func (p *PorkbunProvider) zones(ctx context.Context) ([]*porkbun.Zone, error) {
	zones, err := p.repo.GetZones(ctx)
	if err != nil {
		return nil, err
	}

	return utils.Filter(zones, func(zone *porkbun.Zone) bool {
		return p.domainFilter.Match(zone.Domain)
	}), nil
}

func (p *PorkbunProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	zones, err := p.zones(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*endpoint.Endpoint, 0)

	for _, zone := range zones {
		records, err := p.repo.GetRecords(ctx, zone.Domain)
		if err != nil {
			return nil, err
		}

		endpoints := utils.MapFilter(records, func(item *porkbun.Record) **endpoint.Endpoint {
			if !provider.SupportedRecordType(item.Type) {
				return nil
			}
			ep := endpoint.NewEndpoint(item.Name, item.Type, item.Content)

			ttl, err := strconv.ParseInt(item.TTL, 10, 64)
			if err != nil {
				return nil
			}
			ep.RecordTTL = endpoint.TTL(ttl)
			return &ep
		})
		result = append(result, endpoints...)
	}

	return mergeEndpointsByNameType(result), nil
}

func mergeEndpointsByNameType(endpoints []*endpoint.Endpoint) []*endpoint.Endpoint {
	endpointsByNameType := make(map[string][]*endpoint.Endpoint)

	for _, ep := range endpoints {
		key := fmt.Sprintf("%s-%s", ep.DNSName, ep.RecordType)
		endpointsByNameType[key] = append(endpointsByNameType[key], ep)
	}

	if len(endpointsByNameType) == len(endpoints) {
		return endpoints
	}

  result := make([]*endpoint.Endpoint, 0)
  for _, endpoints := range endpointsByNameType {
    e := utils.Reduce(endpoints, endpoints[0], func(e *endpoint.Endpoint, acc *endpoint.Endpoint) *endpoint.Endpoint {
      acc.Targets = append(acc.Targets, e.Targets...)
      return acc
    })

    result = append(result, e)
  }

  return endpoints
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
