import { hash } from "node:crypto";

import { indent } from "./format.js";
import { File } from "./file.js";

class ParseError extends Error {
	constructor(message) {
		super(message);
		this.name = "Link ParseError";
	}
}

class Link {
	constructor({ from, to } = {}) {
		this.from = from;
		this.to = to;
	}

	toString(short = false) {
		if (short) {
			return `${this.from.toString(true)} -> ${this.to.toString(true)}`;
		}

		return [
			"from:",
			indent(this.from.toString()),
			"to:",
			indent(this.to.toString()),
			`needs update: ${this.needsUpdate}`,
		].join("\n");
	}

	get SHA256() {
		return hash(
			"sha256",
			[
				this.from.repo.owner,
				this.from.repo.repo,
				this.from.path,
				this.to.repo.owner,
				this.to.repo.repo,
				this.to.path,
			].join(" "),
			"hex",
		);
	}

	parse(raw) {
		if (!raw || typeof raw !== "object") {
			throw new ParseError("`links` must be an array of objects");
		}

		if (!("from" in raw)) {
			throw new ParseError("`from` must be present");
		}

		if (!("to" in raw)) {
			throw new ParseError("`to` must be present");
		}

		this.from = new File().parse(raw.from);
		this.to = new File().parse(raw.to);

		return this;
	}

	get needsUpdate() {
		if (!this.from || !this.to) {
			throw new ParseError("`from` and `to` must be defined");
		}

		if (!this.from.content) {
			throw new ParseError("`from` must have a content");
		}

		if (!this.to.content) {
			return true;
		}

		return this.from.content !== this.to.content;
	}
}

export { Link, ParseError };
