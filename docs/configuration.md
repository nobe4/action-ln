# Configuration

The configuration lives in a YAML file that follows the following schema:

```yaml
links:
  - from:
      repo: nobe4/repo-1
      path: README.md
    to:
      repo: nobe4/repo-2
      path: README.md
  # ... and more
```

## Link

A link is composed of two parts:
- `from` is the _source_ of the link, where the file is _read_.
- `to` is the _destination_ of the link, where the file is _written_.

Both of those are [files](#file).

If `from` to is omitted, then it's equivalent to `from.path` in the current
repo.

## File

A file is the logical representation of a file on GitHub.

It is composed of 3 parts:

- `repo`: the full name of a repository, with `owner` and `repo` parts.

    If the `owner` part is omitted, it defaults to the current owner.
    If the `owner` and `repo` parts are omitted, it defaults to the current
    owner and repo.

    E.g.

    ```yaml
      # => nobe4/action-ln (TBD #33)
      # repo key not set

      # => nobe4/action-ln (TBD #33)
      repo: action-ln

      # => cli/cli
      repo: cli/cli

      # => cli/go-gh (TBD #33)
      owner: cli
      repo: go-gh

      # => nobe4/gh-not (TBD #33)
      repo: gh-not

      # etc.
    ```

- `path`: the path relative to the root of the repository

    E.g.

    ```yaml
      path: README.md
      path: path/to/file
      # etc.
    ```

- `ref`: a valid git commit, tag, or branch (TBD #34)

    It defaults to the default branch of the targeted repository.

    E.g.

    ```yaml
    ref: main
    ref: v0.1.2
    ref: sha123456
    # etc.
    ```

There are multiple ways to references a file (using `from` here, `to` works
similarly):

-
  ```yaml
  from:
    repo: owner/repo
    path: path/to/file
    ref: ref
  ```

- Verbose (TBD #33)
  ```yaml
  from:
    owner: owner
    repo: owner
    path: path/to/file
    ref: ref
  ```

- Short path (TBD #33)
  ```yaml
  from: owner/repo/ref/path/to/file
  ```

- Alternative short path (TBD #33)
  ```yaml
  from: owner/repo:path/to/file@ref
  ```

- Alternative short path with implicit repo and ref (TBD #33)
  ```yaml
  from: path/to/file
  ```

- GitHub full URL (TBD #33)
  ```yaml
  from: https://github.com/owner/repo/blob/ref/path/to/file
  ```
