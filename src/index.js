const core = require("@actions/core");
const context = require("@actions/github").context;
const { main } = require("./main");

try {
	const configPath = core.getInput("config-path", { required: true });
	let token = core.getInput("token", { required: true });
	let noop = core.getInput("noop", { required: true }) == "true";

	main({
		configConfig: {
			repo: context.repo,
			path: configPath,
		},
		token: token,
		noop: noop,
	});
} catch (e) {
	core.error(e);
	core.error(e.stack);
	core.setFailed(e.message);
}
