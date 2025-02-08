const currentRepo = { owner: "owner", repo: "repo" };
jest.mock("@actions/github", () => ({ context: { repo: currentRepo } }));

const yaml = require("js-yaml");
jest.mock("js-yaml");

const { dedent } = require("../src/utils");

const { Config, ParseError } = require("../src/config");

const { Link } = require("../src/link");
const { File } = require("../src/file");

describe("Config", () => {
	let c = new Config();

	beforeEach(() => {
		c.gh = { getContent: jest.fn() };
	});

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
				),
			);
		});
	});

	describe("load", () => {
		describe("fails", () => {
			it("cannot read", () => {
				c.gh.getContent.mockRejectedValue(new Error("ENOENT"));
				return expect(c.load()).rejects.toThrow(/ENOENT/);
			});

			it("cannot load YAML", () => {
				c.gh.getContent.mockResolvedValue("content");
				yaml.load.mockRejectedValue(new Error("Invalid YAML"));
				return expect(c.load()).rejects.toThrow(/Invalid YAML/);
			});

			it("cannot parse", () => {
				c.gh.getContent.mockResolvedValue("content");
				yaml.load.mockResolvedValue("yaml");
				jest
					.spyOn(Config.prototype, "parse")
					.mockRejectedValue(new Error("Invalid config"));
				return expect(c.load()).rejects.toThrow(/Invalid config/);
			});

			it("cannot getContents", () => {
				c.gh.getContent.mockResolvedValue("content");
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
				c.gh.getContent.mockResolvedValue("content");
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
		let files = [];

		beforeEach(() => {
			files = [
				new File({ path: "0", repo: "0" }),
				new File({ path: "1", repo: "1" }),
				new File({ path: "2", repo: "2" }),
				new File({ path: "3", repo: "3" }),
			];
			c.data = {
				links: [
					new Link({ from: files[0], to: files[1] }),
					new Link({ from: files[0], to: files[2] }),
					new Link({ from: files[1], to: files[3] }),
				],
			};
		});

		describe("fails", () => {
			it("getContents fails for one file", async () => {
				c.gh.getContent.mockImplementation((repo, path) => {
					return new Promise((resolve) => {
						if (path == "1") {
							throw new Error("Error getting contents");
						}
						resolve({ content: repo + path, sha: 123 });
					});
				});

				await expect(() => c.getContents()).rejects.toThrow(
					/Error getting contents/,
				);
				files.forEach((f) =>
					expect(c.gh.getContent).toHaveBeenCalledWith(f.repo, f.path),
				);
			});
		});

		describe("succeeds", () => {
			it("fills all the links correctly", async () => {
				c.gh.getContent.mockImplementation((repo, path) =>
					Promise.resolve({ content: repo + path, sha: 123 }),
				);

				await expect(c.getContents()).resolves.toEqual(c);
				expect(c.data.links[0].from).toEqual(files[0]);

				files.forEach((f) => {
					expect(f.content).toEqual(f.repo + f.path);
					expect(f.sha).toEqual(123);
					expect(c.gh.getContent).toHaveBeenCalledWith(f.repo, f.path);
				});
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
				expect(() => c.parse()).toThrow(ParseError);
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
