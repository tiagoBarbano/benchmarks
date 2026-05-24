import asyncio
import orjson as json

from typing import Annotated
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Path
from fastapi.responses import ORJSONResponse, Response
from aiohttp import ClientSession
from redis.asyncio import Redis

session: ClientSession | None = None
redis: Redis | None = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global session, redis
    try:
        session = ClientSession()
        redis = Redis(host="localhost", port=6379, password="redis1234", db=1, decode_responses=True)
        yield
    finally:
        await session.close()
        await redis.close()


app = FastAPI(lifespan=lifespan, docs_url=None, redoc_url=None, openapi_url=None)


@app.get("/cep/{cep}")
async def get_address_by_cep(
    cep: Annotated[
        str,
        Path(
            description="CEP principal do segurado",
            pattern=r"^\d{8}$",
            example="01001000",
        ),
    ],
):
    cache_key = f"cep:{cep}"

    cached = await redis.get(cache_key)
    if cached:
        return ORJSONResponse(content=json.loads(cached))

    url = f"https://viacep.com.br/ws/{cep}/json/"
    try:
        async with session.get(url) as response:
            if response.status != 200:
                raise HTTPException(status_code=502, detail="Erro ao acessar ViaCEP")

            data = await response.json()

            if "erro" in data:
                raise HTTPException(status_code=404, detail="CEP não encontrado")

            await redis.set(cache_key, json.dumps(data), ex=600)

            return ORJSONResponse(content=json.dumps(data))

    except asyncio.TimeoutError as e:
        raise HTTPException(
            status_code=504, detail="Timeout na requisição ao ViaCEP"
        ) from e
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Erro inesperado: {str(e)}") from e
