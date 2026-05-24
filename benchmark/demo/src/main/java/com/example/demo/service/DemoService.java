package com.example.demo.service;

import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.data.redis.core.ReactiveValueOperations;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.server.ServerRequest;
import org.springframework.web.reactive.function.server.ServerResponse;
import org.springframework.beans.factory.annotation.Autowired;

import com.example.demo.model.User;
import com.example.demo.repository.UserRepository;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.time.Duration;
import lombok.extern.java.Log;
import reactor.core.publisher.Mono;

@Log
@Component
public class DemoService {
    @Autowired
    private UserRepository userRepository;

    @Autowired
    private ObjectMapper objectMapper;
 
    @Autowired
    private ReactiveRedisTemplate<String, String> redisTemplate;    


    public Mono<ServerResponse> hello(ServerRequest request) {
        return ServerResponse.ok().bodyValue("Hello, WebFlux!");
    }

    // @Cacheavel(chave = "user:id:{id}", ttl = 300, type = User.class)
    @SuppressWarnings("null")
    public Mono<User> findUserByIdTest(String id) {
        return userRepository.findById(id);
    }


    @SuppressWarnings("null")
    public Mono<User> findUserById(String id) {
        String cacheKey = "user:" + id;
        ReactiveValueOperations<String, String> ops = redisTemplate.opsForValue();

        return ops.get(cacheKey)
            .flatMap(json -> {
                try {
                    User user = objectMapper.readValue(json, User.class);
                    return Mono.just(user);
                } catch (Exception e) {
                    return Mono.error(new RuntimeException("Erro ao desserializar do cache", e));
                }
            })
            .switchIfEmpty(
                userRepository.findById(id)
                    .flatMap(user -> {
                        try {
                            String json = objectMapper.writeValueAsString(user);
                            return ops.set(cacheKey, json, Duration.ofMinutes(5))
                                    .thenReturn(user);
                        } catch (Exception e) {
                            return Mono.error(new RuntimeException("Erro ao serializar para o cache", e));
                        }
                    })
            );
    }

    public Mono<ServerResponse> findAllUsers() {
        return ServerResponse.ok().body(userRepository.findAll(), User.class);
    }
}
