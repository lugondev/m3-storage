# Amazon S3 Storage Provider Configuration

## Overview

Amazon S3 (Simple Storage Service) is Amazon's highly scalable, reliable, and cost-effective object storage service. It's the industry standard for cloud storage and provides excellent integration with other AWS services.

**When to use Amazon S3 Provider:**
- Production applications requiring high availability and reliability
- Applications with global user base needing worldwide accessibility
- Integration with other AWS services (CloudFront, Lambda, etc.)
- Applications requiring advanced features (versioning, lifecycle policies, cross-region replication)
- Enterprise applications with compliance requirements

**When to consider alternatives:**
- Cost-sensitive applications with high egress traffic (consider Cloudflare R2)
- Applications hosted entirely on other cloud platforms (Azure, Google Cloud)
- Self-hosted or on-premise requirements (consider MinIO)

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Amazon S3 Configuration
s3:
    accessKeyID: 'AKIAIOSFODNN7EXAMPLE'     # AWS Access Key ID
    secretAccessKey: 'wJalrXUtnFEMI/...'    # AWS Secret Access Key
    region: 'us-east-1'                     # AWS Region
    bucketName: 'your-s3-bucket-name'       # S3 Bucket Name
    endpoint: ''                            # Leave empty for AWS S3
    disableSSL: false                       # Use SSL/TLS (recommended: true)
    forcePathStyle: false                   # Use virtual-hosted style URLs
```

## Environment Variables

You can also configure Amazon S3 using environment variables (recommended for production):

- `S3_ACCESS_KEY_ID`: AWS Access Key ID
- `S3_SECRET_ACCESS_KEY`: AWS Secret Access Key
- `S3_REGION`: AWS Region
- `S3_BUCKET_NAME`: S3 Bucket Name
- `S3_ENDPOINT`: Custom endpoint (leave empty for AWS S3)
- `S3_DISABLE_SSL`: Disable SSL (not recommended for production)
- `S3_FORCE_PATH_STYLE`: Force path-style addressing

## Configuration Examples

### Production Environment
```yaml
s3:
    accessKeyID: '${AWS_ACCESS_KEY_ID}'      # Use environment variables
    secretAccessKey: '${AWS_SECRET_ACCESS_KEY}'
    region: 'us-east-1'
    bucketName: 'production-media-storage'
    endpoint: ''                             # Use AWS S3 endpoints
    disableSSL: false
    forcePathStyle: false
```

### Multi-Region Setup (US East)
```yaml
s3:
    accessKeyID: '${AWS_ACCESS_KEY_ID}'
    secretAccessKey: '${AWS_SECRET_ACCESS_KEY}'
    region: 'us-east-1'                      # Virginia region
    bucketName: 'media-us-east'
    endpoint: ''
    disableSSL: false
    forcePathStyle: false
```

### Multi-Region Setup (EU West)
```yaml
s3:
    accessKeyID: '${AWS_ACCESS_KEY_ID}'
    secretAccessKey: '${AWS_SECRET_ACCESS_KEY}'
    region: 'eu-west-1'                      # Ireland region
    bucketName: 'media-eu-west'
    endpoint: ''
    disableSSL: false
    forcePathStyle: false
```

### Development Environment
```yaml
s3:
    accessKeyID: 'your-dev-access-key'
    secretAccessKey: 'your-dev-secret-key'
    region: 'us-east-1'
    bucketName: 'dev-media-storage'
    endpoint: ''
    disableSSL: false
    forcePathStyle: false
```

### S3-Compatible Services (using S3 provider)
```yaml
s3:
    accessKeyID: 'your-access-key'
    secretAccessKey: 'your-secret-key'
    region: 'us-east-1'
    bucketName: 'your-bucket'
    endpoint: 'https://s3.digitaloceanspaces.com'  # DigitalOcean Spaces
    disableSSL: false
    forcePathStyle: true                     # Required for some S3-compatible services
```

## AWS Regions

Popular AWS regions for S3:

| Region Code | Region Name | Location |
|-------------|-------------|----------|
| us-east-1 | US East (N. Virginia) | United States |
| us-west-2 | US West (Oregon) | United States |
| eu-west-1 | Europe (Ireland) | Europe |
| eu-central-1 | Europe (Frankfurt) | Europe |
| ap-southeast-1 | Asia Pacific (Singapore) | Asia |
| ap-northeast-1 | Asia Pacific (Tokyo) | Asia |

## Features

The S3 provider supports all standard storage operations:

- **Upload**: Upload files to S3 buckets with multipart upload for large files
- **Download**: Download files from S3 buckets
- **Delete**: Remove files from S3 buckets
- **GetURL**: Get public URLs for files (if bucket allows public access)
- **GetSignedURL**: Generate pre-signed URLs for secure, time-limited access
- **GetObject**: Retrieve file metadata and information
- **CheckHealth**: Verify connection to S3 and bucket access

## Advanced S3 Features

### Storage Classes
S3 offers different storage classes for cost optimization:
- **Standard**: Frequently accessed data
- **Standard-IA**: Infrequently accessed data
- **One Zone-IA**: Infrequently accessed, single AZ
- **Glacier**: Archive storage for rarely accessed data
- **Deep Archive**: Lowest cost archive storage

### Bucket Policies
Configure bucket policies for security and access control:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": "arn:aws:iam::YOUR-ACCOUNT:user/YOUR-USER"},
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    }
  ]
}
```

### CORS Configuration
For web applications accessing S3 directly:

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

- AWS Account with S3 access
- IAM user with appropriate S3 permissions
- Valid access credentials (Access Key ID and Secret Access Key)
- Existing S3 bucket or permissions to create buckets
- Internet connectivity to AWS S3 endpoints

## IAM Permissions

Minimum required IAM policy for the S3 provider:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:GetObjectVersion",
                "s3:GetBucketLocation"
            ],
            "Resource": [
                "arn:aws:s3:::your-bucket-name",
                "arn:aws:s3:::your-bucket-name/*"
            ]
        }
    ]
}
```

## Security Considerations

1. **Credentials Management**:
   - Use IAM roles when running on EC2
   - Store credentials as environment variables
   - Rotate access keys regularly
   - Use least privilege principle

2. **Bucket Security**:
   - Enable bucket versioning
   - Configure bucket policies
   - Enable access logging
   - Use MFA delete for critical buckets

3. **Encryption**:
   - Enable server-side encryption (SSE-S3, SSE-KMS, or SSE-C)
   - Use HTTPS for all transfers
   - Enable bucket encryption by default

4. **Network Security**:
   - Use VPC endpoints for private connectivity
   - Configure proper security groups
   - Consider using AWS PrivateLink

## Performance Optimization

1. **Multipart Upload**: Automatically used for large files
2. **Transfer Acceleration**: Enable for global users
3. **CloudFront CDN**: Use for content delivery
4. **Request Rate**: Optimize prefix patterns for high request rates
5. **Connection Pooling**: Handled automatically by AWS SDK

## Cost Optimization

1. **Storage Classes**: Use appropriate storage class for your use case
2. **Lifecycle Policies**: Automatically transition or delete old objects
3. **Intelligent Tiering**: Automatic cost optimization
4. **CloudWatch Metrics**: Monitor usage and costs
5. **S3 Select**: Reduce data transfer costs for queries

## Troubleshooting

### Common Issues

1. **Access Denied (403)**:
   - Verify IAM permissions
   - Check bucket policies
   - Ensure correct region configuration
   - Verify access key and secret key

2. **Bucket Not Found (404)**:
   - Verify bucket name spelling
   - Check bucket region
   - Ensure bucket exists

3. **Slow Upload/Download**:
   - Check network connectivity
   - Consider using Transfer Acceleration
   - Verify region proximity
   - Monitor CloudWatch metrics

4. **SSL/TLS Issues**:
   - Ensure disableSSL is set to false
   - Check certificate configuration
   - Verify endpoint URLs

### Health Check

Verify S3 connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=s3"
```

### AWS CLI Testing

Test your S3 configuration using AWS CLI:

```bash
# Configure AWS CLI
aws configure set aws_access_key_id YOUR_ACCESS_KEY
aws configure set aws_secret_access_key YOUR_SECRET_KEY
aws configure set region us-east-1

# List buckets
aws s3 ls

# Test upload
aws s3 cp test-file.txt s3://your-bucket-name/test-file.txt

# Test download
aws s3 cp s3://your-bucket-name/test-file.txt downloaded-file.txt
```

## Monitoring and Logging

### CloudWatch Metrics
Monitor these key S3 metrics:
- BucketRequests
- BucketBytes
- AllRequests
- GetRequests
- PutRequests

### S3 Access Logging
Enable server access logging:
```yaml
# Bucket logging configuration
LoggingEnabled:
  TargetBucket: your-log-bucket
  TargetPrefix: access-logs/
```

### CloudTrail Integration
Enable CloudTrail for API-level logging of S3 operations.

## Provider Type

When using the storage factory or API endpoints, use provider type: `"s3"`

## Migration

### From Other Providers to S3

1. **Assessment**: Analyze current storage usage and requirements
2. **Bucket Setup**: Create S3 buckets with appropriate configuration
3. **Data Transfer**: Use AWS DataSync, S3 Transfer Family, or custom scripts
4. **Application Update**: Update configuration to use S3 provider
5. **Testing**: Thoroughly test the migration
6. **Cutover**: Switch traffic to S3 storage

### S3 Cross-Region Replication

Set up cross-region replication for disaster recovery:

```json
{
  "Role": "arn:aws:iam::YOUR-ACCOUNT:role/replication-role",
  "Rules": [
    {
      "Status": "Enabled",
      "Prefix": "",
      "Destination": {
        "Bucket": "arn:aws:s3:::destination-bucket",
        "StorageClass": "STANDARD_IA"
      }
    }
  ]
}
```