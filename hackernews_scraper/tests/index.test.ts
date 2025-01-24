import { main } from '../src/index';

describe('main', () => {
  it('processes and saves articles with mock dependencies', async () => {
    const mockDb = {
      saveArticles: jest.fn().mockResolvedValue(undefined),
      updateArticleSummary: jest.fn().mockResolvedValue(undefined),
      closeDatabaseConnection: jest.fn().mockResolvedValue(undefined),
      connection: {},
    };

    const mockBrowser = { close: jest.fn().mockResolvedValue(undefined) };

    const mockProcessor = {
      processTopStories: jest.fn().mockResolvedValue([
        {
          title: 'Test Article',
          content: 'Content',
          link: 'http://hackernews.com',
          summary: '',
          source: 'Hacker News',
          upvotes: 100,
          comment_count: 50,
          comment_link: 'http://hackernews.com/comments',
        },
      ]),
    };

    const mockSummarizer = {
      summarizeArticles: jest.fn().mockResolvedValue([
        {
          title: 'Test Article',
          content: 'Content',
          link: 'http://hackernews.com',
          summary: 'Summary',
          source: 'Hacker News',
          upvotes: 100,
          comment_count: 50,
          comment_link: 'http://hackernews.com/comments',
        },
      ]),
    };

    await main({
      db: mockDb as any,
      browser: mockBrowser as any,
      articleProcessor: mockProcessor,
      articleSummarizer: mockSummarizer,
    });

    // Use Jest's expect assertions
    expect(mockProcessor.processTopStories).toHaveBeenCalledTimes(1);
    expect(mockSummarizer.summarizeArticles).toHaveBeenCalledTimes(1);
    expect(mockDb.saveArticles).toHaveBeenCalledWith([
      {
        title: 'Test Article',
        content: 'Content',
        link: 'http://hackernews.com',
        summary: 'Summary',
        source: 'Hacker News',
        upvotes: 100,
        comment_count: 50,
        comment_link: 'http://hackernews.com/comments',
      },
    ]);
    expect(mockDb.closeDatabaseConnection).toHaveBeenCalledTimes(1);
    expect(mockBrowser.close).toHaveBeenCalledTimes(1);
  });
});
