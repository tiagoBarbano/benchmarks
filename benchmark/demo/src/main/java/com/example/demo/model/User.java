package com.example.demo.model;

import lombok.*;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;


@Document(collection = "users")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class User {

    @Field("_id")
    private String id;

    private String empresa;

    private Double cotacao_final;

    private String created_at;

    private String updated_at;

    private Boolean deleted;
}