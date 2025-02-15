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
	let headBranch = {
		needsUpdate: false,
	};
	let toRepo = group[0].to.repo;

	return gh
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
			if (headBranch.new) {
				return (headBranch.needsUpdate = true);
			}

			const promises = [];
			for (let link of group) {
				if (!link.needsUpdate) {
					core.info("diff checking not needed for ${link.toString(true)}");
					continue;
				}

				core.info(`checking for diff on branch for ${link.toString(true)}`);

				promises.push(() => {
					return gh
						.getContent(toRepo, link.to.path, headBranch.name)
						.then(
							(c) => (headBranch.needsUpdate = link.from.content !== c.content),
						)
						.catch((e) => {
							if (e.status === 404) {
								return (headBranch.needsUpdate = true);
							}

							throw e;
						});
				});
			}

			return Promise.all(promises);
		})

		.then(() => {
			if (!headBranch.needsUpdate) {
				core.info(`update not needed for ${toRepo}:${headBranch.name}`);
				return;
			}

			const promises = [];
			for (let link of group) {
				if (!link.needsUpdate) {
					core.info(`update not needed for ${link.toString(true)}`);
					continue;
				}

				core.info(`updating: ${link.toString(true)}`);

				promises.push(() => {
					return gh.createOrUpdateFileContents(
						toRepo,
						link.to.path,
						link.to.sha,
						headBranch.name,
						link.from.content,
						commitMessage(link),
					);
				});
			}

			return Promise.all(promises);
		})

		.then(() =>
			gh.getOrCreatePullRequest(
				toRepo,
				headBranch.name,
				baseBranch.name,
				pullTitle(),
				pullBody(group, config, context),
			),
		)

		.catch((e) => {
			core.setFailed(`failed to create PR for ${group}: ${p(e)}`);
		});
}

main();
