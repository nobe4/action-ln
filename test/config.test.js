const config = require("../src/config");
const yaml = require("js-yaml");

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
			])("%p", (c) => {
				return expect(config.validate(c)).rejects.toBeInstanceOf(
					config.ValidationError,
				);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					links: [],
				},
				{
					links: [
						{
							from: { path: "a/b" },
							to: { path: "b/c" },
						},
					],
				},
				{
					links: [
						{
							from: {
								path: "a/b",
								repo: "x",
							},
							to: {
								path: "b/c",
								repo: "y",
							},
						},
					],
				},
			])("%p", (c) => {
				return expect(config.validate(c)).resolves.toBe(c);
			});
		});
	});
});
