import { jest } from "@jest/globals";

import * as core from "../__mocks__/@actions/core.js";
jest.unstable_mockModule("@actions/core", () => core);

import { github } from "../__mocks__/@actions/github.js";
jest.unstable_mockModule("@actions/github", () => github);

const { GitHub } = await import("../src/github.js");

const repo = { owner: "owner", repo: "repo" };
const path = "path";
const branch = "branch";
const sha = "sha";
const content = "content";
const prettyRepo = `${repo.owner}/${repo.repo}:${path}@${branch}`;
const prettyBranch = `${repo.owner}/${repo.repo}@${branch}`;

describe("GitHub", () => {
	let g = undefined;
	let octokit = undefined;

	beforeEach(() => {
		octokit = {
			rest: {
				git: {
					createCommit: jest.fn(),
					createRef: jest.fn(),
					getRef: jest.fn(),
				},
				pulls: { create: jest.fn(), list: jest.fn() },
				repos: {
					getContent: jest.fn(),
					get: jest.fn(),
					createOrUpdateFileContents: jest.fn(),
				},
			},
		};

		g = new GitHub(octokit);
	});

	describe("constructor", () => {
		it("sets up the octokit client", () => {
			expect(octokit).toEqual(octokit);
		});
	});

	describe("getContent", () => {
		const expectedcalls = () => {
			expect(core.debug).toHaveBeenCalledWith(`fetching ${prettyRepo}`);

			expect(octokit.rest.repos.getContent).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
				ref: branch,
			});
		};

		it("catches a 404", async () => {
			octokit.rest.repos.getContent.mockRejectedValue({ status: 404 });

			await expect(g.getContent(repo, path, branch)).resolves.not.toBeDefined();

			expect(core.warning).toHaveBeenCalledWith(`${prettyRepo} not found`);
			expectedcalls();
		});

		it("fails to decode a base64 string", async () => {
			octokit.rest.repos.getContent.mockResolvedValue({
				data: { content: content },
			});
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					throw new Error("Error");
				}),
			};

			await expect(g.getContent(repo, path, branch)).rejects.toThrow(/Error/);

			expect(global.Buffer.from).toHaveBeenCalledWith(content, "base64");
			expect(core.setFailed).toHaveBeenCalledWith(
				expect.stringContaining(`failed to fetch ${prettyRepo}`),
			);
			expectedcalls();
		});

		it("succeeds", async () => {
			octokit.rest.repos.getContent.mockResolvedValue({
				data: { content: content, sha: 123 },
			});
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					return { toString: () => content };
				}),
			};

			await expect(g.getContent(repo, path, branch)).resolves.toEqual({
				content: content,
				sha: 123,
			});

			expect(global.Buffer.from).toHaveBeenCalledWith(content, "base64");
			expect(core.debug).toHaveBeenCalledWith(
				expect.stringContaining(`fetched ${prettyRepo}`),
			);
			expectedcalls();
		});

		it("succeeds without a ref", async () => {
			const prettyRepo = `${repo.owner}/${repo.repo}:${path}@undefined`;
			octokit.rest.repos.getContent.mockResolvedValue({
				data: { content: content, sha: 123 },
			});
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					return { toString: () => content };
				}),
			};

			await expect(g.getContent(repo, path)).resolves.toEqual({
				content: content,
				sha: 123,
			});

			expect(global.Buffer.from).toHaveBeenCalledWith(content, "base64");
			expect(core.debug).toHaveBeenCalledWith(
				expect.stringContaining(`fetched ${prettyRepo}`),
			);
			expect(core.debug).toHaveBeenCalledWith(`fetching ${prettyRepo}`);

			expect(octokit.rest.repos.getContent).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
			});
		});
	});

	describe("getDefaultBranch", () => {
		beforeEach(() => {
			g.getDefaultBranchName = jest.fn().mockResolvedValue("main");
		});

		const expectedcalls = () => {
			expect(g.getDefaultBranchName).toHaveBeenCalledWith(repo);
			expect(g.getBranch).toHaveBeenCalledWith(repo, "main");
		};

		it("fetches the default branch", async () => {
			g.getBranch = jest.fn().mockResolvedValue({ object: { sha: 123 } });

			await expect(g.getDefaultBranch(repo)).resolves.toEqual({
				name: "main",
				sha: 123,
			});

			expectedcalls();
		});

		it("defaults to nothing", async () => {
			g.getBranch = jest.fn().mockResolvedValue(undefined);

			await expect(g.getDefaultBranch(repo)).resolves.toEqual({
				name: "main",
				sha: undefined,
			});

			expectedcalls();
		});
	});

	describe("getOrCreateBranch", () => {
		it("gets an existing branch", async () => {
			g.getBranch = jest
				.fn()
				.mockResolvedValue({ object: { sha: "sha_new_branch" } });
			await expect(g.getOrCreateBranch(repo, branch, sha)).resolves.toEqual({
				name: branch,
				sha: "sha_new_branch",
				new: false,
			});
			expect(g.getBranch).toHaveBeenCalledWith(repo, branch);
		});

		it("creates a new branch", async () => {
			g.getBranch = jest.fn().mockRejectedValue({ status: 404 });
			g.createBranch = jest
				.fn()
				.mockResolvedValue({ object: { sha: "sha_new_branch" } });
			await expect(g.getOrCreateBranch(repo, branch, sha)).resolves.toEqual({
				name: branch,
				sha: "sha_new_branch",
				new: true,
			});
			expect(g.getBranch).toHaveBeenCalledWith(repo, branch);
			expect(g.createBranch).toHaveBeenCalledWith(repo, branch, sha);
		});

		it("fails on non-404", async () => {
			g.getBranch = jest.fn().mockRejectedValue(new Error("Error"));
			await expect(() =>
				g.getOrCreateBranch(repo, branch, sha),
			).rejects.toThrow(/Error/);
			expect(g.getBranch).toHaveBeenCalledWith(repo, branch);
		});
	});

	// All tests after this one are just checking for proper calling. Since they
	// do very little more than calling octokit and returning the data.
	describe("getDefaultBranchName", () => {
		it("fetches the default branch", async () => {
			octokit.rest.repos.get.mockResolvedValue({
				data: { default_branch: "main" },
			});
			await expect(g.getDefaultBranchName(repo)).resolves.toEqual("main");
			expect(octokit.rest.repos.get).toHaveBeenCalledWith(repo);
		});
	});

	describe("getBranch", () => {
		it("fetches the branch", async () => {
			octokit.rest.git.getRef.mockResolvedValue({ data: branch });
			await expect(g.getBranch(repo, branch)).resolves.toEqual(branch);
			expect(octokit.rest.git.getRef).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				ref: `heads/${branch}`,
			});
		});
	});

	describe("createBranch", () => {
		it("creates the branch", async () => {
			octokit.rest.git.createRef.mockResolvedValue({ data: branch });
			await expect(g.createBranch(repo, branch, sha)).resolves.toEqual(branch);
			expect(octokit.rest.git.createRef).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				ref: `refs/heads/${branch}`,
				sha: sha,
			});
		});
	});

	describe("createOrUpdateFileContents", () => {
		beforeEach(() => {
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					return { toString: () => "base64'ed content" };
				}),
			};
			octokit.rest.repos.createOrUpdateFileContents.mockResolvedValue({
				data: "ok",
			});
		});

		it("creates a new file", async () => {
			g.getContent = jest.fn().mockResolvedValue(undefined);

			await expect(
				g.createOrUpdateFileContents(repo, path, branch, content, "message"),
			).resolves.toEqual("ok");

			expect(global.Buffer.from).toHaveBeenCalledWith(content);
			expect(g.getContent).toHaveBeenCalledWith(repo, path, branch);
			expect(
				octokit.rest.repos.createOrUpdateFileContents,
			).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
				sha: undefined,
				branch: branch,
				content: "base64'ed content",
				message: "message",
			});
		});

		it("updates an existing file", async () => {
			g.getContent = jest.fn().mockResolvedValue({ sha: 123 });

			await expect(
				g.createOrUpdateFileContents(repo, path, branch, content, "message"),
			).resolves.toEqual("ok");

			expect(g.getContent).toHaveBeenCalledWith(repo, path, branch);
			expect(global.Buffer.from).toHaveBeenCalledWith(content);
			expect(
				octokit.rest.repos.createOrUpdateFileContents,
			).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
				sha: 123,
				branch: branch,
				content: "base64'ed content",
				message: "message",
			});
		});
	});

	describe("getOrCreatePullRequest", () => {
		const expectedcalls = () => {
			expect(octokit.rest.pulls.list).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				head: `${repo.owner}:${branch}`,
				per_page: 2,
			});
		};
		it("gets an existing pull request", async () => {
			octokit.rest.pulls.list.mockResolvedValue({ data: ["pull"] });

			await expect(
				g.getOrCreatePullRequest(repo, branch, "base", "title", "body"),
			).resolves.toEqual("pull");
			expectedcalls();
		});

		it("gets an existing pull request and warns for duplicates", async () => {
			octokit.rest.pulls.list.mockResolvedValue({ data: ["pull1", "pull2"] });

			await expect(
				g.getOrCreatePullRequest(repo, branch, "base", "title", "body"),
			).resolves.toEqual("pull1");

			expect(core.warning).toHaveBeenCalledWith(
				`found 2 PRs for ${prettyBranch}`,
			);
			expectedcalls();
		});

		it("creates a new pull request", async () => {
			octokit.rest.pulls.list.mockResolvedValue({ data: [] });
			octokit.rest.pulls.create.mockResolvedValue({ data: "pull" });

			await expect(
				g.getOrCreatePullRequest(repo, branch, "base", "title", "body"),
			).resolves.toEqual("pull");

			expectedcalls();
			expect(octokit.rest.pulls.create).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				head: branch,
				base: "base",
				title: "title",
				body: "body",
			});
		});
	});

	describe("compareFileContent", () => {
		const expectedcalls = () => {
			expect(GitHub.prototype.getContent).toHaveBeenCalledWith(
				repo,
				path,
				branch,
			);
		};

		it("doesn't find the file", async () => {
			jest.spyOn(GitHub.prototype, "getContent").mockResolvedValue(undefined);
			await expect(
				g.compareFileContent(repo, path, branch, content),
			).resolves.toEqual({
				found: false,
			});
			expectedcalls();
		});

		it("doesn't find the file with a 404", async () => {
			jest
				.spyOn(GitHub.prototype, "getContent")
				.mockRejectedValue({ status: 404 });
			await expect(
				g.compareFileContent(repo, path, branch, content),
			).resolves.toEqual({
				found: false,
			});
			expectedcalls();
		});

		it("find a different content", async () => {
			jest
				.spyOn(GitHub.prototype, "getContent")
				.mockResolvedValue({ content: "different" });
			await expect(
				g.compareFileContent(repo, path, branch, content),
			).resolves.toEqual({
				found: true,
				equal: false,
			});
			expectedcalls();
		});

		it("find the same content", async () => {
			jest
				.spyOn(GitHub.prototype, "getContent")
				.mockResolvedValue({ content: content });
			await expect(
				g.compareFileContent(repo, path, branch, content),
			).resolves.toEqual({
				found: true,
				equal: true,
			});
			expectedcalls();
		});

		it("throw other errors", async () => {
			const error = { status: 500 };
			jest.spyOn(GitHub.prototype, "getContent").mockRejectedValue(error);
			await expect(() =>
				g.compareFileContent(repo, path, branch, content),
			).rejects.toEqual(error);
			expectedcalls();
		});
	});
});
