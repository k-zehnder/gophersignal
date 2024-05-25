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
   MYSQL_ROOT_USER=root
   MYSQL_ROOT_PASSWORD=password
   MYSQL_PASSWORD=password
   MYSQL_DATABASE=gophersignal
   MYSQL_USER=user
   MYSQL_PORT=3306
   MYSQL_HOST=mysql

   # Third Party
   HUGGING_FACE_API_KEY=
   ```

   **Obtain API Keys:**

   Visit [Hugging Face](https://huggingface.co/) to create an account and obtain your API key. Set it as `HUGGING_FACE_API_KEY` in your `.env` file.

3. **Ensure Docker is Installed and Running:**

   Make sure Docker is installed and running on your host machine. You can download Docker from [here](https://www.docker.com/products/docker-desktop).

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
- **Swagger UI:** Access the API documentation at `http://localhost:8080/swagger/index.html`.
