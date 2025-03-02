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

	var b github.Branch

	if b, err = g.GetBranch(ctx, e.Repo, "main"); err != nil {
		fmt.Fprintln(os.Stderr, "Error getting branch", "main", err)
	} else {
		fmt.Fprintln(os.Stdout, "Branch main:", b)
	}

	if nb, err := g.CreateBranch(ctx, e.Repo, "test", b.Commit.SHA); err != nil {
		fmt.Fprintln(os.Stderr, "Error creating branch", "test", err)
	} else {
		fmt.Fprintln(os.Stdout, "Created branch:", nb)
	}
}
