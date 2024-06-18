package http

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"sigs.k8s.io/external-dns/provider"
	"sigs.k8s.io/external-dns/provider/webhook/api"
)

type WebhookServer struct {
	port     int
	cfg      config.ServerConfig
	provider provider.Provider
}

func NewWebhookServer(prov provider.Provider, port int, cfg config.ServerConfig) *WebhookServer {
	return &WebhookServer{port: port, cfg: cfg, provider: prov}
}

func (s *WebhookServer) Run(ctx context.Context, startedChan chan<- bool) error {
	apiSrv := api.WebhookServer{Provider: s.provider}

	m := http.NewServeMux()
	m.HandleFunc("/", apiSrv.NegotiateHandler)
	m.HandleFunc("/records", apiSrv.RecordsHandler)
	m.HandleFunc("/adjustendpoints", apiSrv.AdjustEndpointsHandler)

	logger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)

	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	srv := &http.Server{
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.ReadTimeout,
		Addr:         addr,
		Handler:      m,
		ErrorLog:     logger,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("server could not listen on port", slog.String("addr", addr), slog.String("err", err.Error()))
		startedChan <- false
		return err
	}

	startedChan <- true

	return srv.Serve(l)
}
