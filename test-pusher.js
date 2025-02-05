const github = require("@actions/github");
const pusher = require("./src/pusher");

const octokit = github.getOctokit(process.env.GITHUB_TOKEN, {
	log: console,
});

const repo = { owner: "frozen-fishsticks", repo: "action-ln-test-0" };
const headBranch = "test";
let baseBranch = "";
const tree = pusher.createTreeHash("a.js", "console.log('hello')");

pusher
	.getBaseBranch(octokit, repo)
	.then((b) => (baseBranch = b))
	.then((b) => pusher.getBranch(octokit, repo, b))
	.then((b) => pusher.createBranch(octokit, repo, headBranch, b.object.sha))
	.then((b) => pusher.getCommit(octokit, repo, b.object.sha))
	.then((c) => pusher.createCommit(octokit, repo, tree, c))
	.then((c) => pusher.updateRef(octokit, repo, headBranch, c.sha))
	.then(() =>
		pusher.createPullRequest(
			octokit,
			repo,
			headBranch,
			baseBranch,
			"body test",
			"title test",
		),
	)
	.then(console.log);
