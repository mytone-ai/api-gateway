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

# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs (run this in api-gateway directory)
swag init
