# Local Storage Provider Configuration

## Overview

The Local Storage provider stores files directly on the server's filesystem. This is the simplest storage provider and is ideal for development environments, small-scale deployments, or when you want full control over your file storage.

**When to use Local Storage Provider:**
- Development and testing environments
- Small-scale applications with limited storage needs
- Applications where you have full control over the server
- Cost-sensitive deployments without cloud storage budget
- Single-server deployments

**When NOT to use Local Storage Provider:**
- Production applications requiring high availability
- Distributed or multi-server deployments
- Applications requiring CDN or global file distribution
- When server storage capacity is limited

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Local Storage Configuration
localStorage:
    path: './uploads'                        # Directory path for storing files
    baseURL: '/files'                        # Base URL for accessing files publicly
    signedUrlExpiry: '24h'                   # Signed URL expiration time
    signedUrlSecret: 'your-secret-key'       # Secret key for signing URLs
```

## Environment Variables

You can also configure Local Storage using environment variables:

- `LOCAL_STORAGE_PATH`: Directory path for storing files
- `LOCAL_STORAGE_BASE_URL`: Base URL for accessing files publicly
- `LOCAL_STORAGE_SIGNED_URL_EXPIRY`: Signed URL expiration time
- `LOCAL_STORAGE_SIGNED_URL_SECRET`: Secret key for signing URLs

## Configuration Examples

### Development Environment
```yaml
localStorage:
    path: './uploads'
    baseURL: '/files'
    signedUrlExpiry: '1h'
    signedUrlSecret: 'dev-secret-key'
```

### Production Environment
```yaml
localStorage:
    path: '/var/www/storage'
    baseURL: 'https://yourdomain.com/files'
    signedUrlExpiry: '24h'
    signedUrlSecret: '${LOCAL_STORAGE_SIGNED_URL_SECRET}'  # Use environment variable
```

### Docker Container Setup
```yaml
localStorage:
    path: '/app/uploads'                     # Path inside container
    baseURL: 'https://api.example.com/files'
    signedUrlExpiry: '12h'
    signedUrlSecret: '${LOCAL_STORAGE_SIGNED_URL_SECRET}'
```

### Custom Storage Directory
```yaml
localStorage:
    path: '/mnt/storage/media-files'         # Custom mount point
    baseURL: 'https://cdn.example.com/media'
    signedUrlExpiry: '7d'                    # 7 days for long-term access
    signedUrlSecret: 'production-secret-key'
```

## Features

The Local Storage provider supports all standard storage operations:

- **Upload**: Store files in the specified local directory
- **Download**: Serve files directly from the filesystem
- **Delete**: Remove files from the local directory
- **GetURL**: Generate public URLs for file access
- **GetSignedURL**: Create time-limited signed URLs for secure access
- **GetObject**: Retrieve file metadata and information
- **CheckHealth**: Verify directory access and permissions

## Directory Structure

Files are organized in the following structure:
```
uploads/
├── {user-id}/
│   ├── image/
│   │   ├── {file-id}.jpg
│   │   └── {file-id}.png
│   ├── video/
│   │   └── {file-id}.mp4
│   └── document/
│       └── {file-id}.pdf
```

## Requirements

- Write permissions to the storage directory
- Sufficient disk space for your files
- Web server configuration to serve files (if using public URLs)

## Security Considerations

1. **File Permissions**: Ensure proper file system permissions
   ```bash
   chmod 755 /path/to/uploads
   chown www-data:www-data /path/to/uploads  # For web servers
   ```

2. **Directory Traversal**: The provider automatically prevents directory traversal attacks

3. **Signed URLs**: Use signed URLs for secure, time-limited access to private files

4. **File Validation**: Always validate file types and sizes before storage

## Performance Considerations

- **Disk I/O**: Performance depends on your disk speed (SSD recommended)
- **Concurrent Access**: File system handles concurrent reads/writes
- **Backup**: Implement regular backup strategies for important data
- **Monitoring**: Monitor disk space usage to prevent storage full errors

## Backup and Recovery

### Regular Backup
```bash
# Simple backup script
tar -czf backup-$(date +%Y%m%d).tar.gz /path/to/uploads/
```

### Rsync Backup
```bash
# Sync to remote server
rsync -avz /path/to/uploads/ user@backup-server:/backup/uploads/
```

### Docker Volume Backup
```bash
# Backup Docker volume
docker run --rm -v storage_volume:/data -v $(pwd):/backup alpine tar czf /backup/storage-backup.tar.gz /data
```

## Troubleshooting

### Common Issues

1. **Permission Denied**:
   - Check directory permissions
   - Ensure the application user has write access
   - Verify SELinux/AppArmor policies if applicable

2. **Disk Space Full**:
   - Monitor disk usage regularly
   - Implement cleanup policies for old files
   - Consider file compression for archives

3. **File Not Found**:
   - Verify file paths and directory structure
   - Check if files were moved or deleted externally

4. **Slow Performance**:
   - Use SSD storage for better I/O performance
   - Monitor disk usage and fragmentation
   - Consider file caching strategies

### Health Check

You can verify the Local Storage configuration using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=local"
```

### Testing Configuration

To test your Local Storage configuration:

1. Start your application with Local Storage configuration
2. Check the logs for initialization messages
3. Verify the storage directory exists and is writable
4. Use the health check API endpoint
5. Try uploading a test file

## Web Server Configuration

### Nginx Configuration
```nginx
location /files {
    alias /path/to/uploads;
    expires 1y;
    add_header Cache-Control "public, immutable";
    
    # Security headers
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
}
```

### Apache Configuration
```apache
Alias /files /path/to/uploads
<Directory "/path/to/uploads">
    Options -Indexes
    AllowOverride None
    Require all granted
    
    # Cache control
    ExpiresActive On
    ExpiresDefault "access plus 1 year"
</Directory>
```

## Provider Type

When using the storage factory or API endpoints, use provider type: `"local"`

## Migration

### From Local to Cloud Storage

When migrating from local storage to cloud providers:

1. **Export file list**:
   ```bash
   find /path/to/uploads -type f > file-list.txt
   ```

2. **Bulk upload to new provider** (example with AWS CLI):
   ```bash
   aws s3 sync /path/to/uploads s3://your-bucket/
   ```

3. **Update database records** to reflect new provider and URLs

4. **Verify migration** before removing local files

### From Cloud to Local Storage

1. **Download files** from cloud provider
2. **Organize in local directory structure**
3. **Update configuration** to use local provider
4. **Update database records** with new local paths