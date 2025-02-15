const core = require("@actions/core");
const context = require("@actions/github").context;

const { Config } = require("./config");
const { GitHub } = require("./github");
const {
	branchName,
	commitMessage,
	pullBody,
	pullTitle,
	prettify: p,
} = require("./format");

async function main() {
	const configPath = core.getInput("config-path", { required: true });
	let token = core.getInput("token", { required: true });
	let noop = core.getInput("noop", { required: true }) == "true";

	const gh = new GitHub(token);
	const config = new Config(context.repo, configPath, gh);
	await config.load();

	for (let i in config.data.groups) {
		core.info(`group ${i}`);
		const group = config.data.groups[i];

		let toRepo = group[0].to.repo;

		let baseBranch = await gh.getDefaultBranch(toRepo);
		core.info(`${toRepo.owner}/${toRepo.repo} base branch: ${p(baseBranch)}`);

		let headBranch = await gh.getOrCreateBranch(
			toRepo,
			branchName(),
			baseBranch.sha,
		);
		core.info(`${toRepo.owner}/${toRepo.repo} head branch: ${p(headBranch)}`);

		// Need to run those updates in sequence.
		for (let j in group) {
			const link = group[j];
			core.startGroup(link.toString(true));

			await checkIfLinkNeedsUpdate(link, gh, toRepo, headBranch).then(
				(needsUpdate) => {
					if (needsUpdate) {
						if (noop) {
							core.info("noop: would have created or updated file");
							return;
						}

						return gh.createOrUpdateFileContents(
							toRepo,
							link.to.path,
							headBranch.name,
							link.from.content,
							commitMessage(link),
						);
					}
				},
			);

			core.endGroup();
		}

		core.info(
			`creating PR for ${toRepo.owner}/${toRepo.repo}:${headBranch.name} -> ${baseBranch.name}`,
		);

		if (noop) {
			core.info("noop: would have created PR");
			return;
		}

		const pr = await gh.getOrCreatePullRequest(
			toRepo,
			headBranch.name,
			baseBranch.name,
			pullTitle(),
			pullBody(group, config, context),
		);
		core.info(`PR created: ${p(pr)}`);
	}
}

async function checkIfLinkNeedsUpdate(link, gh, toRepo, headBranch) {
	return new Promise((resolve) => {
		if (headBranch.new) {
			core.info("branch is new");
			resolve(true);
			return;
		}

		if (!link.needsUpdate) {
			core.info("link doesn't need update");
			resolve(false);
			return;
		}

		core.info("checking for diff for by getting content");

		gh.getContent(toRepo, link.to.path, headBranch.name)
			.then((c) => {
				if (!c) {
					core.info("file not found");
					resolve(true);
					return;
				}

				if (link.from.content !== c.content) {
					core.info("content is different");
					resolve(true);
					return;
				}
			})
			.catch((e) => {
				if (e.status !== 404) {
					reject(e);
					return;
				}

				core.info("file not found");
				resolve(true);
			});

		resolve(false);
	});
}

main().catch((e) => {
	core.error(e);
	core.error(e.stack);
	core.setFailed(e.message);
});
