const currentRepo = require("@actions/github").context.repo;

class ParseError extends Error {
	constructor(message) {
		super(message);
		this.name = "File ParseError";
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
			throw new ParseError("location must not be null");
		}

		this.parsePath(raw);
		this.parseRepo(raw);

		return this;
	}

	parsePath(raw) {
		if (!("path" in raw) || !raw.path) {
			throw new ParseError("`path` must be present");
		}

		const path = raw.path.trim();
		if (!path) {
			throw new ParseError("`path` must be not be empty");
		}

		return (this.path = path);
	}

	parseRepo(raw) {
		if (!("repo" in raw) || !raw.repo) {
			return (this.repo = currentRepo);
		}

		if (typeof raw.repo === "object") {
			if (
				!("repo" in raw.repo) ||
				!raw.repo.repo ||
				!("owner" in raw.repo) ||
				!raw.repo.owner
			) {
				throw new ParseError("`repo` object must have `owner` and `repo`");
			}

			return (this.repo = raw.repo);
		}

		const [owner, repo] = raw.repo.split("/");
		if (!owner || !repo) {
			throw new ParseError("`repo` must be in the format `owner/repo`");
		}

		return (this.repo = { owner, repo });
	}
}

module.exports = { File, ParseError };
