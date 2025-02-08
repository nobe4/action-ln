const core = require("@actions/core");
const { getOctokit } = require("@actions/github");
const { prettify: _ } = require("./utils");

class GitHub {
	constructor(token) {
		this.octokit = getOctokit(token, {
			log: console,
		});
	}

	normalizeBranch(branch) {
		return branch.replace(/[^a-zA-Z0-9]/g, "-");
	}

	async getContent({ owner, repo }, path, ref = undefined) {
		const prettyRepo = `${owner}/${repo}:${path}@${ref}`;

		core.debug(`fetching ${prettyRepo}`);

		return this.octokit.rest.repos
			.getContent({
				owner: owner,
				repo: repo,
				path: path,
				ref: ref,
			})
			.then(({ data: { content, sha } }) => ({
				content: Buffer.from(content, "base64").toString("utf-8"),
				sha: sha,
			}))
			.then((c) => {
				core.debug(`fetched ${prettyRepo}: ${JSON.stringify(c)}`);
				return c;
			})
			.catch((e) => {
				// This can fail if the file is missing, or if the repo is not
				// accessible. There's no way to differentiate that here.
				if (e.status === 404) {
					core.warning(`${prettyRepo} not found`);
					return;
				}

				// However, any non-404 error is a real problem.
				core.setFailed(`failed to fetch ${prettyRepo}: ${_(e)}`);

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
		// Note: I wanted to change this code multiple times to just do
		// createBranch and catch the 422. It won't work since we need the new
		// SHA for the branch if it exists and was updated.
		return this.getBranch({ owner, repo }, name)
			.then((b) => {
				core.debug(`${owner}/${repo}@${name} found`);
				return { name: name, sha: b.object.sha, new: false };
			})
			.catch(async (e) => {
				if (e.status === 404) {
					core.debug(`${owner}/${repo}@${name} not found, creating branch`);
					return this.createBranch({ owner, repo }, name, sha).then((b) => ({
						name: name,
						sha: b.object.sha,
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

	async createOrUpdateFileContents(
		{ owner, repo },
		path,
		sha,
		branch,
		content,
		message,
	) {
		const base64Content = Buffer.from(content).toString("base64");

		return this.octokit.rest.repos
			.createOrUpdateFileContents({
				owner: owner,
				repo: repo,
				path: path,
				sha: sha,
				branch: branch,
				content: base64Content,
				message: message,
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
