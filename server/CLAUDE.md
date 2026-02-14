# Server Claude Code Guide

## APIテスト

serverのAPI（BFF grpc-gatewayが公開するRESTエンドポイント）を修正した場合は、`server/test/` 配下のhurlテストファイルも合わせて更新してください。

### テストファイル構成

```
server/test/
├── feed_health.hurl          # Feed HealthCheck
├── feed_list.hurl             # ListFeeds / ListArticles
├── collector_health.hurl      # Collector HealthCheck
└── collector_operations.hurl  # AddFeed / SyncFeeds
```

### テスト実行

```bash
mise run server:test
```
