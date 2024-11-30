import { expect } from 'chai';
import nock from 'nock';
import sinon from 'sinon';
import Instructor from '@instructor-ai/instructor';
import OpenAI from 'openai';
import { createArticleSummarizer } from '../src/services/articleSummarizer';
import { SummarySchema } from '../src/types/index';

describe('articleSummarizer', () => {
  let consoleLogStub: sinon.SinonStub;
  let consoleErrorStub: sinon.SinonStub;

  before(() => {
    // Stub console methods
    consoleLogStub = sinon.stub(console, 'log');
    consoleErrorStub = sinon.stub(console, 'error');
  });

  after(() => {
    // Restore console methods
    consoleLogStub.restore();
    consoleErrorStub.restore();
  });

  afterEach(() => {
    // Clean up nock after each test
    nock.cleanAll();
  });

  const initTestClients = () => {
    const mockConfig = {
      apiKey: 'test-api-key',
      baseUrl: 'http://mock.api',
      model: 'mockModel',
      maxContentLength: 2000,
      maxSummaryLength: 500,
    };

    const openaiClient = new OpenAI({
      apiKey: mockConfig.apiKey,
      baseURL: mockConfig.baseUrl,
    });

    const instructorClient = Instructor({
      client: openaiClient,
      mode: 'JSON',
    });

    return { mockConfig, instructorClient };
  };

  it('summarizes articles correctly', async () => {
    const { mockConfig, instructorClient } = initTestClients();

    // Mock OpenAI API response in the expected format
    const mockSummaryResponse = {
      id: 'chatcmpl-123',
      object: 'chat.completion',
      created: 1677652288,
      choices: [
        {
          index: 0,
          message: {
            role: 'assistant',
            content: '{"summary":"This is a test summary."}',
          },
          finish_reason: 'stop',
        },
      ],
      usage: {
        prompt_tokens: 9,
        completion_tokens: 12,
        total_tokens: 21,
      },
    };

    // Set up nock to intercept the HTTP POST request to the correct endpoint
    nock(mockConfig.baseUrl)
      .post('/chat/completions')
      .reply(200, mockSummaryResponse, {
        'Content-Type': 'application/json',
      });

    // Create an instance of the summarizer with mock dependencies
    const summarizer = createArticleSummarizer(
      instructorClient,
      mockConfig,
      SummarySchema
    );

    // Mock article input
    const articles = [
      {
        title: 'Test Article',
        content: 'Test content',
        link: 'http://example.com',
        summary: '', // Initialize with no summary
      },
    ];

    // Call the summarizeArticles method
    const result = await summarizer.summarizeArticles(articles);

    // Assert the result matches the expected summary
    expect(result[0].summary).to.equal('This is a test summary.');

    // Ensure the mock request was called
    expect(nock.isDone()).to.be.true;
  });

  it('handles API errors gracefully', async () => {
    const { mockConfig, instructorClient } = initTestClients();

    // Set up nock to simulate an API error on the correct endpoint
    nock(mockConfig.baseUrl).post('/chat/completions').reply(500, {
      error: 'Internal Server Error',
    });

    // Create an instance of the summarizer with mock dependencies
    const summarizer = createArticleSummarizer(
      instructorClient,
      mockConfig,
      SummarySchema
    );

    // Mock article input
    const articles = [
      {
        title: 'Error Article',
        content: 'This content will cause an API error.',
        link: 'http://example.com',
        summary: '', // Initialize with no summary
      },
    ];

    // Call the summarizeArticles method
    const result = await summarizer.summarizeArticles(articles);

    // Assert the summary is empty or indicates an error
    expect(result[0].summary).to.equal('');
    expect(nock.isDone()).to.be.true;
  });

  it('truncates content exceeding max length', async () => {
    const { mockConfig, instructorClient } = initTestClients();

    // Update the mock configuration to enforce truncation
    mockConfig.maxContentLength = 20;

    // Mock OpenAI API response in the expected format
    const mockSummaryResponse = {
      id: 'chatcmpl-124',
      object: 'chat.completion',
      created: 1677652290,
      choices: [
        {
          index: 0,
          message: {
            role: 'assistant',
            content: '{"summary":"This is a truncated test summary."}',
          },
          finish_reason: 'stop',
        },
      ],
      usage: {
        prompt_tokens: 9,
        completion_tokens: 12,
        total_tokens: 21,
      },
    };

    // Set up nock to intercept the HTTP POST request to the correct endpoint
    nock(mockConfig.baseUrl)
      .post('/chat/completions')
      .reply(200, mockSummaryResponse, {
        'Content-Type': 'application/json',
      });

    // Create an instance of the summarizer with mock dependencies
    const summarizer = createArticleSummarizer(
      instructorClient,
      mockConfig,
      SummarySchema
    );

    // Mock article input
    const articles = [
      {
        title: 'Long Content Article',
        content:
          'This content is much longer than the maximum allowed length and should be truncated.',
        link: 'http://example.com',
        summary: '', // Initialize with no summary
      },
    ];

    // Call the summarizeArticles method
    const result = await summarizer.summarizeArticles(articles);

    // Assert the result matches the expected summary
    expect(result[0].summary).to.equal('This is a truncated test summary.');

    // Ensure the mock request was called
    expect(nock.isDone()).to.be.true;
  });
});
