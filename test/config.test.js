const config = require("../src/config");
const yaml = require("js-yaml");

describe("config", () => {
	describe("read", () => {
		test("missing file", async () => {
			expect.assertions(1);

			return config.read("./test/fixtures/config/not_a_file").catch((e) => {
				expect(e.message).toMatch(/ENOENT: no such file or directory, open /);
			});
		});

		test("not a YAML file", async () => {
			expect.assertions(1);

			return config.read("./test/fixtures/config/not_yaml.txt").catch((e) => {
				expect(e).toBeInstanceOf(yaml.YAMLException);
			});
		});

		test("invalid YAML config", async () => {
			expect.assertions(1);

			return config
				.read("./test/fixtures/config/invalid_config.yaml")
				.then((e) => {
					expect(e).not.toBeNull();
				});
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
							to: {},
						},
					],
				},
			])("%p", (c) => {
				expect.assertions(1);
				config.validate(c).catch((e) => {
					expect(e).toBeInstanceOf(config.ValidationError);
				});
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
				expect.assertions(1);
				config.validate(c).then((c) => {
					expect(c).not.toBeNull();
				});
			});
		});
	});
});
