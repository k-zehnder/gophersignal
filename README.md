# gophersignal

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

   # Third Party
   HUGGING_FACE_API_KEY=
   ```

   **Obtain API Key:**

   To summarize articles, you need a Hugging Face API key. Visit [Hugging Face](https://huggingface.co/) to create an account and obtain your API key. Then, add it to your `.env` file:

   ```dotenv
   HUGGING_FACE_API_KEY=your-huggingface-api-key
   ```

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
