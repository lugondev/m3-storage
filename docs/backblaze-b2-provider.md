# Backblaze B2 Storage Provider Configuration

## Overview

Backblaze B2 Cloud Storage is a cost-effective, S3-compatible object storage service that provides reliable storage at a fraction of the cost of traditional cloud providers. It's ideal for backup, archival, and cost-conscious applications.

**When to use Backblaze B2 Provider:**
- Cost-sensitive applications requiring affordable storage
- Backup and archival storage use cases
- Applications with predictable storage growth
- Small to medium-scale deployments
- Applications requiring S3 compatibility without AWS costs

**When to consider alternatives:**
- Applications requiring advanced features not available in B2
- High-performance applications with strict latency requirements
- Applications heavily integrated with other cloud ecosystems
- Use cases requiring global edge distribution (consider Cloudflare R2)

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Backblaze B2 Configuration
backblaze:
    keyID: 'your-application-key-id'            # Backblaze Application Key ID
    applicationKey: 'your-application-key'      # Backblaze Application Key
    bucketID: 'your-bucket-id'                  # Backblaze Bucket ID
    bucketName: 'your-bucket-name'              # Backblaze Bucket Name
    region: 'us-west-002'                       # Optional: Region (default us-west-002)
    endpoint: ''                                # Optional: Custom endpoint URL
```

## Environment Variables

You can also configure Backblaze B2 using environment variables (recommended for production):

- `BACKBLAZE_KEY_ID`: Backblaze Application Key ID
- `BACKBLAZE_APPLICATION_KEY`: Backblaze Application Key
- `BACKBLAZE_BUCKET_ID`: Backblaze Bucket ID
- `BACKBLAZE_BUCKET_NAME`: Backblaze Bucket Name
- `BACKBLAZE_REGION`: Backblaze Region (optional)
- `BACKBLAZE_ENDPOINT`: Custom endpoint URL (optional)

## Configuration Examples

### Production Environment
```yaml
backblaze:
    keyID: '${BACKBLAZE_KEY_ID}'                # Use environment variables
    applicationKey: '${BACKBLAZE_APPLICATION_KEY}'
    bucketID: '${BACKBLAZE_BUCKET_ID}'
    bucketName: 'production-media-storage'
    region: 'us-west-002'
    endpoint: ''                                # Use default B2 endpoints
```

### Development Environment
```yaml
backblaze:
    keyID: 'your-dev-key-id'
    applicationKey: 'your-dev-application-key'
    bucketID: 'your-dev-bucket-id'
    bucketName: 'dev-media-storage'
    region: 'us-west-002'
    endpoint: ''
```

### Multi-Environment Setup
```yaml
# config.production.yaml
backblaze:
    keyID: '${BACKBLAZE_PROD_KEY_ID}'
    applicationKey: '${BACKBLAZE_PROD_APP_KEY}'
    bucketID: '${BACKBLAZE_PROD_BUCKET_ID}'
    bucketName: 'production-media'
    region: 'us-west-002'
    endpoint: ''

# config.staging.yaml
backblaze:
    keyID: '${BACKBLAZE_STAGING_KEY_ID}'
    applicationKey: '${BACKBLAZE_STAGING_APP_KEY}'
    bucketID: '${BACKBLAZE_STAGING_BUCKET_ID}'
    bucketName: 'staging-media'
    region: 'us-west-002'
    endpoint: ''
```

### Cost-Optimized Configuration
```yaml
backblaze:
    keyID: '${BACKBLAZE_KEY_ID}'
    applicationKey: '${BACKBLAZE_APPLICATION_KEY}'
    bucketID: '${BACKBLAZE_BUCKET_ID}'
    bucketName: 'cost-optimized-storage'
    region: 'us-west-002'
    endpoint: 'https://s3.us-west-002.backblazeb2.com'  # Direct region endpoint
```

### Custom Endpoint Configuration
```yaml
backblaze:
    keyID: '${BACKBLAZE_KEY_ID}'
    applicationKey: '${BACKBLAZE_APPLICATION_KEY}'
    bucketID: '${BACKBLAZE_BUCKET_ID}'
    bucketName: 'custom-bucket'
    region: 'eu-central-003'                    # European region
    endpoint: 'https://s3.eu-central-003.backblazeb2.com'
```

## Backblaze B2 Regions

Available Backblaze B2 regions:

| Region Code | Region Name | Location | Endpoint |
|-------------|-------------|----------|----------|
| us-west-002 | US West | California, USA | s3.us-west-002.backblazeb2.com |
| us-west-001 | US West | California, USA | s3.us-west-001.backblazeb2.com |
| us-east-005 | US East | Miami, USA | s3.us-east-005.backblazeb2.com |
| eu-central-003 | EU Central | Amsterdam, Netherlands | s3.eu-central-003.backblazeb2.com |

## Features

The Backblaze B2 provider supports all standard storage operations:

- **Upload**: Upload files to B2 buckets with large file support
- **Download**: Download files from B2 buckets
- **Delete**: Remove files from B2 buckets
- **GetURL**: Get public URLs for files (if bucket allows public access)
- **GetSignedURL**: Generate pre-signed URLs for secure, time-limited access
- **GetObject**: Retrieve file metadata and properties
- **CheckHealth**: Verify connection to B2 and bucket access

## Application Key Setup

### Creating Application Keys in Backblaze Console

1. **Go to Backblaze B2 Console**: https://secure.backblaze.com/b2_buckets.htm
2. **Navigate to App Keys section**
3. **Create New Application Key**
4. **Configure Key Permissions**

### Key Types and Permissions

#### Master Application Key
- Full access to all buckets and operations
- Can create, delete, and manage buckets
- Not recommended for production applications

#### Restricted Application Keys
- Limited to specific buckets or operations
- Recommended for production use
- Can be configured with these capabilities:

```json
{
  "capabilities": [
    "listBuckets",
    "listFiles", 
    "readFiles",
    "shareFiles",
    "writeFiles",
    "deleteFiles"
  ],
  "bucketId": "your-bucket-id",
  "bucketName": "your-bucket-name"
}
```

### Capability Descriptions
- **listBuckets**: List buckets in the account
- **listFiles**: List files in allowed buckets
- **readFiles**: Download files from allowed buckets
- **shareFiles**: Create secure download URLs
- **writeFiles**: Upload files to allowed buckets
- **deleteFiles**: Delete files from allowed buckets

## Bucket Configuration

### Bucket Types

Backblaze B2 offers different bucket types:

| Type | Visibility | Use Case |
|------|------------|----------|
| **Private** | Files are private, require authentication | Secure storage |
| **Public** | Files are publicly accessible | Static websites, CDN |

### Bucket Settings
```json
{
  "bucketType": "allPrivate",
  "bucketInfo": {
    "application": "m3-storage",
    "environment": "production"
  },
  "lifecycleRules": [
    {
      "daysFromHidingToDeleting": 7,
      "daysFromUploadingToHiding": null,
      "fileNamePrefix": "temp/"
    }
  ]
}
```

## Requirements

- Backblaze B2 account
- Application Key with appropriate permissions
- Existing bucket or permissions to create buckets
- Network connectivity to Backblaze endpoints

## Security Considerations

1. **Application Key Security**:
   - Store keys securely using environment variables
   - Use restricted keys with minimal permissions
   - Rotate application keys regularly
   - Monitor key usage in B2 console

2. **Bucket Security**:
   - Use private buckets for sensitive data
   - Configure appropriate CORS policies
   - Implement proper access controls in your application
   - Monitor access patterns

3. **Network Security**:
   - All traffic is encrypted in transit (HTTPS)
   - Use secure endpoints for API calls
   - Implement proper error handling

4. **Access Control**:
   - Use bucket-specific application keys
   - Implement application-level access controls
   - Log and monitor file access patterns

## Performance Optimization

1. **Large File Uploads**: B2 automatically uses large file API for files > 100MB
2. **Connection Pooling**: Handled by S3-compatible client
3. **Regional Selection**: Choose region closest to your users
4. **Parallel Operations**: Support for concurrent uploads/downloads
5. **Retry Logic**: Built-in retry mechanisms for transient failures

## Cost Optimization

### Backblaze B2 Pricing (Competitive Advantages)

- **Storage**: $0.005 per GB per month (very competitive)
- **Downloads**: First 1GB free per day, then $0.01 per GB
- **API Calls**: Very low cost per operation
- **No Hidden Fees**: Transparent pricing model

### Cost Optimization Strategies

1. **Lifecycle Rules**: Automatically delete temporary files
2. **Efficient Operations**: Minimize unnecessary API calls
3. **Download Patterns**: Optimize for free daily download allowance
4. **Compression**: Compress files before storage
5. **Monitoring**: Track usage through B2 console

### Lifecycle Rules Example
```json
{
  "lifecycleRules": [
    {
      "daysFromUploadingToHiding": null,
      "daysFromHidingToDeleting": 30,
      "fileNamePrefix": "temp/"
    },
    {
      "daysFromUploadingToHiding": 365,
      "daysFromHidingToDeleting": 7,
      "fileNamePrefix": "archive/"
    }
  ]
}
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**:
   - Verify Application Key ID and Application Key
   - Check key permissions and capabilities
   - Ensure key is not expired or disabled
   - Verify account is in good standing

2. **Bucket Not Found**:
   - Verify bucket ID and bucket name
   - Check if bucket exists in your account
   - Ensure application key has access to the bucket

3. **Access Denied**:
   - Verify application key capabilities
   - Check bucket type (private vs public)
   - Ensure proper permissions for the operation

4. **Upload Failures**:
   - Check file size limits
   - Verify network connectivity
   - Check for special characters in file names
   - Ensure bucket has sufficient quota

### Health Check

Verify Backblaze B2 connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=backblaze"
```

### B2 CLI Testing

Test your B2 configuration using Backblaze CLI:

```bash
# Install B2 CLI
pip install b2

# Authorize with your account
b2 authorize-account YOUR_KEY_ID YOUR_APPLICATION_KEY

# List buckets
b2 list-buckets

# Upload file
b2 upload-file your-bucket-name ./test.txt test.txt

# Download file
b2 download-file-by-name your-bucket-name test.txt ./downloaded-test.txt

# List files in bucket
b2 ls your-bucket-name
```

### S3-Compatible CLI Testing

Since B2 is S3-compatible, you can also use AWS CLI:

```bash
# Configure AWS CLI for B2
aws configure set aws_access_key_id YOUR_B2_KEY_ID
aws configure set aws_secret_access_key YOUR_B2_APPLICATION_KEY
aws configure set region us-west-002

# Set B2 endpoint
export AWS_ENDPOINT_URL=https://s3.us-west-002.backblazeb2.com

# List buckets
aws s3 ls --endpoint-url $AWS_ENDPOINT_URL

# Upload file
aws s3 cp test.txt s3://your-bucket/test.txt --endpoint-url $AWS_ENDPOINT_URL

# Download file
aws s3 cp s3://your-bucket/test.txt downloaded.txt --endpoint-url $AWS_ENDPOINT_URL
```

## Monitoring and Analytics

### Backblaze B2 Console
Monitor through Backblaze B2 Console:
- Storage usage and costs
- Download statistics
- API call usage
- Bucket activity
- Application key usage

### Key Metrics to Track
- **Storage Usage**: Total stored data per bucket
- **Download Volume**: Data downloaded and associated costs
- **API Usage**: Number of API calls and their types
- **Cost Trends**: Monthly spending patterns
- **Error Rates**: Failed operations and their causes

### Custom Monitoring
Implement application-level monitoring:
```python
# Example monitoring code
import logging
from datetime import datetime

def monitor_b2_operation(operation, success, response_time, file_size=None):
    log_data = {
        'timestamp': datetime.utcnow(),
        'operation': operation,
        'success': success,
        'response_time_ms': response_time,
        'file_size_bytes': file_size
    }
    
    if success:
        logging.info(f"B2 {operation} successful", extra=log_data)
    else:
        logging.error(f"B2 {operation} failed", extra=log_data)
```

## Provider Type

When using the storage factory or API endpoints, use provider type: `"backblaze"` (note: this may be mapped to use the S3 provider internally with B2-specific configuration)

## Migration

### From Other Providers to Backblaze B2

1. **Cost Analysis**: Compare current storage costs with B2 pricing
2. **B2 Setup**: Create buckets and configure application keys
3. **Data Transfer**: Use B2 CLI or S3-compatible tools
4. **Application Update**: Update configuration to use B2
5. **Testing**: Verify functionality and performance
6. **Cost Monitoring**: Track actual costs vs. projections

### Migration Tools

#### Using B2 CLI
```bash
# Sync from local storage to B2
b2 sync ./local-storage/ b2://your-bucket-name/

# Copy specific files
b2 upload-file your-bucket ./large-file.zip large-file.zip
```

#### Using rclone
```bash
# Configure B2 destination
rclone config create b2-dest b2 \
    account=YOUR_KEY_ID \
    key=YOUR_APPLICATION_KEY

# Configure source (e.g., S3)
rclone config create s3-source s3 \
    access_key_id=AWS_ACCESS_KEY \
    secret_access_key=AWS_SECRET_KEY \
    region=us-east-1

# Perform migration
rclone sync s3-source:source-bucket b2-dest:destination-bucket --progress
```

#### Using AWS CLI with B2
```bash
# Sync from S3 to B2
aws s3 sync s3://source-bucket s3://destination-bucket \
    --source-region us-east-1 \
    --endpoint-url https://s3.us-west-002.backblazeb2.com
```

## Integration Examples

### With CDN (Cloudflare)
```yaml
# Use B2 as origin for Cloudflare CDN
backblaze:
    keyID: '${BACKBLAZE_KEY_ID}'
    applicationKey: '${BACKBLAZE_APPLICATION_KEY}'
    bucketID: '${BACKBLAZE_BUCKET_ID}'
    bucketName: 'cdn-origin-storage'
    region: 'us-west-002'
    endpoint: ''

cdn:
    provider: 'cloudflare'
    originUrl: 'https://f002.backblazeb2.com/file/cdn-origin-storage'
```

### Backup Strategy
```yaml
# Primary storage on other provider, B2 for backup
primary_storage:
    provider: 's3'
    # ... S3 configuration

backup_storage:
    provider: 'backblaze'
    keyID: '${BACKBLAZE_BACKUP_KEY_ID}'
    applicationKey: '${BACKBLAZE_BACKUP_APP_KEY}'
    bucketID: '${BACKBLAZE_BACKUP_BUCKET_ID}'
    bucketName: 'backup-storage'
    region: 'us-west-002'
```

### Archive Storage
Use B2 for long-term archival with lifecycle rules:
```json
{
  "lifecycleRules": [
    {
      "fileNamePrefix": "active/",
      "daysFromUploadingToHiding": 90
    },
    {
      "fileNamePrefix": "archive/",
      "daysFromUploadingToHiding": null,
      "daysFromHidingToDeleting": 2555  // ~7 years
    }
  ]
}
```