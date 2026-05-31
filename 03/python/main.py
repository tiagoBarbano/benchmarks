import asyncio

import msgspec
from fastapi import FastAPI, Request
from fastapi.responses import Response
from pydantic import BaseModel


class BenchmarkRequest(BaseModel):
    id: int
    name: str
    payload: str

app = FastAPI()

@app.get("/health")
async def root():
    return {"message": "Hello World"}


@app.get("/ping")
async def ping():
    # return "OK"
    return Response(msgspec.json.encode("OK"), media_type="application/json")

@app.post("/small", response_model=BenchmarkRequest)
async def small(request: BenchmarkRequest):
    return request

@app.post("/")
async def root_body_msgspec(request: Request):
    raw = bytearray()

    async for chunk in request.stream():
        raw.extend(chunk)

    body = msgspec.json.decode(memoryview(raw))

    res = await controller_business_manager_rules_handler(body=body)

    return Response(msgspec.json.encode(res), media_type="application/json")


@app.post("/payload", response_model=BenchmarkRequest)
async def root_body_pydantic(request: BenchmarkRequest):
    res = await controller_business_manager_rules_handler(body=request.model_dump())

    return BenchmarkRequest(**res)

async def controller_business_manager_rules_handler(body):
    await asyncio.sleep(0.1)
    return body
