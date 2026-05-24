FROM node:24.6.0-slim AS base

ENV NODE_ENV=production \
    TZ=America/Sao_Paulo \
    PM2_DISABLE_INTERACTION=true

WORKDIR /app

COPY package*.json ./

RUN apt-get update && \
    apt-get install -y --no-install-recommends dumb-init && \
    npm install --omit=dev && \
    npm install -g pm2 && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

COPY . .

USER node

EXPOSE 3000

ENTRYPOINT ["dumb-init", "--"]
# Substitua o CMD por:
# CMD ["npm", "run", "start"]
# CMD ["pm2-runtime", "index.js"]
CMD ["pm2-runtime", "--no-daemon", "--instances", "2", "index.js"]
