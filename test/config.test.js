const github = require("@actions/github");
jest.mock("@actions/github");

const fs = require("fs/promises");
jest.mock("fs/promises");

const yaml = require("js-yaml");
jest.mock("js-yaml");

const { Config, Link, File, ValidationError } = require("../src/config");
const { dedent } = require("../src/utils");

const { GitHub } = require("../src/github");
jest.mock("../src/github");

describe("Config", () => {
	let c = new Config();

	describe("toString", () => {
		it("formats correctly", () => {
			const l1 = new Link({
				from: new File({
					repo: { repo: "repo", owner: "owner" },
					path: "path",
					content: "content",
				}),
				to: new File({
					repo: { repo: "repo", owner: "owner" },
					path: "path",
					content: "content",
				}),
			});
			const l2 = new Link({
				from: new File({
					repo: { repo: "repo", owner: "owner" },
					path: "path",
					content: "content",
				}),
				to: new File({
					repo: { repo: "repo", owner: "owner" },
					path: "path",
					content: "other content",
				}),
			});

			c.data.links = [l1, l2];
			c.path = "path";

			expect(c.toString()).toStrictEqual(
				dedent(
					`
					path: path
					links:
					  -
					    from:
					        owner/repo:path
					        content
					    to:
					        owner/repo:path
					        content
					    needs update: false
					  -
					    from:
					        owner/repo:path
					        content
					    to:
					        owner/repo:path
					        other content
					    needs update: true
					`,
				).trim(),
			);
		});
	});

	describe("load", () => {
		describe("fails", () => {
			it("cannot read", () => {
				fs.readFile.mockRejectedValue(new Error("ENOENT"));
				return expect(c.load()).rejects.toThrow(/ENOENT/);
			});

			it("cannot load YAML", () => {
				fs.readFile.mockResolvedValue("content");
				yaml.load.mockRejectedValue(new Error("Invalid YAML"));
				return expect(c.load()).rejects.toThrow(/Invalid YAML/);
			});

			it("cannot parse", () => {
				fs.readFile.mockResolvedValue("content");
				yaml.load.mockResolvedValue("yaml");
				jest
					.spyOn(Config.prototype, "parse")
					.mockRejectedValue(new Error("Invalid config"));
				return expect(c.load()).rejects.toThrow(/Invalid config/);
			});

			it("cannot getContents", () => {
				fs.readFile.mockResolvedValue("content");
				yaml.load.mockResolvedValue("yaml");
				jest.spyOn(Config.prototype, "parse").mockResolvedValue("data");
				jest
					.spyOn(Config.prototype, "getContents")
					.mockRejectedValue(new Error("Error getting contents"));
				return expect(c.load()).rejects.toThrow(/Error getting contents/);
			});
		});

		describe("succeeds", () => {
			it("read, load, parse, and getContents", async () => {
				fs.readFile.mockResolvedValue("content");
				yaml.load.mockResolvedValue("yaml");
				const mockParse = jest
					.spyOn(Config.prototype, "parse")
					.mockResolvedValue("data");
				const mockGetContents = jest
					.spyOn(Config.prototype, "getContents")
					.mockResolvedValue("data");
				await expect(c.load()).resolves.toEqual("data");
				expect(mockParse).toHaveBeenCalled();
				expect(mockGetContents).toHaveBeenCalled();
			});
		});
	});

	describe("getContents", () => {
		const mockGithub = { getContents: jest.fn() };
		let files = [];

		beforeEach(() => {
			files = [
				new File({ content: 0 }),
				new File({ content: 1 }),
				new File({ content: 2 }),
				new File({ content: 3 }),
			];
			c.github = mockGithub;
			c.data = {
				links: [
					new Link({ from: files[0], to: files[1] }),
					new Link({ from: files[0], to: files[2] }),
					new Link({ from: files[1], to: files[3] }),
				],
			};
		});

		afterEach(() => {
			files.forEach((f) =>
				expect(mockGithub.getContents).toHaveBeenCalledWith(f),
			);
		});

		describe("fails", () => {
			it("getContents fails for one file", async () => {
				mockGithub.getContents.mockImplementation((file) => {
					return new Promise((resolve) => {
						if (file == files[1]) {
							throw new Error("Error getting contents");
						}
						resolve("content");
					});
				});

				await expect(() => c.getContents()).rejects.toThrow(
					/Error getting contents/,
				);
			});
		});

		describe("succeeds", () => {
			it("fills all the links correctly", async () => {
				mockGithub.getContents.mockImplementation((file) =>
					Promise.resolve(file.content),
				);

				await expect(c.getContents()).resolves.toEqual(c);
				files.forEach((f, i) => expect(f.content).toEqual(i));
			});
		});
	});

	describe("parse", () => {
		describe("fails", () => {
			it.each([
				null,
				{},
				{ links: null },
				{ links: "not a list" },
				{ links: 1 },
			])("%# %j", (data) => {
				c.data = data;
				expect(() => c.parse()).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					data: { links: [0, 1, 2] },
					want: {
						links: ["parsed", "parsed", "parsed"],
					},
				},
			])("%# %j", ({ data, want }) => {
				const mockLink = jest
					.spyOn(Link.prototype, "parse")
					.mockImplementation(() => "parsed");
				c.data = data;
				expect(c.parse().data).toStrictEqual(want);
				expect(mockLink).toHaveBeenCalled();
			});
		});
	});
});

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
					return expect(() => l.parse(raw)).toThrow(ValidationError);
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
				return expect(() => l.needsUpdate).toThrow(ValidationError);
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
				return expect(() => l.parse(raw)).toThrow(ValidationError);
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
				return expect(() => l.parsePath({ path: raw })).toThrow(
					ValidationError,
				);
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
				return expect(() => l.parseRepo(raw)).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			const defaultRepo = {
				owner: "owner",
				repo: "repo",
			};

			beforeEach(() => {
				github.context = { repo: defaultRepo };
			});

			it.each([
				{
					raw: undefined,
					want: defaultRepo,
				},
				{
					raw: "",
					want: defaultRepo,
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
