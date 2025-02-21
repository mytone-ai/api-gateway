# api-gateway

This microservice is part of the MyTone AI project.

## Getting Started
1. Build the Docker image:
   ```
   docker build -t mytone-ai/api-gateway .
   ```
2. Run the container:
   ```
   docker run -p 8080:8080 mytone-ai/api-gateway
   ```

## API
- `GET /health` - Health check endpoint.
