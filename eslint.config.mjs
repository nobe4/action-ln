import globals from "globals";
import js from "@eslint/js";
import prettierRecommended from "eslint-plugin-prettier/recommended";
import jest from "eslint-plugin-jest";

export default [
	js.configs.recommended,
	prettierRecommended,
	jest.configs["flat/recommended"],
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
