const { Config, Link, File, ValidationError } = require("../src/config");
const yaml = require("js-yaml");
const github = require("@actions/github");
const { dedent } = require("../src/utils");

describe("Config", () => {
	let c = new Config();

	describe("toString", () => {
		it("formats correctly", () => {
			const l1 = new Link();
			l1.from = new File();
			l1.from.repo = { repo: "repo", owner: "owner" };
			l1.from.path = "path";
			l1.from.content = "content";
			l1.to = "to";

			const l2 = new Link();
			l2.from = "from2";
			l2.to = "to2";

			c.data.links = [l1, l2];
			c.path = "path";

			expect(c.toString()).toStrictEqual(
				dedent(
					`
					path: path
					links:
					  -
					    from:
					        owner/repo:path
					        content
					    to:
					        to
					  -
					    from:
					        from2
					    to:
					        to2
					`,
				).trim(),
			);
		});
	});

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

	describe("toString", () => {
		it("formats correctly", () => {
			l.from = "from";
			l.to = "to";
			expect(l.toString()).toStrictEqual("from:\n    from\nto:\n    to");
		});
	});

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
				const mockFileParse = jest
					.spyOn(File.prototype, "parse")
					.mockImplementation(() => "parsed");

				l.raw = raw;
				l.parse();
				expect(l.from).toStrictEqual("parsed");
				expect(l.to).toStrictEqual("parsed");
				expect(mockFileParse).toHaveBeenCalled();
			});
		});
	});
});

describe("File", () => {
	let l = new File();

	describe("toString", () => {
		it("formats correctly", () => {
			l.repo = { repo: "repo", owner: "owner" };
			l.path = "path";
			expect(l.toString()).toStrictEqual("owner/repo:path");
		});
		it("formats correctly with a content", () => {
			l.repo = { repo: "repo", owner: "owner" };
			l.path = "path";
			l.content = "some\ncontent";
			expect(l.toString()).toStrictEqual("owner/repo:path\nsome\ncontent");
		});
	});

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
					.spyOn(File.prototype, "parsePath")
					.mockImplementation(() => {});

				const mockParseRepo = jest
					.spyOn(File.prototype, "parseRepo")
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
