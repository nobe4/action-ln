# Examples

This is a collection of various examples that try to showcase the various ways
this action can be used.

Assumptions:
- The default repo is `nobe4/action-ln`.
- Valid [authentication](/docs/authentication.md).

## Update LICENSE once a year

```yaml
# .github/workflows/ln.yaml
name: ln
on:
  schedule:
    - cron: "0 0 1 1 *"
jobs:
  ln:
    runs-on: ubuntu-latest
    steps:
      - uses: nobe4/action-ln@v0
        with:
          token: ${{ secret.GITHUB_ORG_TOKEN }}
          config-path: .action-ln-config.yaml

# .action-ln-config.yaml
links:
  - from:
      path: LICENSE
    to:
      repo: gh-not
  - from:
      path: LICENSE
    to:
      repo: dotfiles

# nobe4/action-ln:LICENSE@main => nobe4/gh-not:LICENSE@main
# nobe4/action-ln:LICENSE@main => nobe4/dotfiles:LICENSE@main
```

## Pull config from public repositories on manual request

```yaml
# .github/workflows/ln.yaml
name: ln
on:
  schedule:
    - workflows_dispatch:
jobs:
  ln:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: nobe4/action-ln@v0

# .github/ln-config.yaml
links:
  - from: cli/go-gh:pkg/text/text.go
    to: internal/text/text.go
  - from: cli/cli:.golangci.yml

# cli/go-gh:pkg/text/text.go@trunk => nobe4/action-ln:internal/text/text.go@main
# cli/cli:.golangci.yml@trunk      => nobe4/action-ln:.golangci.yml@main
```

## Link between two other repositories on push

The `ln-config.yaml` lives in a repo, but only read from/write to other repositories.

```yaml
# .github/workflows/ln.yaml
name: ln
on:
  schedule:
    - push: [main]
jobs:
  ln:
    runs-on: ubuntu-latest
    steps:
      - uses: nobe4/action-ln@v0
        with:
          app-id: ${{ secrets.ACTION_LN_APP_ID }}
          app-private-key: ${{ secrets.ACTION_LN_APP_PRIVATE_KEY }}
          app-install-id: ${{ secrets.ACTION_LN_APP_INSTALL_ID }}

# .github/ln-config.yaml
links:
  - from:
      owner: ccoVeille
      repo: golangci-lint-config-examples
      path: 90/daredevil/.golangci.yml
      ref: v1.1.0
    to:
      owner: nobe4
      repo: gh-not
      path: .golangci.yaml
      ref: edge

  - from:
      repo: gh-not
      path: .goreleaser.yaml
    to:
      repo: safe

# ccoVeille/golangci-lint-config-examples:90/daredevil/.golangci.yml@v1.1.0 => nobe4/action-ln:.golangci.yaml@edge
# nobe4/gh-not:.goreleaser.yaml@main                                        => nobe4/safe:.goreleaser.yaml@main
```
