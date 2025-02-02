const core = require("@actions/core");
const config = require("./config");

try {
	const configPath = core.getInput("config-path", { required: true });
	const c = config.Load(configPath);
	console.log(c);
} catch (error) {
	core.setFailed(error.message);
}
