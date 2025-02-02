const core = require("@actions/core");
const yaml = require("js-yaml");
const fs = require("fs");

function Load(path) {
	return new Promise((resolve, reject) => {
		core.notice(`Using config file: ${path}`);

		fs.promises
			.readFile(path, "utf8")
			.then(yaml.load)
			.then(validate)
			.then(resolve)
			.catch(reject);
	});
}

function validate(config) {
	if (!("from" in config)) {
		throw "missing `from` in config";
	}

	return config;
}

export { Load };
