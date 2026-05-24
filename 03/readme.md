# FastAPI vs WebFlux vs Virtual Threads em Payloads Dinâmicos: Um Benchmark Mais Próximo do Mundo Real

## Introdução

Nos últimos anos, benchmarks entre runtimes web passaram a focar quase exclusivamente em cenários sintéticos:

* hello world
* plaintext
* DTOs perfeitamente tipados
* payloads pequenos
* máquinas com muitos cores
* sem limites reais de CPU

Embora esses testes sejam úteis para medir potencial bruto, eles frequentemente se afastam bastante do comportamento encontrado em aplicações enterprise reais.

Este benchmark buscou avaliar um cenário mais próximo do que normalmente encontramos em produção:

* payload JSON dinâmico
* uso de `Map<String,Object>`
* ausência de DTOs fixos
* workloads IO-bound
* limite de 1 vCPU
* payloads médios e grandes
* concorrência fixa
* execução em ambiente restrito

O objetivo não foi descobrir “a linguagem mais rápida”, mas entender:

* impacto do framework overhead
* comportamento sob payload grande
* eficiência de runtime
* estabilidade de latência
* custo de abstração

---

# Cenário do Benchmark

## Hardware

* 1 vCPU
* limite real de CPU
* ambiente Linux

---

# Tecnologias Avaliadas

## Python

Stack:

* FastAPI
* Granian
* uvloop
* msgspec
* parsing manual do body
* Response direta

Execução:

```bash
gunicorn main:app \
        --worker-class asgi \
        --asgi-loop uvloop \
        --bind 0.0.0.0:8000 \
        --workers 1 \
        --worker-connections 75000 \
        --keep-alive 2 \
        --timeout 30 \
        --graceful-timeout 5 \
        --log-level warning \
        --preload \
        --reuse-port"
```

### Endpoint Python

```python
import asyncio

import msgspec
from fastapi import FastAPI, Request
from fastapi.responses import Response

app = FastAPI(
    docs_url=None,
    redoc_url=None,
    openapi_url=None,
)

@app.post("/")
async def root_body_msgspec(request: Request):

    raw = bytearray()

    async for chunk in request.stream():
        raw.extend(chunk)

    body = msgspec.json.decode(memoryview(raw))

    res = await controller_business_manager_rules_handler(body=body)

    return Response(
        msgspec.json.encode(res),
        media_type="application/json"
    )


async def controller_business_manager_rules_handler(body):
    await asyncio.sleep(0.1)
    return body
```

---

## Spring WebFlux

### Endpoint WebFlux

```java
@RestController
public class DemoController {

    @PostMapping(
        consumes = MediaType.APPLICATION_JSON_VALUE,
        produces = MediaType.APPLICATION_JSON_VALUE
    )
    public Mono<ResponseEntity<Map<String, Object>>> process(
            @RequestBody Mono<Map<String, Object>> bodyMono) {

        return bodyMono
                .flatMap(this::processBody)
                .map(ResponseEntity::ok);
    }

    private Mono<Map<String, Object>> processBody(
            Map<String, Object> body) {

        return Mono.just(body)
                .delayElement(Duration.ofMillis(100));
    }
}
```

---

## Spring MVC com Virtual Threads

### Endpoint Virtual Threads

```java
@RestController
public class Controller {

    private final ObjectMapper objectMapper = new ObjectMapper();

    @SuppressWarnings("unchecked")
    @PostMapping("/")
    public ResponseEntity<Map<String, Object>> process(
            HttpServletRequest request) throws IOException {

        InputStream inputStream = request.getInputStream();

        Map<String, Object> body =
                objectMapper.readValue(inputStream, Map.class);

        Map<String, Object> res = processBody(body);

        return ResponseEntity.ok(res);
    }

    private Map<String, Object> processBody(
            Map<String, Object> body) {

        try {
            Thread.sleep(100);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        return body;
    }
}
```

---

# Ferramenta de Benchmark

## wrk

```bash
wrk -t4 -c200 -d30s --latency -s post.lua http://localhost:8000/
```

### Script Lua

```lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local payload = string.rep("abcdefghij", PAYLOAD)

wrk.body = [[
{
  "id": 1,
  "name": "benchmark",
  "payload": "]] .. payload .. [["
}
]]
```

---

# Cenário 1 — Payload ~50 KB

## Resultados

| Stack           | Req/s | P50    | P99    | CPU    |
| --------------- | ----- | ------ | ------ | ------ |
| Python          | ~1789 | ~101ms | ~368ms | ~40%   |
| WebFlux         | ~1791 | ~100ms | ~136ms | 50–60% |
| Virtual Threads | ~1500 | ~122ms | ~166ms | 50–60% |

---

# Análise — Payload Médio

## Python

O resultado do Python foi extremamente impressionante.

Mesmo utilizando FastAPI, o pipeline foi bastante reduzido:

```text
socket
 → uvloop
 → bytearray
 → msgspec
 → coroutine
 → msgspec encode
```

O uso combinado de:

* msgspec
* memoryview
* uvloop
* parsing manual
* Response direta

reduziu drasticamente o overhead do framework.

O throughput ficou praticamente empatado com WebFlux utilizando menos CPU.

Por outro lado, o P99 permaneceu significativamente pior, indicando:

* stalls ocasionais do event loop
* bursts de alloc
* crescimento de buffers
* pausas ocasionais do GC

---

## WebFlux

O WebFlux apresentou o comportamento mais equilibrado.

Embora tenha utilizado mais CPU, apresentou:

* melhor estabilidade
* melhor P99
* menor jitter
* latência mais previsível

O Netty mostrou forte capacidade de estabilizar latência mesmo sob concorrência.

---

## Virtual Threads

As Virtual Threads apresentaram um comportamento intermediário.

O modelo simplifica bastante o desenvolvimento, porém:

* não resolve alloc pressure
* não reduz parsing JSON
* não reduz GC

Quando o payload cresce, parte da vantagem do modelo diminui.

---

# Cenário 2 — Payload ~150 KB

## Resultados

| Stack           | Req/s | P50    | P99    | CPU     |
| --------------- | ----- | ------ | ------ | ------- |
| Python          | ~1794 | ~101ms | ~405ms | 60–70%  |
| WebFlux         | ~1236 | ~145ms | ~204ms | 90–100% |
| Virtual Threads | ~846  | ~225ms | ~311ms | 90–100% |

---

# Análise — Payload Grande

Aqui o benchmark mudou completamente.

O gargalo deixou de ser apenas:

* scheduler
* coordenação concorrente
* waiting model

E passou a ser dominado por:

* parsing JSON
* allocs
* cópias de memória
* GC
* buffers

---

# O Resultado Mais Interessante

O throughput do Python praticamente não mudou.

| Payload | Req/s Python |
| ------- | ------------ |
| ~50 KB  | ~1789        |
| ~150 KB | ~1794        |

Isso sugere fortemente que:

* msgspec é extremamente eficiente
* o pipeline Python é muito curto
* o custo incremental do payload é baixo
* há pouca coordenação intermediária

---

# O Que Mais Pesou no Java

O benchmark mostrou claramente o custo de payload dinâmico no ecossistema Java.

Especialmente:

```java
Map<String, Object>
```

Isso gera:

* LinkedHashMap
* reflection
* boxing
* String allocations
* type coercion
* object graphs grandes

O resultado foi:

* CPU próxima de 100%
* aumento significativo de latência
* queda de throughput

---

---

# Cenário 3 — Payload ~250 KB

## Resultados

| Stack           | Req/s | P50    | P99    | CPU   |
| --------------- | ----- | ------ | ------ | ----- |
| Python          | ~1805 | ~102ms | ~506ms | ~85%  |
| WebFlux         | ~724  | ~259ms | ~337ms | ~100% |
| Virtual Threads | ~556  | ~322ms | ~504ms | ~100% |

---

# Análise — Payload Muito Grande

Neste cenário o benchmark passou a ser fortemente dominado por:

* parsing JSON
* cópias de memória
* alloc pressure
* crescimento de buffers
* pressão de GC

O comportamento das stacks mudou drasticamente.

---

# O Resultado Mais Surpreendente

O Python passou a entregar quase 3x mais throughput que WebFlux.

| Stack           | Req/s |
| --------------- | ----- |
| Python          | ~1805 |
| WebFlux         | ~724  |
| Virtual Threads | ~556  |

Esse resultado foi bastante significativo.

---

# O Que Isso Mostra

O benchmark sugere fortemente que o gargalo principal deixou de ser:

* scheduler
* coordenação concorrente
* runtime HTTP

E passou a ser:

* serialização
* parsing genérico
* allocs
* object graphs
* pressão de memória

---

# O Impacto de `Map<String,Object>`

O uso de:

```java
Map<String,Object>
```

começou a impactar fortemente o throughput da JVM.

Especialmente por gerar:

* muitos objetos intermediários
* LinkedHashMaps
* boxing
* type coercion
* String allocations
* graphs dinâmicos grandes

Com payloads muito grandes, esse custo tornou-se dominante.

---

# O Papel do msgspec

O resultado também mostrou a eficiência do:

* msgspec
* memoryview
* bytearray
* pipeline reduzido do Python

O custo incremental do payload foi significativamente menor.

Mesmo com payload muito grande, o throughput permaneceu elevado.

---

# CPU

Outro ponto extremamente interessante:

| Stack           | CPU   |
| --------------- | ----- |
| Python          | ~85%  |
| WebFlux         | ~100% |
| Virtual Threads | ~100% |

Mesmo entregando throughput muito maior, o Python ainda utilizou menos CPU.

Isso reforça fortemente a hipótese de que:

* o gargalo principal estava no pipeline/framework Java
* e não necessariamente na JVM em si.

---

# P99

Apesar do throughput extremamente alto, o Python continuou apresentando pior tail latency.

| Stack           | P99    |
| --------------- | ------ |
| Python          | ~506ms |
| WebFlux         | ~337ms |
| Virtual Threads | ~504ms |

Isso sugere:

* stalls ocasionais do event loop
* bursts de alloc
* crescimento de buffers
* pausas ocasionais do GC

Enquanto isso, o WebFlux continuou demonstrando maior estabilidade de latência.

---

# O Que Este Cenário Mostra

Payloads muito grandes mudam completamente o benchmark.

Nesses cenários:

* serialização
* allocs
* cópias
* representação dinâmica de objetos

passam a dominar completamente o comportamento da aplicação.

Isso aproxima bastante o benchmark de workloads enterprise reais envolvendo:

* integração
* eventos
* motores de regras
* orquestração
* documentos JSON extensos
* metadata dinâmica

---

# Por Que Não Foram Utilizados DTOs?

O objetivo foi simular aplicações enterprise reais.

Em muitos projetos:

* payloads são parcialmente dinâmicos
* schemas variam
* integrações possuem metadata flexível
* eventos não possuem contratos rígidos

Especialmente em:

* gateways
* BRMS
* orquestração
* seguros
* eventos
* integrações corporativas

Nesses cenários:

```java
Map<String,Object>
```

continua sendo extremamente comum.

Portanto, o benchmark buscou representar:

* comportamento real
* overhead real
* custo real de abstração

E não apenas o potencial máximo idealizado de cada runtime.

---

# Principais Conclusões

## 1. Framework Overhead Importa Muito

Em workloads IO-bound com payload médio/grande:

* abstrações internas
* allocs
* cópias de memória
* parsing genérico

passam a dominar o benchmark.

---

## 2. Python Moderno Mudou Muito

O ecossistema Python moderno evoluiu bastante:

* uvloop
* msgspec
* Granian
* parsing otimizado
* menos abstração

mudaram significativamente o teto de performance.

---

## 3. WebFlux Continua Extremamente Forte em Estabilidade

O WebFlux apresentou:

* melhor P99
* menor jitter
* excelente previsibilidade

Mesmo consumindo mais CPU.

---

## 4. Virtual Threads Simplificam Muito o Desenvolvimento

As Virtual Threads entregam:

* excelente produtividade
* modelo simples
* código linear
* boa escalabilidade

Mas payload grande continua pressionando:

* allocs
* parsing
* GC

---

## 5. O Benchmark Não Mostra “Linguagem Mais Rápida”

O benchmark mostrou algo mais interessante:

> Pipeline minimalista + serializer extremamente otimizado pode superar frameworks complexos em workloads específicos.

---

# Conclusão Final

Os resultados mostram que a discussão moderna de performance mudou bastante.

Hoje, frequentemente:

* arquitetura
* pipeline
* serialização
* overhead do framework
* comportamento do runtime

pesam mais do que a própria linguagem.

Os benchmarks também reforçam que:

* workloads reais importam
* payloads reais importam
* limites reais de CPU importam
* abstrações enterprise têm custo significativo

Em especial, o benchmark mostrou que um stack Python moderno e extremamente enxuto pode ser surpreendentemente competitivo em workloads IO-bound reais, mesmo quando comparado a stacks reativas maduras da JVM.

Ao mesmo tempo, WebFlux continua demonstrando enorme maturidade em estabilidade e previsibilidade de latência, enquanto Virtual Threads aparecem como uma excelente alternativa para simplificar aplicações enterprise modernas.
