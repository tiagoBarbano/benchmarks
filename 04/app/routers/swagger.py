from fastapi import APIRouter, HTTPException, Request
from fastapi.responses import JSONResponse
from fastapi.templating import Jinja2Templates
from urllib.parse import urlparse

import httpx

router = APIRouter()

templates = Jinja2Templates(
    directory="static/templates"
)

@router.get("/swagger")
async def swagger(
    request: Request,
    url: str
):
    return templates.TemplateResponse(
        request=request,
        name="components/swagger.html",
        context={"url": url}
    )
    
@router.get("/proxy/openapi")
async def proxy_openapi(url: str):

    try:

        async with httpx.AsyncClient() as client:

            response = await client.get(
                url,
                timeout=10
            )

            response.raise_for_status()

            return JSONResponse(
                content=response.json()
            )

    except Exception as ex:

        raise HTTPException(
            status_code=500,
            detail=str(ex)
        )