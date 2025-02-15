const core = require("@actions/core");
const yaml = require("js-yaml");
const { indent } = require("./format");
const { Link } = require("./link");

class ParseError extends Error {
	constructor(message) {
		super(message);
		this.name = "Config ParseError";
	}
}

class Config {
	constructor(repo, path, gh) {
		this.path = path;
		this.gh = gh;
		this.data = {};
		this.repo = repo;
		this.sha = undefined;
	}

	toString() {
		return [
			`path: ${this.path}`,
			`links:`,
			...this.data.links.map((l) => "  -\n" + indent(l.toString())),
		].join("\n");
	}

	get URL() {
		return `https://github.com/${this.repo.owner}/${this.repo.repo}/blob/${this.sha}/${this.path}`;
	}

	async load() {
		core.notice(
			`Using config file: ${this.repo.owner}/${this.repo.repo}:${this.path}@${this.sha}`,
		);

		return Promise.all([
			this.gh
				.getContent(this.repo, this.path)

				.then(({ content }) => yaml.load(content))
				.then((data) => (this.data = data)),

			this.gh.getDefaultBranch(this.repo).then(({ sha }) => (this.sha = sha)),
		])
			.then(() => this.parse())
			.then(() => this.getContents())
			.then(() => this.groupLinks());
	}

	async getContents() {
		const promises = [];

		for (let i in this.data.links) {
			let link = this.data.links[i];

			promises.push(
				this.gh
					.getContent(link.from.repo, link.from.path)
					.then(({ content, sha } = {}) => {
						this.data.links[i].from.content = content;
						this.data.links[i].from.sha = sha;
					}),
			);
			promises.push(
				this.gh
					.getContent(link.to.repo, link.to.path)
					.then(({ content, sha } = {}) => {
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

	groupLinks() {
		this.data.groups = {};

		for (let link of this.data.links) {
			if (link.to.repo in this.data.groups) {
				this.data.groups[link.to.repo].push(link);
			} else {
				this.data.groups[link.to.repo] = [link];
			}
		}

		return this;
	}
}

module.exports = { Config, ParseError };
