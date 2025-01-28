import { orchestrateWorkflow } from '../src/index';

describe('orchestrateWorkflow', () => {
  const STATUS_SUCCESS = 0;
  const STATUS_FAILURE = 1;

  it('returns status code 0 on success and processes articles correctly', async () => {
    // Arrange
    const mockDb = {
      saveArticles: jest.fn().mockResolvedValue(undefined),
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

    // Act
    const statusCode = await orchestrateWorkflow({
      db: mockDb as any,
      browser: mockBrowser as any,
      articleProcessor: mockProcessor,
      articleSummarizer: mockSummarizer,
    });

    // Assert
    expect(statusCode).toBe(STATUS_SUCCESS);
    expect(mockProcessor.processTopStories).toHaveBeenCalledTimes(1);
    expect(mockSummarizer.summarizeArticles).toHaveBeenCalledWith([
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
    ]);
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

  it('returns status code 1 when an error occurs in saving articles', async () => {
    // Arrange
    const mockDb = {
      saveArticles: jest.fn().mockRejectedValue(new Error('DB error')),
      closeDatabaseConnection: jest.fn().mockResolvedValue(undefined),
      connection: {},
    };

    const mockBrowser = { close: jest.fn().mockResolvedValue(undefined) };

    const mockProcessor = {
      processTopStories: jest.fn().mockResolvedValue([]),
    };

    const mockSummarizer = {
      summarizeArticles: jest.fn(),
    };

    // Act
    const statusCode = await orchestrateWorkflow({
      db: mockDb as any,
      browser: mockBrowser as any,
      articleProcessor: mockProcessor,
      articleSummarizer: mockSummarizer,
    });

    // Assert
    expect(statusCode).toBe(STATUS_FAILURE);
    expect(mockDb.closeDatabaseConnection).toHaveBeenCalledTimes(1);
    expect(mockBrowser.close).toHaveBeenCalledTimes(1);
  });

  it('handles exceptions during browser closure gracefully', async () => {
    // Arrange
    const mockDb = {
      saveArticles: jest.fn().mockResolvedValue(undefined),
      closeDatabaseConnection: jest.fn().mockResolvedValue(undefined),
      connection: {},
    };

    const mockBrowser = {
      close: jest.fn().mockRejectedValue(new Error('Browser error')),
    };

    const mockProcessor = {
      processTopStories: jest.fn().mockResolvedValue([]),
    };

    const mockSummarizer = {
      summarizeArticles: jest.fn(),
    };

    // Act
    const statusCode = await orchestrateWorkflow({
      db: mockDb as any,
      browser: mockBrowser as any,
      articleProcessor: mockProcessor,
      articleSummarizer: mockSummarizer,
    });

    // Assert
    expect(statusCode).toBe(STATUS_SUCCESS);
    expect(mockDb.closeDatabaseConnection).toHaveBeenCalledTimes(1);
    expect(mockBrowser.close).toHaveBeenCalledTimes(1);
  });
});
