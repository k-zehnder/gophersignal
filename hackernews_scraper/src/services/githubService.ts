// GitHubService provides methods to fetch Git-related info.

import { execSync } from 'child_process';
import { Octokit } from '@octokit/rest';
import { GitHubConfig } from '../types/config';
import { GitHubService } from '../types';

export function createGitHubService(
  client: Octokit,
  cfg: GitHubConfig
): GitHubService {
  const getCommitHash = async (): Promise<string> => {
    let sha: string | undefined;

    // Env override
    if (process.env.COMMIT_HASH) {
      sha = process.env.COMMIT_HASH;
    }

    // GitHub API
    if (!sha) {
      try {
        const { data } = await client.rest.repos.getCommit({
          owner: cfg.owner,
          repo: cfg.repo,
          ref: cfg.branch,
        });
        sha = data.sha.slice(0, 7);
      } catch (err) {
        console.warn('GitHub API failed:', (err as Error).message);
      }
    }

    // Local git
    if (!sha) {
      console.log('Trying local git…');
      try {
        sha = execSync('git rev-parse --short HEAD', {
          encoding: 'utf-8',
        }).trim();
      } catch (err) {
        console.error('Local git failed:', (err as Error).message);
      }
    }

    // Final fallback
    if (!sha) {
      sha = 'unknown';
      console.log('Falling back to unknown SHA');
    }

    return sha;
  };

  // Expose the service API
  return {
    getCommitHash,
  };
}
