# Imagem base Python
FROM python:3.13-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1

WORKDIR /app

# Instala dependências do projeto
RUN pip install --upgrade pip && pip install \
    fastapi \
    httpx \
    jinja2 \
    msgspec \
    pyyaml \
    uvicorn

# Copia os artefatos do projeto
COPY app/ ./app/
COPY static/ ./static/
COPY config/ ./config/

EXPOSE 8000

ENTRYPOINT [ "uvicorn", "app.main:app", \
    "--host", "0.0.0.0", \
    "--port", "8000", \
    "--workers", "1", \
    "--log-level", "warning" ]