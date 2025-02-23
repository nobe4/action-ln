import { jest } from "@jest/globals";

import { github } from "../__fixtures__/@actions/github.js";
jest.unstable_mockModule("@actions/github", () => github);

const crypto = { hash: jest.fn() };
jest.unstable_mockModule("node:crypto", () => crypto);

const { Link, ParseError } = await import("../src/link.js");
const { File } = await import("../src/file.js");

import { dedent } from "../src/format.js";

describe("Link", () => {
	let l = new Link();

	describe("toString", () => {
		l.from = new File({
			repo: { repo: "repo", owner: "owner" },
			path: "path",
			content: "content",
			sha: 123,
		});
		l.to = new File({
			repo: { repo: "repo", owner: "owner" },
			path: "path",
			content: "content",
			sha: 123,
		});

		it("formats correctly", () => {
			expect(l.toString()).toStrictEqual(
				dedent(`
				from:
				    owner/repo:path@123
				    content
				to:
				    owner/repo:path@123
				    content
				needs update: false
				`),
			);
		});

		it("formats correctly in short format", () => {
			expect(l.toString(true)).toStrictEqual(
				"owner/repo:path@123 -> owner/repo:path@123",
			);
		});
	});

	describe("SHA256", () => {
		it("calculates the hash", () => {
			l.from = new File({
				repo: { repo: "repo", owner: "owner" },
				path: "path",
			});
			l.wto = new File({
				repo: { repo: "repo", owner: "owner" },
				path: "path",
			});
			crypto.hash.mockReturnValue("hash");

			expect(l.SHA256).toStrictEqual("hash");
			expect(crypto.hash).toHaveBeenCalledWith(
				"sha256",
				"owner repo path owner repo path",
				"hex",
			);
		});
	});

	describe("parse", () => {
		describe("fails", () => {
			it.each([undefined, "a", 1, {}, { from: {} }, { to: {} }])(
				"%# %p",
				(raw) => {
					return expect(() => l.parse(raw)).toThrow(ParseError);
				},
			);
		});

		describe("succeeds", () => {
			it.each([{ from: {}, to: {} }])("%# %p", (raw) => {
				const mockFileParse = jest
					.spyOn(File.prototype, "parse")
					.mockImplementation(() => "parsed");

				l.parse(raw);
				expect(l.from).toStrictEqual("parsed");
				expect(l.to).toStrictEqual("parsed");
				expect(mockFileParse).toHaveBeenCalled();
			});
		});
	});

	describe("needsUpdate", () => {
		describe("fails", () => {
			it.each([
				{ from: {} },
				{ to: {} },
				{ from: {}, to: {} },
				{ from: {}, to: { content: "a" } },
				{ from: { content: "" }, to: {} },
			])("%# %p", ({ from, to }) => {
				l.from = from;
				l.to = to;
				return expect(() => l.needsUpdate).toThrow(ParseError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{ from: { content: "a" }, to: {}, want: true },
				{ from: { content: "a" }, to: { content: "a" }, want: false },
				{ from: { content: "a" }, to: { content: "b" }, want: true },
			])("%# %p", ({ from, to, want }) => {
				l.from = from;
				l.to = to;
				return expect(l.needsUpdate).toBe(want);
			});
		});
	});
});
