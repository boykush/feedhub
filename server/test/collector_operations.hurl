POST http://localhost:{{bff_port}}/api/v1/collector/feeds
Content-Type: application/json
{"url": "https://example.com/feed.xml"}
HTTP 200

POST http://localhost:{{bff_port}}/api/v1/collector/sync
Content-Type: application/json
{}
HTTP 200
