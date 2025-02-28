package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
)

const (
	endpoint = "https://api.github.com"
)

func main() {
	e, err := environment.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Environment:", e)

	g := github.New(e.Token, endpoint)
	ctx := context.TODO()

	u, err := g.GetUser(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting user:", err)
	} else {
		fmt.Fprintln(os.Stdout, "User:", u.Login)
	}

	c, err := g.GetContent(
		ctx,
		e.Repo,
		"README.md",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting contents:", err)
	} else {
		fmt.Fprintln(os.Stdout, "Content:\n", c.Content)
	}

	b, err := g.GetDefaultBranch(ctx, e.Repo)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting default branch:", err)
	} else {
		fmt.Fprintln(os.Stdout, "Default branch:", b)
	}
}
