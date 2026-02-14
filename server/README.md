# Server

Foreseeアプリケーションのバックエンドサービス群です。Go + gRPC で構築されています。

## アーキテクチャ概要

```
    ┌────────┐
    │  BFF   │ ← エントリポイント (grpc-gateway)
    └───┬────┘
        │ gRPC
        ▼
   ┌──────────────────────┐
   │  Backend Services    │
   ├──────────────────────┤
   │ • feed               │
   │ • collector           │
   └──────────┬───────────┘
              │
              ▼
        ┌──────────┐
        │ PostgreSQL│
        └──────────┘
```

## サービス一覧

### [bff/](bff/)
**Backend for Frontend** - クライアントアプリケーションのAPIゲートウェイ

- 役割: REST APIエンドポイントを提供し、内部のgRPCマイクロサービスへリクエストをルーティング
- 技術: grpc-gateway を使用してgRPCをHTTP/RESTに変換
- ポート: 8080

### [feed/](feed/)
**フィード読み取りサービス** - フィードと記事データの読み取りAPI

- 役割: フィード一覧・記事一覧の取得
- プロトコル: gRPC
- ポート: 50052

### [collector/](collector/)
**RSS収集サービス** - RSSフィードの収集と保存

- 役割: RSSフィードの登録、全フィードの同期（記事取得・保存）
- プロトコル: gRPC
- ポート: 50053

## 開発

### 利用可能なコマンド
```bash
mise tasks | grep go
```

### Protocol Buffers コード生成
```bash
mise run proto:generate feed
mise run proto:generate collector
mise run proto:generate bff
```
