const core = require("@actions/core");
const github = require("@actions/github");
const { Config } = require("./config");
const { fetchAll } = require("./loader");

try {
	const configPath = core.getInput("config-path", { required: true });
	let token = core.getInput("token", { required: true });

	const octokit = github.getOctokit(token, {
		log: console,
	});

	const config = new Config(configPath);

	config
		.load()
		.then((c) => fetchAll(octokit, c))
		.then((c) => {
			core.info(`parsed and enriched config:\n${c}`);
		})
		.catch((e) => {
			core.error(e);
			core.error(e.stack);
			core.setFailed(e.message);
		});
} catch (e) {
	core.error(e);
	core.error(e.stack);
	core.setFailed(e.message);
}
