module.exports = {
  transform: {
    '^.+\\.(js|jsx|ts|tsx)$': 'babel-jest', // Transform JS, JSX, TS, and TSX files
  },
  testEnvironment: 'jsdom', // Use jsdom as the test environment
};
