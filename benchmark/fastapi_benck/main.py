
import logging
import re
from contextlib import asynccontextmanager
from fastapi import Body, FastAPI, Request
from fastapi.responses import Response
import msgspec
from light_health.asgi.health import HealthASGIApp

from light_health.registry import AsyncHealthRegistry, HealthCheckResult, HealthState

from light_health.checks.mongo import mongo_health_check
from light_health.checks.redis import redis_health_check

from redis.asyncio import Redis
from pymongo import AsyncMongoClient

import asyncio

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

app = FastAPI(
    title="",  # evita validações e headers extras
    version="",
    docs_url=None,
    redoc_url=None,
    openapi_url=None,  # desativa geração automática de OpenAPI
)

# --- Redis (async) ---
redis = Redis(host="localhost", port=6379, password="redis1234", db=0, decode_responses=True)

# --- MongoDB (async) ---
mongo_client = AsyncMongoClient(
    "mongodb://localhost:27017",
    maxPoolSize=500,      # ajuste recomendável
    minPoolSize=50,       # mantém conexões prontas
    connectTimeoutMS=500,
    serverSelectionTimeoutMS=500,
    maxConnecting=150,
)

db = mongo_client["cotador"]
collection = db["users"]

mongo_check = mongo_health_check(mongo_client)
redis_check = redis_health_check(redis)

async def process_alive():
    return HealthCheckResult(status=HealthState.UP)

registry = AsyncHealthRegistry()

registry.register_liveness("process", process_alive)
registry.register_readiness("mongo", mongo_health_check(mongo_client))
registry.register_readiness("redis", redis_health_check(redis))

# @asynccontextmanager
# async def lifespan(app: FastAPI):
#     try:
#         pong = await redis.ping()
#         logging.info(f"Redis conectado: {pong}")
#         await mongo_client.server_info()  # Falha se não conectar
#         logging.info("MongoDB conectado")
#     except Exception as e:
#         logging.error(f"Erro ao conectar: {str(e)}")
#         yield
#     finally:
#         await redis.close()
#         await mongo_client.close()


app.mount("/actuator", HealthASGIApp(registry)) 
# app.mount("/management", ManagementASGIApp())

@app.get("/users/{user_id}")
async def mongo_get_id(user_id: str):  
    result = await collection.find_one({"_id": user_id})
    return Response(msgspec.json.encode(result), media_type="application/json")

@app.post("/users")
async def mongo_insert(user: dict = Body()):
    result = await collection.insert_one(user)

    return {"inserted_id": str(result.inserted_id)}

@app.post("/body-msgspec")
async def root_body_msgspec(request: Request):
    raw = await request.body()
    body = msgspec.json.decode(memoryview(raw))

    res = await controller_business_manager_rules_handler(body=body)

    return Response(msgspec.json.encode(res), media_type="application/json")


async def controller_business_manager_rules_handler(body):
    await asyncio.sleep(0.1)  # Simula operação assíncrona
    return body


class ASGIApp:
    def __init__(self):
        self.route_root = re.compile(r"^/users-asgi$")
        self.route_with_id = re.compile(r"^/users-asgi/(?P<id>[^/]+)$")
        self.json_encoder = msgspec.json.Encoder()
        self.json_decoder = msgspec.json.Decoder()

    async def __call__(self, scope, receive, send):
        if scope["type"] != "http":
            return

        method = scope["method"]
        path = self._internal_path(scope)

        if method == "GET" and self.route_root.match(path):
            await self.transformar_asgi(scope, receive, send)
            return

        if method == "GET":
            if match := self.route_with_id.match(path):
                scope["path_params"] = match.groupdict()
                await self.transformar_asgi_path(scope, receive, send)
                return

        await self._send(send, 404, {"error": "Not Found"})

    # ---------- handlers ----------

    async def transformar_asgi(self, scope, receive, send):
        value = self.get_query_param(scope, b"id")  

        result = await collection.find_one({"_id": value})

        await self._send(send, 200, result) 
        
    async def transformar_asgi_path(self, scope, receive, send):
        value = self.get_path_param(scope, 3)  

        result = await collection.find_one({"_id": value})

        await self._send(send, 200, result) 
    # ---------- infra ----------

    def get_path_param(self, scope, index: int):
        parts = scope["path"].split("/")
        try:
            return parts[index]
        except IndexError:
            return None

    def get_query_param(self, scope, key: bytes):
        qs = scope.get("query_string", b"")
        return next(
            (
                part[len(key) + 1 :].decode()
                for part in qs.split(b"&")
                if part.startswith(key + b"=")
            ),
            None,
        )
    
    async def _read_body(self, receive) -> bytes:
        chunks = []
        while True:
            msg = await receive()
            chunks.append(msg["body"])
            if not msg["more_body"]:
                break

        return msgspec.json.decode(b"".join(chunks))

    def _internal_path(self, scope) -> str:
        path = scope["path"]
        root = scope.get("root_path", "")
        return path[len(root) :] if path.startswith(root) else path

    async def _send(self, send, status: int, body):
        await send(
            {
                "type": "http.response.start",
                "status": status,
                "headers": [
                    (b"content-type", b"application/json"),
                    (b"cache-control", b"no-store"),
                ],
            }
        )
        await send(
            {
                "type": "http.response.body",
                "body": msgspec.json.encode(body)
,
            }
        )


app.mount("/", ASGIApp())