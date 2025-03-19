package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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
		fmt.Fprintf(os.Stdout, "Error parsing environment\n%v\n", err)
		os.Exit(1)
	}

	o := log.Options{Level: slog.LevelInfo}
	if e.Debug {
		o.Level = slog.LevelDebug
	}

	var h slog.Handler
	if e.OnAction {
		h = glog.New(os.Stdout, o)
	} else {
		h = plain.New(os.Stdout, o)
	}

	slog.SetDefault(slog.New(h))

	log.Debug("Environment", "parsed", e)

	g := github.New(e.Endpoint)

	if err = g.Auth(ctx,
		e.Token,
		e.App.ID,
		e.App.PrivateKey,
		e.App.InstallID,
	); err != nil {
		fmt.Fprintf(os.Stdout, "Error authenticating\n%v\n", err)
		os.Exit(1)
	}

	if err := ln.Run(ctx, e, g); err != nil {
		fmt.Fprintf(os.Stdout, "Error running action-ln\n%v\n", err)
		os.Exit(1)
	}
}
