import globals from "globals";
import js from "@eslint/js";
import prettierRecommended from "eslint-plugin-prettier/recommended";

export default [
	js.configs.recommended,
	prettierRecommended,
	{
		ignores: ["**/coverage", "**/dist", "**/linter", "**/node_modules"],
	},
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node,
			},
		},
	},
];
