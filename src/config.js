const core = require("@actions/core");
const yaml = require("js-yaml");
const fs = require("fs");

class ValidationError extends Error {
	constructor(message) {
		super(message);
		this.name = "ValidationError";
	}
}

async function load(path) {
	core.notice(`Using config file: ${path}`);

	return read(path).then(validate);
}

async function read(path) {
	return fs.promises.readFile(path, "utf8").then(yaml.load);
}

async function validate(config) {
	return new Promise((resolve, reject) => {
		if (config == null) {
			reject(new ValidationError("config must not be null"));
		}

		if (!("links" in config)) {
			reject(new ValidationError("`links` must be present"));
		}

		if (!Array.isArray(config.links)) {
			reject(new ValidationError("`links` must be an array"));
		}

		config.links.forEach((link) => {
			if (typeof link !== "object") {
				reject(new ValidationError("`links` must be an array of objects"));
			}

			if (!("from" in link)) {
				reject(new ValidationError("`from` must be present"));
			}

			if (!("path" in link.from)) {
				reject(new ValidationError("`path` must be present in `from`"));
			}

			if (!("to" in link)) {
				reject(new ValidationError("`to` must be present"));
			}

			if (!("path" in link.to)) {
				reject(new ValidationError("`path` must be present in `to`"));
			}
		});

		resolve(config);
	});
}

module.exports = { read, load, validate, ValidationError };
