# Azure Blob Storage Provider Configuration

## Overview

Azure Blob Storage is Microsoft's object storage solution for the cloud. It's optimized for storing massive amounts of unstructured data and provides excellent integration with the Microsoft Azure ecosystem.

**When to use Azure Blob Storage Provider:**
- Applications hosted on Microsoft Azure
- Enterprise environments using Microsoft ecosystem
- Applications requiring integration with Azure services (Azure Functions, Logic Apps, etc.)
- Compliance requirements specific to Azure regions
- Organizations with existing Azure subscriptions and credits

**When to consider alternatives:**
- Applications hosted entirely on AWS (consider S3)
- Cost-sensitive applications with high egress traffic (consider Cloudflare R2)
- Self-hosted requirements (consider MinIO)
- Applications primarily using Google services (consider Firebase Storage)

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Azure Blob Storage Configuration
azure:
    accountName: 'yourstorageaccount'        # Azure Storage account name
    accountKey: 'your-account-key...'        # Azure Storage account key
    containerName: 'media-container'         # Azure Blob container name
    serviceUrl: ''                           # Optional: Custom service URL (e.g., Azurite)
```

## Environment Variables

You can also configure Azure Blob Storage using environment variables (recommended for production):

- `AZURE_ACCOUNT_NAME`: Azure Storage account name
- `AZURE_ACCOUNT_KEY`: Azure Storage account key
- `AZURE_CONTAINER_NAME`: Azure Blob container name
- `AZURE_SERVICE_URL`: Custom service URL (optional)

## Configuration Examples

### Production Environment
```yaml
azure:
    accountName: '${AZURE_STORAGE_ACCOUNT}'  # Use environment variables
    accountKey: '${AZURE_STORAGE_KEY}'
    containerName: 'production-media'
    serviceUrl: ''                           # Use Azure's default endpoints
```

### Development with Azurite Emulator
```yaml
azure:
    accountName: 'devstoreaccount1'          # Default Azurite account
    accountKey: 'Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=='
    containerName: 'dev-media'
    serviceUrl: 'http://127.0.0.1:10000/devstoreaccount1'  # Azurite endpoint
```

### Multi-Environment Setup
```yaml
# config.production.yaml
azure:
    accountName: 'prodstorageaccount'
    accountKey: '${AZURE_PROD_STORAGE_KEY}'
    containerName: 'production-media'
    serviceUrl: ''

# config.staging.yaml  
azure:
    accountName: 'stagingstorageaccount'
    accountKey: '${AZURE_STAGING_STORAGE_KEY}'
    containerName: 'staging-media'
    serviceUrl: ''
```

### Regional Configurations

#### US East
```yaml
azure:
    accountName: 'usmediastore'              # Storage account in US East
    accountKey: '${AZURE_STORAGE_KEY}'
    containerName: 'media-us-east'
    serviceUrl: ''
```

#### Europe West
```yaml
azure:
    accountName: 'eumediastore'              # Storage account in Europe West
    accountKey: '${AZURE_STORAGE_KEY}'
    containerName: 'media-eu-west'
    serviceUrl: ''
```

## Azure Storage Account Types

Choose the appropriate storage account type for your needs:

| Type | Use Case | Performance |
|------|----------|-------------|
| **Standard_LRS** | Low cost, local redundancy | Standard |
| **Standard_GRS** | Geographic redundancy | Standard |
| **Standard_ZRS** | Zone redundancy | Standard |
| **Premium_LRS** | High performance | Premium |

## Features

The Azure Blob Storage provider supports all standard storage operations:

- **Upload**: Upload files to Azure Blob containers with block blob support
- **Download**: Download files from Azure Blob containers
- **Delete**: Remove files from Azure Blob containers
- **GetURL**: Get public URLs for files (if container allows public access)
- **GetSignedURL**: Generate Shared Access Signature (SAS) URLs for secure access
- **GetObject**: Retrieve blob metadata and properties
- **CheckHealth**: Verify connection to Azure Storage and container access

## Storage Tiers

Azure Blob Storage offers different access tiers for cost optimization:

- **Hot**: Frequently accessed data, highest storage cost, lowest access cost
- **Cool**: Infrequently accessed data, lower storage cost, higher access cost
- **Archive**: Rarely accessed data, lowest storage cost, highest access and retrieval cost

## Container Configuration

### Public Access Levels
- **Private**: No anonymous access (default)
- **Blob**: Anonymous read access for blobs only
- **Container**: Anonymous read access for containers and blobs

### Container Properties
```json
{
  "publicAccessLevel": "None",
  "metadata": {
    "application": "m3-storage",
    "environment": "production"
  }
}
```

## Requirements

- Azure subscription with Storage Account
- Storage Account with appropriate access tier
- Container created in the storage account
- Valid account name and account key
- Network connectivity to Azure endpoints

## Authentication Methods

### Account Key (Current Implementation)
```yaml
azure:
    accountName: 'mystorageaccount'
    accountKey: 'base64-encoded-key...'      # Primary or secondary key
```

### Connection String (Alternative)
```yaml
azure:
    connectionString: 'DefaultEndpointsProtocol=https;AccountName=mystorageaccount;AccountKey=key==;EndpointSuffix=core.windows.net'
```

### Managed Identity (For Azure-hosted applications)
```yaml
azure:
    accountName: 'mystorageaccount'
    # Use managed identity - no key required
    useManagedIdentity: true
```

## Security Considerations

1. **Access Keys Management**:
   - Store account keys as environment variables
   - Rotate storage account keys regularly
   - Use Azure Key Vault for sensitive credentials
   - Consider using Managed Identity for Azure-hosted apps

2. **Network Security**:
   - Configure storage account firewalls
   - Use private endpoints for internal traffic
   - Enable secure transfer (HTTPS only)
   - Configure allowed IP ranges

3. **Container Security**:
   - Set appropriate public access levels
   - Use Shared Access Signatures (SAS) for granular access
   - Enable container-level access policies
   - Configure CORS policies for web applications

4. **Encryption**:
   - Enable encryption at rest (enabled by default)
   - Use customer-managed keys (BYOK) if required
   - Enable encryption in transit (HTTPS)

## Performance Optimization

1. **Blob Types**:
   - Use Block Blobs for most scenarios
   - Consider Page Blobs for VHD files
   - Use Append Blobs for logging scenarios

2. **Parallel Uploads**: Automatically handled for large files
3. **CDN Integration**: Use Azure CDN for global content delivery
4. **Hot/Cool/Archive Tiers**: Choose appropriate tier for access patterns
5. **Request Optimization**: Batch operations when possible

## Cost Optimization

1. **Storage Tiers**: Use appropriate tier (Hot/Cool/Archive)
2. **Lifecycle Management**: Automatically transition blobs between tiers
3. **Data Redundancy**: Choose appropriate redundancy level
4. **Monitoring**: Use Azure Cost Management to track expenses
5. **Reserved Capacity**: Consider reserved capacity for predictable workloads

### Lifecycle Management Policy Example
```json
{
  "rules": [
    {
      "name": "mediaLifecycle",
      "type": "Lifecycle",
      "definition": {
        "filters": {
          "blobTypes": ["blockBlob"],
          "prefixMatch": ["media/"]
        },
        "actions": {
          "baseBlob": {
            "tierToCool": {
              "daysAfterModificationGreaterThan": 30
            },
            "tierToArchive": {
              "daysAfterModificationGreaterThan": 90
            },
            "delete": {
              "daysAfterModificationGreaterThan": 365
            }
          }
        }
      }
    }
  ]
}
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**:
   - Verify storage account name and key
   - Check if account key has been rotated
   - Ensure proper permissions on storage account
   - Verify connection string format

2. **Container Not Found**:
   - Verify container name spelling and case
   - Ensure container exists in the storage account
   - Check container permissions

3. **Access Denied**:
   - Verify storage account firewall settings
   - Check public access level configuration
   - Ensure proper SAS token permissions
   - Verify IP allowlist configuration

4. **Network Connectivity**:
   - Check network connectivity to Azure endpoints
   - Verify DNS resolution for storage endpoints
   - Check firewall and proxy settings

### Health Check

Verify Azure Blob Storage connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=azure"
```

### Azure CLI Testing

Test your Azure Storage configuration using Azure CLI:

```bash
# Login to Azure
az login

# Set subscription
az account set --subscription "your-subscription-id"

# List storage accounts
az storage account list

# Test container access
az storage container list --account-name yourstorageaccount --account-key yourkey

# Upload test file
az storage blob upload \
    --account-name yourstorageaccount \
    --account-key yourkey \
    --container-name your-container \
    --name test.txt \
    --file ./test.txt

# Download test file
az storage blob download \
    --account-name yourstorageaccount \
    --account-key yourkey \
    --container-name your-container \
    --name test.txt \
    --file ./downloaded-test.txt
```

## Monitoring and Logging

### Azure Monitor Metrics
Monitor these key Azure Storage metrics:
- **Transactions**: Total number of requests
- **Ingress**: Data uploaded to storage account
- **Egress**: Data downloaded from storage account
- **Availability**: Service availability percentage
- **Success E2E Latency**: End-to-end latency for successful requests

### Storage Analytics Logging
Enable logging for detailed request information:
- **Read requests**: GET, HEAD operations
- **Write requests**: PUT, POST operations
- **Delete requests**: DELETE operations

### Application Insights Integration
For detailed application-level monitoring and debugging.

## Azurite Emulator (Development)

For local development, use Azurite emulator:

```bash
# Install Azurite
npm install -g azurite

# Start Azurite
azurite --silent --location c:\azurite --debug c:\azurite\debug.log

# Default connection details
# Account name: devstoreaccount1
# Account key: Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==
# Blob service: http://127.0.0.1:10000/devstoreaccount1
```

## Provider Type

When using the storage factory or API endpoints, use provider type: `"azure"`

## Migration

### From Other Providers to Azure Blob Storage

1. **Assessment**: Analyze current storage usage and requirements
2. **Storage Account Setup**: Create Azure Storage Account with appropriate configuration
3. **Container Setup**: Create containers with proper access levels
4. **Data Transfer**: Use Azure Data Factory, AzCopy, or custom migration scripts
5. **Application Update**: Update configuration to use Azure provider
6. **Testing**: Thoroughly test the migration
7. **Cutover**: Switch traffic to Azure Blob Storage

### AzCopy Migration Example
```bash
# Copy from AWS S3 to Azure Blob Storage
azcopy copy 'https://s3.amazonaws.com/mybucket/*' 'https://mystorageaccount.blob.core.windows.net/mycontainer' --recursive

# Copy from local storage to Azure Blob Storage
azcopy copy '/local/path/*' 'https://mystorageaccount.blob.core.windows.net/mycontainer' --recursive
```

## Integration with Azure Services

### Azure CDN
Configure Azure CDN for global content delivery:
```yaml
# CDN endpoint configuration
cdnProfile: "my-cdn-profile"
cdnEndpoint: "https://myendpoint.azureedge.net"
```

### Azure Functions
Integrate with Azure Functions for serverless processing:
```csharp
[FunctionName("ProcessBlob")]
public static void Run(
    [BlobTrigger("media-container/{name}")] Stream blob,
    string name,
    ILogger log)
{
    log.LogInformation($"Processing blob: {name}");
}
```

### Logic Apps
Use Logic Apps for workflow automation with blob triggers and actions.