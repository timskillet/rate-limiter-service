# Rate Limiter

## Description

The Rate Limiter project is a Go-based implementation of a rate limiting system designed to control the number of requests a client can make to a service within a specified time window. This project utilizes the Tollbooth package for efficient rate limiting, allowing for both fixed window and more advanced sliding window algorithms.

### Features

- **Rate Limiting**: Control the rate of incoming requests to prevent abuse and ensure fair usage of resources.
- **Redis Integration**: Leverage Redis as a backend store for maintaining request counts and managing rate limits across distributed systems.
- **Flexible Configuration**: Easily configure the rate limit and time window for different clients or endpoints.
- **Error Handling**: Provides informative messages when requests exceed the allowed limit.

### Use Cases

- **API Rate Limiting**: Protect your APIs from excessive usage and ensure that all clients have fair access to resources.
- **Microservices**: Manage request rates in microservices architectures to prevent service overload.
- **User Management**: Implement user-specific rate limits to enhance user experience and system stability.

### Getting Started

To get started with the Rate Limiter project, clone the repository and install the necessary dependencies. You can then configure the rate limiter according to your requirements and integrate it into your application.

### Installation

```
go get github.com/adamryman/tollbooth
go get github.com/redis/go-redis/v9
go
```

### Example of creating a new TollboothLimiter

```
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

limiter := NewTollboothLimiter(redisClient, 10, time.Minute)

// Check if a request is allowed
if limiter.Allow("clientID") {
    // Proceed with the request
} else {
    // Handle rate limit exceeded
}
```

### Example of creating a new FixedWindowLimiter

```
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
limiter := NewFixedWindowLimiter(redisClient, 10, time.Minute)

// Check if a request is allowed
if limiter.Allow("clientID") {
    // Proceed with the request
} else {
    // Handle rate limit exceeded
}

```
