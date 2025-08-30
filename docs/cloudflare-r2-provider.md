# Cloudflare R2 Storage Provider Configuration

## Overview

Cloudflare R2 Storage is an S3-compatible object storage service that eliminates egress bandwidth fees. It provides global distribution through Cloudflare's edge network and offers competitive pricing for storage and operations.

**When to use Cloudflare R2 Provider:**
- Applications requiring global CDN and edge distribution
- High-traffic applications with significant egress bandwidth costs
- Applications already using Cloudflare services (DNS, CDN, etc.)
- Cost-conscious deployments seeking to eliminate egress fees
- Applications requiring S3 compatibility without AWS

**When to consider alternatives:**
- Applications requiring advanced AWS S3 features not available in R2
- Applications heavily integrated with AWS ecosystem
- Use cases requiring features still in development for R2

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Cloudflare R2 Configuration
cloudflare:
    accountID: 'your-cloudflare-account-id'     # Cloudflare Account ID
    accessKeyID: 'your-r2-access-key-id'       # R2 API Token with appropriate permissions
    secretAccessKey: 'your-r2-secret-key'      # R2 Secret Access Key
    bucketName: 'your-r2-bucket-name'          # R2 Bucket Name
    publicDomain: 'cdn.example.com'            # Optional: Custom domain for public access
```

## Environment Variables

You can also configure Cloudflare R2 using environment variables (recommended for production):

- `CLOUDFLARE_ACCOUNT_ID`: Cloudflare Account ID
- `CLOUDFLARE_ACCESS_KEY_ID`: R2 Access Key ID
- `CLOUDFLARE_SECRET_ACCESS_KEY`: R2 Secret Access Key
- `CLOUDFLARE_BUCKET_NAME`: R2 Bucket Name
- `CLOUDFLARE_PUBLIC_DOMAIN`: Custom public domain (optional)

## Configuration Examples

### Production Environment
```yaml
cloudflare:
    accountID: '${CLOUDFLARE_ACCOUNT_ID}'       # Use environment variables
    accessKeyID: '${CLOUDFLARE_ACCESS_KEY_ID}'
    secretAccessKey: '${CLOUDFLARE_SECRET_ACCESS_KEY}'
    bucketName: 'production-media-storage'
    publicDomain: 'media.yourdomain.com'        # Custom domain with CNAME
```

### Development Environment
```yaml
cloudflare:
    accountID: 'abc123def456ghi789'
    accessKeyID: 'your-dev-access-key'
    secretAccessKey: 'your-dev-secret-key'
    bucketName: 'dev-media-storage'
    publicDomain: ''                            # Use R2.dev domain for development
```

### Multi-Region Setup (Global)
```yaml
# R2 automatically provides global distribution
cloudflare:
    accountID: '${CLOUDFLARE_ACCOUNT_ID}'
    accessKeyID: '${CLOUDFLARE_ACCESS_KEY_ID}'
    secretAccessKey: '${CLOUDFLARE_SECRET_ACCESS_KEY}'
    bucketName: 'global-media-storage'
    publicDomain: 'cdn.example.com'             # Global CDN domain
```

### Multiple Buckets Configuration
```yaml
# Production bucket
cloudflare:
    accountID: '${CLOUDFLARE_ACCOUNT_ID}'
    accessKeyID: '${CLOUDFLARE_ACCESS_KEY_ID}'
    secretAccessKey: '${CLOUDFLARE_SECRET_ACCESS_KEY}'
    bucketName: 'production-media'
    publicDomain: 'media.example.com'

# Backup bucket (different configuration file)
# cloudflare:
#     accountID: '${CLOUDFLARE_ACCOUNT_ID}'
#     accessKeyID: '${CLOUDFLARE_ACCESS_KEY_ID}'
#     secretAccessKey: '${CLOUDFLARE_SECRET_ACCESS_KEY}'
#     bucketName: 'backup-media'
#     publicDomain: 'backup.example.com'
```

### Custom Domain with SSL
```yaml
cloudflare:
    accountID: '${CLOUDFLARE_ACCOUNT_ID}'
    accessKeyID: '${CLOUDFLARE_ACCESS_KEY_ID}'
    secretAccessKey: '${CLOUDFLARE_SECRET_ACCESS_KEY}'
    bucketName: 'media-files'
    publicDomain: 'files.yourdomain.com'        # Must be configured in Cloudflare DNS
```

## R2 API Token Setup

### Creating R2 API Token

1. **Go to Cloudflare Dashboard**: https://dash.cloudflare.com/
2. **Navigate to R2 Object Storage**
3. **Go to "Manage R2 API Tokens"**
4. **Create API Token** with appropriate permissions

### Required Permissions

Configure your R2 API Token with these permissions:

```json
{
  "policies": [
    {
      "effect": "allow",
      "resources": {
        "com.cloudflare.api.account.*": "*"
      },
      "permission_groups": [
        {
          "id": "c8fed203ed3043cba015a93ad1616f1f",
          "name": "Zone.Zone Settings:Edit"
        },
        {
          "id": "82a0d7a8cc3c41aaa3ae1a53d7ba2f33",
          "name": "Zone.Zone:Read"
        }
      ]
    }
  ]
}
```

Or use these specific permissions:
- **Object Read**: Read objects from R2 buckets
- **Object Write**: Write objects to R2 buckets
- **Object Delete**: Delete objects from R2 buckets
- **Bucket Read**: List buckets and read bucket metadata

## Features

The Cloudflare R2 provider supports all standard storage operations:

- **Upload**: Upload files to R2 buckets (S3-compatible API)
- **Download**: Download files from R2 buckets
- **Delete**: Remove files from R2 buckets
- **GetURL**: Get public URLs for files using R2.dev domains or custom domains
- **GetSignedURL**: Generate pre-signed URLs for secure, time-limited access
- **GetObject**: Retrieve object metadata and properties
- **CheckHealth**: Verify connection to R2 and bucket access

## Key Features of R2

### Zero Egress Fees
- No charges for data transfer out (egress)
- Only pay for storage and operations
- Significant cost savings for high-traffic applications

### Global Distribution
- Automatically distributed across Cloudflare's global network
- Low latency access from anywhere in the world
- Built-in CDN capabilities

### S3 Compatibility
- Compatible with existing S3 tools and SDKs
- Easy migration from Amazon S3
- Familiar API and operations

## Custom Domain Configuration

### Setting up Custom Domain

1. **Add CNAME Record in Cloudflare DNS**:
   ```
   Name: files.yourdomain.com
   Type: CNAME
   Content: your-bucket-name.your-account-id.r2.cloudflarestorage.com
   ```

2. **Configure SSL/TLS**:
   - Cloudflare automatically provides SSL certificates
   - Configure appropriate SSL/TLS mode in Cloudflare

3. **Update Configuration**:
   ```yaml
   cloudflare:
       # ... other settings ...
       publicDomain: 'files.yourdomain.com'
   ```

### Public Access Configuration

Configure bucket for public access through Cloudflare Dashboard:

1. **Go to R2 Object Storage**
2. **Select your bucket**
3. **Configure Public Access**
4. **Set up Custom Domain (optional)**

## Requirements

- Cloudflare account with R2 enabled
- R2 API token with appropriate permissions
- Existing R2 bucket or permissions to create buckets
- Network connectivity to Cloudflare R2 endpoints

## Security Considerations

1. **API Token Security**:
   - Store API tokens securely using environment variables
   - Use tokens with minimal required permissions
   - Rotate API tokens regularly
   - Monitor token usage through Cloudflare Dashboard

2. **Bucket Security**:
   - Configure appropriate public/private access
   - Use signed URLs for sensitive content
   - Implement proper CORS policies for web access
   - Monitor access logs

3. **Network Security**:
   - All traffic is encrypted in transit (HTTPS)
   - Leverage Cloudflare's security features (DDoS protection, etc.)
   - Configure appropriate firewall rules if needed

4. **Access Control**:
   - Use IAM-like policies through Cloudflare's API token system
   - Implement application-level access controls
   - Monitor and log access patterns

## Performance Optimization

1. **Global Distribution**: Automatic through Cloudflare's edge network
2. **Caching**: Built-in CDN caching capabilities
3. **Compression**: Enable Brotli/Gzip compression in Cloudflare
4. **HTTP/2 and HTTP/3**: Automatic support for modern protocols
5. **Connection Pooling**: Handled by the S3-compatible client

## Cost Optimization

### R2 Pricing (as of 2024)

- **Storage**: $0.015 per GB per month
- **Class A Operations**: $4.50 per million (PUT, LIST, etc.)
- **Class B Operations**: $0.36 per million (GET, SELECT, etc.)
- **No Egress Fees**: $0 for data transfer out

### Cost Optimization Strategies

1. **Lifecycle Policies**: Automatically manage object lifecycle
2. **Compression**: Compress files before storage
3. **Efficient Operations**: Minimize unnecessary API calls
4. **Monitoring**: Use Cloudflare Analytics to track usage

### Lifecycle Management
```yaml
# Example lifecycle configuration
lifecycle:
  rules:
    - id: "DeleteTempFiles"
      status: "Enabled"
      filter:
        prefix: "temp/"
      expiration:
        days: 7
    - id: "ArchiveOldMedia"
      status: "Enabled"
      filter:
        prefix: "media/"
      transition:
        days: 90
        storage_class: "INFREQUENT_ACCESS"  # When available
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**:
   - Verify Cloudflare Account ID
   - Check R2 API token validity and permissions
   - Ensure API token has R2 access enabled
   - Verify account has R2 service enabled

2. **Bucket Not Found**:
   - Verify bucket name spelling and case sensitivity
   - Ensure bucket exists in the specified account
   - Check bucket region/location

3. **Access Denied**:
   - Verify API token permissions
   - Check bucket access policies
   - Ensure account has sufficient R2 quota/limits

4. **Network Issues**:
   - Check connectivity to Cloudflare endpoints
   - Verify DNS resolution for R2 endpoints
   - Check firewall and proxy settings

### Health Check

Verify Cloudflare R2 connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=cloudflare"
```

### Cloudflare CLI Testing

Test your R2 configuration using wrangler CLI:

```bash
# Install Wrangler
npm install -g wrangler

# Login to Cloudflare
wrangler login

# List R2 buckets
wrangler r2 bucket list

# Upload file to R2
wrangler r2 object put your-bucket/test.txt --file ./test.txt

# Download file from R2
wrangler r2 object get your-bucket/test.txt --file ./downloaded-test.txt

# List objects in bucket
wrangler r2 object list your-bucket
```

### S3-Compatible CLI Testing

Since R2 is S3-compatible, you can use AWS CLI:

```bash
# Configure AWS CLI for R2
aws configure set aws_access_key_id YOUR_R2_ACCESS_KEY
aws configure set aws_secret_access_key YOUR_R2_SECRET_KEY
aws configure set region auto

# Set R2 endpoint
export AWS_ENDPOINT_URL_S3=https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com

# List buckets
aws s3 ls --endpoint-url $AWS_ENDPOINT_URL_S3

# Upload file
aws s3 cp test.txt s3://your-bucket/test.txt --endpoint-url $AWS_ENDPOINT_URL_S3

# Download file
aws s3 cp s3://your-bucket/test.txt downloaded.txt --endpoint-url $AWS_ENDPOINT_URL_S3
```

## Monitoring and Analytics

### Cloudflare Dashboard Analytics
Monitor through Cloudflare Dashboard:
- Request volume and patterns
- Bandwidth usage
- Storage usage
- Geographic distribution of requests
- Error rates and response times

### R2 Metrics
Track these key R2 metrics:
- **Storage Usage**: Total stored data
- **Operations**: PUT, GET, DELETE operations
- **Bandwidth**: Data transfer (ingress only for billing)
- **Requests**: Total number of requests

### Custom Monitoring
Implement application-level monitoring:
- Upload/download success rates
- Response times
- Error tracking
- Usage patterns by user/application

## Provider Type

When using the storage factory or API endpoints, use provider type: `"cloudflare"` (note: this may be mapped to use the S3 provider internally with R2-specific configuration)

## Migration

### From AWS S3 to Cloudflare R2

1. **Assessment**: Analyze S3 usage patterns and costs
2. **R2 Setup**: Create R2 buckets and configure API tokens
3. **Data Transfer**: Use S3-compatible tools for migration
4. **Configuration Update**: Update application configuration
5. **Testing**: Thoroughly test functionality
6. **DNS/CDN Update**: Update custom domains if needed
7. **Monitoring**: Monitor migration progress and performance

### Migration Tools

#### Using rclone
```bash
# Install rclone
curl https://rclone.org/install.sh | sudo bash

# Configure S3 source
rclone config create s3-source s3 \
    access_key_id=YOUR_AWS_ACCESS_KEY \
    secret_access_key=YOUR_AWS_SECRET_KEY \
    region=us-east-1

# Configure R2 destination
rclone config create r2-dest s3 \
    access_key_id=YOUR_R2_ACCESS_KEY \
    secret_access_key=YOUR_R2_SECRET_KEY \
    endpoint=https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com \
    region=auto

# Perform migration
rclone sync s3-source:source-bucket r2-dest:destination-bucket --progress
```

#### Using AWS CLI with R2
```bash
# Sync from S3 to R2
aws s3 sync s3://source-bucket s3://destination-bucket \
    --source-region us-east-1 \
    --endpoint-url https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com
```

## Integration Examples

### With Cloudflare Workers
```javascript
// Cloudflare Worker to process R2 uploads
export default {
  async fetch(request, env) {
    if (request.method === 'PUT') {
      const object = await env.MY_BUCKET.put(key, request.body);
      return new Response(`Put ${key} successfully!`);
    }
    
    if (request.method === 'GET') {
      const object = await env.MY_BUCKET.get(key);
      return new Response(object.body);
    }
  }
}
```

### Custom Domain with Cloudflare Pages
Serve R2 content through Cloudflare Pages for additional features like serverless functions and advanced caching rules.