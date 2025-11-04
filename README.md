# S3 Clone Golang

A lightweight, S3-compatible object storage service built with Go, featuring a modern web interface for file management.

## Features

- **S3-Compatible API**: RESTful API endpoints that mimic AWS S3 functionality
- **Bucket Management**: Create, list, and delete buckets
- **Object Operations**: Upload, download, list, and delete files
- **Web Interface**: Modern, responsive web UI for easy file management
- **Health Monitoring**: Built-in health check endpoint
- **File Type Detection**: Automatic file extension handling
- **Drag & Drop**: Intuitive file upload via web interface

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check and API documentation |
| `GET` | `/buckets` | List all buckets |
| `POST` | `/buckets` | Create a new bucket |
| `DELETE` | `/buckets/{bucket}` | Delete a bucket |
| `PUT` | `/upload/{bucket}` | Upload a file to bucket |
| `GET` | `/download/{bucket}/{key}` | Download a file |
| `GET` | `/list/{bucket}` | List files in bucket |
| `DELETE` | `/delete/{bucket}/{key}` | Delete a file |

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd s3-clone-go
   ```

2. **Install dependencies**
   ```bash
   cd backend
   go mod download
   ```

3. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

4. **Access the web interface**
   Open your browser and navigate to `http://localhost:3000`


## Configuration

### Environment Variables

- `PORT`: Server port (default: 3000)

### Storage Configuration

Files are stored in the `backend/uploads/` directory by default. Each bucket creates a subdirectory within this path.

## Usage Examples

### Using the Web Interface

1. Open `http://localhost:3000` in your browser
2. Create a new bucket using the sidebar
3. Select a bucket to view its contents
4. Drag and drop files or click to upload
5. Manage files with download and delete options

### Using the API

#### Create a bucket
```bash
curl -X POST http://localhost:3000/buckets \
  -H "Content-Type: application/json" \
  -d '{"name": "my-bucket"}'
```

#### Upload a file
```bash
curl -X PUT http://localhost:3000/upload/my-bucket \
  -F "file=@example.txt" \
  -F "key=documents/example.txt"
```

#### List bucket contents
```bash
curl http://localhost:3000/list/my-bucket
```

#### Download a file
```bash
curl http://localhost:3000/download/my-bucket/documents/example.txt \
  -o downloaded-file.txt
```

#### Delete a file
```bash
curl -X DELETE http://localhost:3000/delete/my-bucket/documents/example.txt
```

## Dependencies

- **[Gorilla Mux](https://github.com/gorilla/mux)**: HTTP router and URL matcher
- **[godotenv](https://github.com/joho/godotenv)**: Environment variable loader
- **[CORS](https://github.com/rs/cors)**: Cross-origin resource sharing middleware

## Project Structure

```
s3-clone-go/
├── backend/
│   ├── cmd/server/main.go           # Application entry point
│   ├── internal/
│   │   ├── handlers/s3.go           # S3 API handlers
│   │   ├── middleware/cors.go       # CORS middleware
│   │   └── services/                # Storage service logic
│   ├── uploads/                     # File storage
│   ├── go.mod                       # Go module definition
│   └── go.sum                       # Dependency checksums
└── public/
    ├── index.html                   # Web interface
    ├── app.js                       # Frontend JavaScript
    └── styles.css                   # Styling
```

## Features in Detail

### Bucket Operations
- Create buckets with alphanumeric names (including dashes and underscores)
- List all available buckets
- Delete empty buckets
- Automatic bucket validation

### File Operations
- Upload files with multipart form data
- Custom key assignment or automatic filename usage
- File type detection and extension handling
- Stream-based file serving for downloads
- Recursive file listing within buckets

### Web Interface
- Responsive design for desktop and mobile
- Real-time server status monitoring
- Drag-and-drop file uploads
- File management with preview and actions
- Toast notifications for user feedback

## Network Access

The server binds to `0.0.0.0`, making it accessible from other machines on your network. Access it using:

- Local: `http://localhost:3000`
- Network: `http://YOUR_IP_ADDRESS:3000`

## Error Handling

The API returns appropriate HTTP status codes and JSON error responses:

- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Bucket or file not found
- `409 Conflict`: Bucket already exists or not empty when deleting
- `500 Internal Server Error`: Server-side errors

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Roadmap

- [ ] Authentication and authorization
- [ ] Metadata storage for objects
- [ ] Versioning support
- [ ] Compression options
- [ ] Docker containerization
- [ ] Kubernetes deployment manifests
- [ ] Metrics and monitoring
- [ ] Backup and replication features
