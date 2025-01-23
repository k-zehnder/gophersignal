#!/bin/bash

# Start the Ollama server
echo "Starting Ollama server..."
ollama serve &

# Wait for the Ollama server to be available
until curl -s http://localhost:11434 > /dev/null; do
  echo "Waiting for Ollama server..."
  sleep 5
done

echo "Ollama server is running."

# Pull the required Ollama model
echo "Pulling Ollama model..."
if ! ollama pull llama3:instruct; then
  echo "Failed to pull Ollama model. Exiting."
  exit 1
fi

# Run the Ollama model
echo "Running Ollama model..."
if ! ollama run llama3:instruct; then
  echo "Failed to run Ollama model. Exiting."
  exit 1
fi

# Keep the container running
echo "Ollama is ready. Keeping the container alive..."
tail -f /dev/null
