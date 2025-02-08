const core = require("@actions/core");
const currentRepo = require("@actions/github").context.repo;
const yaml = require("js-yaml");
const { indent } = require("./utils");
const { Link } = require("./link");

class ParseError extends Error {
	constructor(message) {
		super(message);
		this.name = "Config ParseError";
	}
}

class Config {
	constructor(path, gh) {
		this.path = path;
		this.gh = gh;
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
		core.notice(
			`Using config file: ${currentRepo.owner}/${currentRepo.repo}:${this.path}`,
		);

		return this.gh
			.getContent(currentRepo, this.path)
			.then(({ content }) => yaml.load(content))
			.then((data) => (this.data = data))
			.then(() => this.parse())
			.then(() => this.getContents());
	}

	async getContents() {
		const promises = [];

		for (let i in this.data.links) {
			let link = this.data.links[i];

			promises.push(
				this.gh
					.getContent(link.from.repo, link.from.path)
					.then(({ content, sha }) => {
						this.data.links[i].from.content = content;
						this.data.links[i].from.sha = sha;
					}),
			);
			promises.push(
				this.gh
					.getContent(link.to.repo, link.to.path)
					.then(({ content, sha }) => {
						this.data.links[i].to.content = content;
						this.data.links[i].to.sha = sha;
					}),
			);
		}

		return Promise.all(promises).then(() => this);
	}

	parse() {
		if (!this.data) {
			throw new ParseError("config must not be null");
		}

		if (!("links" in this.data)) {
			throw new ParseError("`links` must be present");
		}

		if (!Array.isArray(this.data.links)) {
			throw new ParseError("`links` must be an array");
		}

		this.data.links.forEach((raw, i) => {
			this.data.links[i] = new Link().parse(raw);
		});

		return this;
	}
}

module.exports = { Config, ParseError };
