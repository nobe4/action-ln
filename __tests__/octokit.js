/* eslint-disable jest/no-mocks-import */

import { jest } from "@jest/globals";

import { github } from "../__mocks__/@actions/github.js";
jest.unstable_mockModule("@actions/github", () => github);

const rest = { Octokit: jest.fn() };
jest.unstable_mockModule("@octokit/rest", () => rest);

const authApp = { createAppAuth: jest.fn() };
jest.unstable_mockModule("@octokit/auth-app", () => authApp);

const { createOctokit } = await import("../src/octokit.js");

describe("createOctokit", () => {
	it("throws an error if all the auth items are missing", () => {
		expect(() => createOctokit({})).toThrow(
			"either token or app_* should be provided",
		);
	});

	it("create an octokit instance with app authentication", () => {});

	it("create an octokit instance with token authentication", () => {
		github.getOctokit.mockReturnValue("octokit");
		expect(createOctokit({ token: "TOKEN" })).toStrictEqual("octokit");
	});
});
