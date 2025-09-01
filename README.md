# Media Management System with Multiple Providers

A enterprise-grade, scalable media management system built with Go (Fiber) and Next.js, featuring support for multiple cloud storage providers, comprehensive observability, and a modern responsive web interface.

## 🚀 Features

### Storage Providers
- **Cloud Storage**: Azure Blob Storage, AWS S3, Firebase Storage
- **S3-Compatible**: Cloudflare R2, Scaleway Object Storage, Backblaze B2, MinIO
- **Alternative**: Discord CDN, Local Storage
- **Unified API**: Single interface for all storage providers

### Core Features
- 🔐 **Authentication & Authorization**: Firebase Authentication integration
- 📁 **File Management**: Upload, download, delete, and organize media files
- 🔄 **Multi-Provider Support**: Seamlessly switch between storage providers
- 📊 **Observability**: OpenTelemetry integration with SigNoz monitoring
- 🚀 **Performance**: Redis caching and optimized file handling
- 🌐 **Internationalization**: Multi-language support (EN, VI)
- 📱 **Responsive Design**: Modern UI with shadcn/ui components
- 🔔 **Notification System**: Email and messaging notifications
- 🐳 **Containerized**: Full Docker support with Docker Compose

## 🛠 Technologies Used

### Backend
- **Language:** Go 1.23+
- **Framework:** Fiber (High-performance web framework)
- **Database:** PostgreSQL with PostGIS
- **Cache:** Redis
- **Authentication:** Firebase Authentication + JWT
- **API Documentation:** Swagger/OpenAPI 3.0
- **Configuration:** Viper
- **Observability:** OpenTelemetry + SigNoz
- **Containerization:** Docker & Docker Compose
- **ORM:** GORM with PostgreSQL driver

### Frontend
- **Framework:** Next.js 15.3+
- **Language:** TypeScript
- **Package Manager:** Bun (Ultra-fast JavaScript runtime)
- **UI Library:** shadcn/ui + Radix UI
- **Styling:** Tailwind CSS
- **State Management:** Zustand
- **Forms:** React Hook Form + Zod validation
- **HTTP Client:** Axios with React Query
- **Authentication:** Firebase Authentication
- **Icons:** Lucide React + Radix Icons

### DevOps & Infrastructure
- **Monitoring:** SigNoz (Open-source observability platform)
- **Database:** ClickHouse (for metrics and logs)
- **Build Tools:** Make, Docker Multi-stage builds
- **Cloud Deployment:** Google Cloud Platform ready

## 📋 System Requirements

- **Go:** 1.23+ (with modules enabled)
- **Node.js:** Latest LTS (20.x recommended)
- **Bun:** Latest version (for frontend package management)
- **Docker:** 24.x+ and Docker Compose v2
- **Make:** For build automation (optional but recommended)
- **Git:** For version control

## 🚀 Quick Start Guide

### 1. Clone and Setup
```bash
git clone <repository-url>
cd media-management-multiple-providers
```

### 2. Configuration
```bash
# Copy example configuration
cp config/config.example.yaml config/config.yaml

# Edit configuration with your settings
# Add your storage provider credentials, database settings, etc.
```

### 3. Start Infrastructure Services
```bash
# Start PostgreSQL, Redis, ClickHouse, and SigNoz
docker-compose up -d

# Wait for services to be healthy
docker-compose ps
```

### 4. Database Setup
```bash
# Run database migrations
make migrate

# Optional: Seed with test data
make seed-test
```

### 5. Install Frontend Dependencies
```bash
cd next
bun install
cd ..
```

### 6. Generate API Documentation (Optional)
```bash
make swag
```

## 🏃 Running the Application

### Start the Backend Server
```bash
# Development mode
make run

# Or directly with Go
go run cmd/server/main.go
```
- Backend runs on port specified in `config.yaml` (default: `8083`)
- Swagger UI: `http://localhost:8083/swagger/index.html`
- Health check: `http://localhost:8083/health`

### Start the Frontend Application
```bash
cd next
bun dev
```
- Frontend runs on `http://localhost:3033`
- Hot reload enabled for development

### Available Make Commands
```bash
make run          # Start the backend server
make client       # Start the frontend (Next.js)
make db-up        # Start PostgreSQL and Redis
make db-down      # Stop database services
make migrate      # Run database migrations
make seed         # Seed database with all data
make seed-test    # Seed with test data only
make build        # Build the Go application
make build-linux  # Build for Linux deployment
make swag         # Generate Swagger documentation
make gcp          # Deploy to Google Cloud Platform
```

## 📁 Project Architecture

Following **Clean Architecture** and **Domain-Driven Design** principles:

```
.
├── cmd/server/              # Application entry point
├── config/                  # Configuration files and examples
├── docs/                   # API documentation and provider guides
├── internal/               # Private application code
│   ├── adapters/          # Storage provider implementations
│   │   ├── azure/         # Azure Blob Storage
│   │   ├── discord/       # Discord CDN
│   │   ├── firebase/      # Firebase Storage
│   │   ├── local/         # Local file system
│   │   ├── minio/         # MinIO/S3-compatible
│   │   └── s3/            # AWS S3 and variants
│   ├── application/       # Application services and use cases
│   ├── infra/            # Infrastructure concerns
│   │   ├── cache/        # Redis caching
│   │   ├── config/       # Configuration management
│   │   ├── database/     # Database connections
│   │   ├── jwt/          # JWT authentication
│   │   └── tracer/       # OpenTelemetry tracing
│   ├── modules/          # Feature modules (DDD bounded contexts)
│   │   ├── app/          # Application-wide concerns
│   │   ├── auth/         # Authentication & authorization
│   │   ├── media/        # Media file management
│   │   └── storage/      # Storage provider abstraction
│   ├── presentation/     # API layer
│   │   ├── grpc/         # gRPC endpoints (future)
│   │   └── http/         # HTTP/REST endpoints
│   └── shared/           # Shared utilities and common code
│       ├── constants/    # Application constants
│       ├── errors/       # Error definitions
│       ├── utils/        # Utility functions
│       └── validator/    # Request validation
├── locales/              # Internationalization files
├── next/                 # Frontend Next.js application
│   ├── src/
│   │   ├── app/         # Next.js App Router pages
│   │   ├── components/  # Reusable UI components
│   │   ├── contexts/    # React contexts
│   │   ├── lib/         # Frontend utilities
│   │   └── services/    # API service clients
│   └── public/          # Static assets
├── test/                 # Test files and fixtures
├── uploads/              # Local storage directory
└── vendor/               # Go dependencies (vendored)
```

## 🔧 Configuration

### Storage Provider Setup
Configure your preferred storage providers in `config/config.yaml`. See detailed guides:

- [Azure Blob Storage](docs/azure-provider.md)
- [AWS S3](docs/s3-provider.md)
- [Firebase Storage](docs/firebase-provider.md)
- [Cloudflare R2](docs/cloudflare-r2-provider.md)
- [Scaleway Object Storage](docs/scaleway-provider.md)
- [Backblaze B2](docs/backblaze-b2-provider.md)
- [MinIO](docs/minio-provider.md)
- [Discord CDN](docs/discord-provider.md)
- [Local Storage](docs/local-storage-provider.md)

Complete configuration reference: [Storage Providers Documentation](docs/storage-providers.md)

## 📊 Monitoring & Observability

The system includes comprehensive monitoring with **SigNoz**:

- **Metrics Dashboard**: `http://localhost:3301`
- **Distributed Tracing**: Full request tracing across services
- **Logs Aggregation**: Centralized logging with ClickHouse
- **Performance Monitoring**: Real-time performance metrics
- **Error Tracking**: Automatic error detection and alerting

## 📚 API Documentation

### Interactive API Explorer
- **Swagger UI**: `http://localhost:8083/swagger/index.html`
- **OpenAPI Spec**: Available in `docs/swagger.json` and `docs/swagger.yaml`

### Key Endpoints
- `POST /api/v1/auth/login` - User authentication
- `POST /api/v1/media/upload` - File upload to specified provider
- `GET /api/v1/media/list` - List uploaded files
- `DELETE /api/v1/media/{id}` - Delete media file
- `GET /health` - Health check endpoint

### Example Usage
```bash
# Login and get access token
TOKEN=$(curl -s http://localhost:8083/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | \
  jq -r '.access_token')

# Upload file to MinIO
curl -X POST http://localhost:8083/api/v1/media/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@example.jpg" \
  -F "provider=minio"
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific module tests
go test ./internal/modules/media/...

# Frontend tests
cd next
bun test
```

## 🐳 Docker Deployment

### Development
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Build
```bash
# Build production image
make build-linux

# Deploy to Google Cloud Platform
make gcp
```

## 🔐 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Firebase Integration**: Enterprise-grade authentication provider
- **Input Validation**: Comprehensive request validation
- **CORS Configuration**: Proper cross-origin resource sharing
- **Rate Limiting**: API rate limiting (configurable)
- **Secure Headers**: Security headers for web protection

## 🌐 Internationalization

Supported languages:
- **English** (`en.toml`)
- **Vietnamese** (`vi.toml`)

Add new languages by creating locale files in the `locales/` directory.

## 🚀 Performance Features

- **Redis Caching**: High-performance caching layer
- **Connection Pooling**: Optimized database connections
- **Async Processing**: Non-blocking file operations
- **CDN Integration**: Multiple CDN provider support
- **Optimized Builds**: Multi-stage Docker builds
- **Bun Runtime**: Ultra-fast JavaScript runtime for frontend

## 📈 Scalability

- **Horizontal Scaling**: Stateless design for easy scaling
- **Multiple Storage Backends**: Distribute load across providers  
- **Database Optimization**: Efficient queries with GORM
- **Microservices Ready**: Modular architecture for service extraction
- **Cloud Native**: Kubernetes and container orchestration ready

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Follow coding standards**: See [general.instructions.md](vscode-userdata:/Users/lugon/Library/Application%20Support/Code%20-%20Insiders/User/prompts/general.instructions.md)
4. **Write tests**: Ensure your changes are tested
5. **Commit changes**: `git commit -m 'Add amazing feature'`
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Development Guidelines
- Files should not exceed 500 lines
- Follow Clean Architecture principles
- Use appropriate naming conventions
- Add comprehensive comments
- Ensure test coverage

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Express-inspired web framework for Go
- [Next.js](https://nextjs.org/) - The React framework for production
- [SigNoz](https://signoz.io/) - Open-source observability platform
- [shadcn/ui](https://ui.shadcn.com/) - Beautiful UI components
- [Firebase](https://firebase.google.com/) - Authentication and storage
- All the amazing open-source libraries that make this project possible

## 📞 Support

- **Documentation**: Check the `docs/` directory
- **Issues**: [GitHub Issues](../../issues)
- **Discussions**: [GitHub Discussions](../../discussions)

---

<div align="center">
  <strong>Built with ❤️ using Go and Next.js</strong>
</div>
