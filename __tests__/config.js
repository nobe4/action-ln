import { jest } from "@jest/globals";

import { github } from "../__fixtures__/@actions/github.js";
jest.unstable_mockModule("@actions/github", () => github);

import * as core from "../__fixtures__/@actions/core.js";
jest.unstable_mockModule("@actions/core", () => core);

const fs = { readFile: jest.fn() };
jest.unstable_mockModule("node:fs/promises", () => fs);

const yaml = { load: jest.fn() };
jest.unstable_mockModule("js-yaml", () => yaml);

import { Link } from "../__fixtures__/src/link.js";
jest.unstable_mockModule("../src/link.js", () => ({ Link: Link }));

const { Config, ParseError } = await import("../src/config.js");

import { dedent } from "../src/format.js";

const repo = { owner: "owner", repo: "repo" };

describe("Config", () => {
	let c = new Config({});

	beforeEach(() => {
		c = new Config(
			{
				repo: repo,
				path: "path",
				useFS: false,
			},
			{
				getDefaultBranch: jest.fn(),
				getContent: jest.fn(),
			},
		);
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
		it("formats correctly for GitHub", () => {
			c.useFS = false;
			c.sha = "sha";
			expect(c.URL).toEqual("https://github.com/owner/repo/blob/sha/path");
		});

		it("formats correctly for FS", () => {
			c.useFS = true;
			expect(c.URL).toEqual("file://path");
		});
	});

	describe("load", () => {
		let mocks = {};

		beforeEach(() => {
			c.data = undefined;

			mocks = {
				loadFromFS: jest.spyOn(Config.prototype, "loadFromFS"),
				loadFromGitHub: jest.spyOn(Config.prototype, "loadFromGitHub"),
				parse: jest.spyOn(Config.prototype, "parse"),
				getContents: jest.spyOn(Config.prototype, "getContents"),
				groupLinks: jest.spyOn(Config.prototype, "groupLinks"),
			};
		});

		describe("from FS", () => {
			beforeEach(() => {
				c.useFS = true;
			});
			it("loads", async () => {
				mocks.loadFromFS.mockResolvedValue("data");
				yaml.load.mockResolvedValue("yaml");
				mocks.parse.mockResolvedValue(c);
				mocks.getContents.mockResolvedValue(c);
				mocks.groupLinks.mockResolvedValue(c);

				await expect(c.load()).resolves.toEqual(c);

				expect(c.data).toEqual("yaml");
				expect(mocks.loadFromFS).toHaveBeenCalled();
				expect(yaml.load).toHaveBeenCalledWith("data");
				expect(mocks.parse).toHaveBeenCalled();
				expect(mocks.getContents).toHaveBeenCalled();
				expect(mocks.groupLinks).toHaveBeenCalled();
			});

			it("fails to load", async () => {
				mocks.loadFromFS.mockRejectedValue("error");

				await expect(c.load()).rejects.toEqual("error");

				expect(c.data).toEqual(undefined);
				expect(mocks.loadFromFS).toHaveBeenCalled();
				expect(yaml.load).not.toHaveBeenCalled();
				expect(mocks.parse).not.toHaveBeenCalled();
				expect(mocks.getContents).not.toHaveBeenCalled();
				expect(mocks.groupLinks).not.toHaveBeenCalled();
			});

			it("fails to parse YAML", async () => {
				mocks.loadFromFS.mockResolvedValue("data");
				yaml.load.mockRejectedValue("error");

				await expect(c.load()).rejects.toEqual("error");

				expect(c.data).toEqual(undefined);
				expect(mocks.loadFromFS).toHaveBeenCalled();
				expect(yaml.load).toHaveBeenCalledWith("data");
				expect(mocks.parse).not.toHaveBeenCalled();
				expect(mocks.getContents).not.toHaveBeenCalled();
				expect(mocks.groupLinks).not.toHaveBeenCalled();
			});

			it("fails to process the data", async () => {
				mocks.loadFromFS.mockResolvedValue("data");
				yaml.load.mockResolvedValue("yaml");
				mocks.parse.mockRejectedValue("error");
				mocks.getContents.mockResolvedValue(c);
				mocks.groupLinks.mockResolvedValue(c);

				await expect(c.load()).rejects.toEqual("error");

				expect(c.data).toEqual("yaml");
				expect(mocks.loadFromFS).toHaveBeenCalled();
				expect(yaml.load).toHaveBeenCalledWith("data");
				expect(mocks.parse).toHaveBeenCalled();
				expect(mocks.getContents).not.toHaveBeenCalled();
				expect(mocks.groupLinks).not.toHaveBeenCalled();
			});

			// Ignoring further failures, as they would just be mocking the call
			// to the config's function.
		});

		describe("from GitHub", () => {
			beforeEach(() => {
				c.useFS = false;
			});

			it("loads", async () => {
				mocks.loadFromGitHub.mockResolvedValue("data");
				yaml.load.mockResolvedValue("yaml");
				mocks.parse.mockResolvedValue(c);
				mocks.getContents.mockResolvedValue(c);
				mocks.groupLinks.mockResolvedValue(c);

				await expect(c.load()).resolves.toEqual(c);

				expect(c.data).toEqual("yaml");
				expect(mocks.loadFromGitHub).toHaveBeenCalled();
				expect(yaml.load).toHaveBeenCalledWith("data");
				expect(mocks.parse).toHaveBeenCalled();
				expect(mocks.getContents).toHaveBeenCalled();
				expect(mocks.groupLinks).toHaveBeenCalled();
			});

			it("fails to load", async () => {
				mocks.loadFromGitHub.mockRejectedValue("error");

				await expect(c.load()).rejects.toEqual("error");

				expect(c.data).toEqual(undefined);
				expect(mocks.loadFromGitHub).toHaveBeenCalled();
				expect(yaml.load).not.toHaveBeenCalled();
				expect(mocks.parse).not.toHaveBeenCalled();
				expect(mocks.getContents).not.toHaveBeenCalled();
				expect(mocks.groupLinks).not.toHaveBeenCalled();
			});

			// Ignoring further failures as they are in the 'from FS' block as
			// well;
		});
	});

	describe("loadFromFS", () => {
		it("loads", async () => {
			const mockReadFile = jest.spyOn(fs, "readFile").mockResolvedValue("data");

			await expect(c.loadFromFS()).resolves.toEqual("data");

			expect(mockReadFile).toHaveBeenCalledWith(c.path, { encoding: "utf-8" });
			expect(c.sha).toEqual("runninglocally123");
		});

		it("fails to load", async () => {
			const mockReadFile = jest
				.spyOn(fs, "readFile")
				.mockRejectedValue("error");

			await expect(c.loadFromFS()).rejects.toEqual("error");

			expect(mockReadFile).toHaveBeenCalledWith(c.path, { encoding: "utf-8" });
			expect(c.sha).toEqual("runninglocally123");
		});
	});

	describe("loadFromGitHub", () => {
		it("loads", async () => {
			c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });
			c.gh.getContent.mockResolvedValue({ content: "data" });

			await expect(c.loadFromGitHub()).resolves.toEqual("data");
			expect(c.gh.getDefaultBranch).toHaveBeenCalledWith(c.repo);
			expect(c.gh.getContent).toHaveBeenCalledWith(c.repo, c.path);
			expect(c.sha).toEqual("sha");
		});

		it("fails to get the default branch", async () => {
			c.gh.getDefaultBranch.mockRejectedValue("error");

			await expect(c.loadFromGitHub()).rejects.toEqual("error");
			expect(c.gh.getDefaultBranch).toHaveBeenCalledWith(c.repo);
			expect(c.gh.getContent).not.toHaveBeenCalledWith(c.repo, c.path);
			expect(c.sha).toEqual(undefined);
		});

		it("fails to get the content", async () => {
			c.gh.getDefaultBranch.mockResolvedValue({ sha: "sha" });
			c.gh.getContent.mockRejectedValue("error");

			await expect(c.loadFromGitHub()).rejects.toEqual("error");
			expect(c.gh.getDefaultBranch).toHaveBeenCalledWith(c.repo);
			expect(c.gh.getContent).toHaveBeenCalledWith(c.repo, c.path);
			expect(c.sha).toEqual("sha");
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
				Link.parse.mockImplementation(() => "parsed");

				c.data = data;

				expect(c.parse().data).toStrictEqual(want);
				expect(Link.parse).toHaveBeenCalledTimes(data.links.length);
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
