const core = require("@actions/core");
const github = require("@actions/github");
const config = require("./config");
const { fetchAll } = require("./loader");

try {
	const configPath = core.getInput("config-path", { required: true });
	let token = core.getInput("token", { required: false });

	if (token === "") {
		core.debug("Input `token` is empty, using default GITHUB_TOKEN.");
		token = github.token;
	}

	const octokit = github.getOctokit(token);

	config
		.load(configPath)
		.then((c) => {
			core.info(`config: ${JSON.stringify(c)}`);
			return fetchAll(octokit, c);
		})
		.then((c) => {
			core.info(`config after fetch: ${JSON.stringify(c)}`);
		})
		.catch((e) => {
			core.setFailed(e);
		});
} catch (error) {
	core.setFailed(error.message);
}
