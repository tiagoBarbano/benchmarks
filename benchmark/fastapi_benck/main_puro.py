async def app(scope, receive, send):
    assert scope["type"] == "http"

    if scope["method"] == "GET" and scope["path"] == "/":
        body = b'{"message": "Hello, world!"}'

        await send({
            "type": "http.response.start",
            "status": 200,
            "headers": [
                (b"content-type", b"application/json"),
                (b"content-length", str(len(body)).encode())
            ]
        })

        await send({
            "type": "http.response.body",
            "body": body,
        })
    else:
        await send({
            "type": "http.response.start",
            "status": 404,
            "headers": [(b"content-type", b"text/plain")]
        })
        await send({
            "type": "http.response.body",
            "body": b"Not Found"
        })
