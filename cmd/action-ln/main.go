package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/nobe4/gh-ln/pkg/client"
	"github.com/nobe4/gh-ln/pkg/client/noop"
	"github.com/nobe4/gh-ln/pkg/github"
	"github.com/nobe4/gh-ln/pkg/ln"
	"github.com/nobe4/gh-ln/pkg/log"

	"github.com/nobe4/action-ln/internal/environment"
	handler "github.com/nobe4/action-ln/internal/log"
)

func main() {
	ctx := context.TODO()

	e, err := environment.Parse()
	if err != nil {
		log.Error("Environment parsing failed", "reason", err)
		os.Exit(1)
	}

	o := log.Options{Level: slog.LevelInfo}
	if e.Debug {
		o.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(handler.New(os.Stdout, o)))

	log.Info("Environment", "parsed", e)

	var c client.Doer = &http.Client{}
	if e.Noop {
		c = noop.New()
	}

	g := github.New(c, e.Endpoint)

	if err = g.Auth(ctx,
		e.Token,
		e.App.ID,
		e.App.PrivateKey,
		e.App.InstallID,
	); err != nil {
		log.Error("Authentication failed", "err", err)
		os.Exit(1)
	}

	if err := ln.Run(ctx, e, g); err != nil {
		log.Error("Running action-ln failed", "err", err)
		os.Exit(1)
	}
}
