# action-ln

Link files between repositories.

## Config

Use `ln-config.yaml` at the root of your repository, or specify a custom path
with the `config-path` action input.

The format is:

```yaml
links:
  # Local files
  - from:
      path: path/to/file
    to:
      path: path/to/other/file

  # Remote files
  - from:
      repo: org/repo
      path: path/to/file
    to:
      repo: org/repo
      path: path/to/other/file
```

## Random notes

- If within a repo, at least the following permissions are needed:

  ```yaml
  permissions:
    contents: write
    pull-requests: write
  ```

- If using from an org, you need to enable `Allow GitHub Actions to create and
approve pull requests` from
  `https://github.com/organizations/<org>/settings/actions`

- For Classic tokens `repo` scope is needed, assuming you have write access to
  all the updated repositories.
