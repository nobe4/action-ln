package main

import (
	"context"
	"fmt"
	"os"

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

	f, err := g.GetFile(ctx, github.Repo{
		Owner: github.User{Login: "frozen-fishsticks"},
		Repo:  "action-ln-test-2",
	}, "README.md")
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, f.Content)
}
