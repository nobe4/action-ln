<div align="center">
  <img width="300" src="https://github.com/nobe4/action-ln/blob/main/docs/logo.png" /> <br>
  <sub>Logo by <a href="https://www.instagram.com/malohff">@malohff</a></sub>
</div>

# `action-ln`

> Link files between repositories.

This action creates a _link_ between files in various places. When the source is
updated, the destination is as well.

It works by using the GitHub API to read files and create Pull Requests where an
update is needed. You can specify the source, destination, and schedule for the
synchronization.

> [!TIP]
> The authentication for this can be rather tricky, make sure you read
> [authentication](/docs/authentication.md) to get familiar with the various
> methods.

## Quickstart

1. Create a config file in `.github/ln-config.yaml`.

    E.g. [`ln-config.yaml`](.github/ln-config.yaml)

2. Create a workflow.

    ```yaml
    uses: nobe4/action-ln@v0
    ```

    E.g. [`ln.yaml`](.github/workflows/ln.yaml)

## Further readings

- [Authentication](/docs/authentication.md)
- [Configuration](/docs/configuration.md)
- [Development](/docs/development.md)
- [Examples](/docs/examples.md)
