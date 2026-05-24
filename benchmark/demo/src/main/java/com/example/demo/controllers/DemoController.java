package com.example.demo.controllers;

import lombok.RequiredArgsConstructor;

import java.time.Duration;
import java.util.Map;

import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.reactive.function.server.ServerResponse;

import com.example.demo.model.User;
import com.example.demo.service.DemoService;

import reactor.core.publisher.Mono;


@RestController
@RequestMapping("/")
@RequiredArgsConstructor
public class DemoController {

    private final DemoService handler;
    
    @GetMapping
    public Mono<String> hello() {
        return Mono.just("Hello World!");
    }

    @GetMapping("/user/{id}")
    public Mono<User> findByUser(@PathVariable String id) {
        return this.handler.findUserByIdTest(id);
    }

    @GetMapping("/users")
    public Mono<ServerResponse> findAll() {
        return this.handler.findAllUsers();
    }

    @PostMapping(consumes = MediaType.APPLICATION_JSON_VALUE, produces = MediaType.APPLICATION_JSON_VALUE)
    public Mono<ResponseEntity<Map<String, Object>>> process(@RequestBody Mono<Map<String, Object>> bodyMono) {
                return bodyMono
                .flatMap(this::processBody)
                .map(ResponseEntity::ok);
    }

    private Mono<Map<String, Object>> processBody(Map<String, Object> body) {
        return Mono.just(body)
                .delayElement(Duration.ofMillis(100));
    }    
}