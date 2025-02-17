const core = require("@actions/core");
const { parseArgs } = require("node:util");

// required by @action/github, imported by ./main
// TODO: find a way to set the context only once, and mock it from locally.
// Should be able to move it to its own file, in a global var, and from here
// fix its value.
process.env.GITHUB_REPOSITORY = "nobe4/action-ln";

const { main } = require("./main");

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
