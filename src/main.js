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

			const promises = [];
			for (let link of c.data.links) {
				core.debug(`link: ${link}`);

				if (!link.needsUpdate) {
					continue;
				}

				core.info(`updating: ${link.toString(true)}`);
				promises.push(gh.createPRForLink(link));
			}

			return Promise.all(promises);
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
