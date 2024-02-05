# gophersignal 

## Getting Started

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
   To use third-party services, you need to obtain API keys. Follow these steps to get your API keys:
   - **Hugging Face API Key:** Visit [Hugging Face](https://huggingface.co/) to create an account and obtain your API key.
   - **OpenAI API Key:** Visit [OpenAI](https://openai.com/) to create an account and obtain your API key. Set it as `OPEN_AI_API_KEY` in your `.env` file.

3. **Launch Services with Docker:**
   ```bash
   make dev_env
   ```

4. **Access the MySQL Docker Container:**
   Find the running MySQL container ID using the following command:
   ```bash
   docker ps | grep mysql | awk '{print $1}'
   ```

5. **Access the MySQL Container:**
   Replace `<CONTAINER_ID>` with the ID from the previous command, then execute the following command to access the container's bash shell:
   ```bash
   docker exec -it <CONTAINER_ID> mysql -u user -p
   ```

6. **Create the Database and Tables:**
   Inside the MySQL shell, you can directly create the database and tables using the following SQL commands:
   ```sql
   CREATE DATABASE IF NOT EXISTS gophersignal;
   USE gophersignal;

   CREATE TABLE IF NOT EXISTS articles (
       id INT AUTO_INCREMENT PRIMARY KEY,
       title VARCHAR(255) NOT NULL,
       link VARCHAR(512) NOT NULL,
       content TEXT,
       summary VARCHAR(2000),
       source VARCHAR(100) NOT NULL,
       created_at TIMESTAMP NOT NULL,
       updated_at TIMESTAMP NOT NULL
   );
   ```

7. **Exit the MySQL Shell:**
   Once you have set up your database, you can exit the MySQL shell by typing:
   ```sql
   exit
   ```

Your development environment should now be running.

## Accessing the Application

- **Frontend:** Visit `http://localhost:3000` to view the frontend.
- **Swagger UI:** Access the API documentation at `http://localhost:8080/swagger/index.html`.
