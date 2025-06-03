'use client'

import { WebTracerProvider, SimpleSpanProcessor, InMemorySpanExporter } from '@opentelemetry/sdk-trace-web';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { trace, SpanStatusCode } from '@opentelemetry/api';
import { resourceFromAttributes } from '@opentelemetry/resources';


// Create a Span Processor and Exporter
const exporter = new InMemorySpanExporter();
const spanProcessor = new SimpleSpanProcessor(exporter);

// Create a Resource
const resource = resourceFromAttributes({
	'service.name': 'auth3-web',
});

// Create a Web Tracer Provider with span processors passed during instantiation
const provider = new WebTracerProvider({
	resource: resource,
	spanProcessors: [spanProcessor], // Pass processors during instantiation
});

// 4. Register the Provider and Context Manager globally
// This makes the provider and context manager active for all traces
provider.register({
	contextManager: new ZoneContextManager(),
});

// 5. Register instrumentations
// Pass the provider to ensure instrumentations use this specific provider
registerInstrumentations({
	instrumentations: [
		new DocumentLoadInstrumentation(), // DocumentLoadInstrumentation is web-specific
	],
	tracerProvider: provider,
});

// Export the tracer for use in the application
// trace.getTracer will now use the globally registered provider
export const tracer = trace.getTracer('authentication-system-web');

// Utility function to create a span
export const createSpan = <T>(name: string, fn: () => Promise<T>): Promise<T> => {
	return tracer.startActiveSpan(name, async (span) => {
		try {
			const result = await fn();
			span.setStatus({ code: SpanStatusCode.OK });
			return result;
		} catch (error) {
			span.recordException(error as Error);
			span.setStatus({ code: SpanStatusCode.ERROR, message: (error as Error).message });
			throw error;
		} finally {
			span.end();
		}
	});
};

// Utility function to create a synchronous span
export const createSyncSpan = <T>(name: string, fn: () => T): T => {
	return tracer.startActiveSpan(name, (span) => {
		try {
			const result = fn();
			span.setStatus({ code: SpanStatusCode.OK });
			return result;
		} catch (error) {
			span.recordException(error as Error);
			span.setStatus({ code: SpanStatusCode.ERROR, message: (error as Error).message });
			throw error;
		} finally {
			span.end();
		}
	});
};

// Hook to wrap API calls with tracing
export const withTracing = async <T>(
	name: string,
	operation: () => Promise<T>,
	attributes: Record<string, string | number | boolean> = {}
): Promise<T> => {
	return tracer.startActiveSpan(name, async (span) => {
		try {
			// Add custom attributes
			Object.entries(attributes).forEach(([key, value]) => {
				span.setAttribute(key, value);
			});

			const result = await operation();
			span.setStatus({ code: SpanStatusCode.OK });
			return result;
		} catch (error) {
			span.recordException(error as Error);
			span.setStatus({ code: SpanStatusCode.ERROR, message: (error as Error).message });
			throw error;
		} finally {
			span.end();
		}
	});
};
