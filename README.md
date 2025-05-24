# Media Management System with Multiple Providers

A scalable media management system built with Go (Fiber) and Next.js, featuring support for multiple storage providers and a modern web interface.

## Features

- Multiple storage provider support:
  - Azure Blob Storage
  - Discord CDN
  - Amazon S3
  - Firebase Storage
  - Local Storage
- Unified media management interface
- Authentication and authorization
- File upload and management
- Notification system
- Caching support

## Technologies Used

### Backend
- **Language:** Go
- **Framework:** Fiber
- **Database:** PostgreSQL
- **Cache:** Redis
- **Authentication:** Firebase Authentication
- **API Documentation:** Swagger (OpenAPI)
- **Configuration:** Viper
- **Observability:** OpenTelemetry
- **Containerization:** Docker, Docker Compose

### Frontend
- **Framework:** Next.js
- **Language:** TypeScript
- **Package Manager:** Bun
- **UI Components:** shadcn/ui
- **Styling:** Tailwind CSS
- **State Management:** Zustand
- **Forms:** React Hook Form, Zod
- **API Client:** Axios

## System Requirements

- Go 1.20 or higher
- Node.js (latest LTS)
- Bun (latest)
- Docker and Docker Compose
- Make (optional)

## Installation Guide

1. **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd media-management-multiple-providers
    ```

2. **Configure the environment:**
    ```bash
    cp config/config.example.yaml config/config.yaml
    ```
    Edit `config.yaml` with your storage provider credentials and other settings.

3. **Start dependent services:**
    ```bash
    docker-compose up -d
    ```

4. **Install frontend dependencies:**
    ```bash
    cd next
    bun install
    cd ..
    ```

5. **Generate Swagger documentation (optional):**
    ```bash
    make swag
    ```

## Usage Guide

### Running the Backend Server
```bash
make run
```
The backend runs on the port specified in `config.yaml` (default: `8080`).  
Swagger UI: `http://localhost:8080/swagger/index.html`

### Running the Frontend Application
```bash
cd next
bun dev
```
Frontend runs on port `3000` by default.

## Project Structure

```
.
├── cmd/
│   └── server/           # Main application entry point
├── config/              # Configuration files
├── docs/               # API documentation and storage provider docs
├── internal/
│   ├── adapters/       # Storage provider implementations
│   │   ├── azure/      # Azure Blob Storage adapter
│   │   ├── discord/    # Discord CDN adapter
│   │   ├── firebase/   # Firebase Storage adapter
│   │   ├── local/      # Local storage adapter
│   │   └── s3/         # Amazon S3 adapter
│   ├── dependencies/   # Application dependencies setup
│   ├── domain/        # Core domain models/logic
│   ├── infra/         # Infrastructure setup
│   ├── interfaces/    # HTTP handlers and validators
│   ├── modules/       # Feature modules
│   │   ├── media/     # Media management
│   │   ├── notify/    # Notification system
│   │   ├── storage/   # Storage abstraction
│   │   └── user/      # User management
│   ├── router/        # API routes
│   └── shared/        # Shared utilities
├── next/              # Frontend Next.js application
├── uploads/           # Local storage directory
└── docker-compose.yml # Docker services configuration
```

## Storage Provider Configuration

See [docs/storage-providers.md](docs/storage-providers.md) for detailed configuration instructions for each storage provider.

## API Reference

API documentation is available at:  
`http://localhost:8080/swagger/index.html`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[License Type] - Please specify your license
