package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nobe4/dent.go"

	"github.com/nobe4/action-ln/internal/config"
)

const (
	endpoint = "https://api.github.com"
)

func main() {
	src := dent.DedentString(`
	links:
	  - from:
	      repo: x
	      path: y
	    to:
	      repo: a
	      path: b
	`)

	c, err := config.Parse(strings.NewReader(src))
	fmt.Fprintf(os.Stdout, "Config: %+v\n", c)
	fmt.Fprintf(os.Stdout, "Err: %+v\n", err)
}
