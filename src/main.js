const core = require("@actions/core");
const { Config } = require("./config");
const { GitHub } = require("./github");

try {
	const configPath = core.getInput("config-path", { required: true });
	let token = core.getInput("token", { required: true });

	const gh = new GitHub(token);
	const config = new Config(configPath, gh);

	config
		.load()
		.then((c) => {
			core.info(`config:\n${c}`);
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
