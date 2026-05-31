import httpx


class HealthService:

    async def status(
        self,
        url: str
    ) -> str:

        try:

            async with httpx.AsyncClient() as client:

                response = await client.get(
                    url,
                    timeout=2
                )

                return (
                    "UP"
                    if response.status_code == 200
                    else "DOWN"
                )

        except Exception:
            return "DOWN"