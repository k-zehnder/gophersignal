import { createWorkflow } from '../src/workflow';
import { Services } from '../src/services/createServices';

describe('Workflow', () => {
  let mockServices: jest.Mocked<Services>;
  let workflow: ReturnType<typeof createWorkflow>;

  beforeEach(() => {
    // Set up a mocked version of the Services interface
    mockServices = {
      scraper: {
        scrapeFront: jest.fn().mockResolvedValue([]),
        scrapeTopStories: jest.fn().mockResolvedValue([]),
      },
      articleProcessor: {
        helpers: {
          categorizeArticles: jest
            .fn()
            .mockReturnValue({ flagged: [], dead: [], dupe: [] }),
          getTopArticlesWithContent: jest.fn().mockReturnValue([]),
        },
        processArticles: jest.fn().mockResolvedValue([]),
      },
      articleSummarizer: {
        summarizeArticles: jest.fn().mockResolvedValue([]),
      },
      db: {
        saveArticles: jest.fn().mockResolvedValue(undefined),
        closeDatabaseConnection: jest.fn().mockResolvedValue(undefined),
      },
      timeUtil: {
        today: '2024-02-07',
        yesterday: '2024-02-06',
      },
      browser: {
        close: jest.fn().mockResolvedValue(undefined),
      },
    } as unknown as jest.Mocked<Services>;

    // Create the workflow using the factory function
    workflow = createWorkflow(mockServices);
  });

  it('should execute the workflow successfully', async () => {
    await expect(workflow.run()).resolves.not.toThrow();

    // Verify that the expected service functions were called
    expect(mockServices.scraper.scrapeFront).toHaveBeenCalledTimes(1);
    expect(mockServices.scraper.scrapeTopStories).toHaveBeenCalledTimes(1);
    expect(mockServices.articleProcessor.processArticles).toHaveBeenCalledTimes(
      1
    );
    expect(
      mockServices.articleSummarizer.summarizeArticles
    ).toHaveBeenCalledTimes(2);
    expect(mockServices.db.saveArticles).toHaveBeenCalledTimes(1);
  });

  it('should handle shutdown gracefully', async () => {
    await expect(workflow.shutdown()).resolves.not.toThrow();

    // Verify that resources are properly closed
    expect(mockServices.db.closeDatabaseConnection).toHaveBeenCalledTimes(1);
    expect(mockServices.browser.close).toHaveBeenCalledTimes(1);
  });
});
