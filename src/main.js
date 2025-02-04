const core = require("@actions/core");
const github = require("@actions/github");
const config = require("./config");
const { fetchAll } = require("./loader");

try {
	const configPath = core.getInput("config-path", { required: true });
	let token = core.getInput("token", { required: true });

	const octokit = github.getOctokit(token);

	config
		.load(configPath)
		.then((c) => {
			core.info(`config: ${JSON.stringify(c, null, "  ")}`);
			return fetchAll(octokit, c);
		})
		.then((c) => {
			core.info(`config after fetch: ${JSON.stringify(c, null, "  ")}`);
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
