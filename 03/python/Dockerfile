# Imagem base Python
FROM python:3.13.13-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 

WORKDIR /app

# Instala bibliotecas diretamente (sem requirements.txt)
RUN pip install --upgrade pip && pip install \
    fastapi \
    uvloop \
    gunicorn \
    msgspec

# Copia o app
COPY main.py .

EXPOSE 8000

ENTRYPOINT [ "sh", "-c", "\
    exec gunicorn main:app \
        --worker-class asgi \
        --asgi-loop uvloop \
        --bind 0.0.0.0:8000 \
        --workers 1 \
        --worker-connections 75000 \
        --keep-alive 2 \
        --timeout 30 \
        --graceful-timeout 5 \
        --log-level warning \
        --preload \
        --reuse-port" ]