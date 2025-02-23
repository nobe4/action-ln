import * as core from "@actions/core";

import { prettify as _ } from "./format.js";

class GitHub {
	constructor(octokit) {
		this.octokit = octokit;
	}

	async getContent({ owner, repo }, path, ref = undefined) {
		const prettyFile = `${owner}/${repo}:${path}@${ref}`;

		core.debug(`fetching ${prettyFile}`);

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
				core.debug(`fetched ${prettyFile}: ${JSON.stringify(c)}`);
				return c;
			})
			.catch((e) => {
				// This can fail if the file is missing, or if the repo is not
				// accessible. There's no way to differentiate that here.
				if (e.status === 404) {
					core.warning(`${prettyFile} not found`);
					return;
				}

				// However, any non-404 error is a real problem.
				core.setFailed(`failed to fetch ${prettyFile}: ${_(e)}`);

				throw e;
			});
	}

	async getDefaultBranch(repo) {
		let name = "";

		return this.getDefaultBranchName(repo)
			.then((n) => this.getBranch(repo, (name = n)))
			.then(({ object: { sha } } = { object: { sha: undefined } }) => ({
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
		branch,
		content,
		message,
	) {
		return this.getContent({ owner: owner, repo: repo }, path, branch)
			.then((c) =>
				this.octokit.rest.repos.createOrUpdateFileContents({
					owner: owner,
					repo: repo,
					path: path,
					sha: (c || {}).sha,
					branch: branch,
					content: Buffer.from(content).toString("base64"),
					message: message,
				}),
			)
			.then(({ data }) => data);
	}

	async getOrCreatePullRequest({ owner, repo }, head, base, title, body) {
		// Notes: Similar to the createBranch function, we wanted to just do a
		// pulls.create, because it fails with 422 when it exists. But since it
		// doesn't return the PR info, we need to still to search for it.
		return this.octokit.rest.pulls
			.list({
				owner: owner,
				repo: repo,
				head: `${owner}:${head}`,

				// We'll use only the first item anyway, and we'll warn if more
				// than one PR exist for such a head.
				per_page: 2,
			})
			.then(({ data }) => {
				// There really should not be more than one Pull matching this
				// head, because creating more fails with 422.
				if (data.length > 1) {
					core.warning(`found ${data.length} PRs for ${owner}/${repo}@${head}`);
				}

				if (data.length > 0) {
					return data[0];
				}

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
			});
	}

	async compareFileContent(repo, path, branch, content) {
		return this.getContent(repo, path, branch)
			.then(({ content: c } = { undefined }) => {
				if (!c) {
					core.info("file not found");
					return { found: false };
				}

				if (content !== c) {
					core.info("content is different");
					return { found: true, equal: false };
				}

				return { found: true, equal: true };
			})
			.catch((e) => {
				if (e.status !== 404) {
					throw e;
				}

				core.info("file not found");
				return { found: false };
			});
	}
}

export { GitHub };
