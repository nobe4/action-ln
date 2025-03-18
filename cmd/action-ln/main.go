package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/nobe4/action-ln/internal/log"
	"github.com/nobe4/action-ln/internal/log/plain"
)

func main() {
	debug := true

	o := log.Options{
		Level: slog.LevelInfo,
	}
	if debug {
		o.Level = slog.LevelDebug
	}

	handler := plain.New(os.Stdout, o)
	// handler := github.New(os.Stdout, o)
	slog.SetDefault(slog.New(handler))

	log.Group("test")
	log.Info("message", "a", 1, "b", "c")
	log.Debug("message")
	log.Error("message", "a", []string{"x", "y", "z"})
	log.Warn("message")
	log.Notice("message", "a", time.Now())
	log.GroupEnd()
}

// func main() {
// 	ctx := context.TODO()
//
// 	e, err := environment.Parse()
// 	if err != nil {
// 		fmt.Fprintf(os.Stdout, "Error parsing environment\n%v\n", err)
// 		os.Exit(1)
// 	}
//
// 	fmt.Fprintln(os.Stdout, "Environment:", e)
//
// 	g := github.New(e.Endpoint)
//
// 	if err = g.Auth(ctx,
// 		e.Token,
// 		e.App.ID,
// 		e.App.PrivateKey,
// 		e.App.InstallID,
// 	); err != nil {
// 		fmt.Fprintf(os.Stdout, "Error authenticating\n%v\n", err)
// 		os.Exit(1)
// 	}
//
// 	if err := ln.Run(ctx, e, g); err != nil {
// 		fmt.Fprintf(os.Stdout, "Error running action-ln\n%v\n", err)
// 		os.Exit(1)
// 	}
// }
