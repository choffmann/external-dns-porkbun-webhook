package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/choffmann/external-dns-porkbun-webhook/config"
	"github.com/choffmann/external-dns-porkbun-webhook/internal/entities"
)

type HealthServer struct {
	port   int
	cfg    config.ServerConfig
	status *entities.HealthStatus
}

func NewHealthServer(healthStatus *entities.HealthStatus, port int, serverConfig config.ServerConfig) *HealthServer {
	return &HealthServer{port: port, cfg: serverConfig, status: healthStatus}
}

func (s *HealthServer) Run(ctx context.Context) error {
  addr := fmt.Sprintf("127.0.0.1:%d", s.port)

	m := http.NewServeMux()
	m.HandleFunc("/healthz", s.healthHandler)
	m.HandleFunc("/ready", s.readyHandler)

	srv := &http.Server{ReadTimeout: s.cfg.ReadTimeout, WriteTimeout: s.cfg.ReadTimeout, Addr: addr, Handler: m}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	l, err := net.Listen("tcp", addr)
	if err != nil {
    slog.Error("server could not listen on port", slog.String("addr", addr), slog.String("err", err.Error()))
		return err
	}

	return srv.Serve(l)
}

func (s *HealthServer) healthHandler(w http.ResponseWriter, _ *http.Request) {
	health := s.status.GetHealth()
	var err error
	if health {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(http.StatusText(http.StatusOK)))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err = w.Write([]byte(http.StatusText(http.StatusServiceUnavailable)))
	}

	if err != nil {
		log.Println("error while sending response", err)
	}
}

func (s *HealthServer) readyHandler(w http.ResponseWriter, _ *http.Request) {
	ready := s.status.GetReady()
	var err error
	if ready {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(http.StatusText(http.StatusOK)))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err = w.Write([]byte(http.StatusText(http.StatusServiceUnavailable)))
	}

	if err != nil {
		log.Println("error while sending response", err)
	}
}
