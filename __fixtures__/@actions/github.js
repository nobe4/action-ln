import { jest } from "@jest/globals";

export const github = {
	context: {
		repo: { owner: "owner", repo: "repo" },
	},
	getOctokit: jest.fn(),
};
