# Storage Providers

This document describes how to configure and use different storage providers in the media management system.

## Available Providers

- **Local Storage** - Store files on local filesystem
- **AWS S3** - Store files on Amazon S3 or S3-compatible services
- **Cloudflare R2** - Store files on Cloudflare R2 (S3-compatible)
- **Firebase Storage** - Store files on Google Firebase Cloud Storage
- **Azure Blob Storage** - Store files on Microsoft Azure Blob Storage
- **Discord** - Store files using Discord channels (unique approach)

## Discord Storage Provider

The Discord storage provider is a unique solution that uses Discord channels as a storage backend. This can be useful for:

- Prototype projects
- Small applications with limited storage needs
- Free storage solution (within Discord's limits)
- Applications that already integrate with Discord

### Configuration

Add the following to your `config.yaml`:

```yaml
# Discord Configuration (File Storage)
discord:
    botToken: 'YOUR_DISCORD_BOT_TOKEN'      # Discord Bot Token
    channelID: 'YOUR_CHANNEL_ID'            # Discord Channel ID for file storage
    webhookURL: ''                          # Optional: Discord Webhook URL for notifications
```

### Setting up Discord Bot

1. **Create a Discord Application**
   - Go to https://discord.com/developers/applications
   - Click "New Application"
   - Give it a name (e.g., "Media Storage Bot")

2. **Create a Bot**
   - In your application, go to the "Bot" section
   - Click "Add Bot"
   - Copy the bot token and add it to your config

3. **Set Bot Permissions**
   The bot needs the following permissions:
   - Read Messages
   - Send Messages
   - Attach Files
   - Manage Messages (for deletion)
   - View Channel

4. **Invite Bot to Server**
   - Go to OAuth2 → URL Generator
   - Select "bot" scope
   - Select the required permissions
   - Use the generated URL to invite the bot to your server

5. **Get Channel ID**
   - Enable Developer Mode in Discord settings
   - Right-click on the channel you want to use
   - Select "Copy ID"

### Usage Examples

#### Upload a file using Discord storage

```go
// Using the media service with Discord provider
mediaService := // ... get your media service instance

// Upload file specifying Discord as the provider
media, err := mediaService.UploadFile(
    ctx,
    userID,
    fileHeader,
    "discord", // Provider name
    "image",   // Media type hint
)
if err != nil {
    log.Error("Failed to upload file", err)
    return
}

log.Info("File uploaded successfully", 
    "mediaID", media.ID,
    "url", media.URL,
    "provider", media.Provider)
```

#### Direct usage of Discord provider

```go
import (
    "github.com/lugondev/m3-storage/internal/adapters/discord"
    "github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// Create Discord provider directly
config := map[string]interface{}{
    "token":      "YOUR_BOT_TOKEN",
    "channel_id": "YOUR_CHANNEL_ID",
}

provider, err := discord.NewDiscordProvider(config)
if err != nil {
    log.Fatal("Failed to create Discord provider", err)
}

// Upload a file
file, _ := os.Open("example.jpg")
defer file.Close()

fileObject, err := provider.Upload(
    ctx,
    "user123/images/20240123/example.jpg", // Storage key
    file,
    1024*1024, // File size
    &port.UploadOptions{
        ContentType: "image/jpeg",
    },
)
if err != nil {
    log.Fatal("Failed to upload file", err)
}

log.Info("File uploaded", "url", fileObject.URL)
```

### Limitations

1. **File Size Limits**
   - Discord has a file size limit (8MB for free servers, 50MB for Nitro servers)
   - Large files will be rejected

2. **Rate Limits**
   - Discord API has rate limits for message posting
   - High-frequency uploads may be throttled

3. **Persistence**
   - Files are stored as Discord messages
   - If messages are deleted, files are lost
   - Channel history limits may affect old files

4. **Public Access**
   - All files uploaded to Discord get public URLs
   - URLs are long-lived but not permanent
   - Discord may change URL format in the future

5. **Search Performance**
   - Finding files requires searching through channel messages
   - Performance degrades with large numbers of files
   - Limited to 100 messages per search request

### Best Practices

1. **Use Dedicated Channels**
   - Create dedicated channels only for file storage
   - Don't mix storage with regular Discord conversations

2. **File Organization**
   - Use meaningful file names and organize by user/date
   - Consider the storage key structure: `{userID}/{type}/{date}/{filename}`

3. **Backup Strategy**
   - Discord storage should not be used for critical data
   - Consider it as a temporary or cache storage solution
   - Implement backup strategies for important files

4. **Monitor Usage**
   - Keep track of file counts and sizes
   - Monitor Discord API usage and rate limits
   - Set up alerts for storage failures

### Error Handling

The Discord provider handles several error scenarios:

- **Channel Access**: Verifies bot has access to the specified channel
- **File Size**: Returns appropriate errors for oversized files
- **Rate Limiting**: Handles Discord API rate limits gracefully
- **Network Issues**: Provides meaningful error messages for connection problems

### Security Considerations

1. **Bot Token Security**
   - Store bot tokens in environment variables, not config files
   - Rotate tokens periodically
   - Use minimal required permissions

2. **Channel Privacy**
   - Use private channels for sensitive data
   - Consider who has access to the storage channel
   - Monitor channel membership

3. **Data Sensitivity**
   - Discord storage is not suitable for highly sensitive data
   - Consider encryption for sensitive files before upload
   - Be aware of Discord's data retention policies

## Other Storage Providers

### Local Storage

```yaml
localStorage:
    path: './uploads'
    baseURL: '/files'
```

### AWS S3

```yaml
s3:
    accessKeyID: 'YOUR_ACCESS_KEY'
    secretAccessKey: 'YOUR_SECRET_KEY'
    region: 'us-east-1'
    bucketName: 'your-bucket-name'
```

### Cloudflare R2

```yaml
cloudflare:
    accountID: 'YOUR_ACCOUNT_ID'
    accessKeyID: 'YOUR_R2_ACCESS_KEY'
    secretAccessKey: 'YOUR_R2_SECRET_KEY'
    bucketName: 'your-r2-bucket'
```

### Firebase Storage

```yaml
firestore:
    projectID: 'your-project-id'
    credentialsFile: 'path/to/serviceAccountKey.json'
    bucketName: 'your-project-id.appspot.com'
```

## Provider Selection

You can specify which provider to use when uploading files:

```go
// Upload to Discord
media, err := mediaService.UploadFile(ctx, userID, fileHeader, "discord", "image")

// Upload to S3
media, err := mediaService.UploadFile(ctx, userID, fileHeader, "s3", "image")

// Upload to default provider (local)
media, err := mediaService.UploadFile(ctx, userID, fileHeader, "", "image")
```

## Configuration via Environment Variables

All configuration can be overridden using environment variables:

```bash
# Discord
export DISCORD_BOT_TOKEN="your_bot_token"
export DISCORD_CHANNEL_ID="your_channel_id"

# S3
export S3_ACCESS_KEY_ID="your_access_key"
export S3_SECRET_ACCESS_KEY="your_secret_key"
export S3_BUCKET_NAME="your_bucket"

# Local Storage
export LOCAL_STORAGE_PATH="/var/uploads"
export LOCAL_STORAGE_BASE_URL="/files"
