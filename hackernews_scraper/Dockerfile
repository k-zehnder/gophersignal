# Use an official Node runtime as the base image
FROM node:latest

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Bundle app source
COPY . .

# Keep the container running without doing anything
# This facilitates manually triggering the scraper as needed
CMD ["tail", "-f", "/dev/null"]
