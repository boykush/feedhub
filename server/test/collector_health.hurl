GET http://localhost:{{bff_port}}/api/v1/collector/health
HTTP 200
[Asserts]
jsonpath "$.status" == "SERVING"
