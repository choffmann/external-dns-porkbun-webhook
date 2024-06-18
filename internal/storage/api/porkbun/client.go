package porkbun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/entities/porkbun"
)

const baseURL string = "https://api.porkbun.com/api/json/v3"

type apiCredentials struct {
	Apikey       string `json:"apikey,omitempty"`
	Secretapikey string `json:"secretapikey,omitempty"`
}

type PorkbunRepository struct {
	cred   apiCredentials
	client *http.Client
}

type porkbunResponse struct {
	Status string `json:"status"`
}

func NewPorkbunRepository(cfg *config.Config, client *http.Client) (*PorkbunRepository, error) {
	repo := &PorkbunRepository{
		cred: apiCredentials{
			cfg.ApiKey,
			cfg.ApiSecret,
		},
		client: client,
	}

	valid, err := repo.checkValidCredentials()
	if err != nil {
		slog.Error("", slog.String("error", err.Error()))
		return nil, err
	}
	if !valid {
		return nil, errors.New("credentials are not valid!")
	}

	return repo, nil
}

func makeApiRequest[T, K any](ctx context.Context, r *PorkbunRepository, path string, reqBody T, resBody K) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, &buf)
	if err != nil {
		slog.Error("", slog.String("error", err.Error()))
		return err
	}

	res, err := r.client.Do(req)
	if err != nil {
		slog.Error("", slog.String("error", err.Error()))
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Debug("porkbun credentials are not valid")
		return errors.New("response was not status ok")
	}

	if err := json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return err
	}

	return nil
}

func (r *PorkbunRepository) checkValidCredentials() (bool, error) {
	slog.Debug("checking porkbun credentials")

	if err := makeApiRequest(context.Background(), r, "/ping", &r.cred, &struct{}{}); err != nil {
		return false, err
	}

	slog.Debug("porkbun credentials are valid")
	return true, nil
}

func (r *PorkbunRepository) GetZones(ctx context.Context) ([]*porkbun.Zone, error) {
	type response struct {
		porkbunResponse
		Domains []*porkbun.Zone `json:"domains"`
	}

	var res response
	if err := makeApiRequest(ctx, r, "/domain/listAll", &r.cred, &res); err != nil {
		return nil, err
	}

	return res.Domains, nil
}

func (r *PorkbunRepository) GetRecords(ctx context.Context, domain string) ([]*porkbun.Record, error) {
	type response struct {
		porkbunResponse
		Records []*porkbun.Record `json:"records"`
	}

	var res response
	if err := makeApiRequest(ctx, r, "/dns/retrieve/"+domain, &r.cred, &res); err != nil {
		return nil, err
	}

	return res.Records, nil
}

func (r *PorkbunRepository) CreateRecord(ctx context.Context, domain string, record *porkbun.Record) (*porkbun.Record, error) {
	type request struct {
		apiCredentials
		*porkbun.Record
	}

	req := request{
		r.cred,
		record,
	}
	if err := makeApiRequest(ctx, r, "/dns/create/"+domain, &req, &struct{}{}); err != nil {
		return nil, err
	}

	return record, nil
}

func (r *PorkbunRepository) UpdateRecord(ctx context.Context, domain string, record *porkbun.Record) (*porkbun.Record, error) {
	type request struct {
		apiCredentials
		*porkbun.Record
	}

	req := request{
		r.cred,
		record,
	}
	if err := makeApiRequest(ctx, r, fmt.Sprintf("/dns/editByNameType/%s/%s/%s", domain, record.Type, record.Name), &req, &struct{}{}); err != nil {
		return nil, err
	}

	return record, nil
}

func (r *PorkbunRepository) DeleteRecord(ctx context.Context, domain string, record *porkbun.Record) error {
	type request struct {
		apiCredentials
		*porkbun.Record
	}

	req := request{
		r.cred,
		record,
	}
	if err := makeApiRequest(ctx, r, fmt.Sprintf("/dns/deleteByNameType/%s/%s/%s", domain, record.Type, record.Name), &req, &struct{}{}); err != nil {
		return err
	}

	return nil
}
