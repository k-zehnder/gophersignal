# Use an official Node runtime as the base image
FROM node:latest

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install Node dependencies and additional system dependencies needed for scraping
RUN npm install ts-node && apt-get update && apt-get install -y \
    libnss3 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxrandr2 \
    libgbm1 \
    libasound2 \
    libpangocairo-1.0-0 \
    libpango-1.0-0 \
    libcairo2 \
    libatspi2.0-0 \
    libgtk-3-0 \
    libdbus-1-3

# Bundle app source
COPY . .

# Keep the container running without doing anything
# This facilitates manually triggering the scraper as needed
CMD ["tail", "-f", "/dev/null"]
