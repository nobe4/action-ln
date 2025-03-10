package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/ln"
)

func main() {
	ctx := context.TODO()

	e, err := environment.Parse()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error parsing environment\n%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "Environment:", e)

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
