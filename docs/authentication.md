# Authentication

This action supports various authentication methods.

See this table for a quick comparison

| method       | ease | public | private | cross-org | trigger CI |
| ---          | ---  | ---    | ---     | ---       | ---        |
| Action token | 🟩   | ✅     | ❌      | ❌        | ❌         |
| Custom token | 🟨   | ✅     | ✅      | ✅        | ✅         |
| GitHub app   | 🟥   | ✅     | ✅      | ❌        | ✅         |

## GitHub Action token (default)

GitHub Action generates a token[^automatic-token-authentication] for you.

You can use it like so:

```yaml
jobs:
  ln:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
    steps:
      - uses: nobe4/action-ln@v0
```

> [!NOTE]
> The `permissions` set is required, as well as allowing GitHub action to
> [create and approve Pull
> Requests](#allowing-github-action-to-create-pull-requests)

> [!NOTE]
> GitHub action's provided token don't trigger CI checks[^automatic token ci
> checks trigger].
> You need to use another authentication method if you need that.

## Custom GitHub token

You can use a custom token if you require more permissions, want to act on
behalf of a user, etc.

The permissions needed for the token are:

- Classic token: `repo`

  The owner of the token also needs at least `write` access to all the linked
  repositories.

- Fine-grained token: `TBD`

Store it in your repository's secret and use it like so:

```yaml
jobs:
  ln:
    runs-on: ubuntu-latest
    steps:
      - uses: nobe4/action-ln@v0
        with:
          token: ${{ secret.CUSTOM_GITHUB_TOKEN }}
```

## GitHub application installation

> [!IMPORTANT]
> This section needs some love, I'm new in this.

```yaml
jobs:
  ln:
    runs-on: ubuntu-latest
    steps:
      - uses: nobe4/action-ln@v0
        with:
          app-id: ${{ secrets.ACTION_LN_APP_ID }}
          app-private-key: ${{ secrets.ACTION_LN_APP_PRIVATE_KEY }}
          app-install-id: ${{ secrets.ACTION_LN_APP_INSTALL_ID }}
```

> [!NOTE]
> The application installation cannot work across organizations. This is a know
> limitation; use a custom token if you need that.

## Allowing GitHub Action to create pull requests

For `action-ln` to create Pull Requests automatically, you need to authorized
GitHub Action to do so. It can be done at the organization or repository level:

Go to the applicable page and check `Allow GitHub Actions to create and approve
pull requests`:
- `https://github.com/organizations/<org>/settings/actions`
- `https://github.com/<owner>/<repo>/settings/actions`


[^automatic token authentication]: https://docs.github.com/en/actions/security-for-github-actions/security-guides/automatic-token-authentication
[^automatic token ci checks trigger]: https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/triggering-a-workflow#triggering-a-workflow-from-a-workflow
