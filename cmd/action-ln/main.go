package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/nobe4/action-ln/internal/client"
	"github.com/nobe4/action-ln/internal/client/noop"
	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/ln"
	"github.com/nobe4/action-ln/internal/log"
	glog "github.com/nobe4/action-ln/internal/log/github"
	"github.com/nobe4/action-ln/internal/log/plain"
)

func main() {
	ctx := context.TODO()

	e, err := environment.Parse()
	if err != nil {
		log.Error("Environment parsing failed", "reason", err)
		os.Exit(1)
	}

	setLogger(e.Debug, e.OnAction)

	e.PrintDebug()

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

//nolint:revive // debug here is expected.
func setLogger(debug, onAction bool) {
	o := log.Options{Level: slog.LevelInfo}
	if debug {
		o.Level = slog.LevelDebug
	}

	var h slog.Handler
	if onAction {
		h = glog.New(os.Stdout, o)
	} else {
		h = plain.New(os.Stdout, o)
	}

	slog.SetDefault(slog.New(h))
}
