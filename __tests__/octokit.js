import { jest } from "@jest/globals";

import * as core from "../__mocks__/@actions/core.js";
jest.unstable_mockModule("@actions/core", () => core);

import { github } from "../__mocks__/@actions/github.js";
jest.unstable_mockModule("@actions/github", () => github);

const retry = { retry: jest.fn() };
jest.unstable_mockModule("@octokit/plugin-retry", () => retry);

const rest = { Octokit: class Octokit {} };
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

	it("create an octokit instance with app authentication", () => {
		expect(createOctokit({ appId: "id", appPrivKey: "key" })).toBeInstanceOf(
			rest.Octokit,
		);
		// TODO: not really sure how to check that the constructor was called.
	});

	it("create an octokit instance with token authentication", () => {
		github.getOctokit.mockReturnValue("octokit");
		expect(createOctokit({ token: "TOKEN" })).toStrictEqual("octokit");
		expect(github.getOctokit).toHaveBeenCalledWith("TOKEN", {
			userAgent: "nobe4/action-ln",
			additionalPlugins: [expect.any(Function)],
			log: console,
		});
	});
});
