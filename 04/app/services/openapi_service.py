import httpx


class OpenApiService:

    async def version(
        self,
        url: str
    ) -> str | None:

        try:

            async with httpx.AsyncClient() as client:

                response = await client.get(
                    url,
                    timeout=5
                )

                response.raise_for_status()

                return (
                    response
                    .json()
                    .get("info", {})
                    .get("version")
                )

        except Exception:

            return None