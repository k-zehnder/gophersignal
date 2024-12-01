import { expect } from 'chai';
import sinon from 'sinon';
import { main } from '../src/index';

describe('main', () => {
  it('processes and saves articles with mock dependencies', async () => {
    const mockDb = {
      saveArticles: sinon.stub().resolves(),
      updateArticleSummary: sinon.stub().resolves(),
      closeDatabaseConnection: sinon.stub().resolves(),
      connection: {},
    };

    const mockBrowser = { close: sinon.stub().resolves() };

    const mockProcessor = {
      processTopStories: sinon.stub().resolves([
        {
          title: 'Test Article',
          content: 'Content',
          link: 'http://hackernews.com',
        },
      ]),
    };

    const mockSummarizer = {
      summarizeArticles: sinon
        .stub()
        .resolves([{ title: 'Test Article', summary: 'Summary' }]),
    };

    await main({
      db: mockDb as any,
      browser: mockBrowser as any,
      articleProcessor: mockProcessor,
      articleSummarizer: mockSummarizer,
    });

    expect(mockProcessor.processTopStories.calledOnce).to.be.true;
    expect(mockSummarizer.summarizeArticles.calledOnce).to.be.true;
    expect(
      mockDb.saveArticles.calledOnceWithExactly([
        { title: 'Test Article', summary: 'Summary' },
      ])
    ).to.be.true;
    expect(mockDb.closeDatabaseConnection.calledOnce).to.be.true;
    expect(mockBrowser.close.calledOnce).to.be.true;
  });
});
