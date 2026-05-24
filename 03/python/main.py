import asyncio

import msgspec
from fastapi import FastAPI, Request
from fastapi.responses import Response


app = FastAPI(
    title="",
    version="",
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

    return Response(msgspec.json.encode(res), media_type="application/json")


async def controller_business_manager_rules_handler(body):
    await asyncio.sleep(0.1)
    return body
