// Provides a GitHub API client using Octokit.

import { Octokit } from '@octokit/rest';

type GitHubClient = Octokit;

export const createGitHubClient = (token?: string): GitHubClient => {
  const authToken = token ?? process.env.GH_TOKEN;
  return new Octokit({ auth: authToken });
};
