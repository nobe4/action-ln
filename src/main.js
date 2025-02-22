import * as core from "@actions/core";
import { context } from "@actions/github";

import { Config } from "./config.js";
import { GitHub } from "./github.js";

import {
	branchName,
	commitMessage,
	pullBody,
	pullTitle,
	prettify as p,
} from "./format.js";

async function main({ configConfig, token, noop }) {
	const gh = new GitHub(token);
	const config = new Config(configConfig, gh);
	await config.load();

	for (let i in config.data.groups) {
		core.info(`group ${i}`);
		const group = config.data.groups[i];

		let toRepo = group[0].to.repo;

		let baseBranch = await gh.getDefaultBranch(toRepo);
		core.info(`${toRepo.owner}/${toRepo.repo} base branch: ${p(baseBranch)}`);

		let headBranch = await gh.getOrCreateBranch(
			toRepo,
			branchName,
			baseBranch.sha,
		);
		core.info(`${toRepo.owner}/${toRepo.repo} head branch: ${p(headBranch)}`);

		// Need to run those updates in sequence.
		for (let j in group) {
			const link = group[j];
			core.startGroup(link.toString(true));

			await checkIfLinkNeedsUpdate(link, gh, toRepo, headBranch).then(
				(needsUpdate) => {
					core.info(`needs update: ${needsUpdate}`);

					if (needsUpdate) {
						if (noop) {
							core.info("noop: would have created or updated file");
							return;
						}

						core.info("creating or updating file");

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
			pullTitle,
			pullBody(group, config, context),
		);
		core.info(`PR created: ${p(pr)}`);
	}
}

async function checkIfLinkNeedsUpdate(link, gh, toRepo, headBranch) {
	return new Promise((resolve) => {
		core.info("checking if the link needs an update");

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

		gh.compareFileContent(
			toRepo,
			link.to.path,
			headBranch.name,
			link.from.content,
		).then(({ found, equal }) => resolve(!found || !equal));
	});
}

export { main };
