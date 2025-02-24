<div align="center">
  <img width="300" src="https://github.com/nobe4/action-ln/blob/main/docs/logo.png" /> <br>
  <sub>Logo by <a href="https://www.instagram.com/malohff">@malohff</a></sub>
</div>

# `action-ln`

Link files between repositories.

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
- [Examples](/docs/examples.md)

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
