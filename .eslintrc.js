module.exports = {
  extends: [
    'standard-with-typescript'
  ],
  parserOptions: {
    project: './tsconfig.json'
  },
  rules: {
    '@typescript-eslint/no-non-null-assertion': 'off'
  }
}
