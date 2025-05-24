package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppConfig stores application-specific configuration.
type AppConfig struct {
	Env       string `mapstructure:"env"`
	Name      string `mapstructure:"name"`
	Port      string `mapstructure:"port"`
	Secret    string `mapstructure:"secret"`
	ClientURL string `mapstructure:"clientUrl"`
	Origins   string `mapstructure:"origins"`
}

// DBConfig stores database-specific configuration.
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SslMode  string `mapstructure:"sslMode"`
	LogLevel string `mapstructure:"logLevel"`
}

// RedisConfig stores Redis-specific configuration.
type RedisConfig struct {
	Url  string `mapstructure:"url"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Pass string `mapstructure:"password"`
	User string `mapstructure:"username"`
	DB   int    `mapstructure:"db"`
}

// LogConfig stores logging-specific configuration.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// FireStoreConfig holds Firestore specific configuration.
type FireStoreConfig struct {
	ProjectID       string `mapstructure:"projectID"`
	CredentialsFile string `mapstructure:"credentialsFile"`
	BucketName      string `mapstructure:"bucketName"`
}

// S3Config holds S3 specific configuration.
type S3Config struct {
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	Region          string `mapstructure:"region"`
	BucketName      string `mapstructure:"bucketName"`
	Endpoint        string `mapstructure:"endpoint"`
	DisableSSL      bool   `mapstructure:"disableSSL"`
	ForcePathStyle  bool   `mapstructure:"forcePathStyle"`
}

// CloudflareConfig holds Cloudflare R2 specific configuration.
type CloudflareConfig struct {
	AccountID       string `mapstructure:"accountID"`
	AccessKeyID     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
	BucketName      string `mapstructure:"bucketName"`
	PublicDomain    string `mapstructure:"publicDomain"`
}

// ToS3Config converts CloudflareConfig to S3Config for use with S3-compatible APIs
func (c CloudflareConfig) ToS3Config() S3Config {
	return S3Config{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		BucketName:      c.BucketName,
		Endpoint:        fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID),
		Region:          "auto", // R2 uses "auto" as region
		ForcePathStyle:  true,   // R2 requires path-style addressing
	}
}

// DiscordConfig holds Discord specific configuration.
type DiscordConfig struct {
	BotToken   string `mapstructure:"botToken"`
	ChannelID  string `mapstructure:"channelID"`
	WebhookURL string `mapstructure:"webhookURL"`
}

// LocalStorageConfig holds local adapters specific configuration.
type LocalStorageConfig struct {
	Path            string        `mapstructure:"path"`
	BaseURL         string        `mapstructure:"baseURL"`
	SignedURLExpiry time.Duration `mapstructure:"signedUrlExpiry"` // Default expiry for signed URLs
	SignedURLSecret string        `mapstructure:"signedUrlSecret"` // Secret key for signing URLs (basic example)
}

type AzureConfig struct {
	AccountName   string `mapstructure:"accountName"`
	AccountKey    string `mapstructure:"accountKey"`
	ContainerName string `mapstructure:"containerName"`
	ServiceURL    string `mapstructure:"serviceUrl"`
}

// ScalewayConfig holds Scaleway Object Storage specific configuration
type ScalewayConfig struct {
	AccessKeyID     string `mapstructure:"accessKeyID"`     // API Access Key ID
	SecretAccessKey string `mapstructure:"secretAccessKey"` // API Secret Access Key
	Region          string `mapstructure:"region"`          // Region (e.g., fr-par, nl-ams)
	BucketName      string `mapstructure:"bucketName"`      // The bucket name
	Endpoint        string `mapstructure:"endpoint"`        // Optional: Custom endpoint URL
}

// BackBlazeConfig holds Backblaze B2 specific configuration
type BackBlazeConfig struct {
	KeyID          string `mapstructure:"keyID"`          // Application Key ID
	ApplicationKey string `mapstructure:"applicationKey"` // Application Key
	BucketID       string `mapstructure:"bucketID"`       // Bucket ID
	BucketName     string `mapstructure:"bucketName"`     // Bucket Name
	Region         string `mapstructure:"region"`         // Optional: Region (e.g., us-west-002)
	Endpoint       string `mapstructure:"endpoint"`       // Optional: Custom endpoint URL
}

// ToS3Config converts BackBlazeConfig to S3Config for use with S3-compatible API
func (c BackBlazeConfig) ToS3Config() S3Config {
	endpoint := c.Endpoint
	if endpoint == "" && c.Region != "" {
		endpoint = fmt.Sprintf("https://s3.%s.backblazeb2.com", c.Region)
	}

	return S3Config{
		AccessKeyID:     c.KeyID,
		SecretAccessKey: c.ApplicationKey,
		BucketName:      c.BucketName,
		Endpoint:        endpoint,
		ForcePathStyle:  true, // BackBlaze requires path-style addressing
	}
}

// ToS3Config converts ScalewayConfig to S3Config for use with S3-compatible API
func (c ScalewayConfig) ToS3Config() S3Config {
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://%s.s3.%s.scw.cloud", c.BucketName, c.Region)
	}

	return S3Config{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		Region:          c.Region,
		BucketName:      c.BucketName,
		Endpoint:        endpoint,
		ForcePathStyle:  true, // Scaleway requires path-style addressing
	}
}

// Config stores all configuration of the application.
type Config struct {
	App          AppConfig          `mapstructure:"app"`
	DB           DBConfig           `mapstructure:"db"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Log          LogConfig          `mapstructure:"log"`
	Telegram     TelegramConfig     `mapstructure:"telegram"`
	Adapter      AdapterConfig      `mapstructure:"adapter"`
	RateLimiter  RateLimiterConfig  `mapstructure:"rateLimiter"`
	Signoz       SignozConfig       `mapstructure:"signoz"`
	FireStore    FireStoreConfig    `mapstructure:"firestore"`
	S3           S3Config           `mapstructure:"s3"`
	Cloudflare   CloudflareConfig   `mapstructure:"cloudflare"`
	Discord      DiscordConfig      `mapstructure:"discord"`
	LocalStorage LocalStorageConfig `mapstructure:"localStorage"`
	Azure        AzureConfig        `mapstructure:"azure"`
	Scaleway     ScalewayConfig     `mapstructure:"scaleway"`
	BackBlaze    BackBlazeConfig    `mapstructure:"backblaze"`
}

// RateLimiterConfig holds rate limiter specific configuration.
type RateLimiterConfig struct {
	Max               int `mapstructure:"max"`               // Max requests per expiration window
	ExpirationSeconds int `mapstructure:"expirationSeconds"` // Window duration in seconds
}

type NotifyChannel string

const (
	NotifyMock     NotifyChannel = "mock"
	NotifyTelegram NotifyChannel = "telegram"
	NotifySlack    NotifyChannel = "slack"
)

type EmailProvider string

const (
	EmailSendGrid EmailProvider = "sendgrid"
	EmailBrevo    EmailProvider = "brevo"
	EmailMock     EmailProvider = "mock"
)

type SMSProvider string

const (
	SMSProviderTwilio SMSProvider = "twilio"
	SMSProviderBrevo  SMSProvider = "brevo"
	SMSProviderMock   SMSProvider = "mock"
)

// AdapterConfig holds configuration for different notification adapters.
type AdapterConfig struct {
	Notify NotifyChannel `mapstructure:"notify"`
	Email  EmailProvider `mapstructure:"email"`
	SMS    SMSProvider   `mapstructure:"sms"`
}

// TelegramConfig holds Telegram specific configuration.
type TelegramConfig struct {
	BotToken string `mapstructure:"botToken"`
	ChatID   string `mapstructure:"chatId"`
	Debug    bool   `mapstructure:"debug"`
}

// SignozConfig holds Signoz specific configuration.
type SignozConfig struct {
	CollectorURL string            `mapstructure:"collectorUrl"`
	Insecure     string            `mapstructure:"insecure"`
	Headers      map[string]string `mapstructure:"headers"`
}

// LoadConfig reads configuration from a YAML file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config") // Look for config.yaml or config.json etc.
	viper.SetConfigType("yaml")   // Specify YAML format

	// Attempt to read the config file first
	if readErr := viper.ReadInConfig(); readErr != nil {
		if _, ok := readErr.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if it's just not found.
			log.Println("Viper: config.yaml not found. Relying on environment variables and defaults.")
		} else {
			// Config file was found but another error was produced
			err = fmt.Errorf("failed to read config file (config.yaml): %w", readErr)
			return
		}
	} else {
		log.Println("Using configuration file:", viper.ConfigFileUsed())
	}

	// Allow overriding config values with environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Replace dots with underscores for nested keys in env vars (e.g., FIREBASE.PROJECT_ID -> FIREBASE_PROJECT_ID)

	// Unmarshal the config into the struct using the mapstructure tags
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %x", err)
	}

	// Load AppSecret separately from environment variable for better security
	// Env var should be APP_SECRET (uppercase snake case)
	if config.App.Secret == "your_strong_secret_key" {
		log.Println("Warning: APP_SECRET environment variable not set, and 'app.secret' in config.yaml is missing or using default. Set APP_SECRET environment variable.")
		log.Println("Current value:", config.App.Secret)
		// Depending on requirements, you might want to return an error here in production
		// return config, fmt.Errorf("APP_SECRET must be set via environment variable")
	}

	return
}
