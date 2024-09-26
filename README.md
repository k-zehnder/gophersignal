# gophersignal

[![CI/CD Pipeline](https://github.com/k-zehnder/gophersignal/actions/workflows/workflow.yml/badge.svg)](https://github.com/k-zehnder/gophersignal/actions/workflows/workflow.yml)

## Quickstart

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/k-zehnder/gophersignal.git
   cd gophersignal
   ```

2. **Configure Environment:**

   Copy the example environment file:

   ```bash
   cp .env.example .env
   ```

   **Edit the `.env` file:**

   Open the `.env` file and set your environment variables:

   ```dotenv
   # Frontend
   NEXT_PUBLIC_ENV=development

   # Backend
   GO_ENV=development
   SERVER_ADDRESS=0.0.0.0:8080

   # MySQL
   MYSQL_HOST=mysql
   MYSQL_PORT=3306
   MYSQL_DATABASE=gophersignal
   MYSQL_USER=user
   MYSQL_PASSWORD=password
   MYSQL_ROOT_PASSWORD=password

   # Ollama Configuration
   OLLAMA_BASE_URL=http://localhost:11434/api/generate
   OLLAMA_MODEL=llama3:instruct

   # Additional Configuration
   MAX_CONTENT_LENGTH=5000
   ```

   **Set Up Ollama:**

   The project uses [Ollama](https://ollama.ai/) to summarize articles. Follow these steps to set up Ollama:

   - **Install Ollama:**

     Visit the [Ollama installation page](https://ollama.ai/download) and download the installer for your operating system. Alternatively, for macOS users with Homebrew:

     ```bash
     brew install ollama/tap/ollama
     ```

   - **Start the Ollama Server:**

     Ollama runs as a local server. Start it by running:

     ```bash
     ollama serve
     ```

   - **Download the Required Model:**

     The default model used is `llama3:instruct`. Pull the model using:

     ```bash
     ollama pull llama3:instruct
     ```

     If you specified a different model in your `.env` file under `OLLAMA_MODEL`, make sure to pull that model instead.

   **Note:** Ensure that Ollama is running whenever you run the scraper or the application components that require summarization.

3. **Ensure Docker is Installed and Running:**

   Make sure Docker is installed and running on your host machine. You can download Docker Desktop from [here](https://www.docker.com/products/docker-desktop).

   Alternatively, you can install Docker via the command line:

   For **Ubuntu**:

   ```bash
   sudo apt-get update
   sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
   sudo apt-get update
   sudo apt-get install -y docker-ce
   sudo systemctl status docker
   ```

   For **Mac**:

   ```bash
   brew install docker
   brew install docker-compose
   ```

   For **Windows**:

   You can download Docker Desktop for Windows from [here](https://www.docker.com/products/docker-desktop) and follow the installation instructions provided on the website.

4. **Launch Development Environment with Docker:**

   This will build and start all necessary services:

   ```bash
   make dev
   ```

5. **Populate the Database with Data by Running the Scraper:**

   ```bash
   cd hackernews_scraper
   make scrape
   cd ..
   ```

   Your development environment should now be running.

## Accessing the Application

- **Frontend:** Visit `http://localhost:3000` to view the frontend.
- **Swagger UI:** Access the API documentation at `http://localhost:8080/swagger/index.html`
