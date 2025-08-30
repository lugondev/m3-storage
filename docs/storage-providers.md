# Storage Providers

This document provides an overview of all supported storage providers in the M3 Storage system. For detailed configuration instructions for each provider, please refer to the individual provider documentation files.

## Available Providers

### Local Storage
- **Local Filesystem** - Store files directly on the server's filesystem
  - 📚 **Documentation**: [Local Storage Provider Guide](./local-storage-provider.md)

### S3-Compatible Storage
- **Amazon S3** - Amazon Simple Storage Service
  - 📚 **Documentation**: [Amazon S3 Provider Guide](./s3-provider.md)
- **MinIO** - High-performance object storage (self-hosted or cloud)
  - 📚 **Documentation**: [MinIO Provider Guide](./minio-provider.md)

### Cloud Platform Storage
- **Firebase Storage** - Google Firebase Cloud Storage service
  - 📚 **Documentation**: [Firebase Storage Provider Guide](./firebase-provider.md)
- **Azure Blob Storage** - Microsoft Azure Blob Storage service
  - 📚 **Documentation**: [Azure Blob Storage Provider Guide](./azure-provider.md)

### Cost-Effective Storage
- **Cloudflare R2** - Cloudflare's S3-compatible storage solution with zero egress fees
  - 📚 **Documentation**: [Cloudflare R2 Provider Guide](./cloudflare-r2-provider.md)
- **Backblaze B2** - Backblaze B2 Cloud Storage (S3-compatible)
  - 📚 **Documentation**: [Backblaze B2 Provider Guide](./backblaze-b2-provider.md)
- **Scaleway Object Storage** - Scaleway's S3-compatible European storage service
  - 📚 **Documentation**: [Scaleway Provider Guide](./scaleway-provider.md)

### Alternative Storage
- **Discord** - Store files using Discord channels (experimental/educational use)
  - 📚 **Documentation**: [Discord Provider Guide](./discord-provider.md)
  - ⚠️ **Note**: For experimental/educational use only

---

## Quick Provider Comparison

| Provider | Use Case | Pros | Cons | Best For |
|----------|----------|------|------|----------|
| **Local Storage** | Development, small deployments | Simple setup, no costs, full control | Not scalable, single point of failure | Development, testing |
| **Amazon S3** | Production, enterprise | Highly reliable, feature-rich, global | Can be expensive with high egress | Enterprise applications |
| **MinIO** | Self-hosted, hybrid cloud | S3-compatible, self-hosted option | Requires infrastructure management | On-premise, hybrid setups |
| **Firebase Storage** | Firebase ecosystem | Great mobile integration, real-time features | Limited to Google ecosystem | Mobile apps, Firebase projects |
| **Azure Blob Storage** | Microsoft ecosystem | Excellent Azure integration | Best for Azure-hosted apps | Enterprise Microsoft environments |
| **Cloudflare R2** | High-traffic applications | Zero egress fees, global CDN | Newer service, fewer features | High-bandwidth applications |
| **Backblaze B2** | Cost-conscious applications | Very low cost, reliable | Fewer advanced features | Backup, archival storage |
| **Scaleway** | European applications | GDPR compliant, competitive pricing | Limited to European regions | EU-based applications |
| **Discord** | Experimental projects | Creative solution, no setup cost | Not reliable, ToS concerns | Educational, experiments only |

## Provider Selection Guide

### For Development & Testing
1. **Local Storage** - Quick setup, no external dependencies
2. **MinIO** - Self-hosted S3-compatible option

### For Production Applications
1. **Amazon S3** - Industry standard with extensive features
2. **Azure Blob Storage** - Best for Azure-hosted applications
3. **Firebase Storage** - Excellent for mobile/web apps using Firebase

### For Cost-Conscious Applications
1. **Backblaze B2** - Extremely competitive storage pricing
2. **Cloudflare R2** - Zero egress fees, great for high-traffic
3. **Scaleway** - Good European pricing options

### For Specific Requirements
- **Global Distribution**: Cloudflare R2, Amazon S3
- **GDPR Compliance**: Scaleway, Azure (EU regions)
- **Self-Hosted**: Local Storage, MinIO
- **Mobile Integration**: Firebase Storage
- **Microsoft Ecosystem**: Azure Blob Storage

## Configuration Overview

Each provider requires specific configuration parameters. Here's a quick reference:

### Environment Variables Pattern
All providers support environment variable configuration using this pattern:
- `{PROVIDER}_ACCESS_KEY_ID` - Access credentials
- `{PROVIDER}_SECRET_ACCESS_KEY` - Secret credentials  
- `{PROVIDER}_BUCKET_NAME` - Storage container name
- `{PROVIDER}_REGION` - Region/location (where applicable)

### Configuration Priority
Configuration values are loaded in this order (highest to lowest priority):
1. **Environment variables** (recommended for production)
2. **Configuration file values**
3. **Default values**

## Health Checks

All providers support health checks through the API:

### Check Specific Provider
```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=PROVIDER_NAME"
```

### Check All Providers
```bash
curl -X GET "http://localhost:8083/api/v1/storage/health/all"
```

### Available Provider Types
- `local` - Local Storage
- `s3` - Amazon S3
- `minio` - MinIO
- `firebase` - Firebase Storage  
- `azure` - Azure Blob Storage
- `cloudflare` - Cloudflare R2
- `backblaze` - Backblaze B2
- `scaleway` - Scaleway Object Storage
- `discord` - Discord Storage

## Getting Started

1. **Choose your provider** based on your requirements
2. **Read the detailed documentation** for your chosen provider
3. **Set up your account** and obtain necessary credentials
4. **Configure your application** using environment variables
5. **Test the connection** using health check endpoints
6. **Start uploading files** through the API

## Migration Between Providers

The M3 Storage system is designed to make provider migration straightforward:

1. **Configure new provider** alongside existing one
2. **Test new provider** functionality
3. **Migrate data** using appropriate tools (rclone, AWS CLI, etc.)
4. **Update application configuration**
5. **Monitor and verify** the migration

Each provider documentation includes specific migration guidance and tool recommendations.

## Support and Troubleshooting

- Check the individual provider documentation for specific troubleshooting guides
- Use health check endpoints to verify connectivity
- Monitor application logs for detailed error messages
- Refer to each provider's official documentation for service-specific issues

For detailed configuration instructions, troubleshooting guides, and best practices, please refer to the individual provider documentation linked above.
