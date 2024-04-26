import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Description from '../../components/Description';

describe('Description Component', () => {
  it('renders the welcome message correctly', () => {
    render(<Description />);
    // Check if the heading contains the expected text
    expect(screen.getByRole('heading', { level: 3 })).toHaveTextContent(
      'Welcome to Gopher Signal',
    );
  });

  it('renders the explanation text correctly', () => {
    render(<Description />);
    // Use a function matcher to find the text across multiple elements
    const explanationText = screen.getByText((_, node: any) => {
      // Helper function to match the node's text content accurately
      const hasText = (node: any) =>
        node.textContent ===
        'Gopher Signal uses smart technology to quickly summarize important points from Hacker News articles, giving you brief and useful updates.';
      const nodeHasText = hasText(node);
      // Ensure none of the child elements contain the same text
      const childrenDontHaveText = Array.from(node.children).every(
        (child) => !hasText(child),
      );

      return nodeHasText && childrenDontHaveText;
    });
    expect(explanationText).toBeInTheDocument();
  });

  it('contains a link to Hacker News with the correct attributes', () => {
    render(<Description />);
    // Verify that the link to Hacker News is correctly set up
    const hackerNewsLink = screen.getByRole('link', { name: 'Hacker News' });
    expect(hackerNewsLink).toHaveAttribute(
      'href',
      'https://news.ycombinator.com',
    );
    expect(hackerNewsLink).toHaveAttribute('target', '_blank');
    expect(hackerNewsLink).toHaveAttribute('rel', 'noopener noreferrer');
  });
});
