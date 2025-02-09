const core = require("@actions/core");
const { Config } = require("./config");
const { GitHub } = require("./github");
const { prettify: p } = require("./utils");

function main() {
	try {
		const configPath = core.getInput("config-path", { required: true });
		let token = core.getInput("token", { required: true });

		const gh = new GitHub(token);
		const config = new Config(configPath, gh);

		config
			.load()
			.then((c) => {
				core.info(`config:\n${c}`);

				const promises = [];
				for (let link of c.data.links) {
					core.debug(`link: ${link}`);

					if (!link.needsUpdate) {
						continue;
					}

					core.info(`updating: ${link.toString(true)}`);
					promises.push(createPRForLink(gh, link));
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

async function createPRForLink(gh, link) {
	let baseBranch = {};
	let headBranchName = gh.normalizeBranch(
		`link-${link.from.repo.owner}-${link.from.repo.repo}-${link.from.path}`,
	);
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
			gh.getOrCreateBranch(link.to.repo, headBranchName, baseBranch.sha),
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
				"TODO commit message",
			);
		})

		.then(() =>
			gh.getOrCreatePullRequest(
				link.to.repo,
				headBranch,
				baseBranch.name,
				"TODO title",
				"TODO body",
			),
		)

		.catch((e) => {
			core.setFailed(`failed to create PR for ${link.toString(true)}: ${p(e)}`);
		});
}

main();

module.exports = { createPRForLink };
