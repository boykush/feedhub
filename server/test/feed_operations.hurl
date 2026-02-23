# AddFeed - success
POST http://localhost:{{bff_port}}/api/v1/feeds
Content-Type: application/json
{"url": "https://github.blog/feed/"}
HTTP 200
[Asserts]
jsonpath "$.feedId" exists
jsonpath "$.title" exists

# AddFeed - duplicate URL
POST http://localhost:{{bff_port}}/api/v1/feeds
Content-Type: application/json
{"url": "https://github.blog/feed/"}
HTTP 409
