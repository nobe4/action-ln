/*
 Entrypoint to run action-ln from GitHub Actions.
*/

import * as core from "@actions/core";
import { context } from "@actions/github";

import { main } from "./main.js";

try {
	const configPath = core.getInput("config-path", { required: true });
	const noop = core.getBooleanInput("noop", { required: true });
	const token = core.getInput("token", { required: true });
	const appId = core.getInput("app-id");
	const appPrivKey = core.getInput("app-private-key");
	const appInstallId = core.getInput("app-install-id");

	main({
		configConfig: {
			repo: context.repo,
			path: configPath,
		},
		auth: {
			token: token,
			appId: appId,
			appPrivKey: appPrivKey,
			appInstallId: appInstallId,
		},
		noop: noop,
	});
} catch (e) {
	core.error(e);
	core.error(e.stack);
	core.setFailed(e.message);
}
