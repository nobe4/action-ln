const core = require("@actions/core");
const github = require("@actions/github");
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
	if (config == null) {
		throw new ValidationError("config must not be null");
	}

	if (!("links" in config)) {
		throw new ValidationError("`links` must be present");
	}

	if (!Array.isArray(config.links)) {
		throw new ValidationError("`links` must be an array");
	}

	config.links.forEach((link, i) => {
		config.links[i] = validateLink(link);
	});

	return config;
}

function validateLink(link) {
	if (typeof link !== "object") {
		throw new ValidationError("`links` must be an array of objects");
	}

	if (!("from" in link)) {
		throw new ValidationError("`from` must be present");
	}

	if (!("to" in link)) {
		throw new ValidationError("`to` must be present");
	}

	link.from = validateLocation(link.from);
	link.to = validateLocation(link.to);

	return link;
}

function validateLocation(location) {
	if (!("path" in location) || location.path == "") {
		throw new ValidationError("`path` must be present");
	}

	if ("repo" in location) {
		const [owner, name] = location.repo.split("/");
		console.log(owner, name);
		if (
			owner === "" ||
			name === "" ||
			owner === undefined ||
			name === undefined
		) {
			throw new ValidationError("`repo` must be in the format `owner/name`");
		}
		location.repo = { owner, name };
	} else {
		location.repo = github.context.repo;
	}

	return location;
}

module.exports = { read, load, validate, ValidationError };
