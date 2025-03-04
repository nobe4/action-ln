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

	// if u, err := g.GetUser(ctx); err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error getting user:", err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "User:", u.Login)
	// }
	//
	// if b, err := g.GetDefaultBranch(ctx, e.Repo); err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error getting default branch:", err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "Default branch:", b)
	// }

	// branch := "test"
	//
	// var defaultBranch github.Branch
	//
	// if defaultBranch, err = g.GetBranch(ctx, e.Repo, "main"); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error getting default %s %+v", branch, err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "Default branch:", defaultBranch)
	// }
	//
	// if b, err := g.GetOrCreateBranch(ctx, e.Repo, branch, defaultBranch.Commit.SHA); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error getting or creating branch %s %+v\n", branch, err)
	// } else {
	// 	fmt.Fprintln(os.Stdout, "Branch:", b)
	// }

	path := "README.md"

	c, err := g.GetContent(
		ctx,
		e.Repo,
		path,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting contents:", err)
		os.Exit(1)
	}

	c.Content += "\n\nHello, World!"

	fmt.Println(g.CreateOrUpdateContent(ctx, e.Repo, c, "main", "Update README.md"))
}
