package service

import "sigs.k8s.io/external-dns/provider"

type Provider struct {
	PorkbunProvider provider.Provider
}
