import time


class MemoryCache:

    def __init__(self):

        self.data = {}

    def get(self, key):

        item = self.data.get(key)

        if not item:
            return None

        if item["expires"] < time.time():

            del self.data[key]
            return None

        return item["value"]

    def put(
        self,
        key,
        value,
        ttl=60
    ):

        self.data[key] = {

            "value": value,

            "expires":
                time.time() + ttl
        }