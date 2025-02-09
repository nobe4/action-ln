const {
	indent: indentFn,
	dedent,
	branchName,
	commitMessage,
	pullBody,
} = require("../src/format");

describe("indent", () => {
	it.each([
		{ str: "a", indent: "", want: "a" },
		{ str: "a", indent: " ", want: " a" },
		{ str: "a", indent: "--", want: "--a" },
		{ str: "a", indent: undefined, want: "    a" },
	])("%# %j", ({ str, indent, want }) => {
		expect(indentFn(str, indent)).toStrictEqual(want);
	});
});

describe("dedent", () => {
	it.each([
		{
			str: "",
			trim: true,
			want: "",
		},
		{
			str: "a",
			trim: true,
			want: "a",
		},
		{
			str: `
a
 b
  c
`,
			trim: true,
			want: `a
 b
  c`,
		},
		{
			str: `
a
 b
  c
`,
			trim: false,
			want: `
a
 b
  c
`,
		},
		{
			str: `
	  a
	   b
	    c
`,
			trim: true,
			want: `a
 b
  c`,
		},
		{
			str: `
	  a
	   b
	    c
`,
			trim: false,
			want: `
a
 b
  c
`,
		},
		{
			str: `
	    a
	   b
	    c
`,
			trim: false,
			want: `
 a
b
 c
`,
		},
	])("%# %j", ({ str, trim, want }) => {
		expect(dedent(str, trim)).toStrictEqual(want);
	});
});

describe("branchName", () => {
	it("formats the branch name", () => {
		const link = { SHA256: "sha256" };
		expect(branchName(link)).toEqual("link-sha256");
	});
});

describe("commitMessage", () => {
	it("formats the commit message", () => {
		const link = {
			from: { toString: () => "from" },
			to: { toString: () => "to" },
		};
		expect(commitMessage(link)).toEqual(
			dedent(`
			auto(link): update links

			From: from
			To:   to
		`),
		);
	});
});
describe("commitMessage", () => {
	it("formats the commit message", () => {
		const link = {
			from: { toString: () => "from" },
			to: { toString: () => "to" },
		};
		expect(commitMessage(link)).toEqual(
			dedent(`
			auto(link): update links

			From: from
			To:   to
		`),
		);
	});
});

describe("pullBody", () => {
	it("formats the pull request body", () => {
		const config = {
			URL: "URL",
		};
		const link = {
			from: { toString: () => "from" },
			to: { toString: () => "to" },
		};

		const body = pullBody(link, config);
		expect(body).toEqual(
			expect.stringContaining(
				"This automated PR updates the following link as defined in [the configuration](URL).",
			),
		);
		expect(body).toEqual(expect.stringContaining("From: from"));
		expect(body).toEqual(expect.stringContaining("To:   to"));
	});
});
