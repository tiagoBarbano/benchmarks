package com.example.demo.cache;

import com.fasterxml.jackson.databind.ObjectMapper;

import lombok.RequiredArgsConstructor;

import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.stereotype.Component;

import reactor.core.publisher.Mono;

import java.time.Duration;

@Aspect
@Component
@RequiredArgsConstructor
public class CacheavelReactiveAspect {

    private final ReactiveRedisTemplate<String, String> redisTemplate;
    private final ObjectMapper objectMapper;

    @SuppressWarnings("null")
    @Around("@annotation(cacheavel)")
    public Mono<?> interceptar(ProceedingJoinPoint joinPoint, Cacheavel cacheavel) {
        String chave = gerarChave(joinPoint, cacheavel.chave());

        return redisTemplate.opsForValue().get(chave)
                .flatMap(json -> {
                    try {
                        Object valor = objectMapper.readValue(json, cacheavel.type());
                        return Mono.just(valor);
                    } catch (Exception e) {
                        return Mono.error(e);
                    }
                })
                .switchIfEmpty(Mono.defer(() -> {
                    try {
                        Mono<?> resultadoMono = (Mono<?>) joinPoint.proceed();
                        return resultadoMono.flatMap(obj -> {
                            try {
                                String json = objectMapper.writeValueAsString(obj);
                                return redisTemplate.opsForValue()
                                        .set(chave, json, Duration.ofSeconds(cacheavel.ttl()))
                                        .thenReturn(obj);
                            } catch (Exception e) {
                                return Mono.error(e);
                            }
                        });
                    } catch (Throwable e) {
                        return Mono.error(e);
                    }
                }));
    }

    private String gerarChave(ProceedingJoinPoint joinPoint, String template) {
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        String[] parameterNames = signature.getParameterNames();
        Object[] args = joinPoint.getArgs();

        if (parameterNames == null || args == null) {
            return template; // fallback
        }

        for (int i = 0; i < parameterNames.length; i++) {
            String name = parameterNames[i];
            Object value = args[i];
            template = template.replace("{" + name + "}", String.valueOf(value));
        }

        return template;
    }
}
