/*
Entrypoint to run action-ln from the local environment.

Because it doesn't actually come from any GitHub job workflow, most of the data
will be mocked.

The configuration is loaded from a local file.
*/

import * as core from "@actions/core";
import { parseArgs } from "node:util";

// Required by @action/github, imported by ./main
// TODO: find a way to set the context only once, and mock it from locally.
// Should be able to move it to its own file, in a global var, and from here
// fix its value.
process.env.GITHUB_REPOSITORY = "nobe4/action-ln";

import { main } from "./main.js";

try {
	// TODO: write some simple help if --help is passed
	const { values } = parseArgs({
		options: {
			config: {
				type: "string",
				default: "config.yaml",
			},
			token: {
				type: "string",
				default: process.env.GITHUB_TOKEN,
			},
			noop: {
				type: "boolean",
				default: false,
			},
		},
	});

	core.info(
		`Running action-ln with the following values: ${JSON.stringify(values, null, 2)}`,
	);

	main({
		configConfig: {
			useFS: true,
			path: values.config,
		},
		token: values.token,
		noop: values.noop,
	});
} catch (e) {
	core.error(e);
	core.error(e.stack);
	core.setFailed(e.message);
}
