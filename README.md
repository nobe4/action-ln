<div align="center">
  <img width="300" src="https://github.com/nobe4/action-ln/blob/main/docs/logo.png" /> <br>
  <sub>Logo by <a href="https://www.instagram.com/malohff">@malohff</a></sub>
</div>

# `action-ln`

Link files between repositories.

- [Authentication](/docs/authentication.md)

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

See [`example.config.yaml`](./example.config.yaml) for an example.

## Running locally

```shell
npm start -- [--token="..."] [--config="..."] [--noop]
```

## Development

To test the action from a branch, run `npm run build` and `npm run build:add`.
Then, commit the `dist` folder, it make CI fails but allows you to use the
branch name/commit sha as a version to run the action on.

Once testing is done, run `npm run build:clean` before you merge to the main branch.

> [!NOTE]
> There's no need to push the code, just pushing the dist is enough for testing.
