package com.example.demo.controllers;

import java.time.Duration;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import reactor.core.publisher.Mono;

@RestController
public class DemoController {

    @PostMapping("/payload")
    public Mono<BenchmarkRequest> process_payload(@RequestBody Mono<BenchmarkRequest> bodyMono) {
        return bodyMono.delayElement(Duration.ofMillis(100));
    }

    @PostMapping("/small")
    public Mono<BenchmarkRequest> process_small(@RequestBody Mono<BenchmarkRequest> bodyMono) {
        return bodyMono;
    }

    @GetMapping("/ping")
    public Mono<String> ping() {
        return Mono.just("OK");
    }

    public record BenchmarkRequest(
            Long id,
            String name,
            String payload) {
    }

}