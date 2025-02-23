import { jest } from "@jest/globals";

const parse = jest.fn();
class Link {
	static parse = parse;
	constructor() {
		return { parse: parse };
	}
}

export { Link };
