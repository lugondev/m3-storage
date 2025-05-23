package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"
	"github.com/lugondev/m3-storage/internal/modules/storage/port"
)

// discordProvider implements the port.StorageProvider interface for Discord.
type discordProvider struct {
	config config.DiscordConfig
	client *http.Client
	logger logger.Logger
}

// Discord API constants
const (
	discordAPIBaseURL = "https://discord.com/api/v10"
)

// Discord API response structures
type discordMessage struct {
	ID          string              `json:"id"`
	Content     string              `json:"content"`
	Attachments []discordAttachment `json:"attachments"`
	Timestamp   string              `json:"timestamp"`
}

type discordAttachment struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	ContentType string `json:"content_type"`
}

// NewDiscordProvider creates a new Discord storage provider.
func NewDiscordProvider(config config.DiscordConfig, logger logger.Logger) (port.StorageProvider, error) {
	if config.BotToken == "" {
		return nil, errors.New("discord provider: token is required")
	}

	if config.ChannelID == "" {
		return nil, errors.New("discord provider: channel_id is required")
	}

	// Verify channel exists and is accessible
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/channels/%s", discordAPIBaseURL, config.ChannelID), nil)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bot "+config.BotToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to access channel: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discord provider: failed to access channel, status: %d", resp.StatusCode)
	}

	return &discordProvider{
		config: config,
		client: client,
	}, nil
}

// Upload uploads a file to Discord.
func (p *discordProvider) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *port.UploadOptions) (*port.FileObject, error) {
	// Create a message with the file
	filename := key

	// Create a buffer to read the file
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to read file: %w", err)
	}

	// Create a multipart form for the file upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file content
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("discord provider: failed to write file data: %w", err)
	}

	// Add the message content
	if err := writer.WriteField("content", fmt.Sprintf("File: %s", key)); err != nil {
		return nil, fmt.Errorf("discord provider: failed to write content field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("discord provider: failed to close multipart writer: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/channels/%s/messages", discordAPIBaseURL, p.config.ChannelID),
		body)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+p.config.BotToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("discord provider: failed to upload file, status: %d, response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var message discordMessage
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("discord provider: failed to parse response: %w", err)
	}

	// Get the URL of the uploaded file
	if len(message.Attachments) == 0 {
		return nil, errors.New("discord provider: no attachment URL found after upload")
	}

	fileURL := message.Attachments[0].URL
	fileSize := int64(message.Attachments[0].Size)
	fileContentType := message.Attachments[0].ContentType

	// Parse the timestamp
	lastModified := time.Now()
	if message.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, message.Timestamp); err == nil {
			lastModified = t
		}
	}

	return &port.FileObject{
		Key:          key,
		URL:          fileURL,
		Size:         fileSize,
		ContentType:  fileContentType,
		LastModified: lastModified,
		Provider:     p.ProviderType(),
	}, nil
}

// getMessages retrieves messages from a Discord channel
func (p *discordProvider) getMessages(ctx context.Context, limit int) ([]discordMessage, error) {
	url := fmt.Sprintf("%s/channels/%s/messages?limit=%d", discordAPIBaseURL, p.config.ChannelID, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+p.config.BotToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discord provider: failed to get messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discord provider: failed to get messages, status: %d", resp.StatusCode)
	}

	var messages []discordMessage
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("discord provider: failed to parse messages: %w", err)
	}

	return messages, nil
}

// findMessageWithFile finds a message containing a file with the given key
func (p *discordProvider) findMessageWithFile(ctx context.Context, key string) (*discordMessage, error) {
	messages, err := p.getMessages(ctx, 100)
	if err != nil {
		return nil, err
	}

	filePrefix := fmt.Sprintf("File: %s", key)
	for i := range messages {
		if strings.Contains(messages[i].Content, filePrefix) && len(messages[i].Attachments) > 0 {
			return &messages[i], nil
		}
	}

	return nil, errors.New("discord provider: file not found")
}

// GetURL returns the URL for a file.
func (p *discordProvider) GetURL(ctx context.Context, key string) (string, error) {
	message, err := p.findMessageWithFile(ctx, key)
	if err != nil {
		return "", err
	}

	return message.Attachments[0].URL, nil
}

// GetSignedURL returns a signed URL for a file.
// Discord doesn't support signed URLs natively, so we just return the regular URL.
func (p *discordProvider) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	// Discord doesn't support signed URLs, so we just return the regular URL
	return p.GetURL(ctx, key)
}

// Delete removes a file from Discord.
func (p *discordProvider) Delete(ctx context.Context, key string) error {
	message, err := p.findMessageWithFile(ctx, key)
	if err != nil {
		// If the file is not found, consider it already deleted
		if err.Error() == "discord provider: file not found" {
			return nil
		}
		return err
	}

	// Delete the message
	url := fmt.Sprintf("%s/channels/%s/messages/%s", discordAPIBaseURL, p.config.ChannelID, message.ID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("discord provider: failed to create delete request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+p.config.BotToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("discord provider: failed to delete message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord provider: failed to delete message, status: %d", resp.StatusCode)
	}

	return nil
}

// GetObject retrieves file information.
func (p *discordProvider) GetObject(ctx context.Context, key string) (*port.FileObject, error) {
	message, err := p.findMessageWithFile(ctx, key)
	if err != nil {
		return nil, err
	}

	attachment := message.Attachments[0]

	// Parse the timestamp
	lastModified := time.Now()
	if message.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, message.Timestamp); err == nil {
			lastModified = t
		}
	}

	return &port.FileObject{
		Key:          key,
		URL:          attachment.URL,
		Size:         int64(attachment.Size),
		ContentType:  attachment.ContentType,
		LastModified: lastModified,
		Provider:     p.ProviderType(),
	}, nil
}

// Download downloads a file from Discord.
func (p *discordProvider) Download(ctx context.Context, key string) (io.ReadCloser, *port.FileObject, error) {
	// Get the file object first
	fileObj, err := p.GetObject(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	// Download the file from the URL
	req, err := http.NewRequestWithContext(ctx, "GET", fileObj.URL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("discord provider: failed to create download request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("discord provider: failed to download file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("discord provider: failed to download file, status: %d", resp.StatusCode)
	}

	return resp.Body, fileObj, nil
}

// ProviderType returns the type of the storage provider.
func (p *discordProvider) ProviderType() port.StorageProviderType {
	return port.ProviderDiscord
}
