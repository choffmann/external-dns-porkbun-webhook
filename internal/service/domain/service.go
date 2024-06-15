package domain

import "github.com/choffmann/external-dns-porkbun-webhook/internal/service"

func NewService() *service.Service {
	return &service.Service{}
}
