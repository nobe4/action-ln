const { indent: indentFn, dedent } = require("../src/utils");

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
