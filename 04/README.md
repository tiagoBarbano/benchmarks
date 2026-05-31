# API Catalog

Dashboard web para centralizar e monitorar APIs internas. Exibe status de saúde e versão de cada API em tempo real, agrupadas por domínio de negócio.

## Funcionalidades

- **Dashboard visual** — lista todas as APIs registradas, agrupadas por domínio
- **Health check em tempo real** — verifica se cada API está `UP` ou `DOWN` via endpoint de healthcheck
- **Versão automática** — lê a versão diretamente do `openapi.json` de cada API
- **Swagger viewer** — abre a documentação interativa de qualquer API registrada
- **Proxy OpenAPI** — rota `/proxy/openapi` contorna restrições de CORS ao buscar specs externas

## Stack

| Camada | Tecnologia |
|---|---|
| Framework web | [FastAPI](https://fastapi.tiangolo.com/) |
| Servidor ASGI | [Uvicorn](https://www.uvicorn.org/) |
| Templates | [Jinja2](https://jinja.palletsprojects.com/) |
| HTTP assíncrono | [httpx](https://www.python-httpx.org/) |
| Serialização | [msgspec](https://jcristharif.com/msgspec/) |
| Configuração | [PyYAML](https://pyyaml.org/) |
| Runtime / deps | [uv](https://docs.astral.sh/uv/) |
| Python | >= 3.13 |

## Estrutura do projeto

```
app/
├── main.py                    # Inicialização do FastAPI, middlewares e rotas
├── cache/
│   └── memory_cache.py        # Cache em memória com TTL
├── model/
│   └── api.py                 # Modelos de dados (ApiConfig, ApiView)
├── repositories/
│   └── catalog_repository.py  # Leitura do arquivo YAML de configuração
├── routers/
│   ├── home.py                # Rota do dashboard (/)
│   ├── api_catalog.py         # Rota REST (/api/catalog)
│   └── swagger.py             # Swagger viewer e proxy OpenAPI
└── services/
    ├── catalog_service.py     # Orquestra enriquecimento das APIs (status + versão)
    ├── health_service.py      # Verifica saúde de cada API
    └── openapi_service.py     # Extrai versão do spec OpenAPI

config/
└── apis.yaml                  # Registro das APIs monitoradas

static/
├── css/app.css
├── js/catalog.js
└── templates/
    ├── index.html
    └── components/
        ├── card.html
        ├── dashboard.html
        └── swagger.html
```

## Configuração das APIs

Edite `config/apis.yaml` para registrar as APIs que devem aparecer no catálogo:

```yaml
apis:
  - id: minha-api
    nome: Minha API
    dominio: Comercial
    descricao: Descrição da API

    tags:
      - fastapi
      - python

    openapi: http://localhost:8000/openapi.json
    healthcheck: http://localhost:8000/health
```

| Campo | Descrição |
|---|---|
| `id` | Identificador único |
| `nome` | Nome de exibição |
| `dominio` | Domínio de negócio (usado para agrupamento) |
| `descricao` | Descrição curta |
| `tags` | Lista de tags livres |
| `openapi` | URL do `openapi.json` |
| `healthcheck` | URL do endpoint de health (deve retornar HTTP 200 quando saudável) |

## Como executar

**Pré-requisito:** [uv](https://docs.astral.sh/uv/getting-started/installation/) instalado.

```bash
uv run uvicorn app.main:app --host 0.0.0.0 --port 8080 --reload
```

Acesse em: [http://localhost:8080](http://localhost:8080)

## Rotas disponíveis

| Rota | Método | Descrição |
|---|---|---|
| `/` | GET | Dashboard principal |
| `/api/catalog` | GET | Lista todas as APIs em JSON |
| `/swagger?url=<openapi_url>` | GET | Visualizador Swagger |
| `/proxy/openapi?url=<openapi_url>` | GET | Proxy para specs externas |
