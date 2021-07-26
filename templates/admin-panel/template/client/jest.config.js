module.exports = {
    'projects': [
        {
            'displayName': 'test',
            'moduleFileExtensions': [
                'js',
                'json',
                'sfc'
            ],
            'transform': {
                '^.+\\.sht$': 'sham-ui-templates-jest-preprocessor',
                '^.+\\.sfc$': 'sham-ui-templates-jest-preprocessor',
                '^.+\\.js$': 'sham-ui-macro-jest-preprocessor'
            },
            'transformIgnorePatterns': [],
            'collectCoverageFrom': [
                'src/**/*.js'
            ],
            'coveragePathIgnorePatterns': [
                '^.+\\.sht$',
                '<rootDir>/.babel-plugin-macrosrc.js',
                '<rootDir>/__tests__/setup-jest.js'
            ],
            'testPathIgnorePatterns': [
                '<rootDir>/node_modules/',
                '<rootDir>/__tests__/integration/helpers.js',
                '<rootDir>/__tests__/setup-jest.js'
            ],
            'testMatch': [
                '<rootDir>/__tests__/**/*.js',
            ],
            'setupTestFrameworkScriptFile': '<rootDir>/__tests__/setup-jest.js',
            "testURL": "http://client.example.com"
        },
        {
            'runner': 'jest-runner-eslint',
            'displayName': 'eslint',
            'moduleFileExtensions': [
                'js',
                'json',
                'sfc'
            ],
            'testMatch': [
                '<rootDir>/src/**/*.*',
                '<rootDir>/__tests__/**/*.js',
                '<rootDir>/__mocks__/**/*.js',
                '<rootDir>/jest/**/*.js'
            ],
            'testPathIgnorePatterns': [
                '<rootDir>/dist'
            ]
        },
        {
            'runner': 'jest-runner-stylelint',
            'displayName': 'stylelint',
            'moduleFileExtensions': [
                'scss'
            ],
            'testMatch': [
                '**/*.scss'
            ]
        }
    ]
};
