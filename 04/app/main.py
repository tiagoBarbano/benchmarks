from fastapi import FastAPI


from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import (
    StaticFiles
)

from app.routers.home import router as home
from app.routers.api_catalog import router as catalog
from app.routers.swagger import router as swagger_router

app = FastAPI(
    title="API Catalog"
)


app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.mount(
    "/static",
    StaticFiles(
        directory="static"
    ),
    name="static"
)


app.include_router(home)
app.include_router(catalog)
app.include_router(swagger_router)