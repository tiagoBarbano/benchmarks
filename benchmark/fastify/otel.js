import { resourceFromAttributes } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { HttpInstrumentation } from '@opentelemetry/instrumentation-http';
import { MongoDBInstrumentation } from '@opentelemetry/instrumentation-mongodb';
import { IORedisInstrumentation } from '@opentelemetry/instrumentation-ioredis';
import { NodeTracerProvider, BatchSpanProcessor } from '@opentelemetry/sdk-trace-node';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { FastifyOtelInstrumentation } from '@fastify/otel';
import { TraceIdRatioBasedSampler } from '@opentelemetry/sdk-trace-base';

const traceExporter = new OTLPTraceExporter({
   // url: 'http://tempo:4317', // Replace with your OTLP endpoint
   url: 'http://localhost:4317', // Replace with your OTLP endpoint
})

const provider = new NodeTracerProvider({
   sampler: new TraceIdRatioBasedSampler(0.1),
   resource: resourceFromAttributes({
      [SemanticResourceAttributes.SERVICE_NAME]: 'fastify-service',
      [SemanticResourceAttributes.SERVICE_VERSION]: '1.0.0',
   }),
   spanProcessors: [
      new BatchSpanProcessor(traceExporter)
   ],
});

const fastifyOtelInstrumentation = new FastifyOtelInstrumentation({
   registerOnInitialization: true,
});

provider.register();

registerInstrumentations({
   instrumentations: [
      new HttpInstrumentation(),
      new IORedisInstrumentation(),
      new MongoDBInstrumentation(),
   ],
   tracerProvider: provider,
});