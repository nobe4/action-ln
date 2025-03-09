package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
)

func main() {
	ctx := context.TODO()

	e, err := environment.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Environment:", e)

	g := github.New(e.Endpoint)

	if err = g.Auth(ctx,
		e.Token,
		e.App.ID,
		e.App.PrivateKey,
		e.App.InstallID,
	); err != nil {
		panic(err)
	}

	f := github.File{
		Repo: github.Repo{
			Owner: github.User{Login: "frozen-fishsticks"},
			Repo:  "action-ln-test-0",
		},
		Path: "ln-config.yaml",
	}
	if err = g.GetFile(ctx, &f); err != nil {
		panic(err)
	}

	c, err := config.Parse(strings.NewReader(f.Content))
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Content:", f.Content)
	fmt.Fprintln(os.Stdout, "Config:", c)
}
