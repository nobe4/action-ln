import {
	indent as indentFn,
	dedent,
	commitMessage,
	pullBody,
} from "../src/format.js";

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

describe("commitMessage", () => {
	it("formats the commit message", () => {
		const link = {
			from: { toString: () => "from" },
			to: { toString: () => "to" },
		};
		expect(commitMessage(link)).toEqual(
			dedent(`
			auto(ln): update links

			From: from
			To:   to
		`),
		);
	});
});

describe("pullBody", () => {
	it("formats the pull request body", () => {
		const config = {
			path: "path",
			URL: "URL",
		};
		const group = [
			{
				from: { toString: () => "from0" },
				to: { toString: () => "to0" },
			},
			{
				from: { toString: () => "from1" },
				to: { toString: () => "to1" },
			},
		];
		const context = {
			workflow: "workflow",
			repo: {
				owner: "owner",
				repo: "repo",
			},
			serverUrl: "serverUrl",
			runId: "runId",
		};

		const body = pullBody(group, config, context);
		expect(body).toEqual(expect.stringContaining("[configuration](URL)"));
		expect(body).toEqual(expect.stringContaining("`from0` | `to0`"));
		expect(body).toEqual(expect.stringContaining("`from1` | `to1`"));
	});
});
