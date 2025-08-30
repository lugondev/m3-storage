# Scaleway Object Storage Provider Configuration

## Overview

Scaleway Object Storage is a European S3-compatible object storage service that provides secure, scalable, and cost-effective storage solutions. It's ideal for European applications requiring GDPR compliance and regional data sovereignty.

**When to use Scaleway Object Storage Provider:**
- European applications requiring GDPR compliance
- Applications hosted on Scaleway infrastructure
- Cost-effective European cloud storage needs
- Multi-region European deployments
- Applications requiring data sovereignty in Europe

**When to consider alternatives:**
- Global applications requiring worldwide distribution (consider AWS S3 or Cloudflare R2)
- Applications primarily hosted on other cloud platforms
- Applications requiring advanced features not available in Scaleway

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Scaleway Object Storage Configuration
scaleway:
    accessKeyID: 'SCW1234567890ABCDEF'         # Scaleway Access Key ID
    secretAccessKey: 'your-secret-access-key'   # Scaleway Secret Access Key
    region: 'fr-par'                            # Scaleway Region (fr-par, nl-ams, pl-waw)
    bucketName: 'your-scaleway-bucket'          # Scaleway Bucket Name
    endpoint: ''                                # Optional: Custom endpoint URL
```

## Environment Variables

You can also configure Scaleway Object Storage using environment variables (recommended for production):

- `SCALEWAY_ACCESS_KEY_ID`: Scaleway Access Key ID
- `SCALEWAY_SECRET_ACCESS_KEY`: Scaleway Secret Access Key
- `SCALEWAY_REGION`: Scaleway Region
- `SCALEWAY_BUCKET_NAME`: Scaleway Bucket Name
- `SCALEWAY_ENDPOINT`: Custom endpoint URL (optional)

## Configuration Examples

### Production Environment (Paris)
```yaml
scaleway:
    accessKeyID: '${SCALEWAY_ACCESS_KEY_ID}'    # Use environment variables
    secretAccessKey: '${SCALEWAY_SECRET_ACCESS_KEY}'
    region: 'fr-par'                            # Paris region
    bucketName: 'production-media-storage'
    endpoint: ''                                # Use default Scaleway endpoints
```

### Development Environment
```yaml
scaleway:
    accessKeyID: 'your-dev-access-key'
    secretAccessKey: 'your-dev-secret-key'
    region: 'fr-par'
    bucketName: 'dev-media-storage'
    endpoint: ''
```

### Multi-Region Setup (Amsterdam)
```yaml
scaleway:
    accessKeyID: '${SCALEWAY_ACCESS_KEY_ID}'
    secretAccessKey: '${SCALEWAY_SECRET_ACCESS_KEY}'
    region: 'nl-ams'                            # Amsterdam region
    bucketName: 'media-nl-ams'
    endpoint: ''
```

### Multi-Region Setup (Warsaw)
```yaml
scaleway:
    accessKeyID: '${SCALEWAY_ACCESS_KEY_ID}'
    secretAccessKey: '${SCALEWAY_SECRET_ACCESS_KEY}'
    region: 'pl-waw'                            # Warsaw region
    bucketName: 'media-pl-waw'
    endpoint: ''
```

### Custom Endpoint Configuration
```yaml
scaleway:
    accessKeyID: '${SCALEWAY_ACCESS_KEY_ID}'
    secretAccessKey: '${SCALEWAY_SECRET_ACCESS_KEY}'
    region: 'fr-par'
    bucketName: 'custom-bucket'
    endpoint: 'https://s3.fr-par.scw.cloud'    # Explicit endpoint
```

### Private Network Configuration
```yaml
scaleway:
    accessKeyID: '${SCALEWAY_ACCESS_KEY_ID}'
    secretAccessKey: '${SCALEWAY_SECRET_ACCESS_KEY}'
    region: 'fr-par'
    bucketName: 'private-media'
    endpoint: 'https://s3.internal.fr-par.scw.cloud'  # Internal endpoint for private networks
```

## Scaleway Regions

Available Scaleway regions for Object Storage:

| Region Code | Region Name | Location | Endpoint |
|-------------|-------------|----------|----------|
| fr-par | France (Paris) | Paris, France | s3.fr-par.scw.cloud |
| nl-ams | Netherlands (Amsterdam) | Amsterdam, Netherlands | s3.nl-ams.scw.cloud |
| pl-waw | Poland (Warsaw) | Warsaw, Poland | s3.pl-waw.scw.cloud |

## Features

The Scaleway Object Storage provider supports all standard storage operations:

- **Upload**: Upload files to Scaleway buckets with multipart support
- **Download**: Download files from Scaleway buckets
- **Delete**: Remove files from Scaleway buckets
- **GetURL**: Get public URLs for files (if bucket allows public access)
- **GetSignedURL**: Generate pre-signed URLs for secure, time-limited access
- **GetObject**: Retrieve object metadata and properties
- **CheckHealth**: Verify connection to Scaleway and bucket access

## API Token Setup

### Creating API Tokens in Scaleway Console

1. **Go to Scaleway Console**: https://console.scaleway.com/
2. **Navigate to Identity and Access Management (IAM)**
3. **Go to API Keys section**
4. **Create new API Key**
5. **Set appropriate permissions**

### Required Permissions

Configure your API key with these Object Storage permissions:

```json
{
  "rules": [
    {
      "permission_sets": [
        "ObjectStorageRead",
        "ObjectStorageWrite",
        "ObjectStorageDelete"
      ],
      "scope": {
        "projects": ["your-project-id"]
      }
    }
  ]
}
```

### API Key Permissions Breakdown
- **ObjectStorageRead**: Read objects and list buckets
- **ObjectStorageWrite**: Upload objects and create buckets
- **ObjectStorageDelete**: Delete objects and buckets
- **ObjectStorageAdmin**: Full administrative access

## Bucket Configuration

### Storage Classes

Scaleway offers different storage classes:

| Storage Class | Use Case | Pricing |
|---------------|----------|---------|
| **STANDARD** | Frequently accessed data | Standard pricing |
| **COLD** | Infrequently accessed data | Lower storage cost, higher retrieval cost |

### Bucket Policies

Configure bucket policies for access control:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:scw:s3:::your-bucket/*"
    }
  ]
}
```

### CORS Configuration

For web applications accessing Scaleway directly:

```json
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
    "AllowedOrigins": ["https://yourdomain.com"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]
```

## Requirements

- Scaleway account with Object Storage enabled
- Valid API key with Object Storage permissions
- Existing bucket or permissions to create buckets
- Network connectivity to Scaleway endpoints

## Security Considerations

1. **API Key Security**:
   - Store API keys securely using environment variables
   - Use least privilege principle for API key permissions
   - Rotate API keys regularly
   - Monitor API key usage in Scaleway Console

2. **Bucket Security**:
   - Configure appropriate bucket policies
   - Use private buckets for sensitive data
   - Implement proper CORS policies
   - Enable access logging for audit trails

3. **Network Security**:
   - All traffic is encrypted in transit (HTTPS)
   - Use private networks for internal traffic
   - Configure appropriate firewall rules

4. **GDPR Compliance**:
   - Data stays within European regions
   - Implement proper data retention policies
   - Configure appropriate access controls
   - Maintain audit logs for compliance

## Performance Optimization

1. **Regional Proximity**: Choose region closest to your users/infrastructure
2. **Connection Pooling**: Handled automatically by S3-compatible client
3. **Multipart Upload**: Automatic for large files
4. **CDN Integration**: Use Scaleway's CDN for global distribution
5. **Private Networks**: Use internal endpoints for faster access

## Cost Optimization

### Scaleway Object Storage Pricing

- **Storage**: Competitive European pricing
- **Operations**: Pay-per-request pricing
- **Data Transfer**: Outbound transfer costs apply
- **Storage Classes**: Use COLD storage for infrequently accessed data

### Cost Optimization Strategies

1. **Storage Classes**: Use appropriate storage class for access patterns
2. **Lifecycle Policies**: Automatically transition or delete old objects
3. **Compression**: Compress files before storage
4. **Monitoring**: Track usage through Scaleway Console
5. **Regional Selection**: Choose cost-effective regions

### Lifecycle Management
```json
{
  "Rules": [
    {
      "ID": "TransitionToCold",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "archive/"
      },
      "Transitions": [
        {
          "Days": 30,
          "StorageClass": "COLD"
        }
      ]
    },
    {
      "ID": "DeleteTempFiles",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "temp/"
      },
      "Expiration": {
        "Days": 7
      }
    }
  ]
}
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**:
   - Verify Access Key ID and Secret Access Key
   - Check API key permissions in Scaleway Console
   - Ensure API key is active and not expired
   - Verify project access permissions

2. **Bucket Not Found**:
   - Verify bucket name and region
   - Ensure bucket exists in the specified region
   - Check bucket naming conventions

3. **Access Denied**:
   - Verify API key has required permissions
   - Check bucket policies and ACLs
   - Ensure correct project context

4. **Network Issues**:
   - Check connectivity to Scaleway endpoints
   - Verify DNS resolution
   - Check firewall and proxy settings
   - Try using internal endpoints for private networks

### Health Check

Verify Scaleway Object Storage connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=scaleway"
```

### S3-Compatible CLI Testing

Since Scaleway is S3-compatible, you can use AWS CLI:

```bash
# Configure AWS CLI for Scaleway
aws configure set aws_access_key_id YOUR_SCALEWAY_ACCESS_KEY
aws configure set aws_secret_access_key YOUR_SCALEWAY_SECRET_KEY
aws configure set region fr-par

# Set Scaleway endpoint
export AWS_ENDPOINT_URL=https://s3.fr-par.scw.cloud

# List buckets
aws s3 ls --endpoint-url $AWS_ENDPOINT_URL

# Upload file
aws s3 cp test.txt s3://your-bucket/test.txt --endpoint-url $AWS_ENDPOINT_URL

# Download file
aws s3 cp s3://your-bucket/test.txt downloaded.txt --endpoint-url $AWS_ENDPOINT_URL

# List objects
aws s3 ls s3://your-bucket/ --endpoint-url $AWS_ENDPOINT_URL
```

### Scaleway CLI Testing

Test using Scaleway CLI (scw):

```bash
# Install Scaleway CLI
curl -s https://raw.githubusercontent.com/scaleway/scaleway-cli/master/scripts/get.sh | sh

# Configure CLI
scw init

# List buckets
scw object bucket list

# Create bucket
scw object bucket create name=test-bucket region=fr-par

# Upload object
scw object object put bucket=test-bucket key=test.txt body=./test.txt

# Download object
scw object object get bucket=test-bucket key=test.txt > downloaded-test.txt

# List objects
scw object object list bucket=test-bucket
```

## Monitoring and Analytics

### Scaleway Console Monitoring
Monitor through Scaleway Console:
- Storage usage and trends
- Request statistics and patterns
- Bandwidth usage
- Cost analysis
- Performance metrics

### Key Metrics to Track
- **Storage Usage**: Total stored data per bucket
- **Request Volume**: GET, PUT, DELETE operations
- **Bandwidth**: Ingress and egress traffic
- **Error Rates**: Failed requests and their causes
- **Cost**: Monthly spending and trends

### Custom Monitoring
Implement application-level monitoring:
- Upload/download success rates
- Response times and latency
- Error tracking and alerting
- Usage patterns by region

### Log Analysis
- Access logs for security analysis
- Performance logs for optimization
- Error logs for troubleshooting
- Cost analysis for optimization

## Provider Type

When using the storage factory or API endpoints, use provider type: `"scaleway"` (note: this may be mapped to use the S3 provider internally with Scaleway-specific configuration)

## Migration

### From Other Providers to Scaleway

1. **Assessment**: Analyze current storage requirements and costs
2. **Scaleway Setup**: Create buckets and configure API keys
3. **Data Transfer**: Use S3-compatible tools for migration
4. **Configuration Update**: Update application configuration
5. **Testing**: Thoroughly test functionality and performance
6. **Monitoring**: Monitor migration progress and costs

### Migration Tools

#### Using rclone
```bash
# Configure source (e.g., AWS S3)
rclone config create s3-source s3 \
    access_key_id=AWS_ACCESS_KEY \
    secret_access_key=AWS_SECRET_KEY \
    region=us-east-1

# Configure Scaleway destination
rclone config create scaleway-dest s3 \
    access_key_id=SCALEWAY_ACCESS_KEY \
    secret_access_key=SCALEWAY_SECRET_KEY \
    endpoint=https://s3.fr-par.scw.cloud \
    region=fr-par

# Perform migration
rclone sync s3-source:source-bucket scaleway-dest:destination-bucket --progress
```

#### Using AWS CLI
```bash
# Sync from another S3-compatible service to Scaleway
aws s3 sync s3://source-bucket s3://destination-bucket \
    --source-region us-east-1 \
    --endpoint-url https://s3.fr-par.scw.cloud \
    --region fr-par
```

## GDPR Compliance

### Data Location
- All data stored in European regions (Paris, Amsterdam, Warsaw)
- No cross-border data transfers outside EU
- Clear data residency guarantees

### Compliance Features
- **Right to be Forgotten**: Easy object deletion
- **Data Portability**: S3-compatible exports
- **Access Logging**: Complete audit trails
- **Encryption**: Data encrypted at rest and in transit

### Implementation Guidelines
```yaml
# GDPR-compliant configuration example
scaleway:
    accessKeyID: '${SCALEWAY_ACCESS_KEY_ID}'
    secretAccessKey: '${SCALEWAY_SECRET_ACCESS_KEY}'
    region: 'fr-par'                             # EU region
    bucketName: 'gdpr-compliant-storage'
    endpoint: ''
    
# Additional GDPR considerations
gdpr:
    dataRetentionPeriod: '7 years'               # As per business requirements
    encryptionRequired: true
    accessLoggingEnabled: true
    rightToBeForgettenEnabled: true
```

## Integration with Scaleway Services

### Scaleway Functions
Integrate with Scaleway Functions for serverless processing:

```python
# Example Scaleway Function
def handler(event, context):
    # Process uploaded objects
    bucket = event['Records'][0]['s3']['bucket']['name']
    key = event['Records'][0]['s3']['object']['key']
    
    # Process the object
    return {
        'statusCode': 200,
        'body': f'Processed {key} from {bucket}'
    }
```

### Scaleway Container Registry
Store container images alongside your object data for a complete solution.

### Scaleway Database
Use with Scaleway's managed databases for metadata storage and application data.