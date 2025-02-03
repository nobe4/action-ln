const core = require("@actions/core");
const config = require("./config");

try {
	const configPath = core.getInput("config-path", { required: true });

	config
		.load(configPath)
		.then((c) => {
			core.info(`config: ${c}`);
		})
		.catch((e) => {
			core.setFailed(e);
		});
} catch (error) {
	core.setFailed(error.message);
}
