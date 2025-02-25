> Still in progress, here be 🐉

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
