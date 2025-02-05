#!/bin/bash

# Start the Ollama server in the background
ollama serve &

# Wait for the Ollama server to be available
echo "Starting Ollama server..."
while ! curl -s http://localhost:11434 > /dev/null; do
  echo "Waiting for Ollama server to be ready..."
  sleep 5
done

# Pull the required Ollama model
echo "Pulling the required Ollama model..."
ollama pull llama3:instruct
echo "Model pulled successfully."

# Create a flag file to signal that the model has been pulled
touch /tmp/ollama_model_ready

# Keep the container running
tail -f /dev/null
