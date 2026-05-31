from fastapi import APIRouter

from app.services.catalog_service import (
    CatalogService
)
router = APIRouter(
    prefix="/api/catalog"
)

service = CatalogService()

@router.get("")
async def list_catalog():

    return await service.list()