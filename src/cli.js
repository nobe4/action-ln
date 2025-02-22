/*
Entrypoint to run action-ln from the local environment.

Because it doesn't actually come from any GitHub job workflow, most of the data
will be mocked.

The configuration is loaded from a local file.
*/

import * as core from "@actions/core";
import { parseArgs } from "node:util";
import { readFileSync } from "node:fs";

// Required by @action/github, imported by ./main
// TODO: find a way to set the context only once, and mock it from locally.
// Should be able to move it to its own file, in a global var, and from here
// fix its value.
// FIXME: this doesn't anymore with ESM, since imports are hoisted to the top of
// the file.
// cc https://stackoverflow.com/a/51730422
// Instead, it's setup in package.json, and the value is read from there.
// process.env.GITHUB_REPOSITORY = "nobe4/action-ln";

import { main } from "./main.js";
import { dedent } from "./format.js";
const help = () => {
	console.log(
		dedent(`
			usage: cli.js [FLAGS]

			flags:
				--help|-h               show this help
				--config|-c            path to config file (default: config.yaml)
				--token                GitHub PAT to use for auth
				--app_id               GitHub Application ID to use for auth
				--app_private_key_file path to the private key file for the GitHub Application
				--app_install_id       GitHub Application installation ID
				--noop                 don't actually run the action, just print what it would do
		`),
	);
	return;
};

const options = {
	noop: { type: "boolean" },
	help: { type: "boolean" },

	config: {
		type: "string",
		default: "config.yaml",
	},

	// Auth with a GitHub PAT
	token: {
		type: "string",
		default: process.env.GITHUB_TOKEN,
	},

	// Auth with a GitHub Application
	app_id: { type: "string" },
	app_private_key_file: { type: "string" },
	app_install_id: { type: "string" },
};

try {
	const { values } = parseArgs({ options: options });

	if (values.help) {
		return help();
	}

	if (values.app_private_key_file) {
		values.app_private_key_file = readFileSync(
			values.app_private_key_file,
		).toString();
	}

	main({
		configConfig: {
			useFS: true,
			path: values.config,
		},
		githubConfig: {
			token: values.token,
			appId: values.app_id,
			appPrivKey: values.app_private_key_file,
			appInstallId: values.app_install_id,
		},
		noop: values.noop,
	});
} catch (e) {
	core.error(e);
	core.error(e.stack);
	core.setFailed(e.message);
	help();
}
