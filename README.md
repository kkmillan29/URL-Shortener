# URL Shortener Service

A simple URL shortener service built with Go that provides REST APIs for shortening URLs, redirecting to original URLs, and viewing metrics.

## Features

- Shorten URLs via REST API
- Redirect from shortened URLs to original URLs
- Track and display metrics for the top 3 most shortened domains
- In-memory storage for URL mappings
- Docker support for easy deployment

## API Endpoints

### 1. Shorten URL

```
POST /shorten
```

**Request Body:**
```json
{
  "url": "https://www.example.com"
}
```

**Response:**
```json
{
  "original_url": "https://www.example.com",
  "short_url": "http://localhost:8080/abc123"
}
```

### 2. Redirect to Original URL

```
GET /{shortURL}
```

This endpoint redirects the user to the original URL associated with the provided short URL.

### 3. Get Metrics

```
GET /metrics
```

**Response:**
```json
{
  "top_domains": [
    {
      "domain": "example.com",
      "count": 5
    },
    {
      "domain": "google.com",
      "count": 3
    },
    {
      "domain": "github.com",
      "count": 2
    }
  ]
}
```

## Running the Application

### Prerequisites

- Go 1.21 or higher
- Docker (optional)

### Local Development

1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Run the application:
   ```
   go run main.go
   ```
4. The server will start on port 8080

### Using Docker

1. Build the Docker image:
   ```
   docker build -t urlshortener .
   ```
2. Run the container:
   ```
   docker run -p 8080:8080 urlshortener
   ```

## Testing

Run the tests with:

```
go test ./...
```

## Project Structure

- `main.go` - Application entry point
- `model/` - Data models and structures
- `service/` - Business logic for URL shortening
- `handler/` - HTTP handlers for API endpoints