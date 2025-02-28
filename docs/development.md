> Still in progress, here be 🐉

## Release

Until #70 is solved, the release workflow is gonna be wonky.

1. Make a branch for the release
1. Run `script/tag-release` and choose the next tag
1. Run `script/build.sh TAG`, add the results
1. Remove any older bin version from `dist`
1. Commit, PR, Merge
1. On the merge commit, run `script/tag-release` again

## Development

To test the action from a branch, run `script/build.sh` and commit the results.

You can now use the branch name to test your code.

E.g.

```yaml
      - uses: nobe4/action-ln@branch-name
```

> [!TIP]
> There's no need to push the code, just pushing the dist is enough for testing.
