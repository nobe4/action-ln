// required by File, which is required by Link
jest.mock("@actions/github", () => ({ context: { repo: "repo" } }));

const core = require("@actions/core");
jest.mock("@actions/core");

const yaml = require("js-yaml");
jest.mock("js-yaml");

const { dedent } = require("../src/format");
const { Config, ParseError } = require("../src/config");

const { Link } = require("../src/link");
jest.mock("../src/link");

const repo = { owner: "owner", repo: "repo" };

describe("Config", () => {
	let c = new Config();

	beforeEach(() => {
		c.repo = repo;
		c.path = "path";
		c.sha = "sha";
		c.gh = {
			getDefaultBranch: jest.fn(),
			getContent: jest.fn(),
		};
	});

	describe("toString", () => {
		it("formats correctly", () => {
			const l1 = { toString: () => "l1" };
			const l2 = { toString: () => "l2" };

			c.data.links = [l1, l2];

			expect(c.toString()).toStrictEqual(
				dedent(
					`
					path: path
					links:
					  -
					    l1
					  -
					    l2
					`,
				),
			);
		});
	});

	describe("URL", () => {
		it("formats correctly", () => {
			c.path = "path";
			expect(c.URL).toEqual("https://github.com/owner/repo/blob/sha/path");
		});
	});

	describe("load", () => {
		const expectedcalls = () => {
			expect(core.notice).toHaveBeenCalledWith(
				"Using config file: owner/repo:path@sha",
			);
			expect(c.gh.getDefaultBranch).toHaveBeenCalledWith(c.repo);
			expect(c.gh.getContent).toHaveBeenCalledWith(c.repo, c.path);
		};

		describe("fails", () => {
			it("cannot read", async () => {
				c.gh.getContent.mockRejectedValue(new Error("ENOENT"));
				c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });

				await expect(c.load()).rejects.toThrow(/ENOENT/);

				expectedcalls();
			});

			it("cannot load YAML", async () => {
				c.gh.getContent.mockResolvedValue({ content: "content" });
				c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });
				yaml.load.mockRejectedValue(new Error("Invalid YAML"));

				await expect(c.load()).rejects.toThrow(/Invalid YAML/);

				expectedcalls();
			});

			it("cannot parse", async () => {
				c.gh.getContent.mockResolvedValue({ content: "content" });
				c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });
				yaml.load.mockResolvedValue("yaml");
				jest
					.spyOn(Config.prototype, "parse")
					.mockRejectedValue(new Error("Invalid config"));

				await expect(c.load()).rejects.toThrow(/Invalid config/);

				expectedcalls();
			});

			it("cannot getContents", async () => {
				c.gh.getContent.mockResolvedValue({ content: "content" });
				c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });
				yaml.load.mockResolvedValue("yaml");
				jest.spyOn(Config.prototype, "parse").mockResolvedValue("data");
				jest
					.spyOn(Config.prototype, "getContents")
					.mockRejectedValue(new Error("Error getting contents"));

				await expect(c.load()).rejects.toThrow(/Error getting contents/);

				expectedcalls();
			});
		});

		describe("succeeds", () => {
			it("read, load, parse, and getContents", async () => {
				c.gh.getContent.mockResolvedValue({ content: "content" });
				c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });
				yaml.load.mockResolvedValue("yaml");
				const mockParse = jest
					.spyOn(Config.prototype, "parse")
					.mockResolvedValue("data");
				const mockGetContents = jest
					.spyOn(Config.prototype, "getContents")
					.mockResolvedValue("data");
				const mockGroupLinks = jest
					.spyOn(Config.prototype, "groupLinks")
					.mockResolvedValue("data");

				await expect(c.load()).resolves.toEqual("data");
				expect(c.sha).toEqual("sha");

				expect(mockParse).toHaveBeenCalled();
				expect(mockGetContents).toHaveBeenCalled();
				expect(mockGroupLinks).toHaveBeenCalled();
				expectedcalls();
			});
		});
	});

	describe("getContents", () => {
		let files = [];

		beforeEach(() => {
			files = [
				{ path: "0", repo: "0" },
				{ path: "1", repo: "1" },
				{ path: "2", repo: "2" },
				{ path: "3", repo: "3" },
			];
			c.data = {
				links: [
					{ from: files[0], to: files[1] },
					{ from: files[0], to: files[2] },
					{ from: files[1], to: files[3] },
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
			it("fills all but one the links correctly", async () => {
				c.gh.getContent.mockImplementation((repo, path) => {
					return new Promise((resolve) => {
						if (path == "1") {
							resolve();
						}
						resolve({ content: repo + path, sha: 123 });
					});
				});

				await expect(c.getContents()).resolves.toEqual(c);
				expect(c.data.links[0].from).toEqual(files[0]);

				expect(files[1].content).not.toBeDefined();
				expect(files[1].sha).not.toBeDefined();

				files.forEach((f) => {
					expect(c.gh.getContent).toHaveBeenCalledWith(f.repo, f.path);
				});
			});

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
					data: { links: [] },
					want: {
						links: [],
					},
				},
				{
					data: { links: [0, 1, 2] },
					want: {
						links: ["parsed", "parsed", "parsed"],
					},
				},
			])("%# %j", ({ data, want }) => {
				const mockParse = jest
					.spyOn(Link.prototype, "parse")
					.mockImplementation(() => "parsed");

				c.data = data;

				expect(c.parse().data).toStrictEqual(want);
				expect(mockParse).toHaveBeenCalledTimes(data.links.length);
			});
		});
	});

	describe("groupLinks", () => {
		const links = [
			{ to: { repo: { owner: "o0", repo: "r0" } } },
			{ to: { repo: { owner: "o1", repo: "r1" } } },
			{ to: { repo: { owner: "o2", repo: "r2" } } },
		];
		it.each([
			{
				links: [],
				want: {},
			},
			{
				links: [links[0]],
				want: {
					"o0/r0": [links[0]],
				},
			},
			{
				links: [links[0], links[0], links[0]],
				want: {
					"o0/r0": [links[0], links[0], links[0]],
				},
			},
			{
				links: [links[0], links[1], links[0]],
				want: {
					"o0/r0": [links[0], links[0]],
					"o1/r1": [links[1]],
				},
			},
			{
				links: [links[0], links[1], links[0], links[1], links[2], links[1]],
				want: {
					"o0/r0": [links[0], links[0]],
					"o1/r1": [links[1], links[1], links[1]],
					"o2/r2": [links[2]],
				},
			},
		])("%# %j", ({ links, want }) => {
			c.data = { links: links };

			expect(c.groupLinks().data.groups).toStrictEqual(want);
		});
	});
});
