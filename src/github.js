const core = require("@actions/core");
const { getOctokit } = require("@actions/github");
const { jsonError } = require("./utils");

class GitHub {
	constructor(token) {
		this.octokit = getOctokit(token, {
			log: console,
		});
	}

	normalizeBranch(branch) {
		return branch.replace(/[^a-zA-Z0-9]/g, "-");
	}

	createTree(path, content) {
		return [
			{
				path: path,
				content: content,
				mode: "100644", // TODO: this needs to be set by the `from`
				type: "blob",
			},
		];
	}

	async getContent({ owner, repo }, path) {
		core.debug(`fetching ${owner}/${repo}:${path}`);

		return this.octokit.rest.repos
			.getContent({
				owner: owner,
				repo: repo,
				path: path,
			})
			.then(({ data: { content, sha } }) => ({
				content: Buffer.from(content, "base64").toString("utf-8"),
				sha: sha,
			}))
			.then((c) => {
				core.debug(`fetched ${owner}/${repo}:${path}: ${JSON.stringify(c)}`);
				return c;
			})
			.catch((e) => {
				// This can fail if the file is missing, or if the repo is not
				// accessible. There's no way to differentiate that here.
				if (e.status === 404) {
					core.warning(`${owner}/${repo}:${path} not found`);
					return;
				}

				// However, any non-404 error is a real problem.
				core.setFailed(
					`failed to fetch ${owner}/${repo}:${path}: ${jsonError(e)}`,
				);

				throw e;
			});
	}

	async getDefaultBranch(repo) {
		let name = "";

		return this.getDefaultBranchName(repo)
			.then((n) => this.getBranch(repo, (name = n)))
			.then(({ object: { sha } } = {}) => ({
				name: name,
				sha: sha,
			}));
	}

	async getOrCreateBranch({ owner, repo }, name, sha) {
		return this.getBranch({ owner, repo }, name)
			.then((b) => ({ branch: b, new: false }))
			.catch(async (e) => {
				if (e.status === 404) {
					return this.createBranch({ owner, repo }, name, sha).then((b) => ({
						branch: b,
						new: true,
					}));
				}

				throw e;
			});
	}

	async getDefaultBranchName({ owner, repo }) {
		return this.octokit.rest.repos
			.get({ owner: owner, repo: repo })
			.then(({ data }) => data.default_branch);
	}

	async getBranch({ owner, repo }, name) {
		return this.octokit.rest.git
			.getRef({
				owner: owner,
				repo: repo,
				ref: `heads/${name}`,
			})
			.then(({ data }) => data);
	}

	async getCommit({ owner, repo }, sha) {
		return this.octokit.rest.git
			.getCommit({
				owner: owner,
				repo: repo,
				commit_sha: sha,
			})
			.then(({ data }) => data);
	}

	async createBranch({ owner, repo }, name, sha) {
		return this.octokit.rest.git
			.createRef({
				owner: owner,
				repo: repo,
				ref: `refs/heads/${name}`,
				sha: sha,
			})
			.then(({ data }) => data);
	}

	async createCommit({ owner, repo }, tree, parent) {
		return this.octokit.rest.git
			.createTree({
				owner: owner,
				repo: repo,
				tree: tree,
				base_tree: parent.tree.sha,
			})
			.then(({ data }) =>
				this.octokit.rest.git.createCommit({
					owner: owner,
					repo: repo,
					message: "[test] Update links",
					tree: data.sha,
					parents: [parent.sha],
				}),
			)
			.then(({ data }) => data);
	}

	async updateBranch({ owner, repo }, branch, sha) {
		return this.octokit.rest.git
			.updateRef({
				owner: owner,
				repo: repo,
				ref: `heads/${branch}`,
				sha: sha,
			})
			.then(({ data }) => data);
	}

	async createPullRequest({ owner, repo }, head, base, title, body) {
		return this.octokit.rest.pulls
			.create({
				owner: owner,
				repo: repo,
				head: head,
				base: base,
				title: title,
				body: body,
			})
			.then(({ data }) => data);
	}
}

module.exports = { GitHub };
