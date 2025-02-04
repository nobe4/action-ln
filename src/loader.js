async function fetchAll(octokit, config) {
	const promises = [];

	for (let l of config.links) {
		promises.push(fetch(octokit, l.from).then((c) => (l.from.content = c)));
		promises.push(fetch(octokit, l.to).then((c) => (l.to.content = c)));
	}

	return Promise.all(promises);
}

async function fetch(octokit, { repo: { owner, name }, path }) {
	octokit.rest.repos.getContent({
		owner: owner,
		repo: name,
		path: path,
	});
}

module.exports = { fetchAll };
