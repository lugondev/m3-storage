# MinIO Storage Provider Configuration

## Overview

MinIO is a high-performance object storage system that is compatible with Amazon S3 API. This document describes how to configure the MinIO storage provider in the M3 Storage system.

**When to use MinIO Provider:**
- Self-hosted MinIO servers
- MinIO cloud services 
- Any storage service that explicitly identifies as MinIO
- When you want MinIO-specific optimizations and error handling

**When to use S3 Provider instead:**
- Amazon S3
- Other S3-compatible services (Cloudflare R2, Backblaze B2, etc.)
- Services that implement S3 API but are not MinIO-based

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# MinIO Configuration (Self-hosted or hosted MinIO service)
minio:
    accessKeyID: 'minioadmin'                    # MinIO Access Key ID
    secretAccessKey: 'minioadmin'                # MinIO Secret Access Key
    bucketName: 'm3-storage'                     # MinIO Bucket Name
    endpoint: 'https://play.min.io'              # MinIO endpoint URL
    region: 'us-east-1'                          # Optional: MinIO region
    useSSL: true                                 # Whether to use SSL/TLS
```

## Environment Variables

You can also configure MinIO using environment variables:

- `MINIO_ACCESS_KEY_ID`: MinIO Access Key ID
- `MINIO_SECRET_ACCESS_KEY`: MinIO Secret Access Key  
- `MINIO_BUCKET_NAME`: MinIO Bucket Name
- `MINIO_ENDPOINT`: MinIO endpoint URL
- `MINIO_REGION`: MinIO region (optional)
- `MINIO_USE_SSL`: Whether to use SSL/TLS (true/false)

## Configuration Examples

### Self-hosted MinIO (Local Development)
```yaml
minio:
    accessKeyID: 'minioadmin'
    secretAccessKey: 'minioadmin'
    bucketName: 'dev-storage'
    endpoint: 'http://localhost:9000'
    region: 'us-east-1'
    useSSL: false
```

### MinIO Cloud (Play Server)
```yaml
minio:
    accessKeyID: 'Q3AM3UQ867SPQQA43P2F'
    secretAccessKey: 'zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG'
    bucketName: 'test-bucket'
    endpoint: 'https://play.min.io'
    region: 'us-east-1'
    useSSL: true
```

### Production MinIO Setup
```yaml
minio:
    accessKeyID: '${MINIO_ACCESS_KEY_ID}'
    secretAccessKey: '${MINIO_SECRET_ACCESS_KEY}'
    bucketName: 'production-storage'
    endpoint: 'https://minio.yourdomain.com'
    region: 'us-east-1'
    useSSL: true
```

### Custom Hosted MinIO
```yaml
minio:
    accessKeyID: 'your-access-key'
    secretAccessKey: 'your-secret-key'
    bucketName: 'my-app-storage'
    endpoint: 'https://storage.mycompany.com'  # Custom domain
    region: 'us-west-1'
    useSSL: true
```

### MinIO with Custom Port
```yaml
minio:
    accessKeyID: 'minioadmin'
    secretAccessKey: 'minioadmin'
    bucketName: 'test-bucket'
    endpoint: 'https://minio.example.com:9443'  # Custom port
    region: 'us-east-1'
    useSSL: true
```

## Features

The MinIO provider supports all standard storage operations:

- **Upload**: Upload files to MinIO buckets
- **Download**: Download files from MinIO buckets
- **Delete**: Remove files from MinIO buckets
- **GetURL**: Get public URLs for files (if bucket allows public access)
- **GetSignedURL**: Generate time-limited signed URLs for private file access
- **GetObject**: Retrieve file metadata
- **CheckHealth**: Verify connection to MinIO server

## Requirements

- MinIO server (self-hosted or cloud service)
- Valid access credentials (Access Key ID and Secret Access Key)
- Existing bucket or permissions to create buckets
- Network connectivity to the MinIO endpoint

## Security Considerations

1. **Credentials**: Store access credentials securely using environment variables
2. **SSL/TLS**: Always use SSL/TLS in production (`useSSL: true`)
3. **Bucket Policy**: Configure appropriate bucket policies for your use case
4. **Network**: Ensure proper network security between your application and MinIO server

## Troubleshooting

### Common Issues

1. **Connection Failed**: 
   - Check endpoint URL format and accessibility
   - Verify SSL/TLS settings
   - Ensure firewall allows connection

2. **Authentication Failed**:
   - Verify access key ID and secret access key
   - Check user permissions in MinIO

3. **Bucket Not Found**:
   - Ensure bucket exists or create it manually
   - Verify bucket name spelling and case sensitivity

4. **SSL Certificate Issues**:
   - For self-signed certificates, consider setting `useSSL: false` in development
   - For production, ensure valid SSL certificates

### Health Check

You can verify the MinIO connection using the health check endpoint:

```bash
curl -X POST http://localhost:8083/api/storage/health \
  -H "Content-Type: application/json" \
  -d '{"provider_type": "minio"}'
```

### Testing Configuration

To test your MinIO configuration:

1. Start your application with MinIO configuration
2. Check the logs for MinIO initialization messages
3. Use the health check API endpoint
4. Try uploading a test file

## Performance Notes

- MinIO is designed for high performance and can handle large files efficiently
- Path-style URLs are used (required for S3-compatible storage)
- Connection pooling is handled by the AWS SDK
- Large file uploads use multipart upload automatically

## Provider Type

When using the storage factory or API endpoints, use provider type: `"minio"`