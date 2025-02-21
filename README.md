<div align="center">
  <img width="400" src="https://github.com/nobe4/action-ln/blob/main/docs/logo.png" />
  <small>Logo by <a href="https://www.instagram.com/malohff">@malohff</a></small>
</div>

# `action-ln`

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
