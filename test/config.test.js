const { Config, Link, Location, ValidationError } = require("../src/config");
const yaml = require("js-yaml");
const github = require("@actions/github");

describe("Config", () => {
	let c = new Config();

	describe("load", () => {
		test("calls read and parse", async () => {
			const mockRead = jest
				.spyOn(Config.prototype, "read")
				.mockResolvedValue("read");

			const mockParse = jest
				.spyOn(Config.prototype, "parse")
				.mockResolvedValue("parsed");

			await expect(c.load()).resolves.toStrictEqual("parsed");
			expect(mockRead).toHaveBeenCalled();
			expect(mockParse).toHaveBeenCalled();
		});
	});

	describe("read", () => {
		test("missing file", () => {
			c.path = "./test/fixtures/config/not_a_file";
			return expect(c.read()).rejects.toThrow(
				/ENOENT: no such file or directory, open /,
			);
		});

		test("not a YAML file", () => {
			c.path = "./test/fixtures/config/not_yaml.txt";
			return expect(c.read()).rejects.toThrow(yaml.YAMLException);
		});

		test("invalid YAML config", () => {
			c.path = "./test/fixtures/config/invalid_config.yaml";
			return expect(c.read()).resolves.not.toBeNull();
		});
	});

	describe("parse", () => {
		describe("fails", () => {
			it.each([
				null,
				{},
				{ links: null },
				{ links: "not a list" },
				{ links: 1 },
			])("%# %j", (data) => {
				c.data = data;
				expect(() => c.parse()).toThrow(ValidationError);
			});
		});

		describe("succeeds", () => {
			it.each([
				{
					data: { links: [0, 1, 2] },
					want: {
						links: ["parsed", "parsed", "parsed"],
					},
				},
			])("%# %j", ({ data, want }) => {
				const mockLink = jest
					.spyOn(Link.prototype, "parse")
					.mockImplementation(() => "parsed");
				c.data = data;
				expect(c.parse().data).toStrictEqual(want);
				expect(mockLink).toHaveBeenCalled();
			});
		});
	});
});

describe("Link", () => {
	let l = new Link();

	describe("parse", () => {
		describe("fails", () => {
			it.each([undefined, "a", 1, {}, { from: {} }, { to: {} }])(
				"%# %p",
				(raw) => {
					l.raw = raw;
					return expect(() => l.parse()).toThrow(ValidationError);
				},
			);
		});

		describe("succeeds", () => {
			it.each([{ from: {}, to: {} }])("%# %p", (raw) => {
				const mockLocationParse = jest
					.spyOn(Location.prototype, "parse")
					.mockImplementation(() => "parsed");

				l.raw = raw;
				l.parse();
				expect(l.from).toStrictEqual("parsed");
				expect(l.to).toStrictEqual("parsed");
				expect(mockLocationParse).toHaveBeenCalled();
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

				expect(l.parse()).toStrictEqual(l);
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
			const defaultRepo = {
				owner: "owner",
				repo: "repo",
			};

			jest.mock("@actions/github");
			beforeEach(() => {
				github.context = { repo: defaultRepo };
			});

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
