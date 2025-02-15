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

function main() {
	try {
		const configPath = core.getInput("config-path", { required: true });
		let token = core.getInput("token", { required: true });
		let noop = core.getInput("noop", { required: true }) == "true";

		const gh = new GitHub(token);
		const config = new Config(context.repo, configPath, gh);

		config
			.load()
			.then((c) => {
				core.info(`config:\n${c}`);

				const promises = [];
				for (let groupName in c.data.groups) {
					core.info(`group: ${groupName}`);

					// This is only the first usage of noop, it should go deeper
					// into the creation of the branch and PRs.
					if (!noop) {
						promises.push(
							createPRForGroup(gh, c.data.groups[groupName], config),
						);
					}
				}

				return Promise.all(promises);
			})
			.catch((e) => {
				core.error(e);
				core.error(e.stack);
				core.setFailed(e.message);
			});
	} catch (e) {
		core.error(e);
		core.error(e.stack);
		core.setFailed(e.message);
	}
}

async function createPRForGroup(gh, group, config) {
	let baseBranch = {};
	let headBranch = {};

	let toRepo = group[0].to.repo;

	return (
		gh
			.getDefaultBranch(toRepo)

			.then((b) => {
				core.debug(`default branch: ${p(b)}`);
				baseBranch = b;
			})

			.then(() => gh.getOrCreateBranch(toRepo, branchName(), baseBranch.sha))

			.then((b) => {
				core.debug(`head branch: ${p(b)}`);
				headBranch = b;
			})

			.then(() => {
				const promises = group.map((link) => {
					return new Promise((resolve, reject) => {
						core.info(`checking for diff for ${link.toString(true)}`);

						if (headBranch.new) {
							core.info(
								`diff checking not needed for ${link.toString(true)}: head branch is new`,
							);
							resolve(false);
							return;
						}

						if (!link.needsUpdate) {
							core.info(
								`diff checking not needed for ${link.toString(true)}: links is up to date`,
							);
							resolve(false);
							return;
						}

						core.info(
							`checking for diff for ${link.toString(true)} by getting content`,
						);

						gh.getContent(toRepo, link.to.path, headBranch.name)
							.then((c) => {
								const needsUpdate = link.from.content !== c.content;

								core.info(
									`diff found for ${link.toString(true)}: ${needsUpdate}`,
								);

								resolve(needsUpdate);
							})
							.catch((e) => {
								if (e.status === 404) {
									core.info(`file not found ${link.toString(true)}`);
									resolve(true);
								}

								reject(e);
							});
					});
				});

				Promise.all(promises).then((x) => core.info(`update needs: ${x}`));

				//return Promise.all(promises);
			})

			//.then((updateNeeded) => {
			//	core.info(`update needs: ${updateNeeded}`);
			//
			//const promises = [];
			//
			//group.forEach((link, i) => {
			//	if (!updateNeeded[i]) {
			//		core.info(
			//			`update not needed for ${toRepo.owner}/${toRepo.repo}:${headBranch.name}: branch is up to date`,
			//		);
			//		return;
			//	}
			//
			//	core.info(`updating: ${link.toString(true)}`);
			//
			//	promises.push(
			//		(() => {
			//			return gh.createOrUpdateFileContents(
			//				toRepo,
			//				link.to.path,
			//				link.to.sha,
			//				headBranch.name,
			//				link.from.content,
			//				commitMessage(link),
			//			);
			//		})(),
			//	);
			//});
			//
			//return Promise.all(promises);
			//})

			//.then((values) => {
			//	core.info(`result from updating branches: ${values}`);
			//})
			//
			//.then(() =>
			//	gh.getOrCreatePullRequest(
			//		toRepo,
			//		headBranch.name,
			//		baseBranch.name,
			//		pullTitle(),
			//		pullBody(group, config, context),
			//	),
			//)
			//
			.catch((e) => {
				core.setFailed(`failed to create PR for ${group}: ${p(e)}`);
			})
	);
}

main();
