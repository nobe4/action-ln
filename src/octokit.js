import * as core from "@actions/core";

import { Octokit } from "@octokit/rest";
import { createAppAuth } from "@octokit/auth-app";
import { getOctokit } from "@actions/github";

function createOctokit({ token, appId, appPrivKey, appInstallId }) {
	if (!token && !appId && !appPrivKey) {
		throw new Error("either token or app_* should be provided");
	}

	if (appId && appPrivKey) {
		core.debug("creating octokit from application");

		// TODO: can't this be done with getOctokit?
		return new Octokit({
			authStrategy: createAppAuth,
			auth: {
				appId: appId,
				privateKey: appPrivKey,
				installationId: appInstallId,
			},
		});
	}

	core.debug("creating octokit from token");
	return getOctokit(token, {
		log: console,
	});
}

export { createOctokit };
