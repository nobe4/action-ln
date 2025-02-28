package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nobe4/action-ln/internal/github"
)

const (
	endpoint = "https://api.github.com"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		token = os.Getenv("INPUT_TOKEN")
	}

	if token == "" {
		fmt.Fprintln(os.Stdout, "Environment variables:")

		for _, env := range os.Environ() {
			parts := strings.Split(env, "=")
			fmt.Fprintln(os.Stdout, parts[0])
		}

		panic("GITHUB_TOKEN/input 'token' is required")
	}

	g := github.New(token, endpoint)
	ctx := context.TODO()
	repo := github.Repo{Owner: github.User{Login: "nobe4"}, Repo: "action-ln"}

	u, err := g.GetUser(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting user:", err)
	} else {
		fmt.Fprintln(os.Stdout, "User:", u.Login)
	}

	c, err := g.GetContent(
		ctx,
		repo,
		"go.mod",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting user:", err)
	} else {
		fmt.Fprintln(os.Stdout, "Content:\n", c.Content)
	}

	b, err := g.GetDefaultBranch(ctx, repo)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting default branch:", err)
	} else {
		fmt.Fprintln(os.Stdout, "Default branch:", b)
	}
}
