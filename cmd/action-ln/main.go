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
	e, err := environment.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Environment:", e)

	g := github.New(e.Token, e.Endpoint)
	ctx := context.TODO()

	f, err := g.GetFile(
		ctx,
		github.Repo{
			Owner: github.User{Login: "nobe4"},
			Repo:  "action-ln",
		},
		".github/ln-config.yaml",
	)
	if err != nil {
		panic(err)
	}

	c, err := config.Parse(strings.NewReader(f.Content))

	fmt.Fprintf(os.Stdout, "Config: %+v\n", c)
	fmt.Fprintf(os.Stdout, "Err: %+v\n", err)
}
