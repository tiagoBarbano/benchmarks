// package com.example.demo.config;

// import org.springframework.context.annotation.Bean;
// import org.springframework.context.annotation.Configuration;

// import io.opentelemetry.sdk.trace.samplers.Sampler;


// @Configuration
// public class CustomSamplerConfiguration {

//     @Bean
//     Sampler otelSampler() {
//         // 10% sampling, como TraceIdRatioBased(0.1)
//         return Sampler.traceIdRatioBased(0.1);
//     }
// }