const core = require("@actions/core");

//async function link(octokit, l) {
//	if (!l.needUpdate) {
//		core.info(`No need to update ${l}`);
//		return;
//	}
//}

async function getBaseBranch(octokit, { owner, repo }) {
	return octokit.rest.repos
		.get({ owner: owner, repo: repo })
		.then(({ data }) => data.default_branch);
}

async function getBranch(octokit, { owner, repo }, branch) {
	return octokit.rest.git
		.getRef({
			owner: owner,
			repo: repo,
			ref: `heads/${branch}`,
		})
		.then(({ data }) => data);
}

async function getCommit(octokit, { owner, repo }, sha) {
	return octokit.rest.git
		.getCommit({
			owner: owner,
			repo: repo,
			commit_sha: sha,
		})
		.then(({ data }) => data);
}

async function createBranch(octokit, { owner, repo }, name, sha) {
	return octokit.rest.git
		.createRef({
			owner: owner,
			repo: repo,
			ref: `refs/heads/${name}`,
			sha: sha,
		})
		.then(({ data }) => data);
}

function createTreeHash(path, content) {
	return [
		{
			path: path,
			content: content,
			mode: "100644", // TODO: this needs to be set by the `from`
			type: "blob",
		},
	];
}

async function createCommit(octokit, { owner, repo }, tree, parent) {
	return octokit.rest.git
		.createTree({
			owner: owner,
			repo: repo,
			tree: tree,
			base_tree: parent.tree.sha,
		})
		.then(({ data }) =>
			octokit.rest.git.createCommit({
				owner: owner,
				repo: repo,
				message: "[test] Update links",
				tree: data.sha,
				parents: [parent.sha],
			}),
		)
		.then(({ data }) => data);
}

async function updateBranch(octokit, { owner, repo }, branch, sha) {
	return octokit.rest.git
		.updateRef({
			owner: owner,
			repo: repo,
			ref: `heads/${branch}`, // TODO: make  this consistent
			sha: sha,
		})
		.then(({ data }) => data);
}

async function createPullRequest(
	octokit,
	{ owner, repo },
	head,
	base,
	body,
	title,
) {
	return octokit.rest.pulls
		.create({
			owner: owner,
			repo: repo,
			head: head,
			base: base,
			body: body,
			title: title,
		})
		.then(({ data }) => data);
}

module.exports = {
	getBaseBranch,
	getBranch,
	getCommit,
	createCommit,
	createTreeHash,
	createBranch,
	createPullRequest,
	updateRef: updateBranch,
};
