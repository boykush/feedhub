GET http://localhost:{{bff_port}}/api/v1/feeds/health
HTTP 200
[Asserts]
jsonpath "$.status" == "SERVING"
