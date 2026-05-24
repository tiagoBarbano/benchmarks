package com.example.demo.repository;

import com.example.demo.model.User;

import reactor.core.publisher.Mono;

import org.springframework.data.mongodb.repository.Query;
import org.springframework.data.repository.reactive.ReactiveCrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserRepository extends ReactiveCrudRepository<User, String> {

    @Query("{ '_id': ?0 }")
    Mono<User> buscarPorIdPersonalizado(String id);
}