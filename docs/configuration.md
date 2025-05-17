# Configuration

> [!IMPORTANT]
> See[`all-cases.yaml`](../internal/config/fixtures/all-cases.yaml) for a comprehensive
> list of supported cases in a YAML format. It is the canonical reference for
> configuration.

The configuration lives in a YAML file that follows the following schema:

```yaml
links:
  - from: owner/repo:path/to/file@main
    to: path/to/file

  # ... and more
```

> [!NOTE]
> A file is usually printed by the following string:
> `owner/repo:path/to/file@ref`
> A link is usually printed by the following string:
> `from -> to`

## Link

A link is composed of two [files](#file)
- `from` is the _source_ of the link, where the file is _read_.
- `to` is the _destination_ of the link, where the file is _written_.

## File

A file is the logical representation of a file on GitHub.

It is composed of 3 parts:

- `repo`: the full name of a repository, with `owner` and `repo` parts.
- `path`: the path relative to the root of the repository
- `ref`: a valid git commit, tag, or branch (TBD #34)
    It defaults to the default branch of the targeted repository.

## Defaults

- `link`: a [link](#link) whose values are used if not further specified.
- More TBD
