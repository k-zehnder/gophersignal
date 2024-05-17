# gophersignal

## Quickstart

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/k-zehnder/gophersignal.git
   cd gophersignal
   ```

2. **Configure Environment:**
   Copy the example environment file and edit it with your preferences:

   ```bash
   cp .env.example .env
   ```

   **Obtain API Keys:**
   Visit [Hugging Face](https://huggingface.co/) to create an account and obtain your API key. Set it as `HUGGING_FACE_API_KEY` in your `.env` file.

3. **Launch Development Environment with Docker:**
   This will build and start all necessary services, including setting up the database:

   ```bash
   make dev
   ```

Your development environment should now be running.

## Accessing the Application

- **Frontend:** Visit `http://localhost:3000` to view the frontend.
- **Swagger UI:** Access the API documentation at `http://localhost:8080/swagger/index.html`.
