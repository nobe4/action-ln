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

class Config {
	constructor(path) {
		this.path = path;
		this.data = {};
	}

	async load() {
		core.notice(`Using config file: ${this.path}`);

		return this.read().then(() => this.validate());
	}

	async read() {
		return fs.promises
			.readFile(this.path, "utf8")
			.then(yaml.load)
			.then((data) => (this.data = data));
	}

	validate() {
		if (this.data == null) {
			throw new ValidationError("config must not be null");
		}

		if (!("links" in this.data)) {
			throw new ValidationError("`links` must be present");
		}

		if (!Array.isArray(this.data.links)) {
			throw new ValidationError("`links` must be an array");
		}

		this.data.links.forEach((link, i) => {
			this.data.links[i] = this.validateLink(link);
		});

		return this.data;
	}

	validateLink(link) {
		if (typeof link !== "object") {
			throw new ValidationError("`links` must be an array of objects");
		}

		if (!("from" in link)) {
			throw new ValidationError("`from` must be present");
		}

		if (!("to" in link)) {
			throw new ValidationError("`to` must be present");
		}

		link.from = this.validateLocation(link.from);
		link.to = this.validateLocation(link.to);

		return link;
	}

	validateLocation(location) {
		if (!("path" in location) || location.path == "") {
			throw new ValidationError("`path` must be present");
		}

		if ("repo" in location) {
			const [owner, repo] = location.repo.split("/");
			if (
				owner === "" ||
				repo === "" ||
				owner === undefined ||
				repo === undefined
			) {
				throw new ValidationError("`repo` must be in the format `owner/name`");
			}
			location.repo = { owner, repo };
		} else {
			location.repo = github.context.repo;
		}

		return location;
	}
}

module.exports = { Config, ValidationError };
