const { Config, Location, ValidationError } = require("../src/config");
const yaml = require("js-yaml");

const github = require("@actions/github");
jest.mock("@actions/github");

const defaultRepo = {
	owner: "owner",
	repo: "repo",
};

beforeEach(() => {
	github.context = { repo: defaultRepo };
});

describe("Config", () => {
	describe("read", () => {
		test("missing file", () => {
			return expect(
				new Config("./test/fixtures/config/not_a_file").read(),
			).rejects.toThrow(/ENOENT: no such file or directory, open /);
		});

		test("not a YAML file", () => {
			return expect(
				new Config("./test/fixtures/config/not_yaml.txt").read(),
			).rejects.toThrow(yaml.YAMLException);
		});

		test("invalid YAML config", () => {
			return expect(
				new Config("./test/fixtures/config/invalid_config.yaml").read(),
			).resolves.not.toBeNull();
		});
	});

	describe("validate", () => {
		let c = new Config();

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
			])("%# %j", (data) => {
				c.data = data;
				// Gotcha: c needs to keep its `this`, so wrapping it let it
				// keeps it.
				return expect(() => c.validate()).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{ data: { links: [] }, want: { links: [] } },
				{
					data: {
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
					data: {
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
			])("%# %j", ({ data, want }) => {
				c.data = data;
				return expect(c.validate()).toStrictEqual(want);
			});
		});
	});
});

describe("Location", () => {
	let l = new Location();

	describe("parse", () => {
		describe("fails", () => {
			it.each([null, undefined, ""])("%# %p", (raw) => {
				l.raw = raw;
				return expect(() => l.parse()).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			it.each(["non-nil"])("%# %p", (raw) => {
				const mockParsePath = jest
					.spyOn(Location.prototype, "parsePath")
					.mockImplementation(() => {});

				const mockParseRepo = jest
					.spyOn(Location.prototype, "parseRepo")
					.mockImplementation(() => {});

				l.raw = raw;

				expect(l.parse()).toStrictEqual({});
				expect(mockParsePath).toHaveBeenCalled();
				expect(mockParseRepo).toHaveBeenCalled();
			});
		});
	});

	describe("parsePath", () => {
		describe("fails", () => {
			it.each([null, undefined, "", "\n", "    ", " \t"])("%# %p", (raw) => {
				l.raw = { path: raw };
				return expect(() => l.parsePath()).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					raw: "a",
					want: "a",
				},
			])("%# %p", ({ raw, want }) => {
				l.raw = { path: raw };
				return expect(l.parsePath()).toStrictEqual(want);
			});
		});
	});

	describe("parseRepo", () => {
		describe("fails", () => {
			it.each([
				{ repo: "a" },
				{ repo: "a/" },
				{ repo: "/a" },
				{ repo: "/" },
				{ repo: {} },
				{ repo: { owner: "" } },
				{ repo: { repo: "" } },
				{ repo: { repo: "", owner: undefined } },
				{ repo: { repo: undefined, owner: "" } },
			])("%# %p", (raw) => {
				l.raw = raw;
				return expect(() => l.parseRepo()).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					raw: undefined,
					want: defaultRepo,
				},
				{
					raw: "",
					want: defaultRepo,
				},
				{
					raw: "owner/repo",
					want: {
						repo: "repo",
						owner: "owner",
					},
				},
				{
					raw: {
						repo: "repo",
						owner: "owner",
					},
					want: {
						repo: "repo",
						owner: "owner",
					},
				},
			])("%# %p", ({ raw, want }) => {
				l.raw = { repo: raw };
				return expect(l.parseRepo()).toStrictEqual(want);
			});
		});
	});
});
