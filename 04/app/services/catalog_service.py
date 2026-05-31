import asyncio

from app.repositories.catalog_repository import (
    CatalogRepository
)

from app.services.health_service import (
    HealthService
)

from app.services.openapi_service import (
    OpenApiService
)


class CatalogService:

    def __init__(self):

        self.repo = CatalogRepository()

        self.health = HealthService()

        self.openapi = OpenApiService()

    async def list(self):

        config = self.repo.load()

        apis = config["apis"]

        tasks = [
            self._enrich(api)
            for api in apis
        ]

        return await asyncio.gather(*tasks)

    async def _enrich(
        self,
        api
    ):

        status = await self.health.status(
            api["healthcheck"]
        )

        version = await self.openapi.version(
            api["openapi"]
        )

        return {

            **api,

            "status": status,

            "version": version
        }