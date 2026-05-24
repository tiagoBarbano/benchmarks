import 'dotenv/config';
import pino from 'pino';
import Fastify from 'fastify';
import fastifyRedis from '@fastify/redis';
import fastifyMongo from '@fastify/mongodb';
import fastifySwagger from '@fastify/swagger';
import fastifySwaggerUi from '@fastify/swagger-ui';

const log = pino({ level: 'error' });

const fastify = Fastify({ logger: false });

import fastifyHealthcheck from 'fastify-healthcheck';
fastify.register(fastifyHealthcheck);
fastify.register(fastifySwagger);
fastify.register(fastifySwaggerUi, {
  routePrefix: '/docs',
  uiConfig: {
    docExpansion: 'full',
    deepLinking: false
  },
  uiHooks: {
    onRequest: function (request, reply, next) { next() },
    preHandler: function (request, reply, next) { next() }
  },
  staticCSP: true,
  transformStaticCSP: (header) => header,
  transformSpecification: (swaggerObject, request, reply) => { return swaggerObject },
  transformSpecificationClone: true
})

fastify.register(fastifyRedis, {
  host: process.env.REDIS_HOST,
  port: Number(process.env.REDIS_PORT),
  password: process.env.REDIS_PASSWORD,
  db: 0,
  maxRetriesPerRequest: null, // Evita atraso em caso de erro
  enableOfflineQueue: false,  // Evita uso de memória caso Redis caia
  lazyConnect: false,         // Garante conexão pronta no boot
  retryStrategy: (times) => Math.min(times * 50, 2000),  
});

fastify.register(fastifyMongo, {
  forceClose: true,
  url: process.env.MONGO_URL,
});


fastify.get('/users/:id', async function (req, reply) {
    try {
      const users = fastify.mongo.db.collection('users');
      const user = await users.findOne({ _id: req.params.id });
      if (!user) {
        log.info(`User not found: ${req.params.id}`);
        return reply.code(204).send();
      }
      return user;
    } catch (err) {
      return reply.code(500).send({ error: 'Erro ao buscar usuário' });
    }
  });

// Lista usuários
fastify.get('/users', {
  schema: {
    description: 'Consulta todos os usuários',
    tags: ['user'],
    summary: 'Consulta todos os usuários',
    response: {
      200: {
        description: 'Successful response',
        type: 'array',
        items: {
          type: 'object',
          properties: {
            _id: { type: 'string' },
            empresa: { type: 'string' },
            cotacao_final: { type: 'number' },
            created_at: { type: 'string' },
            updated_at: { type: 'string' },
            deleted: { type: 'boolean' },
          }
        }
      },
      400: {
        description: 'Bad request',
        type: 'object',
        properties: {
          error: { type: 'string' }
        }
      },
      500: {
        description: 'Internal server error',
        type: 'object',
        properties: {
          error: { type: 'string' }
        }
      }
    }
  }
}, async function (req, reply) {
  try {
    const users = this.mongo.db.collection('users');
    return await users.find({}).toArray();
  } catch (err) {
    console.error(err);
    return reply.code(500).send({ error: 'Erro ao buscar usuários' });
  }
});


fastify.get('/', {
  schema: {
    description: 'Hello World',
    tags: ['HelloWorld'],
    summary: 'Hello World'
  }
}, async (request, reply) => {
  return { message: 'Hello World' };

});

// Inicia servidor
const start = async () => {
  try {
    await fastify.listen({ port: Number(process.env.PORT), host: '0.0.0.0' });
    console.log('Server is running at http://localhost:3000');
    console.log('Prometheus metrics at http://localhost:3000/metrics');
    console.log('Swagger-ui at http://localhost:3000/docs');
    console.log('Health Check at http://localhost:3000/health');
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();
