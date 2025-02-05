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

		this.data.links.forEach((l, i) => {
			this.data.links[i] = new Link(l);
		});

		return this.data;
	}
}

class Link {
	constructor(raw) {
		this.raw = raw;
		this.data = {};
	}

	parse() {
		if (!this.raw || typeof this.raw !== "object") {
			throw new ValidationError("`links` must be an array of objects");
		}

		if (!("from" in this.raw)) {
			throw new ValidationError("`from` must be present");
		}

		if (!("to" in this.raw)) {
			throw new ValidationError("`to` must be present");
		}

		this.data.from = new Location(this.raw.from).parse();
		this.data.to = new Location(this.raw.to).parse();

		return this.data;
	}
}

class Location {
	constructor(raw) {
		this.raw = raw;
		this.data = {};
	}

	parse() {
		if (!this.raw) {
			throw new ValidationError("location must not be null");
		}

		this.parsePath();
		this.parseRepo();

		return this.data;
	}

	parsePath() {
		if (!("path" in this.raw) || !this.raw.path) {
			throw new ValidationError("`path` must be present");
		}

		const path = this.raw.path.trim();
		if (!path) {
			throw new ValidationError("`path` must be not be empty");
		}

		return (this.data.path = path);
	}

	parseRepo() {
		if (!("repo" in this.raw) || !this.raw.repo) {
			return (this.data.repo = github.context.repo);
		}

		if (typeof this.raw.repo === "object") {
			if (
				!("repo" in this.raw.repo) ||
				!this.raw.repo.repo ||
				!("owner" in this.raw.repo) ||
				!this.raw.repo.owner
			) {
				throw new ValidationError("`repo` object must have `owner` and `repo`");
			}

			return (this.data.repo = this.raw.repo);
		}

		const [owner, repo] = this.raw.repo.split("/");
		if (!owner || !repo) {
			throw new ValidationError("`repo` must be in the format `owner/repo`");
		}

		return (this.data.repo = { owner, repo });
	}
}

module.exports = { Config, Link, Location, ValidationError };
