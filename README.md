# SafelyYou Monitoring API

A lightweight Go REST API service for monitoring and tracking device's heartbeats and upload statistics. The service calculates device uptime metrics and average upload durations for registered devices.

## Prerequisites

- Go 1.24.3 or higher
- Gorilla Mux router (`github.com/gorilla/mux`)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Dewanshi08/safelyyou-monitoring-fleet.git
cd safelyyou-monitoring-fleet
```

2. Install dependencies:
```bash
go mod init safelyyou-monitoring-fleet
go get github.com/gorilla/mux
```

## Running the Application

Execute the following command from the project root:
```bash
go run main.go
```

The server will start listening on port `6733`


## API Endpoints

### 1. Register Heartbeat
**Endpoint:** `POST /api/v1/devices/{device_id}/heartbeat`

Registers a heartbeat signal from a device. These heartbeats are used to calculate device uptime.

**Request Body:**
```json
{
    "sent_at": "2024-01-20T10:30:00Z"
}
```

**Response:**
- `204 No Content` - Heartbeat registered successfully
- `400 Bad Request` - Invalid JSON body
- `404 Not Found` - Device not found
- `500 Internal Server Error` - Error Response

### 2. Register Stats
**Endpoint:** `POST /api/v1/devices/{device_id}/stats`

Registers the statistics of device with actual device's upload time. These upload times are used to calculate the average upload time of a device.

**Request Body:**
```json
{
    "sent_at": "2024-01-20T10:30:00Z",
    "upload_time": 1500000000
}
```

**Response:**
- `204 No Content` - Stats registered successfully
- `400 Bad Request` - Invalid JSON body
- `404 Not Found` - Device not found
- `500 Internal Server Error` - Error Response

### 3. Get Stats
**Endpoint:** `GET /api/v1/devices/{device_id}/stats`

Retrieves aggregated statistics for a device including uptime and average upload time.

**Response:**
```json
{
    "uptime": 98.75000,
    "avg_upload_time": "3m17.331667813s"
}
```

**Response Codes:**
- `200 OK` - Stats retrieved successfully
- `400 Bad Request` - [ No upload times found, Invalid data: first and last heartbeat has same timestamp]
- `404 Not Found` - Device not found


## Running Tests

1. Run all tests from the project root:
```bash
go test ./...
```

2. Run tests with verbose output:
```bash
go test -v ./...
```

3. Run tests for a specific package:
```bash
go test -v ./handlers
```

4. Run tests with coverage report:
```bash
go test -v -cover ./...
```

## Running the Application Via Docker

1. Build the Docker image:
```bash
docker build -t monitoring .
```

2. Run the container:
```bash
docker run -p 6733:6733 monitoring
```

3. Verify the container is running:
```bash
docker ps
```

4. Stop the running container:
```bash
docker stop <CONTAINER_ID>
```