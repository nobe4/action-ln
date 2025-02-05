const core = require("@actions/core");

async function fetchAll(octokit, config) {
	const promises = [];

	for (let i in config.data.links) {
		promises.push(
			fetch(octokit, config.data.links[i].data.from.data).then((c) => {
				config.data.links[i].data.from.data.content = c;
			}),
		);
		promises.push(
			fetch(octokit, config.data.links[i].data.to.data).then((c) => {
				config.data.links[i].data.to.data.content = c;
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
			if (e.status === 404) {
				core.warning(`${owner}/${repo}:${path} not found`);
				return;
			}

			// However, any non-404 error is a real problem.
			core.error(
				`failed to fetch ${owner}/${repo}:${path}: ${JSON.stringify(e)}`,
			);
		});
}

module.exports = { fetchAll };
