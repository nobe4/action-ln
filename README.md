# action-ln

Link files between repositories.

## Config

Use `ln-config.yaml` at the root of your repository, or specify a custom path
with the `config-path` action input.

The format is:

```yaml
# Local files
links:
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
