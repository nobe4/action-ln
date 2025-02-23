import globals from "globals";
import js from "@eslint/js";
import prettierRecommended from "eslint-plugin-prettier/recommended";
import jest from "eslint-plugin-jest";

export default [
	js.configs.recommended,
	prettierRecommended,
	jest.configs["flat/recommended"],
	{
		rules: {
			// Seems that the support for mocking with ESM is not perfect yet.
			// I found a thing  that worked, so I'll stick with it.
			"jest/no-mocks-import": "off",
		},
	},
	{
		ignores: ["**/coverage", "**/dist", "**/linter", "**/node_modules"],
	},
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node,
				...jest.environments.globals.globals,
			},
			ecmaVersion: "latest",
			sourceType: "module",
			parserOptions: {
				ecmaFeatures: {
					implicitStrict: true,
				},
			},
		},
	},
	{
		rules: {},
	},
];
