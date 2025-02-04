const core = require("@actions/core");

async function fetchAll(octokit, config) {
	const promises = [];

	for (let i in config.links) {
		promises.push(
			fetch(octokit, config.links[i].from).then((c) => {
				config.links[i].from.content = c;
			}),
		);
		promises.push(
			fetch(octokit, config.links[i].to).then((c) => {
				config.links[i].to.content = c;
			}),
		);
	}

	return Promise.all(promises).then(() => config);
}

async function fetch(octokit, { repo: { owner, repo }, path }) {
	core.debug(`fetching ${owner}/${repo}:${path}`);

	return octokit.rest.repos
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
			core.warning(
				`failed to fetch ${owner}/${repo}:${path}: ${JSON.stringify(e)}`,
			);
		});
}

module.exports = { fetchAll };
