const core = require("@actions/core");
const github = require("@actions/github");

class GitHub {
	constructor(token) {
		this.octokit = github.getOctokit(token, {
			log: console,
		});
	}

	// TODO: research what's the best interface for this method.
	async getContents({ repo: { owner, repo }, path }) {
		core.debug(`fetching ${owner}/${repo}:${path}`);

		return this.octokit.rest.repos
			.getContent({
				owner: owner,
				repo: repo,
				path: path,
			})
			.then(({ data: { content } }) =>
				Buffer.from(content, "base64").toString("utf-8"),
			)
			.catch((e) => {
				// This can fail if the file is missing, or if the repo is not
				// accessible. There's no way to differentiate that here.
				if (e.status === 404) {
					core.warning(`${owner}/${repo}:${path} not found`);
					return;
				}

				// However, any non-404 error is a real problem.
				core.setFailed(
					`failed to fetch ${owner}/${repo}:${path}: ${JSON.stringify(e)}`,
				);
			});
	}

	normalizeBranch(branch) {
		return branch.replace(/[^a-zA-Z0-9]/g, "-");
	}

	async getBaseBranch({ owner, repo }) {
		return this.octokit.rest.repos
			.get({ owner: owner, repo: repo })
			.then(({ data }) => data.default_branch);
	}

	async getBranch({ owner, repo }, branch) {
		return this.octokit.rest.git
			.getRef({
				owner: owner,
				repo: repo,
				ref: `heads/${branch}`,
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
				ref: `heads/${branch}`, // TODO: make  this consistent
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
