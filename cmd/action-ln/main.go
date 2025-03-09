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

	f := github.File{
		Repo: github.Repo{
			Owner: github.User{Login: "frozen-fishsticks"},
			Repo:  "action-ln-test-1",
		},
		Path: "README.md",
	}

	if err = g.GetFile(ctx, &f); err != nil {
		panic(err)
	}

	f.Content += "sup\n"

	f2, err := g.UpdateFile(ctx, f, "main", "test #125")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "File: %+v\n", f)
	fmt.Fprintf(os.Stdout, "File2: %+v\n", f2)
}
