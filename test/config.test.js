const config = require("../src/config");
const yaml = require("js-yaml");

const github = require("@actions/github");
jest.mock("@actions/github");

beforeEach(() => {
	github.context = {
		repo: {
			owner: "owner",
			repo: "repo",
		},
	};
});

describe("config", () => {
	describe("read", () => {
		test("missing file", () => {
			return expect(
				config.read("./test/fixtures/config/not_a_file"),
			).rejects.toThrow(/ENOENT: no such file or directory, open /);
		});

		test("not a YAML file", () => {
			return expect(
				config.read("./test/fixtures/config/not_yaml.txt"),
			).rejects.toThrow(yaml.YAMLException);
		});

		test("invalid YAML config", () => {
			return expect(
				config.read("./test/fixtures/config/invalid_config.yaml"),
			).resolves.not.toBeNull();
		});
	});

	describe("validate", () => {
		describe("fails", () => {
			it.each([
				null,
				{},
				{ links: null },
				{ links: "not a list" },
				{
					links: ["a", "b"],
				},
				{
					links: [{}, {}],
				},
				{
					links: [
						{
							from: {},
						},
					],
				},
				{
					links: [
						{
							from: {},
							to: {},
						},
					],
				},
				{
					links: [
						{
							from: { path: "" },
							to: {},
						},
					],
				},
				{
					links: [
						{
							from: { path: "non-empty" },
							to: { path: "" },
						},
					],
				},
				{
					links: [
						{
							from: { path: "a", repo: "x" },
							to: { path: "a" },
						},
					],
				},
				{
					links: [
						{
							from: { path: "a", repo: "x/y" },
							to: { path: "a", repo: "z" },
						},
					],
				},
				{
					links: [
						{
							from: { path: "a", repo: "x/" },
							to: { path: "a", repo: "z" },
						},
					],
				},
			])("%# %j", (c) => {
				return expect(config.validate(c)).rejects.toBeInstanceOf(
					config.ValidationError,
				);
			});
		});

		describe("succeeds", () => {
			it.each([
				{ c: { links: [] }, want: { links: [] } },
				{
					c: {
						links: [
							{
								from: { path: "a/b" },
								to: { path: "b/c" },
							},
						],
					},
					want: {
						links: [
							{
								from: {
									path: "a/b",
									repo: {
										owner: "owner",
										repo: "repo",
									},
								},
								to: {
									path: "b/c",
									repo: {
										owner: "owner",
										repo: "repo",
									},
								},
							},
						],
					},
				},
				{
					c: {
						links: [
							{
								from: {
									path: "a/b",
									repo: "x/y",
								},
								to: {
									path: "b/c",
									repo: "y/z",
								},
							},
						],
					},
					want: {
						links: [
							{
								from: {
									path: "a/b",
									repo: {
										owner: "x",
										repo: "y",
									},
								},
								to: {
									path: "b/c",
									repo: {
										owner: "y",
										repo: "z",
									},
								},
							},
						],
					},
				},
			])("%# %j", ({ c, want }) => {
				return expect(config.validate(c)).resolves.toStrictEqual(want);
			});
		});
	});
});
