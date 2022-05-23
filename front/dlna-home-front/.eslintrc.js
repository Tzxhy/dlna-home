module.exports = {
    env: {
        browser: true,
        es2021: true,
        es6: true,
        node: true,
    },
    extends: ['eslint:recommended', 'plugin:@typescript-eslint/recommended'],

    parserOptions: {
        parser: '@typescript-eslint/parser',
        ecmaVersion: 12,
        sourceType: 'module',
    },
    plugins: ['@typescript-eslint', 'simple-import-sort'],
    rules: {
        'eol-last': [
            'error', 'always'
        ],
        'space-infix-ops': [
            'error',
            {
                int32Hint: false,
            },
        ],
        'no-empty': [
            'error',
            {
                allowEmptyCatch: true,
            },
        ],
        'max-len': [
            'error',
            {
                code: 120,
                comments: 999,
                tabWidth: 4,
            },
        ],
        'array-element-newline': ['error', 'consistent'],
        'comma-dangle': [
            'error',
            {
                arrays: 'only-multiline',
                objects: 'always',
                imports: 'always',
                exports: 'always',
                functions: 'never',
            },
        ],
        'no-multiple-empty-lines': [
            'error',
            {
                max: 2,
            },
        ],
        'no-trailing-spaces': ['error'],
        'simple-import-sort/imports': 'error',
        'simple-import-sort/exports': 'error',
        'object-property-newline': [
            'error',
            {
                allowAllPropertiesOnSameLine: false,
            },
        ],
        'object-curly-newline': [
            'error',
            {
                // multiline: true,
                minProperties: 1,
            },
        ],
        'padding-line-between-statements': [
            'error',
            {
                blankLine: 'always',
                prev: 'var',
                next: 'return',
            },
        ],
        'no-useless-rename': [
            'error',
            {
                ignoreImport: true,
                ignoreExport: true,
                ignoreDestructuring: true,
            },
        ],
        'semi-style': ['error', 'last'],
        indent: ['error', 4],
        'linebreak-style': ['error', 'unix'],
        quotes: ['error', 'single'],
        semi: ['error', 'always'],
        'no-extra-boolean-cast': 'off',
        'no-async-promise-executor': 'off',
        '@typescript-eslint/ban-types': 'off',
        '@typescript-eslint/no-empty-function': 'off',
        '@typescript-eslint/explicit-module-boundary-types': 'off',
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/no-var-requires': 'off',
        '@typescript-eslint/no-non-null-assertion': 'off',
        '@typescript-eslint/no-unused-vars': [
            'error',
            {
                argsIgnorePattern: '^_',
                varsIgnorePattern: '^_',
            },
        ],
        '@typescript-eslint/ban-ts-comment': 0,
        'no-mixed-spaces-and-tabs': 0,
        'no-multi-spaces': ['error'],
    },
    ignorePatterns: [
    ],
};
