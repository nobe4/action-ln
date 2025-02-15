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

					for (let link of c.data.groups[groupName]) {
						core.info(`link: ${link}`);

						if (!link.needsUpdate) {
							continue;
						}

						core.info(`updating: ${link.toString(true)}`);

						// This is only the first usage of noop, it should go deeper
						// into the creation of the branch and PRs.
						if (!noop) {
							promises.push(createPRForLink(gh, link, config));
						}
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

async function createPRForLink(gh, link, config) {
	let baseBranch = {};
	let headBranch = {
		needsUpdate: false,
	};

	return gh
		.getDefaultBranch(link.to.repo)

		.then((b) => {
			core.debug(`default branch: ${p(b)}`);
			baseBranch = b;
		})
		.then(() =>
			gh.getOrCreateBranch(link.to.repo, branchName(link), baseBranch.sha),
		)

		.then((b) => {
			core.debug(`head branch: ${p(b)}`);
			headBranch = b;
		})
		.then(() => {
			if (headBranch.new) {
				return (headBranch.needsUpdate = true);
			}

			return gh
				.getContent(link.to.repo, link.to.path, headBranch.name)
				.then((c) => (headBranch.needsUpdate = link.from.content !== c.content))
				.catch((e) => {
					if (e.status === 404) {
						return (headBranch.needsUpdate = true);
					}

					throw e;
				});
		})

		.then(() => {
			if (!headBranch.needsUpdate) {
				console.log("update not needed");
				return;
			}

			return gh.createOrUpdateFileContents(
				link.to.repo,
				link.to.path,
				link.to.sha,
				headBranch.name,
				link.from.content,
				commitMessage(link),
			);
		})

		.then(() =>
			gh.getOrCreatePullRequest(
				link.to.repo,
				headBranch.name,
				baseBranch.name,
				pullTitle(link),
				pullBody(link, config, context),
			),
		)

		.catch((e) => {
			core.setFailed(`failed to create PR for ${link.toString(true)}: ${p(e)}`);
		});
}

main();
