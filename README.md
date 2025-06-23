# Go Rate Limiter

This project is a rate limiter implemented in Go, designed to control the number of requests to a web server based on IP address or an access token.

## Features

- Rate limiting by IP address.
- Rate limiting by access token.
- Configurable rate limits and block times.
- Redis-based storage for tracking request counts.
- Middleware-based integration with web servers.
- Pluggable storage backend through an interface.

## Configuration

The rate limiter is configured through environment variables or a `.env` file in the root directory.

- `REDIS_ADDR`: The address of the Redis server (e.g., `localhost:6379`).
- `RATE_LIMIT_BY_IP`: The maximum number of requests per second for a single IP address.
- `RATE_LIMIT_BY_TOKEN`: The default maximum number of requests per second for an access token.
- `BLOCK_TIME_IN_SECONDS`: The duration in seconds to block an IP or token after the rate limit is exceeded.
- `TOKEN_<token>_RATE_LIMIT`: The specific rate limit for a given token. For example, `TOKEN_abc123_RATE_LIMIT=100` sets a rate limit of 100 requests per second for the token `abc123`.

## How to Run

1.  **Clone the repository:**
    ```sh
    git clone git@github.com:AugustoSandim/rate-limiter-go.git
    cd rate-limiter-go
    ```

2.  **Create a `.env` file:**
    Create a `.env` file in the root of the project and add the configuration variables as described above.

3.  **Run with Docker Compose:**
    ```sh
    docker-compose up --build
    ```

The server will be running on port `8080`.

## How to Use

-   **IP-based limiting:** Requests are automatically limited by the client's IP address.
-   **Token-based limiting:** Include the access token in the `API_KEY` header of your request.

    ```
    API_KEY: <your-token>
    ```

If the rate limit is exceeded, the server will respond with a `429 Too Many Requests` status code and the message: `you have reached the maximum number of requests or actions allowed within a certain time frame`.
