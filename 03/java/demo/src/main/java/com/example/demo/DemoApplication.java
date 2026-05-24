package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class DemoApplication {

	public static void main(String[] args) {
		System.setProperty("reactor.netty.ioWorkerCount", "1");
        System.setProperty("reactor.netty.ioSelectCount", "1");

        System.setProperty(
            "java.util.concurrent.ForkJoinPool.common.parallelism",
            "1"
        );
		SpringApplication.run(DemoApplication.class, args);
	}

}
