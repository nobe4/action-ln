const core = require("@actions/core");

async function fetchAll(octokit, config) {
	const promises = [];

	for (let l of config.links) {
		promises.push(
			fetch(octokit, l.from).then((c) => {
				core.debug(`fetched ${l.from.path}: ${c}`);
				l.from.content = c;
			}),
		);
		promises.push(
			fetch(octokit, l.to).then((c) => {
				core.debug(`fetched ${l.to.path}: ${c}`);
				l.to.content = c;
			}),
		);
	}

	return Promise.all(promises);
}

async function fetch(octokit, { repo: { owner, name }, path }) {
	core.debug(`fetching ${owner}/${name}:${path}`);

	octokit.rest.repos.getContent({
		owner: owner,
		repo: name,
		path: path,
	});
}

module.exports = { fetchAll };
