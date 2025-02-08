const currentRepo = { owner: "owner", repo: "repo" };
jest.mock("@actions/github", () => ({ context: { repo: currentRepo } }));

const { File, ParseError } = require("../src/file");

describe("File", () => {
	let l = new File();

	describe("toString", () => {
		it("formats correctly", () => {
			l.repo = { repo: "repo", owner: "owner" };
			l.path = "path";
			expect(l.toString()).toStrictEqual("owner/repo:path");
		});
		it("formats correctly with a content", () => {
			l.repo = { repo: "repo", owner: "owner" };
			l.path = "path";
			l.content = "some\ncontent";
			expect(l.toString()).toStrictEqual("owner/repo:path\nsome\ncontent");
		});
	});

	describe("parse", () => {
		describe("fails", () => {
			it.each([null, undefined, ""])("%# %p", (raw) => {
				return expect(() => l.parse(raw)).toThrow(ParseError);
			});
		});

		describe("succeeds", () => {
			it.each(["non-nil"])("%# %p", (raw) => {
				const mockParsePath = jest
					.spyOn(File.prototype, "parsePath")
					.mockImplementation(() => {});

				const mockParseRepo = jest
					.spyOn(File.prototype, "parseRepo")
					.mockImplementation(() => {});

				expect(l.parse(raw)).toStrictEqual(l);
				expect(mockParsePath).toHaveBeenCalled();
				expect(mockParseRepo).toHaveBeenCalled();
			});
		});
	});

	describe("parsePath", () => {
		describe("fails", () => {
			it.each([null, undefined, "", "\n", "    ", " \t"])("%# %p", (raw) => {
				return expect(() => l.parsePath({ path: raw })).toThrow(ParseError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					raw: "a",
					want: "a",
				},
			])("%# %p", ({ raw, want }) => {
				return expect(l.parsePath({ path: raw })).toStrictEqual(want);
			});
		});
	});

	describe("parseRepo", () => {
		describe("fails", () => {
			it.each([
				{ repo: "a" },
				{ repo: "a/" },
				{ repo: "/a" },
				{ repo: "/" },
				{ repo: {} },
				{ repo: { owner: "" } },
				{ repo: { repo: "" } },
				{ repo: { repo: "", owner: undefined } },
				{ repo: { repo: undefined, owner: "" } },
			])("%# %p", (raw) => {
				return expect(() => l.parseRepo(raw)).toThrow(ParseError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					raw: undefined,
					want: currentRepo,
				},
				{
					raw: "",
					want: currentRepo,
				},
				{
					raw: "owner/repo",
					want: {
						repo: "repo",
						owner: "owner",
					},
				},
				{
					raw: {
						repo: "repo",
						owner: "owner",
					},
					want: {
						repo: "repo",
						owner: "owner",
					},
				},
			])("%# %p", ({ raw, want }) => {
				return expect(l.parseRepo({ repo: raw })).toStrictEqual(want);
			});
		});
	});
});
