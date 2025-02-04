const core = require("@actions/core");

async function fetchAll(octokit, config) {
	const promises = [];

	//for (let l of config.links) {
	//	promises.push(
	//		fetch(octokit, l.from).then((c) => {
	//			core.debug(`fetched ${JSON.stringify(l.from)}: ${JSON.stringify(c)}`);
	//			l.from.content = c;
	//			return l.from;
	//		}),
	//	);
	//	promises.push(
	//		fetch(octokit, l.to).then((c) => {
	//			core.debug(`fetched ${JSON.stringify(l.to)}: ${JSON.stringify(c)}`);
	//			l.to.content = c;
	//
	//			return l.to;
	//		}),
	//	);
	//}

	for (let i in config.links) {
		promises.push(
			fetch(octokit, config.links[i].from).then((c) => {
				core.debug(
					`fetched ${JSON.stringify(config.links[i].from)}: ${JSON.stringify(c)}`,
				);
				config.links[i].from.content = c;
			}),
		);
		promises.push(
			fetch(octokit, config.links[i].to).then((c) => {
				core.debug(
					`fetched ${JSON.stringify(config.links[i].to)}: ${JSON.stringify(c)}`,
				);
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
			core.error(`failed to fetch ${owner}/${repo}:${path}: ${e}`);
			return undefined;
		});
}

module.exports = { fetchAll };
