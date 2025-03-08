package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/jwt"
)

func main() {
	ctx := context.TODO()

	e, err := environment.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stdout, "Environment:", e)

	g := github.New(e.Endpoint)
	g.Token = e.Token

	if e.App.Valid() {
		var jwtToken string

		jwtToken, err = jwt.New(time.Now().Unix(), e.App.ID, e.App.PrivateKey)
		if err != nil {
			panic(err)
		}

		if g.Token, err = g.GetAppToken(ctx, e.App.InstallID, jwtToken); err != nil {
			panic(err)
		}
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
