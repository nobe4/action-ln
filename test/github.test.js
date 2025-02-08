const core = require("@actions/core");
jest.mock("@actions/core");

const { getOctokit } = require("@actions/github");
jest.mock("@actions/github");

const { GitHub } = require("../src/github");

const repo = { owner: "owner", repo: "repo" };
const path = "path";
const ref = "ref";
const prettyRepo = `${repo.owner}/${repo.repo}:${path}@${ref}`;

describe("GitHub", () => {
	let g = undefined;

	beforeEach(() => {
		g = new GitHub();
		g.octokit = {
			rest: {
				git: {
					createCommit: jest.fn(),
					createRef: jest.fn(),
					getRef: jest.fn(),
				},
				pulls: { create: jest.fn() },
				repos: { getContent: jest.fn(), get: jest.fn() },
			},
		};
	});

	describe("constructor", () => {
		it("sets up the octokit client", () => {
			const mockOctokit = getOctokit.mockReturnValue("ok");
			expect(new GitHub("token")).toBeDefined();
			expect(mockOctokit).toHaveBeenCalledWith("token", { log: console });
		});
	});

	describe("normalizeBranch", () => {
		it.each([
			{ input: "", output: "" },
			{ input: "abc123", output: "abc123" },
			{ input: " x y z ", output: "-x-y-z-" },
			{ input: "()|]\\ xxx {}", output: "------xxx---" },
		])("%# %s", ({ input, output }) => {
			expect(g.normalizeBranch(input)).toEqual(output);
		});
	});

	describe("createTree", () => {
		// Waiting for the `mode` to be testable.
		it.todo;
	});

	describe("getContent", () => {
		const expectedcalls = () => {
			expect(core.debug).toHaveBeenCalledWith(`fetching ${prettyRepo}`);

			expect(g.octokit.rest.repos.getContent).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
				ref: ref,
			});
		};

		it("catches a 404", async () => {
			g.octokit.rest.repos.getContent.mockRejectedValue({ status: 404 });

			await expect(g.getContent(repo, path, ref)).resolves.not.toBeDefined();

			expect(core.warning).toHaveBeenCalledWith(`${prettyRepo} not found`);
			expectedcalls();
		});

		it("fails to decode a base64 string", async () => {
			g.octokit.rest.repos.getContent.mockResolvedValue({
				data: { content: "content" },
			});
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					throw new Error("Error");
				}),
			};

			await expect(g.getContent(repo, path, ref)).rejects.toThrow(/Error/);

			expect(global.Buffer.from).toHaveBeenCalledWith("content", "base64");
			expect(core.setFailed).toHaveBeenCalledWith(
				expect.stringContaining(`failed to fetch ${prettyRepo}`),
			);
			expectedcalls();
		});

		it("succeeds", async () => {
			g.octokit.rest.repos.getContent.mockResolvedValue({
				data: { content: "content", sha: 123 },
			});
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					return { toString: () => "content" };
				}),
			};

			await expect(g.getContent(repo, path, ref)).resolves.toEqual({
				content: "content",
				sha: 123,
			});

			expect(global.Buffer.from).toHaveBeenCalledWith("content", "base64");
			expect(core.debug).toHaveBeenCalledWith(
				expect.stringContaining(`fetched ${prettyRepo}`),
			);
			expectedcalls();
		});

		it("succeeds without a ref", async () => {
			const prettyRepo = `${repo.owner}/${repo.repo}:${path}@undefined`;
			g.octokit.rest.repos.getContent.mockResolvedValue({
				data: { content: "content", sha: 123 },
			});
			global.Buffer = {
				from: jest.fn().mockImplementation(() => {
					return { toString: () => "content" };
				}),
			};

			await expect(g.getContent(repo, path)).resolves.toEqual({
				content: "content",
				sha: 123,
			});

			expect(global.Buffer.from).toHaveBeenCalledWith("content", "base64");
			expect(core.debug).toHaveBeenCalledWith(
				expect.stringContaining(`fetched ${prettyRepo}`),
			);
			expect(core.debug).toHaveBeenCalledWith(`fetching ${prettyRepo}`);

			expect(g.octokit.rest.repos.getContent).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
			});
		});
	});

	describe("getDefaultBranch", () => {
		it("fetches the default branch", async () => {
			g.getDefaultBranchName = jest.fn().mockResolvedValue("main");
			g.getBranch = jest.fn().mockResolvedValue({ object: { sha: 123 } });

			await expect(g.getDefaultBranch(repo)).resolves.toEqual({
				name: "main",
				sha: 123,
			});

			expect(g.getDefaultBranchName).toHaveBeenCalledWith(repo);
			expect(g.getBranch).toHaveBeenCalledWith(repo, "main");
		});
	});

	describe("getOrCreateBranch", () => {
		it("gets an existing branch", async () => {
			g.getBranch = jest
				.fn()
				.mockResolvedValue({ object: { sha: "sha_new_branch" } });
			await expect(g.getOrCreateBranch(repo, "branch", "sha")).resolves.toEqual(
				{
					name: "branch",
					sha: "sha_new_branch",
					new: false,
				},
			);
			expect(g.getBranch).toHaveBeenCalledWith(repo, "branch");
		});

		it("creates a new branch", async () => {
			g.getBranch = jest.fn().mockRejectedValue({ status: 404 });
			g.createBranch = jest
				.fn()
				.mockResolvedValue({ object: { sha: "sha_new_branch" } });
			await expect(g.getOrCreateBranch(repo, "branch", "sha")).resolves.toEqual(
				{
					name: "branch",
					sha: "sha_new_branch",
					new: true,
				},
			);
			expect(g.getBranch).toHaveBeenCalledWith(repo, "branch");
			expect(g.createBranch).toHaveBeenCalledWith(repo, "branch", "sha");
		});

		it("fails on non-404", async () => {
			g.getBranch = jest.fn().mockRejectedValue(new Error("Error"));
			await expect(() =>
				g.getOrCreateBranch(repo, "branch", "sha"),
			).rejects.toThrow(/Error/);
			expect(g.getBranch).toHaveBeenCalledWith(repo, "branch");
		});
	});

	// All tests after this one are just checking for proper calling. Since they
	// do very little more than calling octokit and returning the data.
	describe("getDefaultBranchName", () => {
		it("fetches the default branch", async () => {
			g.octokit.rest.repos.get.mockResolvedValue({
				data: { default_branch: "main" },
			});
			await expect(g.getDefaultBranchName(repo)).resolves.toEqual("main");
			expect(g.octokit.rest.repos.get).toHaveBeenCalledWith(repo);
		});
	});

	describe("getBranch", () => {
		it("fetches the branch", async () => {
			g.octokit.rest.git.getRef.mockResolvedValue({ data: "branch" });
			await expect(g.getBranch(repo, "branch")).resolves.toEqual("branch");
			expect(g.octokit.rest.git.getRef).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				ref: "heads/branch",
			});
		});
	});

	describe("createBranch", () => {
		it("creates the branch", async () => {
			g.octokit.rest.git.createRef.mockResolvedValue({ data: "branch" });
			await expect(g.createBranch(repo, "branch", "sha")).resolves.toEqual(
				"branch",
			);
			expect(g.octokit.rest.git.createRef).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				ref: "refs/heads/branch",
				sha: "sha",
			});
		});
	});

	describe("createPullRequest", () => {
		it("creates a pull request", async () => {
			g.octokit.rest.pulls.create.mockResolvedValue({ data: "pull" });
			await expect(
				g.createPullRequest(repo, "head", "base", "title", "body"),
			).resolves.toEqual("pull");
			expect(g.octokit.rest.pulls.create).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				head: "head",
				base: "base",
				title: "title",
				body: "body",
			});
		});
	});
});
