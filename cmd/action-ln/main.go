package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/version"
)

const (
	endpoint = "https://api.github.com"
)

func main() {
	fmt.Fprintln(os.Stdout, version.String())

	e, err := environment.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Environment:", e)

	g := github.New(e.Token, endpoint)
	ctx := context.TODO()

	p, err := g.GetOrCreatePull(ctx, e.Repo, "main", "test-1", "title", "body")
	fmt.Fprintf(os.Stdout, "Pull request: %+v\n", p)
	fmt.Fprintf(os.Stdout, "Err: %+v\n", err)
}
