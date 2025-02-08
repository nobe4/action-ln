const core = require("@actions/core");
jest.mock("@actions/core");

const { getOctokit } = require("@actions/github");
jest.mock("@actions/github");

const { GitHub } = require("../src/github");

const repo = { owner: "owner", repo: "repo" };
const path = "path";
const prettyRepo = `${repo.owner}/${repo.repo}:${path}`;

describe("GitHub", () => {
	let g = undefined;

	beforeEach(() => {
		g = new GitHub();
		g.octokit = {
			rest: {
				git: {
					createCommit: jest.fn(),
					createRef: jest.fn(),
					createTree: jest.fn(),
					getCommit: jest.fn(),
					getRef: jest.fn(),
					updateRef: jest.fn(),
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
			expect(core.debug).toHaveBeenCalledWith("fetching owner/repo:path");

			expect(g.octokit.rest.repos.getContent).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				path: path,
			});
		};

		it("catches a 404", async () => {
			g.octokit.rest.repos.getContent.mockRejectedValue({ status: 404 });
			await expect(g.getContent(repo, path)).resolves.not.toBeDefined();
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
			await expect(g.getContent(repo, path)).rejects.toThrow(/Error/);
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
			await expect(g.getContent(repo, path)).resolves.toEqual({
				content: "content",
				sha: 123,
			});
			expect(global.Buffer.from).toHaveBeenCalledWith("content", "base64");
			expect(core.debug).toHaveBeenCalledWith(
				expect.stringContaining("fetched owner/repo:path"),
			);
			expectedcalls();
		});
	});

	describe("getDefaultBranch", () => {
		it("fetches the default branch", async () => {
			g.getDefaultBranchName = jest.fn().mockResolvedValue("main");
			g.getBranch = jest.fn().mockResolvedValue("branch");

			await expect(g.getDefaultBranch(repo)).resolves.toEqual({
				name: "main",
				branch: "branch",
			});

			expect(g.getDefaultBranchName).toHaveBeenCalledWith(repo);
			expect(g.getBranch).toHaveBeenCalledWith(repo, "main");
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

	describe("getCommit", () => {
		it("fetches the commit", async () => {
			g.octokit.rest.git.getCommit.mockResolvedValue({ data: "commit" });
			await expect(g.getCommit(repo, "sha")).resolves.toEqual("commit");
			expect(g.octokit.rest.git.getCommit).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				commit_sha: "sha",
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

	describe("createCommit", () => {
		it("creates the commit", async () => {
			const parent = {
				sha: "sha_123",
				tree: {
					sha: "tree_123",
				},
			};
			const newTree = {
				sha: "tree_456",
			};

			g.octokit.rest.git.createTree.mockResolvedValue({
				data: newTree,
			});

			g.octokit.rest.git.createCommit.mockResolvedValue({ data: "commit" });

			await expect(g.createCommit(repo, "tree", parent)).resolves.toEqual(
				"commit",
			);

			expect(g.octokit.rest.git.createTree).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				tree: "tree",
				base_tree: parent.tree.sha,
			});

			expect(g.octokit.rest.git.createCommit).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				message: expect.any(String),
				tree: newTree.sha,
				parents: [parent.sha],
			});
		});
	});

	describe("updateBranch", () => {
		it("updates the branch", async () => {
			g.octokit.rest.git.updateRef.mockResolvedValue({ data: "branch" });
			await expect(g.updateBranch(repo, "branch", "sha")).resolves.toEqual(
				"branch",
			);
			expect(g.octokit.rest.git.updateRef).toHaveBeenCalledWith({
				owner: repo.owner,
				repo: repo.repo,
				ref: "heads/branch",
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
