const core = require("@actions/core");
const { Config } = require("./config");
const { GitHub } = require("./github");

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
	let headBranch = gh.normalizeBranch(
		`link-${link.from.repo.owner}-${link.from.repo.repo}-${link.from.path}`,
	);
	const newContent = gh.createTree(link.to.path, link.from.content);

	return gh
		.getDefaultBranch(link.to.repo)
		.then((b) => {
			baseBranch = b;
			return gh.createBranch(link.to.repo, headBranch, baseBranch.sha);
		})
		.then((b) => gh.getCommit(link.to.repo, b.object.sha))
		.then((c) => gh.createCommit(link.to.repo, newContent, c))
		.then((c) => gh.updateBranch(link.to.repo, headBranch, c.sha))
		.then(() =>
			gh.createPullRequest(
				link.to.repo,
				headBranch,
				baseBranch.name,
				"TODO",
				"TODO",
			),
		)
		.catch((e) => {
			core.setFailed(
				`failed to create PR for ${link.toString(true)}: ${e} ${JSON.stringify(e)}`,
			);
		});
}

main();
