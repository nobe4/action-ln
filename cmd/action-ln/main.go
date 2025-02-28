package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nobe4/action-ln/internal/github"
)

const (
	endpoint = "https://api.github.com"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		os.Getenv("INPUT_TOKEN")
	}

	if token == "" {
		panic("GITHUB_TOKEN/input 'token' is required")
	}

	g := github.New(token, endpoint)

	u, err := g.GetUser(context.TODO())
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Hello, "+u.Login)
}
