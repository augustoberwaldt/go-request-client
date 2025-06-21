<p align="center">
  <a href="https://github.com/bozd4g/go-http-client">
    <img alt="go-http-client" src="https://req.cool/images/req.png" width="300">
  </a>
</p>


# Go HTTP Request Client Library

A powerful and flexible HTTP client library for Go, inspired by Guzzle PHP. This library provides a simple interface for building HTTP requests, handling responses, and integrating with web services.

## Features

- **Simple Interface**: Easy-to-use API for building HTTP requests
- **Multiple Request Types**: Support for GET, POST, PUT, DELETE, PATCH
- **JSON Support**: Built-in JSON request/response handling
- **Form Data**: Support for form-encoded data
- **Multipart Uploads**: File upload support with multipart/form-data
- **Authentication**: Basic authentication support
- **Query Parameters**: Easy query string building
- **Headers Management**: Flexible header configuration
- **Async Requests**: Asynchronous request support with promises
- **Concurrent Requests**: Send multiple requests concurrently
- **Middleware Support**: Extensible middleware system
- **Timeout Control**: Configurable request timeouts
- **Cookie Support**: HTTP cookie handling
- **Response Handling**: Convenient response methods

## Installation

```bash
go get github.com/http-client/httpclient
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/http-client/httpclient"
)

func main() {
    // Create a new client
    client := httpclient.NewClient(
        httpclient.WithBaseURL("https://api.example.com"),
        httpclient.WithTimeout(10*time.Second),
    )

    // Send a GET request
    resp, err := client.Get("/users", &httpclient.RequestOptions{
        QueryParams: map[string]string{
            "page":  "1",
            "limit": "10",
        },
        Headers: map[string]string{
            "User-Agent": "MyApp/1.0",
        },
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %d\n", resp.GetStatusCode())
    fmt.Printf("Body: %s\n", resp.GetBody())
}
```

### POST with JSON

```go
// Create data to send
data := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
}

// Send POST request with JSON
resp, err := client.Post("/users", &httpclient.RequestOptions{
    JSON: data,
})

if err != nil {
    log.Fatal(err)
}

// Parse JSON response
var result map[string]interface{}
if err := resp.UnmarshalJSON(&result); err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user with ID: %v\n", result["id"])
```

### Form Data

```go
resp, err := client.Post("/login", &httpclient.RequestOptions{
    FormData: map[string]string{
        "username": "john_doe",
        "password": "secret123",
    },
})
```

### File Upload (Multipart)

```go
// Create multipart data
multipartData := httpclient.NewMultipartData()
multipartData.AddField("description", "Profile picture")
multipartData.AddFileFromPath("file", "/path/to/image.jpg")

// Send multipart request
resp, err := client.Post("/upload", &httpclient.RequestOptions{
    Multipart: multipartData,
})
```

### Authentication

```go
// Create client with authentication
client := httpclient.NewClient(
    httpclient.WithBaseURL("https://api.example.com"),
    httpclient.WithAuth("username", "password"),
)

// Or set auth per request
resp, err := client.Get("/protected", &httpclient.RequestOptions{
    Auth: &httpclient.Auth{
        Username: "username",
        Password: "password",
    },
})
```

### Async Requests

```go
// Create async client
asyncClient := httpclient.NewAsyncClient(
    httpclient.WithBaseURL("https://api.example.com"),
)

// Send async request
promise := asyncClient.GetAsync("/users", nil)

// Do other work while request is processing
fmt.Println("Request sent, doing other work...")

// Wait for response
resp, err := promise.Wait()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Async response: %s\n", resp.GetBody())
```

### Using Promises

```go
promise := asyncClient.GetAsync("/users", nil)

promise.Then(func(resp *httpclient.Response) (*httpclient.Response, error) {
    fmt.Printf("Request succeeded: %d\n", resp.GetStatusCode())
    return resp, nil
}).Catch(func(err error) error {
    fmt.Printf("Request failed: %v\n", err)
    return err
})

promise.Wait()
```

### Concurrent Requests

```go
requests := []httpclient.ConcurrentRequest{
    {Method: "GET", Path: "/users", Options: &httpclient.RequestOptions{}},
    {Method: "GET", Path: "/posts", Options: &httpclient.RequestOptions{}},
    {Method: "GET", Path: "/comments", Options: &httpclient.RequestOptions{}},
}

results := asyncClient.SendConcurrent(requests)

for i, result := range results {
    if result.Error != nil {
        fmt.Printf("Request %d failed: %v\n", i, result.Error)
    } else {
        fmt.Printf("Request %d succeeded: %d\n", i, result.Response.GetStatusCode())
    }
}
```

### With Concurrency Limit

```go
// Send requests with a limit of 5 concurrent requests
results := asyncClient.SendConcurrentWithLimit(requests, 5)
```

## Client Configuration

### Client Options

```go
client := httpclient.NewClient(
    // Set base URL for all requests
    httpclient.WithBaseURL("https://api.example.com"),
    
    // Set default timeout
    httpclient.WithTimeout(30*time.Second),
    
    // Set default headers
    httpclient.WithHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
        "Accept":     "application/json",
    }),
    
    // Set default authentication
    httpclient.WithAuth("username", "password"),
)
```

### Request Options

```go
options := &httpclient.RequestOptions{
    // Query parameters
    QueryParams: map[string]string{
        "page":  "1",
        "limit": "10",
    },
    
    // Headers
    Headers: map[string]string{
        "Authorization": "Bearer token123",
    },
    
    // JSON data
    JSON: map[string]interface{}{
        "name": "John Doe",
    },
    
    // Form data
    FormData: map[string]string{
        "username": "john_doe",
    },
    
    // Custom body
    Body: strings.NewReader("custom body"),
    
    // Timeout for this request
    Timeout: 5 * time.Second,
    
    // Authentication for this request
    Auth: &httpclient.Auth{
        Username: "user",
        Password: "pass",
    },
    
    // Cookies
    Cookies: []*http.Cookie{
        {Name: "session", Value: "abc123"},
    },
    
    // Allow redirects
    AllowRedirects: true,
}
```

## Response Handling

```go
resp, err := client.Get("/users", nil)
if err != nil {
    log.Fatal(err)
}

// Get status code
status := resp.GetStatusCode()

// Get specific header
contentType := resp.GetHeader("Content-Type")

// Get response body as string
body := resp.GetBody()

// Get response body as bytes
bodyBytes := resp.GetBodyBytes()

// Parse JSON response
var users []User
if err := resp.UnmarshalJSON(&users); err != nil {
    log.Fatal(err)
}
```

## Middleware

The library includes a middleware system for extending client behavior:

```go
// Logging middleware
logger := &httpclient.SimpleLogger{}
loggingMiddleware := httpclient.LoggingMiddleware(logger)

// Retry middleware with exponential backoff
backoff := &httpclient.ExponentialBackoff{
    BaseDelay: 1 * time.Second,
    MaxDelay:  30 * time.Second,
}
retryMiddleware := httpclient.RetryMiddleware(3, backoff)

// Timeout middleware
timeoutMiddleware := httpclient.TimeoutMiddleware(5 * time.Second)
```

## Error Handling

```go
resp, err := client.Get("/users", nil)
if err != nil {
    // Handle network errors, timeouts, etc.
    log.Printf("Request failed: %v", err)
    return
}

// Check HTTP status code
if resp.GetStatusCode() >= 400 {
    log.Printf("HTTP error: %d - %s", resp.GetStatusCode(), resp.GetBody())
    return
}
```

## Testing

Run the tests:

```bash
go test ./...
```

## Examples

See the `examples/` directory for comprehensive usage examples.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Comparison with Guzzle PHP

This Go library provides similar functionality to Guzzle PHP:

| Feature | Guzzle PHP | Go HTTP Client |
|---------|------------|----------------|
| Simple interface | ✅ | ✅ |
| JSON support | ✅ | ✅ |
| Form data | ✅ | ✅ |
| Multipart uploads | ✅ | ✅ |
| Authentication | ✅ | ✅ |
| Query parameters | ✅ | ✅ |
| Headers management | ✅ | ✅ |
| Async requests | ✅ | ✅ |
| Concurrent requests | ✅ | ✅ |
| Middleware | ✅ | ✅ |
| Timeout control | ✅ | ✅ |
| Cookie support | ✅ | ✅ |
| Response handling | ✅ | ✅ |

The Go library maintains the same ease of use while leveraging Go's strengths in concurrency and performance. 
