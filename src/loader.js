const core = require("@actions/core");

async function fetchAll(octokit, config) {
	const promises = [];

	for (let l of config.links) {
		promises.push(
			fetch(octokit, l.from).then((c) => {
				core.debug(`fetched ${l.from}: ${c}`);
				l.from.content = c;
			}),
		);
		promises.push(
			fetch(octokit, l.to).then((c) => {
				core.debug(`fetched ${l.to}: ${c}`);
				l.to.content = c;
			}),
		);
	}

	return Promise.all(promises);
}

async function fetch(octokit, { repo: { owner, repo }, path }) {
	core.debug(`fetching ${owner}/${repo}:${path}`);

	return octokit.rest.repos.getContent({
		owner: owner,
		repo: repo,
		path: path,
	});
}

module.exports = { fetchAll };
