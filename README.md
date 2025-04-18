# GopherSignal

[![CI/CD Pipeline](https://github.com/k-zehnder/gophersignal/actions/workflows/workflow.yml/badge.svg)](https://github.com/k-zehnder/gophersignal/actions/workflows/workflow.yml)

## Quickstart

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/k-zehnder/gophersignal.git
   cd gophersignal
   ```

2. **Set Up Environment:**

   Copy the example environment file and update any necessary values:

   ```bash
   cp .env.example .env
   ```

3. **Start Development Environment:**

   Build and start all services:

   ```bash
   make dev
   ```

4. **Run the Scraper:**

   Populate the database by running the scraper:

   ```bash
   make scrape
   ```

5. **Access the Application:**

   - **Frontend:** [http://localhost:3000](http://localhost:3000)
   - **API Documentation (Swagger UI):** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
