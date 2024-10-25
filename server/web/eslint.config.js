import jsRecommendedLib from '@eslint/js';
import typescriptPlugin from '@typescript-eslint/eslint-plugin';
import typescriptParser from '@typescript-eslint/parser';
import importPlugin from 'eslint-plugin-import';
import jsxA11yPlugin from 'eslint-plugin-jsx-a11y';
import prettierPlugin from 'eslint-plugin-prettier';
import reactPlugin from 'eslint-plugin-react';
import reactHooksPlugin from 'eslint-plugin-react-hooks';
import sonarjsPlugin from 'eslint-plugin-sonarjs';
import testingLibPlugin from 'eslint-plugin-testing-library';
import { fixupPluginRules } from '@eslint/compat';

// eslint.config.js
export default [
	{
		files: ['**/styles/**', '**/__tests__/**', '**/*.test.tsx', '**/*.test.ts', '*.less', 'src/**/*.tsx'],
		languageOptions: {
			parserOptions: {
				ecmaFeatures: {
					jsx: true
				}
			},
			parser: typescriptParser
		},
		plugins: {
			react: reactPlugin,
			'react-hooks': fixupPluginRules(reactHooksPlugin),
			sonarjs: sonarjsPlugin,
			import: fixupPluginRules(importPlugin),
			'jsx-a11y': jsxA11yPlugin,
			'@typescript-eslint': typescriptPlugin,
			prettier: prettierPlugin,
			js: jsRecommendedLib,
			'testing-library': fixupPluginRules(testingLibPlugin)
		},
		settings: {
			react: {
				version: 'detect'
			},
			'import/resolver': {
				typescript: {},
				node: {
					paths: ['src'],
					extensions: ['.js', '.jsx', '.ts', '.tsx']
				}
			}
		},
		rules: {
			...jsRecommendedLib.configs.recommended.rules,
			...reactPlugin.configs.recommended.rules,
			...reactPlugin.configs['jsx-runtime'].rules,
			...reactHooksPlugin.configs.recommended.rules,
			...importPlugin.configs.recommended.rules,
			...jsxA11yPlugin.flatConfigs.recommended.rules,
			semi: 'error',
			'no-unused-vars': 'off',
			'@typescript-eslint/no-unused-vars': 'error',
			'no-undef': 'off',
			'prefer-const': 'error',
			'testing-library/no-debugging-utils': 'warn',
			'testing-library/no-dom-import': 'off',
			'@typescript-eslint/no-explicit-any': 'off',
			'@typescript-eslint/no-empty-function': 'off',
			'react/react-in-jsx-scope': 'off',
			'react/jsx-uses-react': 'off',
			'no-console': 'warn',
			'no-duplicate-imports': 'warn',
			'jsx-a11y/no-autofocus': 'off',
			'sonarjs/cognitive-complexity': 'off',
			'sonarjs/elseif-without-else': 'off',
			'sonarjs/max-switch-cases': 'error',
			'sonarjs/no-all-duplicated-branches': 'error',
			'sonarjs/no-collapsible-if': 'error',
			'sonarjs/no-collection-size-mischeck': 'error',
			'sonarjs/no-duplicate-string': 'off',
			'sonarjs/no-duplicated-branches': 'error',
			'sonarjs/no-element-overwrite': 'error',
			'sonarjs/no-empty-collection': 'error',
			'sonarjs/no-extra-arguments': 'error',
			'sonarjs/no-gratuitous-expressions': 'error',
			'sonarjs/no-identical-conditions': 'error',
			'sonarjs/no-identical-expressions': 'error',
			'sonarjs/no-identical-functions': 'error',
			'sonarjs/no-ignored-return': 'error',
			'sonarjs/no-inverted-boolean-check': 'error',
			'sonarjs/no-nested-switch': 'error',
			'sonarjs/no-nested-template-literals': 'error',
			'sonarjs/no-one-iteration-loop': 'error',
			'sonarjs/no-redundant-boolean': 'error',
			'sonarjs/no-redundant-jump': 'error',
			'sonarjs/no-same-line-conditional': 'error',
			'sonarjs/no-small-switch': 'error',
			'sonarjs/no-unused-collection': 'error',
			'sonarjs/no-use-of-empty-return-value': 'error',
			'sonarjs/no-useless-catch': 'error',
			'sonarjs/non-existent-operator': 'error',
			'sonarjs/prefer-immediate-return': 'error',
			'sonarjs/prefer-object-literal': 'error',
			'sonarjs/prefer-single-boolean-return': 'error',
			'sonarjs/prefer-while': 'error',
			'import/order': [
				'warn',
				{
					alphabetize: {
						caseInsensitive: true,
						order: 'asc'
					},
					groups: [['builtin', 'external', 'index', 'sibling', 'parent', 'internal']],
					'newlines-between': 'always',
					pathGroups: [
						{
							pattern: '*.less',
							group: 'index',
							patternOptions: {
								matchBase: true
							},
							position: 'before'
						},
						{
							pattern: '*.json',
							group: 'index',
							patternOptions: {
								matchBase: true
							},
							position: 'after'
						}
					]
				}
			]
		}
	}
];
