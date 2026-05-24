# Etapa 1: Build da aplicação usando Maven e JDK 17
FROM maven:3.9-eclipse-temurin-21 AS build

WORKDIR /app

# Copia os arquivos necessários para o build
COPY pom.xml .
COPY src ./src

# Gera o JAR (sem rodar testes para acelerar)
RUN mvn clean package -DskipTests

# Etapa 2: Runtime enxuto com JRE 17 + OpenTelemetry Agent
FROM eclipse-temurin:21-jre-alpine AS runtime

WORKDIR /app

# Copia o JAR gerado na etapa de build
COPY --from=build /app/target/*.jar app.jar

# Copia o agente do OpenTelemetry (você pode baixar antes ou durante o build)
# ADD https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/latest/download/opentelemetry-javaagent.jar .

# Expõe a porta padrão do WebFlux
EXPOSE 8080

# Comando de execução com o agente OpenTelemetry
ENTRYPOINT [ \
"java", \
# "-javaagent:/app/opentelemetry-javaagent.jar", \
# "-Dotel.service.name=demo-vt-app", \
# "-Dotel.exporter.otlp.endpoint=http://tempo:4317", \
# "-Dotel.exporter.otlp.protocol=grpc", \
# "-Dotel.metrics.exporter=none", \
# "-Dotel.logs.exporter=none", \
# "-Dotel.traces.sampler=traceidratio", \
# "-Dotel.traces.sampler.arg=0.1", \
"-jar", "/app/app.jar" \
]
