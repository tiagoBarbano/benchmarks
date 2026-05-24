package com.example.demo.cache;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import java.lang.annotation.ElementType;

@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
public @interface Cacheavel {
    String chave(); // ex: "usuario:{0}"
    long ttl() default 60; // em segundos
    Class<?> type();       // tipo da entidade para serialização    
}