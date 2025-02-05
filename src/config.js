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

		this.data.links.forEach((raw, i) => {
			this.data.links[i] = new Link().parse(raw);
		});

		return this;
	}
}

class Link {
	constructor({ from, to } = {}) {
		this.from = from;
		this.to = to;
	}

	toString() {
		return [
			"from:",
			indent(this.from.toString()),
			"to:",
			indent(this.to.toString()),
			`needs update: ${this.needsUpdate}`,
		].join("\n");
	}

	parse(raw) {
		if (!raw || typeof raw !== "object") {
			throw new ValidationError("`links` must be an array of objects");
		}

		if (!("from" in raw)) {
			throw new ValidationError("`from` must be present");
		}

		if (!("to" in raw)) {
			throw new ValidationError("`to` must be present");
		}

		this.from = new File(raw.from).parse();
		this.to = new File(raw.to).parse();

		return this;
	}

	get needsUpdate() {
		if (!this.from || !this.to) {
			throw new ValidationError("`from` and `to` must be defined");
		}

		if (!this.from.content) {
			throw new ValidationError("`from` must have a content");
		}

		if (!this.to.content) {
			return true;
		}

		return this.from.content !== this.to.content;
	}
}

class File {
	constructor({ repo, path, content } = {}) {
		this.repo = repo;
		this.path = path;
		this.content = content;
	}

	toString() {
		let out = `${this.repo.owner}/${this.repo.repo}:${this.path}`;
		if (this.content) {
			out += `\n${this.content}`;
		}
		return out;
	}

	parse(raw) {
		if (!raw) {
			throw new ValidationError("location must not be null");
		}

		this.parsePath(raw);
		this.parseRepo(raw);

		return this;
	}

	parsePath(raw) {
		if (!("path" in raw) || !raw.path) {
			throw new ValidationError("`path` must be present");
		}

		const path = raw.path.trim();
		if (!path) {
			throw new ValidationError("`path` must be not be empty");
		}

		return (this.path = path);
	}

	parseRepo(raw) {
		if (!("repo" in raw) || !raw.repo) {
			return (this.repo = github.context.repo);
		}

		if (typeof raw.repo === "object") {
			if (
				!("repo" in raw.repo) ||
				!raw.repo.repo ||
				!("owner" in raw.repo) ||
				!raw.repo.owner
			) {
				throw new ValidationError("`repo` object must have `owner` and `repo`");
			}

			return (this.repo = raw.repo);
		}

		const [owner, repo] = raw.repo.split("/");
		if (!owner || !repo) {
			throw new ValidationError("`repo` must be in the format `owner/repo`");
		}

		return (this.repo = { owner, repo });
	}
}

module.exports = { Config, Link, File, ValidationError };
