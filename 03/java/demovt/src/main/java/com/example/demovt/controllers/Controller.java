package com.example.demovt.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestBody;



@RestController
public class Controller {

    @GetMapping("/ping")
    public String process() {
        return "OK";
    }

    @PostMapping("/small")
    public BenchmarkRequest small(@RequestBody BenchmarkRequest req) {
        return req;
    }

    @PostMapping("/payload")
    public BenchmarkRequest processPayload(@RequestBody BenchmarkRequest request) {
        try {
            Thread.sleep(100);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.out.println("Thread interrompida durante sleep");
        }
        return request;
    }

    public record BenchmarkRequest(
            Long id,
            String name,
            String payload) {
    }

}