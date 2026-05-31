import msgspec


class ApiConfig(msgspec.Struct):
    id: str
    nome: str
    dominio: str
    descricao: str

    tags: list[str]

    openapi: str
    healthcheck: str


class ApiView(msgspec.Struct):
    id: str

    nome: str
    dominio: str
    descricao: str

    tags: list[str]

    status: str
    version: str | None

    openapi: str