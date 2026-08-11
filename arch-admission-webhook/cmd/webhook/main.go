// Command webhook runs the arch validating admission webhook. It observes
// pod admissions, classifies each image's arch support, and emits a statsd
// metric to the cluster Datadog agent. It never blocks admission.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carta/arch-admission-webhook/internal/classify"
	"github.com/carta/arch-admission-webhook/internal/metrics"
	"github.com/carta/arch-admission-webhook/internal/webhook"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	addr := env("LISTEN_ADDR", ":8443")
	certFile := env("TLS_CERT_FILE", "/etc/webhook/certs/tls.crt")
	keyFile := env("TLS_KEY_FILE", "/etc/webhook/certs/tls.key")

	// Datadog agent host is injected by the DownwardAPI (status.hostIP).
	ddHost := env("DD_AGENT_HOST", "localhost")
	ddAddr := net.JoinHostPort(ddHost, env("DD_DOGSTATSD_PORT", "8125"))

	emitter, err := metrics.New(ddAddr)
	if err != nil {
		log.Error("statsd init failed", "addr", ddAddr, "err", err)
		os.Exit(1)
	}
	defer emitter.Close()

	classifier := classify.NewCache(classify.StubClassifier{}, time.Hour)
	h := webhook.New(classifier, emitter, log, 8, 1024)

	mux := http.NewServeMux()
	mux.Handle("/validate", h)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("listening", "addr", addr)
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Info("shutdown complete")
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
