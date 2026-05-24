function cacheable({ keyBuilder, ttlSeconds }) {
    return function (handlerFn) {
        return async function (request, reply) {
            const key = keyBuilder(request);
            const result = await request.server.cache.getOrSet(key, ttlSeconds, () => handlerFn(request, reply));
            return result;
        };
    };
}

export default cacheable;
