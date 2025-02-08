// Needed for the import of File
const currentRepo = { owner: "owner", repo: "repo" };
jest.mock("@actions/github", () => ({ context: { repo: currentRepo } }));

const { Link, ParseError } = require("../src/link");
const { File } = require("../src/file");
const { dedent } = require("../src/utils");

describe("Link", () => {
	let l = new Link();

	describe("toString", () => {
		it("formats correctly", () => {
			l.from = new File();
			l.from.repo = { repo: "repo", owner: "owner" };
			l.from.path = "path";
			l.from.content = "content";
			l.to = new File();
			l.to.repo = { repo: "repo", owner: "owner" };
			l.to.path = "path";
			l.to.content = "content";
			expect(l.toString()).toStrictEqual(
				dedent(`
				from:
				    owner/repo:path
				    content
				to:
				    owner/repo:path
				    content
				needs update: false
				`).trim(),
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
