package tracer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials"
)

// OtelInitConfig holds the necessary configuration fields for Otel initialization.
type OtelInitConfig struct {
	ServiceName  string
	CollectorURL string
	Insecure     string // Representing boolean as string
	Headers      map[string]string
}

// InitOtel initializes OpenTelemetry tracing and logging, returning a shutdown function.
// Assumes it should initialize if called; enablement check happens before calling.
func InitOtel(cfg OtelInitConfig) func(context.Context) error {
	log.Println("Initializing OpenTelemetry integration...")

	var secureOption otlptracegrpc.Option // Used for both trace and log exporters
	var logSecureOption otlploggrpc.Option

	// Determine security option based on Insecure flag
	isInsecure := strings.ToLower(cfg.Insecure) == "true" || cfg.Insecure == "1" || strings.ToLower(cfg.Insecure) == "t"
	if !isInsecure {
		// Assuming the same TLS config works for both trace and log
		tlsCreds := credentials.NewClientTLSFromCert(nil, "")
		secureOption = otlptracegrpc.WithTLSCredentials(tlsCreds)
		logSecureOption = otlploggrpc.WithTLSCredentials(tlsCreds)
	} else {
		secureOption = otlptracegrpc.WithInsecure()
		logSecureOption = otlploggrpc.WithInsecure()
	}

	// --- Shared Resource ---
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName), // Use serviceName from cfg
			attribute.String("library.language", "go"),
			// Add other relevant resource attributes here
			// e.g., attribute.String("deployment.environment", "production"),
		),
		resource.WithSchemaURL(semconv.SchemaURL),
	)
	if err != nil {
		// Log fatal error if resource creation fails
		log.Fatalf("Could not set OpenTelemetry resources: %v", err)
	}

	// --- Trace Setup ---
	// Create the OTLP trace exporter
	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.CollectorURL), // Use CollectorURL from cfg
			otlptracegrpc.WithHeaders(cfg.Headers),       // Use Headers from cfg
		),
	)
	if err != nil {
		log.Fatalf("Failed to create OTLP trace exporter: %v", err)
	}

	// Create the TracerProvider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Configure sampler as needed
		sdktrace.WithBatcher(traceExporter),           // Use the created exporter
		sdktrace.WithResource(res),                    // Attach shared resources
	)

	// Set the global TracerProvider
	otel.SetTracerProvider(tracerProvider)

	// --- Log Setup ---
	// Create the OTLP log exporter
	logExporter, err := otlploggrpc.New(
		context.Background(),
		logSecureOption,                            // Use the determined security option
		otlploggrpc.WithEndpoint(cfg.CollectorURL), // Use CollectorURL from cfg
		otlploggrpc.WithHeaders(cfg.Headers),       // Use Headers from cfg
	)
	if err != nil {
		log.Fatalf("Failed to create OTLP log exporter: %v", err)
	}

	// Create the LoggerProvider with a BatchProcessor
	// Note: Using aliased sdklog from "go.opentelemetry.io/otel/sdk/log"
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res), // Attach shared resources
	)

	// Set the global LoggerProvider
	global.SetLoggerProvider(loggerProvider)

	log.Println("OpenTelemetry integration initialized successfully.")

	// --- Combined Shutdown Function ---
	return func(ctx context.Context) error {
		var wg sync.WaitGroup
		var traceErr, logErr error
		shutdownTimeout := 5 * time.Second // Example timeout

		wg.Add(2) // One for tracer, one for logger

		// Shutdown Tracer Provider
		go func() {
			defer wg.Done()
			shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
			defer cancel()
			if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
				fmt.Printf("Failed to shutdown TracerProvider: %x\n", err)
				traceErr = err // Capture error
			}
		}()

		// Shutdown Logger Provider
		go func() {
			defer wg.Done()
			shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
			defer cancel()
			if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
				fmt.Printf("Failed to shutdown LoggerProvider: %x\n", err)
				logErr = err // Capture error
			}
		}()

		wg.Wait() // Wait for both shutdowns to complete

		// Return the first error encountered, if any
		if traceErr != nil {
			return fmt.Errorf("tracer shutdown error: %w", traceErr)
		}
		if logErr != nil {
			return fmt.Errorf("logger shutdown error: %w", logErr)
		}
		log.Println("OpenTelemetry integration shut down gracefully.")
		return nil // Both shut down successfully
	}
}
