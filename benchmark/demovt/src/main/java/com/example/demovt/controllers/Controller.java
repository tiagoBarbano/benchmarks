package com.example.demovt.controllers;

import java.io.IOException;
import java.io.InputStream;
import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

import com.fasterxml.jackson.databind.ObjectMapper;

// import io.swagger.v3.oas.annotations.Operation;
import jakarta.servlet.http.HttpServletRequest;

@RestController
public class Controller {

    private final ObjectMapper objectMapper = new ObjectMapper();

    @SuppressWarnings("unchecked")
    @PostMapping("/process")
    public ResponseEntity<Map<String, Object>> process(HttpServletRequest request) throws IOException {
        InputStream inputStream = request.getInputStream();
        Map<String, Object> body = objectMapper.readValue(inputStream, Map.class);
        Map<String, Object> res = processBody(body);
        return ResponseEntity.ok(res);
    }

    private Map<String, Object> processBody(Map<String, Object> body) {
        try {
            Thread.sleep(1);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.out.println("Thread interrompida durante sleep");
        }
        return body;
    }
}
