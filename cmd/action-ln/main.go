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

	// if u, err := g.GetUser(ctx); err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error getting user:", err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "User:", u.Login)
	// }
	//
	// if c, err := g.GetContent(
	// 	ctx,
	// 	e.Repo,
	// 	"README.md",
	// ); err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error getting contents:", err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "Content:\n", c.Content)
	// }

	// if b, err := g.GetDefaultBranch(ctx, e.Repo); err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error getting default branch:", err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "Default branch:", b)
	// }

	branch := "test"

	var defaultBranch github.Branch

	if defaultBranch, err = g.GetBranch(ctx, e.Repo, "main"); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting default %s %+v", branch, err)
	} else {
		fmt.Fprintln(os.Stdout, "Default branch:", defaultBranch)
	}

	if b, err := g.GetOrCreateBranch(ctx, e.Repo, branch, defaultBranch.Commit.SHA); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting or creating branch %s %+v\n", branch, err)
	} else {
		fmt.Fprintln(os.Stdout, "Branch:", b)
	}
}
