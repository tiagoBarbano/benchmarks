package com.example.demovt;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;
import com.fasterxml.jackson.module.blackbird.BlackbirdModule;



@SpringBootApplication
public class DemovtApplication {

    public static void main(String[] args) {
        SpringApplication.run(DemovtApplication.class, args);
    }

    @Bean
    ObjectMapper objectMapper() {
        return JsonMapper.builder()
                .addModule(new BlackbirdModule())
                .build();
    }

}
