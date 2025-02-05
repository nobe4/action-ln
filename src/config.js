const core = require("@actions/core");
const github = require("@actions/github");
const yaml = require("js-yaml");
const fs = require("fs");
const { indent } = require("./utils");

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

	toString() {
		return [
			`path: ${this.path}`,
			`links:`,
			...this.data.links.map((l) => "  -\n" + indent(l.toString())),
		].join("\n");
	}

	async load() {
		core.notice(`Using config file: ${this.path}`);

		return this.read().then(() => this.parse());
	}

	async read() {
		return fs.promises
			.readFile(this.path, "utf8")
			.then(yaml.load)
			.then((data) => (this.data = data));
	}

	parse() {
		if (!this.data) {
			throw new ValidationError("config must not be null");
		}

		if (!("links" in this.data)) {
			throw new ValidationError("`links` must be present");
		}

		if (!Array.isArray(this.data.links)) {
			throw new ValidationError("`links` must be an array");
		}

		this.data.links.forEach((l, i) => {
			this.data.links[i] = new Link(l).parse();
		});

		return this;
	}
}

class Link {
	constructor(raw) {
		this.raw = raw;
	}

	toString() {
		return `from:\n${indent(this.from.toString())}\nto:\n${indent(this.to.toString())}`;
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

		this.from = new File(this.raw.from).parse();
		this.to = new File(this.raw.to).parse();

		delete this.raw;
		return this;
	}
}

class File {
	constructor(raw) {
		this.raw = raw;
	}

	toString() {
		let out = `${this.repo.owner}/${this.repo.repo}:${this.path}`;
		if (this.content) {
			out += `\n${this.content}`;
		}
		return out;
	}

	parse() {
		if (!this.raw) {
			throw new ValidationError("location must not be null");
		}

		this.parsePath();
		this.parseRepo();

		delete this.raw;
		return this;
	}

	parsePath() {
		if (!("path" in this.raw) || !this.raw.path) {
			throw new ValidationError("`path` must be present");
		}

		const path = this.raw.path.trim();
		if (!path) {
			throw new ValidationError("`path` must be not be empty");
		}

		return (this.path = path);
	}

	parseRepo() {
		if (!("repo" in this.raw) || !this.raw.repo) {
			return (this.repo = github.context.repo);
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

			return (this.repo = this.raw.repo);
		}

		const [owner, repo] = this.raw.repo.split("/");
		if (!owner || !repo) {
			throw new ValidationError("`repo` must be in the format `owner/repo`");
		}

		return (this.repo = { owner, repo });
	}
}

module.exports = { Config, Link, File, ValidationError };
