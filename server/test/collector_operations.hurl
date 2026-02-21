# AddFeed - success
POST http://localhost:{{bff_port}}/api/v1/collector/feeds
Content-Type: application/json
{"url": "https://github.blog/feed/"}
HTTP 200
[Asserts]
jsonpath "$.feedId" exists
jsonpath "$.title" exists

# AddFeed - duplicate URL
POST http://localhost:{{bff_port}}/api/v1/collector/feeds
Content-Type: application/json
{"url": "https://github.blog/feed/"}
HTTP 409

# SyncFeeds
POST http://localhost:{{bff_port}}/api/v1/collector/sync
Content-Type: application/json
{}
HTTP 200
