package main

import (
	"io"
	"log/slog"
	"testing"
)

func newTestApplication(t *testing.T, cfg *config) *application {
	if cfg == nil {
		cfg = &config{
			env:  "dev",
			port: 4000,
		}
	}

	return &application{
		config: *cfg,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}
