import globals from "globals";
import js from "@eslint/js";
import prettierRecommended from "eslint-plugin-prettier/recommended";
import jestRecommended from "eslint-plugin-jest/recommended";

export default [
	js.configs.recommended,
	prettierRecommended,
	jestRecommended,
	{
		ignores: ["**/coverage", "**/dist", "**/linter", "**/node_modules"],
	},
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node,
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
		rules: {

		},
	},
];
