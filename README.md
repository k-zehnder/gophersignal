# GopherSignal

[![CI/CD Pipeline](https://github.com/k-zehnder/gophersignal/actions/workflows/workflow.yml/badge.svg)](https://github.com/k-zehnder/gophersignal/actions/workflows/workflow.yml)

GopherSignal is a full-stack application designed to scrape Hacker News, generate AI-powered summaries of articles, and present them through a web interface, a RESTful API, and an RSS feed. This repository hosts the code powering the live website at [https://gophersignal.com/](https://gophersignal.com/).

## Table of Contents

*   [Features](#features)
*   [Architecture Overview](#architecture-overview)
*   [Technology Stack](#technology-stack)
*   [Configuration](#configuration)
*   [Getting Started](#getting-started)
*   [Usage](#usage)
*   [API Endpoints](#api-endpoints)
*   [RSS Feed](#rss-feed)
*   [Contributing](#contributing)
*   [License](#license)

## Features

*   **Hacker News Scraping:** Periodically scrapes Hacker News top stories and front page articles.
*   **AI Summarization:** Utilizes Ollama and Instructor AI to generate concise, structured summaries of article content.
    *   Default Model: `llama3:instruct`
*   **Web Interface:** Clean, responsive frontend built with Next.js and Material UI (Joy UI) to display articles and summaries. Includes light/dark mode toggle.
*   **RESTful API:** Go-based API for retrieving articles with support for:
    *   Pagination (`limit`, `offset`)
    *   Filtering by status (`flagged`, `dead`, `dupe`)
    *   Filtering by engagement thresholds (`min_upvotes`, `min_comments`)
*   **API Documentation:** Interactive Swagger UI available for exploring API endpoints.
*   **RSS Feed:** Provides RSS feeds of articles, supporting filtering via query parameters similar to the API.
*   **Dockerized:** Fully containerized setup using Docker and Docker Compose for easy development and deployment.
*   **CI/CD:** GitHub Actions workflow for automated building, testing, and pushing Docker images to Docker Hub.
*   **Database Schema:** Defined schema for storing article details, content, summaries, and metadata (e.g., commit hash, model name).
*   **Configurable:** Uses environment variables (`.env` file) for easy configuration.

## Architecture Overview

The project follows a microservices architecture orchestrated using Docker Compose:

1.  **Frontend:** A Next.js (React/TypeScript) application serving the user interface (the live version is at [https://gophersignal.com/](https://gophersignal.com/)). It fetches data from the Backend API. Built using Static Site Generation (`output: 'export'`).
2.  **Backend:** A Go application providing a RESTful API (documented with Swagger) to query stored articles from the database. Implements caching for optimized responses.
3.  **Hackernews Scraper:** A Node.js (TypeScript) service responsible for scraping Hacker News (top stories and front page for flagged/dead/dupe articles) using Puppeteer. It fetches full article content and utilizes an AI model via Ollama (interfaced with Instructor) to generate summaries. Results are stored in the MySQL database.
4.  **RSS:** A Rust service that generates configurable RSS feeds by fetching article data from the Backend API.
5.  **Ollama:** Hosts a local large language model (LLM) for generating article summaries. Defaults to `llama3:instruct`.
6.  **MySQL:** The relational database used for storing article data and summaries.
7.  **Nginx:** Acts as a reverse proxy, routing requests to the appropriate service (Frontend, Backend API, RSS) and handling SSL termination in production.


## Technology Stack

*   **Frontend:** Next.js, React, TypeScript, Material UI (Joy UI), Zod
*   **Backend API:** Go, Gorilla Mux, MySQL Driver, Swaggo
*   **Scraper:** Node.js, TypeScript, Puppeteer, Puppeteer-Extra (Stealth), Instructor AI, Ollama, Zod, MySQL2
*   **RSS Feed:** Rust, Axum, Reqwest, RSS crate, SQLx (implied via dependencies, though direct DB access might not be used in the final RSS implementation)
*   **Database:** MySQL 8.0
*   **LLM Service:** Ollama
*   **Reverse Proxy:** Nginx
*   **Orchestration:** Docker, Docker Compose
*   **CI/CD:** GitHub Actions
*   **Task Runner:** Make

## Configuration

The application uses environment variables for configuration. A template file `.env.example` is provided. Copy it to `.env` and adjust the values as needed.

Key default configuration values:

*   **Environment:**
    *   `NEXT_PUBLIC_ENV`: `development`
    *   `GO_ENV`: `development`
*   **Backend API:**
    *   `SERVER_ADDRESS`: `0.0.0.0:8080`
    *   `CACHE_MAX_AGE`: `5400` (seconds, 1.5 hours - Go backend default fallback if not set or invalid in `.env`)
*   **MySQL:**
    *   `MYSQL_HOST`: `mysql` (Service name in Docker Compose)
    *   `MYSQL_PORT`: `3306`
    *   `MYSQL_DATABASE`: `gophersignal`
    *   `MYSQL_USER`: `user`
    *   `MYSQL_PASSWORD`: `password`
    *   `MYSQL_ROOT_PASSWORD`: `password`
*   **Ollama (Scraper):**
    *   `OLLAMA_BASE_URL`: `http://ollama:11434/api`
    *   `OLLAMA_MODEL`: `llama3:instruct`
    *   `OLLAMA_CONTEXT_LENGTH`: `8192`
    *   `MAX_CONTENT_LENGTH`: `2000` (Characters used for generating summary)
    *   `MAX_SUMMARY_LENGTH`: `500` (Tokens for summary output)
*   **RSS:**
    *   `RSS_PORT`: `9090`
    *   `API_URL`: `http://backend:8080/api/v1/articles` (Default for RSS service to fetch data)
*   **GitHub (Scraper Metadata):**
    *   `GITHUB_OWNER`: `k-zehnder`
    *   `GITHUB_REPO`: `gophersignal`
    *   `GITHUB_BRANCH`: `main`
    *   `GH_TOKEN`: *Not set by default* (Required for reliable commit hash fetching via GitHub API in the scraper)
*   **API Query Defaults:**
    *   `limit`: `30`
    *   `offset`: `0`
    *   `min_upvotes`: `0`
    *   `min_comments`: `0`
*   **Scraper Behavior:**
    *   Max Scraped Top Story Pages: `2` (Hardcoded in `workflow.ts`)
    *   Max Articles Summarized (non-flagged): `30` (Hardcoded in `workflow.ts`)

## Getting Started

These instructions guide you through setting up the project locally using Docker Compose for development.

1.  **Prerequisites:**
    *   Git
    *   Docker
    *   Docker Compose
    *   Make (optional, simplifies commands)

2.  **Clone the Repository:**
    ```bash
    git clone https://github.com/k-zehnder/gophersignal.git
    cd gophersignal
    ```

3.  **Set Up Environment:**
    Copy the example environment file and update any necessary values (especially if you need specific API keys or non-default configurations). The defaults are generally suitable for local development.
    ```bash
    cp .env.example .env
    # Optional: Edit .env if needed (e.g., provide GH_TOKEN)
    ```

4.  **Start Development Environment:**
    This command builds the development Docker images (if they don't exist) and starts all services defined in `docker-compose-dev.yml`.
    ```bash
    make dev
    ```
    Alternatively, without Make:
    ```bash
    docker compose -f docker-compose-dev.yml up -d --build
    ```
    Wait for all services to start, especially `mysql` and `ollama` which have health checks. The Ollama container will also pull the default model (`llama3:instruct`) upon first startup, which may take some time.

5.  **Run the Scraper:**
    After the environment is up, manually trigger the scraper to populate the database with Hacker News articles and summaries.
    ```bash
    make scrape
    ```
    Alternatively, without Make:
    ```bash
    docker compose -f docker-compose-dev.yml run --rm hackernews_scraper npm run start
    ```
    *Note: If the Ollama model pull is still in progress when you first run the scraper, it might fail. Wait a few minutes and try again.*

6.  **Access the Application:**
    Once the scraper has run successfully, you can access the different parts of the application:
    *   **Frontend:** [http://localhost:3000](http://localhost:3000) (Proxied via Nginx at `http://localhost`)
    *   **API Documentation (Swagger UI):** [http://localhost/swagger/index.html#/](http://localhost/swagger/index.html#/)
    *   **RSS Feed:** [http://localhost/rss](http://localhost/rss)

## Usage

*   **Browse Articles:** Navigate to the frontend URL (`http://localhost` by default in dev, or [https://gophersignal.com/](https://gophersignal.com/) for the live version) to view the latest summarized articles.
*   **Use the API:** Interact with the backend API endpoints, documented via Swagger UI at `/swagger/index.html#/` (or `https://gophersignal.com/swagger/index.html#/` for the live API).
*   **Access RSS:** Use an RSS reader to subscribe to the feed available at `/rss` (or `https://gophersignal.com/rss` for the live feed). Filters can be applied via query parameters (see below).
*   **Run Scraper:** Execute `make scrape` (or the corresponding Docker command) whenever you want to fetch and process new articles.

## API Endpoints

The primary API endpoint is provided by the Go backend:

*   `GET /api/v1/articles`: Retrieves a list of articles.
    *   **Query Parameters:**
        *   `limit` (integer, default: `30`, max: `100`): Number of articles per page.
        *   `offset` (integer, default: `0`): Pagination offset.
        *   `flagged` (boolean, optional): Filter by flagged status.
        *   `dead` (boolean, optional): Filter by dead status.
        *   `dupe` (boolean, optional): Filter by duplicate status.
        *   `min_upvotes` (integer, default: `0`): Minimum upvotes threshold.
        *   `min_comments` (integer, default: `0`): Minimum comments threshold.

Refer to the [Swagger UI](http://localhost/swagger/index.html#/) for detailed request/response schemas and interactive testing.

## RSS Feed

The Rust service provides an RSS feed at `/rss`. You can customize the feed content using query parameters:

*   `/rss`: Default feed (latest non-flagged, non-dead, non-dupe articles).
*   `/rss?flagged=true`: Feed containing only flagged articles.
*   `/rss?dead=true`: Feed containing only dead articles.
*   `/rss?dupe=true`: Feed containing only duplicate articles.
*   `/rss?min_upvotes=50`: Feed containing articles with at least 50 upvotes.
*   `/rss?min_comments=10`: Feed containing articles with at least 10 comments.

Parameters can be combined (e.g., `/rss?flagged=true&min_upvotes=100`).

## Contributing

Contributions are welcome! Please follow these general guidelines:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes, adhering to the existing coding style.
4.  Ensure tests pass (`make test` in relevant subdirectories or root).
5.  Update documentation if necessary.
6.  Submit a Pull Request using the provided template (`.github/pull_request_template.md`).

## License

MIT
