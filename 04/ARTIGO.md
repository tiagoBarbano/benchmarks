# Centralizando a documentação de APIs com FastAPI e Swagger UI

**O problema crônico das APIs espalhadas — e como resolvi com ~200 linhas de Python**

---

## O problema que todo dev já sofreu

Se você trabalha numa empresa com mais de cinco times, já passou por isso:

> *"Qual é a URL do Swagger da API de ofertas mesmo?"*
> *"Deixa eu ver... acho que é o Jenkins... ou era o Confluence? Espera, o link está quebrado."*

APIs existem, estão documentadas (quando estão), mas ninguém sabe onde ficam. O Swagger de cada serviço vive em uma URL diferente, em portas diferentes, em ambientes diferentes. Quando você finalmente acha a URL certa, o serviço está caído e você nem sabe.

Esse é um problema crônico em ambientes com muitos microserviços. A documentação existe, mas o acesso a ela é descentralizado, inconsistente e propenso a falhas silenciosas.

---

## A ideia

A solução mais simples possível: um portal central que:

1. Lê uma lista de APIs de um arquivo de configuração
2. Verifica em tempo real se cada API está no ar
3. Lê a versão atual diretamente do `openapi.json` de cada uma
4. Abre a documentação Swagger de qualquer API com um clique

Sem banco de dados. Sem autenticação complexa. Sem infraestrutura adicional. Um único arquivo YAML com as APIs registradas e uma aplicação FastAPI para dar vida ao catálogo.

---

## A estrutura do projeto

```
app/
├── main.py
├── model/api.py
├── repositories/catalog_repository.py
├── routers/
│   ├── home.py          # dashboard
│   ├── api_catalog.py   # endpoint REST
│   └── swagger.py       # viewer + proxy
└── services/
    ├── catalog_service.py
    ├── health_service.py
    └── openapi_service.py

config/
└── apis.yaml            # o único lugar para registrar APIs
```

Simples, direto. Cada camada tem uma responsabilidade única.

---

## O coração do projeto: o arquivo de configuração

Tudo começa no `config/apis.yaml`. Para registrar uma nova API no catálogo, basta adicionar uma entrada:

```yaml
apis:
  - id: oferta-core
    nome: Oferta Core
    dominio: Comercial
    descricao: Motor de ofertas

    tags:
      - fastapi
      - comercial

    openapi: http://localhost:8000/openapi.json
    healthcheck: http://localhost:8000/health
```

Dois campos são especiais: `openapi` e `healthcheck`. É a partir deles que o sistema vai buscar, em tempo real, o estado de cada API.

O `CatalogRepository` carrega esse arquivo com PyYAML:

```python
class CatalogRepository:
    def __init__(self):
        self._file = Path("config/apis.yaml")

    def load(self):
        with open(self._file) as f:
            return yaml.safe_load(f)
```

---

## Enriquecendo os dados em paralelo

A parte mais interessante é o `CatalogService`. Para cada API registrada, ele precisa fazer duas chamadas HTTP: uma para o healthcheck e outra para o `openapi.json`. Se feitas em sequência, com 10 APIs, isso significaria esperar 10 timeouts em caso de falha.

A solução é usar `asyncio.gather` para disparar todas as chamadas em paralelo:

```python
async def list(self):
    config = self.repo.load()
    apis = config["apis"]

    tasks = [self._enrich(api) for api in apis]
    return await asyncio.gather(*tasks)

async def _enrich(self, api):
    status = await self.health.status(api["healthcheck"])
    version = await self.openapi.version(api["openapi"])

    return {
        **api,
        "status": status,
        "version": version
    }
```

Com `asyncio.gather`, todas as APIs são verificadas simultaneamente. O tempo total é determinado pela API mais lenta, não pela soma de todas.

---

## Health check e versão: simples e resilientes

O `HealthService` faz uma chamada GET com timeout curto (2 segundos) e retorna `"UP"` ou `"DOWN"`. Qualquer exceção — timeout, conexão recusada, DNS — vira `"DOWN"` silenciosamente:

```python
class HealthService:
    async def status(self, url: str) -> str:
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(url, timeout=2)
                return "UP" if response.status_code == 200 else "DOWN"
        except Exception:
            return "DOWN"
```

O `OpenApiService` faz o mesmo para a versão: busca o `openapi.json`, extrai `info.version` e retorna `None` se falhar. Nada explode, o catálogo continua funcionando mesmo com metade das APIs fora do ar.

```python
class OpenApiService:
    async def version(self, url: str) -> str | None:
        try:
            async with httpx.AsyncClient() as client:
                response = await client.get(url, timeout=5)
                response.raise_for_status()
                return response.json().get("info", {}).get("version")
        except Exception:
            return None
```

---

## O dashboard: agrupado por domínio

O router `home.py` monta o contexto para o template Jinja2. As APIs são agrupadas por domínio de negócio — Comercial, Produto, Financeiro — o que facilita muito a navegação quando o catálogo cresce:

```python
@router.get("/")
async def home(request: Request):
    apis = await service.list()
    grouped = {}

    for api in apis:
        grouped.setdefault(api["dominio"], []).append(api)

    total = len(apis)
    online = sum(1 for api in apis if api["status"] == "UP")

    return templates.TemplateResponse(
        request=request,
        name="index.html",
        context={
            "apis": apis,
            "grouped": grouped,
            "total": total,
            "online": online,
            "offline": total - online,
            "domains": len(grouped),
        }
    )
```

O dashboard exibe um resumo no topo (total de APIs, quantas online, quantas offline, quantos domínios) e os cards agrupados por domínio logo abaixo.

---

## O Swagger viewer e o problema de CORS

Clicar num card abre o Swagger UI carregando o `openapi.json` da API selecionada. Mas aqui aparece um problema clássico: o navegador bloqueia requisições cross-origin do Swagger UI para APIs que não têm cabeçalhos CORS configurados.

A solução é um proxy simples no próprio backend:

```python
@router.get("/proxy/openapi")
async def proxy_openapi(url: str):
    try:
        async with httpx.AsyncClient() as client:
            response = await client.get(url, timeout=10)
            response.raise_for_status()
            return JSONResponse(content=response.json())
    except Exception as ex:
        raise HTTPException(status_code=500, detail=str(ex))
```

O Swagger UI aponta para `/proxy/openapi?url=<openapi_url>`. O backend busca o spec, e o navegador faz uma requisição para a mesma origem — sem problemas de CORS.

---

## Como rodar

O projeto usa `uv` para gerenciamento de dependências. Com ele instalado:

```bash
uv run uvicorn app.main:app --host 0.0.0.0 --port 8080 --reload
```

Acesse `http://localhost:8080` e o catálogo estará no ar. Para registrar uma nova API, edite `config/apis.yaml` e recarregue a página.

---

## O que fica de lição

**1. YAML como fonte de verdade é suficiente para começar.** Não precisa de banco de dados para um problema de leitura. Um arquivo versionado no repositório é auditável, fácil de editar e funciona com qualquer processo de CI/CD.

**2. Resiliência silenciosa em integrações externas.** Health checks e chamadas ao OpenAPI falham com frequência. O sistema não pode quebrar por causa disso — tratar qualquer exceção como estado degradado (`"DOWN"`, `None`) é a abordagem correta aqui.

**3. `asyncio.gather` é a ferramenta certa para I/O em paralelo.** A diferença de latência entre chamadas sequenciais e paralelas é percebida imediatamente quando você tem mais de 5 APIs.

**4. Um proxy resolve CORS sem alterar as APIs.** Em vez de exigir que cada time configure CORS corretamente, o catálogo assume essa responsabilidade centralmente.

---

## Próximos passos

O projeto tem espaço para evoluir sem perder a simplicidade:

- **Cache com TTL** — a infraestrutura já existe em `app/cache/memory_cache.py`; aplicar no `CatalogService` reduz a carga nas APIs monitoradas
- **Suporte a ambientes** — adicionar `dev`, `staging`, `prod` por API no YAML
- **Autenticação simples** — um `Depends` no FastAPI com token estático resolve o acesso não autorizado
- **Exportar como OpenAPI unificado** — agregar todos os specs num único documento

---

## Conclusão

Documentação descentralizada é um problema de discoverability. As informações existem, mas o custo de encontrá-las é alto demais para ser pago rotineiramente.

Um catálogo central, mesmo simples, muda esse custo. De "preciso perguntar para alguém" para "abro o portal e vejo". Com ~200 linhas de Python e um arquivo YAML, é possível resolver isso para um time inteiro.

O código está disponível em: **[link do repositório]**

---

*Curtiu? Deixa um clap e compartilha com aquele colega que ainda manda URL de Swagger pelo Slack.* 🙂
