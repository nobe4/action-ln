/*
 Entrypoint to run action-ln from GitHub Actions.
*/

import * as core from "@actions/core";
import { context } from "@actions/github";

import { main } from "./main.js";

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
