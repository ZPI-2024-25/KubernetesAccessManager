export default {
    preset: 'ts-jest',
    testEnvironment: 'jest-environment-jsdom',
    transform: {
        // "^.+\\.tsx?$": "ts-jest" // process `*.tsx` files with `ts-jest`
        "^.+\\.tsx?$": ["ts-jest", { tsconfig: 'tsconfig.app.json' }]
    },
    moduleNameMapper: {
        '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
        '\\.(gif|ttf|eot|svg|png)$': '<rootDir>/test/__ mocks __/fileMock.js',
    },
    transformIgnorePatterns: ["<rootDir>/node_modules/"],
}