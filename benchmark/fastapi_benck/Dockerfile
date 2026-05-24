# Imagem base Python
FROM python:3.13.11-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 

WORKDIR /app

# Instala bibliotecas diretamente (sem requirements.txt)
RUN pip install --upgrade pip && pip install \
    fastapi==0.129.0 \
    uvloop==0.22.1 \
    uvicorn==0.40.0 \
    gunicorn==25.1.0 \
    pymongo==4.16.0 \
    redis==7.1.1 \ 
    msgspec==0.20.0 \
    light_health==1.0.1

# Copia o app
COPY main.py .

EXPOSE 8000

ENTRYPOINT [ "sh", "-c", "\
    exec gunicorn main:app \
        -k uvicorn.workers.UvicornWorker \
        --bind 0.0.0.0:8000 \
        --workers 2 \
        --worker-connections 75000 \
        --keep-alive 2 \
        --timeout 30 \
        --graceful-timeout 5 \
        --log-level warning \
        --preload \
        --reuse-port" ]

# ENTRYPOINT ["granian"]
# CMD ["--interface", "asgi", "--host", "0.0.0.0", "--port", "8000", "main:app", "--loop", "uvloop", "--backlog", "2048", "--backpressure", "5000", "--workers-kill-timeout", "10s", "--workers", "2"]
