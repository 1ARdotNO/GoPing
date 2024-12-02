# GoPing

A lightweight REST API written in Go for performing ICMP ping tests. The API allows users to send GET requests to a `/ICMP` endpoint to check the reachability of a specified hostname.
This API is useful for monitoring and diagnosing network connectivity by performing ICMP ping tests to verify the reachability and response times of hosts. It can be integrated into monitoring tools, automated scripts, or dashboards to provide insights into network health and latency.
ie. 
- Bridgehead into an internal network to get ping responses without exposing services directly(uptime monitoring)
- Use with https://upptime.js.org/ to monitor hosts that do not have a TCP/HTTP endpoint already
## Features
- Perform ICMP ping tests to any hostname.
- Authenticate requests using an environment-defined API key.
- Return detailed responses, including average ping time (in milliseconds) or appropriate HTTP error codes for failures.

## Requirements
- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- Docker Compose (optional, for simplified multi-container management)
- A valid hostname to ping

## Environment Variables
- `AUTH_KEY`: A secret key used to authenticate API requests.

## Endpoints

### `GET /ICMP`
**Description**: Perform an ICMP ping to a specified hostname.

#### Query Parameters:
- `hostname` (required): The hostname or IP address to ping.
- `key` (required): The authentication key to verify the request.

#### Responses:
- **200 OK**:
    ```json
    {
      "hostname": "example.com",
      "avgPingTime": 12.34
    }
    ```
- **400 Bad Request**: Missing or invalid parameters.
    ```json
    {"error": "Missing hostname or authentication key"}
    ```
- **401 Unauthorized**: Invalid authentication key.
    ```json
    {"error": "Invalid authentication key"}
    ```
- **408 Request Timeout**: Ping timed out.
    ```json
    {"error": "Ping timed out"}
    ```
- **500 Internal Server Error**: Ping execution failed.
    ```json
    {"error": "Failed to execute ping", "details": "Some error message"}
    ```

## Usage

### Local Setup
1. Clone the repository:
    ```bash
    git clone https://github.com/your-repo/icmp-api.git
    cd icmp-api
    ```

2. Set the `AUTH_KEY` environment variable:
    ```bash
    export AUTH_KEY=my_secret_key
    ```

3. Run the application:
    ```bash
    go run main.go
    ```

4. Test the endpoint:
    ```bash
    curl "http://localhost:8080/ICMP?hostname=example.com&key=my_secret_key"
    ```

### Docker Setup
1. Build the Docker image:
    ```bash
    docker build -t icmp-api .
    ```

2. Run the container:
    ```bash
    docker run -d -p 8080:8080 --env AUTH_KEY=my_secret_key icmp-api
    ```

3. Test the endpoint:
    ```bash
    curl "http://localhost:8080/ICMP?hostname=example.com&key=my_secret_key"
    ```

### Docker Compose Setup
1. Create a `docker-compose.yml` file with the following content:
    ```yaml
    version: '3.8'

    services:
      icmp-api:
        build:
          context: .
          dockerfile: Dockerfile
        ports:
          - "8080:8080"
        environment:
          - AUTH_KEY=my_secret_key
        restart: always
    ```

2. Start the service:
    ```bash
    docker-compose up --build
    ```

3. Test the endpoint:
    ```bash
    curl "http://localhost:8080/ICMP?hostname=example.com&key=my_secret_key"
    ```

### Scaling with Docker Compose (Optional)
To scale the service horizontally:
```bash
docker-compose up --scale icmp-api=3