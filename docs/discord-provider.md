# Discord Storage Provider Configuration

## Overview

The Discord Storage provider is an experimental and unique approach that uses Discord channels for file storage. This provider leverages Discord's file attachment capabilities to store files, making it a creative solution for small-scale projects, prototypes, or educational purposes.

**⚠️ Important Notice**: This provider is intended for experimental, educational, or personal use only. It should NOT be used in production environments or for commercial applications due to Discord's Terms of Service limitations and potential reliability issues.

**When to use Discord Storage Provider:**
- Experimental projects and proof-of-concepts
- Educational purposes to demonstrate creative solutions
- Small personal applications with minimal storage needs
- Prototyping applications where traditional storage isn't available
- Creative projects that want to leverage Discord's infrastructure

**When NOT to use Discord Storage Provider:**
- Production applications
- Commercial applications
- Applications requiring high reliability or uptime guarantees
- Applications with significant storage requirements
- Applications requiring fast or consistent performance

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Discord Storage Configuration
discord:
    botToken: 'your-discord-bot-token'          # Discord Bot Token
    channelID: 'your-discord-channel-id'        # Discord Channel ID for file storage
    webhookURL: 'your-discord-webhook-url'      # Optional: Discord Webhook URL
```

## Environment Variables

You can also configure Discord Storage using environment variables (recommended):

- `DISCORD_BOT_TOKEN`: Discord Bot Token
- `DISCORD_CHANNEL_ID`: Discord Channel ID for file storage
- `DISCORD_WEBHOOK_URL`: Discord Webhook URL (optional)

## Configuration Examples

### Basic Bot Configuration
```yaml
discord:
    botToken: '${DISCORD_BOT_TOKEN}'            # Use environment variables
    channelID: '123456789012345678'             # Channel ID where files will be stored
    webhookURL: ''                              # Not using webhook
```

### Webhook Configuration (Alternative)
```yaml
discord:
    botToken: ''                                # Not using bot token
    channelID: '123456789012345678'
    webhookURL: '${DISCORD_WEBHOOK_URL}'        # Using webhook instead
```

### Development Environment
```yaml
discord:
    botToken: 'your-dev-bot-token'
    channelID: '987654321098765432'             # Development channel
    webhookURL: ''
```

### Multiple Channels (Different Environments)
```yaml
# config.development.yaml
discord:
    botToken: '${DISCORD_DEV_BOT_TOKEN}'
    channelID: '111111111111111111'             # Dev channel
    webhookURL: ''

# config.testing.yaml
discord:
    botToken: '${DISCORD_TEST_BOT_TOKEN}'
    channelID: '222222222222222222'             # Test channel
    webhookURL: ''
```

## Discord Bot Setup

### Creating a Discord Bot

1. **Go to Discord Developer Portal**: https://discord.com/developers/applications
2. **Create New Application**
3. **Go to Bot section**
4. **Create Bot and get Bot Token**
5. **Configure Bot Permissions**

### Required Bot Permissions

Your Discord bot needs these permissions:

```json
{
  "permissions": [
    "VIEW_CHANNEL",          // View the storage channel
    "SEND_MESSAGES",         // Send messages (file uploads)
    "ATTACH_FILES",          // Attach files to messages
    "READ_MESSAGE_HISTORY",  // Read message history to find files
    "MANAGE_MESSAGES"        // Delete messages when files are deleted
  ]
}
```

### Bot Permissions Value
The permission integer for these permissions is: `75776`

### Inviting Bot to Server
Use this URL template to invite your bot:
```
https://discord.com/api/oauth2/authorize?client_id=YOUR_BOT_CLIENT_ID&permissions=75776&scope=bot
```

## Channel Setup

### Creating Storage Channel

1. **Create a dedicated channel** for file storage (e.g., #file-storage)
2. **Make it private** or restrict access as needed
3. **Add your bot** to the channel
4. **Copy the Channel ID** (Enable Developer Mode in Discord settings)

### Channel Recommendations
- Use a private channel to prevent clutter
- Consider using a dedicated server for storage
- Set appropriate permissions to prevent unauthorized access
- Use descriptive channel names (e.g., #app-storage, #media-files)

## Features

The Discord Storage provider supports basic storage operations:

- **Upload**: Upload files as Discord attachments (up to 8MB for regular users, 50MB for Nitro)
- **Download**: Download files from Discord CDN URLs
- **Delete**: Delete files by removing Discord messages
- **GetURL**: Get Discord CDN URLs for files
- **GetSignedURL**: Returns the same as GetURL (Discord URLs are already time-limited)
- **GetObject**: Retrieve basic file metadata from Discord messages
- **CheckHealth**: Verify bot connection and channel access

## Limitations

### File Size Limits
- **Regular Discord**: 8MB per file
- **Discord Nitro**: 50MB per file
- **Server Boosts**: May increase limits

### Rate Limits
- Discord API rate limits apply
- Approximately 5 requests per 5 seconds per channel
- Global rate limits may affect performance

### Terms of Service
- Must comply with Discord's Terms of Service
- Not intended for commercial file storage
- Subject to Discord's usage policies

### Reliability Concerns
- Files depend on Discord's infrastructure
- No guaranteed uptime or availability
- Discord could change policies or rate limits
- Messages could be deleted by server administrators

## Requirements

- Discord account
- Discord server (guild) where you can add bots
- Discord bot with appropriate permissions
- Channel where the bot can send messages and attach files
- Network connectivity to Discord's API and CDN

## Security Considerations

1. **Bot Token Security**:
   - Store bot token securely using environment variables
   - Never commit bot tokens to version control
   - Regenerate tokens if compromised
   - Monitor bot activity in Discord

2. **Channel Security**:
   - Use private channels for sensitive files
   - Restrict channel access to necessary users only
   - Monitor channel activity and access
   - Consider using dedicated servers for storage

3. **File Security**:
   - Discord URLs are publicly accessible if known
   - Files are stored on Discord's CDN
   - No additional encryption provided by Discord
   - Consider encrypting files before upload if needed

4. **Access Control**:
   - Implement application-level access controls
   - Don't rely solely on Discord permissions
   - Log and monitor file operations
   - Validate file types and content

## Performance Considerations

### Upload Performance
- Limited by Discord's rate limits
- Large files take longer to upload
- Network connectivity affects performance
- Consider chunking or queuing for multiple files

### Download Performance
- Files served from Discord's CDN
- Generally fast download speeds
- Subject to Discord's CDN availability
- URLs may have expiration times

### Storage Efficiency
- Each file creates a Discord message
- Metadata stored in message content
- Channel history can become cluttered
- Consider periodic cleanup strategies

## Troubleshooting

### Common Issues

1. **Bot Authentication Failed**:
   - Verify bot token is correct and not expired
   - Ensure bot is added to the server
   - Check bot permissions in the channel
   - Verify bot is not banned or restricted

2. **Channel Not Found**:
   - Verify channel ID is correct
   - Ensure bot has access to the channel
   - Check if channel still exists
   - Verify bot is in the correct server

3. **File Upload Failed**:
   - Check file size limits (8MB/50MB)
   - Verify bot has ATTACH_FILES permission
   - Check Discord rate limits
   - Ensure stable network connection

4. **Rate Limited**:
   - Implement proper rate limit handling
   - Add delays between operations
   - Consider queuing uploads
   - Monitor Discord API responses

### Health Check

Verify Discord connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=discord"
```

### Testing Bot Connection

Test your Discord bot configuration:

```python
# Example Python test script
import discord
import asyncio

async def test_bot():
    client = discord.Client(intents=discord.Intents.default())
    
    @client.event
    async def on_ready():
        print(f'Bot logged in as {client.user}')
        
        # Test channel access
        channel = client.get_channel(YOUR_CHANNEL_ID)
        if channel:
            print(f'Channel found: {channel.name}')
            
            # Test message sending
            message = await channel.send('Test message')
            print(f'Message sent: {message.id}')
            
            # Clean up
            await message.delete()
        
        await client.close()
    
    await client.start('YOUR_BOT_TOKEN')

# Run test
asyncio.run(test_bot())
```

## Monitoring and Logging

### Application Monitoring
Track these metrics:
- Upload/download success rates
- File storage usage (number of messages)
- Rate limit encounters
- Error rates and types
- Response times

### Discord Monitoring
- Bot uptime and connectivity
- Channel activity and message count
- Rate limit status
- Permission issues

### Custom Logging
```go
// Example logging structure
type DiscordOperationLog struct {
    Timestamp   time.Time `json:"timestamp"`
    Operation   string    `json:"operation"`
    FileName    string    `json:"fileName"`
    FileSize    int64     `json:"fileSize"`
    Success     bool      `json:"success"`
    MessageID   string    `json:"messageID,omitempty"`
    ChannelID   string    `json:"channelID"`
    Error       string    `json:"error,omitempty"`
    Duration    int64     `json:"duration_ms"`
}
```

## Provider Type

When using the storage factory or API endpoints, use provider type: `"discord"`

## Best Practices

### File Management
1. **Organize by prefixes**: Use consistent file naming conventions
2. **Clean up regularly**: Delete old or unnecessary files
3. **Monitor channel growth**: Keep track of message count
4. **Backup important files**: Don't rely solely on Discord for critical data

### Error Handling
1. **Implement retries**: Handle temporary Discord outages
2. **Rate limit respect**: Implement proper backoff strategies  
3. **Graceful degradation**: Have fallback options when Discord is unavailable
4. **User feedback**: Inform users of Discord-related limitations

### Development Guidelines
```go
// Example error handling
func (d *DiscordProvider) uploadWithRetry(ctx context.Context, data []byte, fileName string) error {
    const maxRetries = 3
    const baseDelay = time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := d.upload(ctx, data, fileName)
        if err == nil {
            return nil
        }
        
        // Check if it's a rate limit error
        if isRateLimit(err) {
            delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
            time.Sleep(delay)
            continue
        }
        
        return err // Non-retryable error
    }
    
    return fmt.Errorf("upload failed after %d attempts", maxRetries)
}
```

## Migration Considerations

### From Discord to Other Providers

When outgrowing Discord storage:

1. **Export file inventory**: List all stored files and their metadata
2. **Download all files**: Bulk download from Discord CDN
3. **Configure new provider**: Set up proper storage provider
4. **Migrate data**: Upload to new provider
5. **Update application**: Switch storage configuration
6. **Cleanup**: Remove old Discord messages if desired

### Migration Script Example
```python
# Example migration script
import discord
import requests
import os

async def migrate_from_discord(channel_id, destination_path):
    client = discord.Client(intents=discord.Intents.default())
    
    @client.event
    async def on_ready():
        channel = client.get_channel(channel_id)
        
        async for message in channel.history(limit=None):
            for attachment in message.attachments:
                # Download file
                response = requests.get(attachment.url)
                
                # Save to local storage
                file_path = os.path.join(destination_path, attachment.filename)
                with open(file_path, 'wb') as f:
                    f.write(response.content)
                
                print(f'Downloaded: {attachment.filename}')
        
        await client.close()
    
    await client.start('YOUR_BOT_TOKEN')
```

## Educational Use Cases

This provider serves as an excellent example for:

1. **API Integration**: Learning to work with REST APIs
2. **Creative Solutions**: Thinking outside the box for storage solutions
3. **Rate Limiting**: Understanding and handling API rate limits
4. **Error Handling**: Implementing robust error handling strategies
5. **Alternative Architecture**: Exploring unconventional storage approaches

## Disclaimer

**Important**: The Discord Storage Provider is provided as-is for experimental and educational purposes. It is not recommended for production use and may violate Discord's Terms of Service if used inappropriately. Always review and comply with Discord's current Terms of Service and Community Guidelines.

Users are responsible for ensuring their use of this provider complies with all applicable terms and regulations.