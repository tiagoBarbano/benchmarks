from fastapi import APIRouter
from fastapi import Request

from fastapi.templating import (
    Jinja2Templates
)

from app.services.catalog_service import (
    CatalogService
)

router = APIRouter()

templates = Jinja2Templates(
    directory="static/templates"
)

service = CatalogService()


@router.get("/")
async def home(
    request: Request
):

    apis = await service.list()
    grouped = {}

    for api in apis:
        grouped.setdefault(
            api["dominio"],
            []
        ).append(api)

    total = len(apis)

    online = sum(
        1
        for api in apis
        if api["status"] == "UP"
    )

    offline = total - online

    domains = len(grouped)

    return templates.TemplateResponse(
        request=request,
        name="index.html",
        context={
            "apis": apis,
            "grouped": grouped,
            "total": total,
            "online": online,
            "offline": offline,
            "domains": domains,
        }
    )