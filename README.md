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

   > **Note:** The `ollama` service starts by running `ollama serve`, pulling the required model (`llama3.2`), and creating a readiness flag. Other services wait until `ollama` is fully ready before proceeding.

4. **Run the Scraper:**

   Populate the database by running the scraper:

   ```bash
   make scrape
   ```

5. **Access the Application:**

   - **Frontend:** [http://localhost:3000](http://localhost:3000)
   - **API Documentation (Swagger UI):** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## System Architecture

Here’s the flow of the GopherSignal system, illustrating user interactions with the web app and RSS feed:

![GopherSignal System Flow](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgR29waGVyU2lnbmFsIFN5c3RlbSBGbG93CgphY3RvciAiVXNlciAxIChCcm93c2VyKSIgYXMgAAYHCnBhcnRpY2lwYW50IE5naW54AAUNIkZyb250ZW5kIChOZXh0LmpzADgGAA8IAB8OQmFja2VuZCBBUEkgKEdvAGIHAA8GAEkORGF0YWJhc2UgKE15U1FMAIENBgANCAB0DkhhY2tlck5ld3MgU2NyYXBlciAoQmFja2dyb3VuZCBKb2IAgUsGABYHAIExDlJTUyBTZXJ2aWNlIChBeHVtAIF3BlJTUwCCEA0yIChSU1MgUmVhZACCGAhSU1MACQYKCgCCMQcgLT4AghwGOiBTZW5kcyByZXF1ZXN0IGZvciB3ZWIgVUkKAII9BSAtPgCCGAk6IFJvdXRlcyBVSQAoCAoAgksJLT4AghYIOiBRdWVyaWVzIEFQSQoAgkAILT4AggsJOgCBGAVzIGFuZCB3cml0ZXMgZGF0YQAfDABpCgCBGQZkYXRhIHRvAIMiCgBsDACBPA0AgTcGIGMAg2YFdCBiYWNrAIFACgCELwc6IERlbGl2ZXIAgUYFACMIdG8gAIRYBgoKAIMVCACBKA1Qb3B1bGF0ZXMgd2l0aCBuZXcgAIRrBWxlAIEZBihiAINACmpvYikKCgCCbwkgLT4gUlNTOiBSAIJhBgCDMQUgZmVlZApSU1MAgi8NRmV0Y2hlcwBODQCCOQwAQgUAgy4GAHENAIFkBW5vdGUgb3ZlcgBoBlRyYW5zZm9ybQCDBgUgSlNPTgCCQQkAfAVvcm1hdAB7CACEHQk6IFJldHVybgCBHAoAghMJMgoK&s=default)
