// eslint.config.js
import { FlatCompat } from '@eslint/eslintrc';
import js from '@eslint/js';
import path from 'path';
import { fileURLToPath } from 'url';

// プラグインやパーサーは ESM import で取得
import tsPlugin from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import prettierPlugin from 'eslint-plugin-prettier';
import reactPlugin from 'eslint-plugin-react';
import rnPlugin from 'eslint-plugin-react-native';

// ESM で __dirname を再現
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// FlatCompat のインスタンス化
const compat = new FlatCompat({
  baseDirectory: __dirname, // {String} flat-config 読み込みの基準ディレクトリ
  resolvePluginsRelativeTo: __dirname, // {String} プラグイン解決の基準ディレクトリ
  recommendedConfig: js.configs.recommended, // 必須: "eslint:recommended" 相当
  allConfig: js.configs.all, // 必須: "eslint:all" を使う場合
});

export default [
  // 既存の extends は FlatCompat を経由して読み込む
  ...compat.extends(
    'plugin:@typescript-eslint/recommended',
    'plugin:react/recommended',
    'plugin:react-native/all',
    'plugin:prettier/recommended',
  ),

  // node_modules は無視
  { ignores: ['node_modules/**'] },

  // 独自ルールの上書き
  {
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 2021,
        sourceType: 'module',
        ecmaFeatures: { jsx: true },
      },
      globals: {
        React: 'readonly',
      },
    },
    plugins: {
      '@typescript-eslint': tsPlugin,
      react: reactPlugin,
      'react-native': rnPlugin,
      prettier: prettierPlugin,
    },
    settings: {
      react: { version: 'detect' },
    },
    rules: {
      // TypeScript
      '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/no-explicit-any': 'warn',

      // React
      'react/prop-types': 'off',
      'react/react-in-jsx-scope': 'off',
      'react/display-name': 'off',
      'react/jsx-curly-brace-presence': ['error', { props: 'never', children: 'never' }],

      // React Native
      'react-native/no-inline-styles': 'warn',
      'react-native/no-raw-text': ['error', { skip: ['Text'] }],
      'react-native/no-unused-styles': 'error',
      'react-native/split-platform-components': 'warn',

      // Prettier
      'prettier/prettier': ['error', {}, { usePrettierrc: true }],

      // その他
      'no-console': 'warn',
      'no-debugger': 'error',
    },
  },
];
