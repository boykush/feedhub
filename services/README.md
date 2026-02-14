# Services

このディレクトリには、Foreseeアプリケーションを構成するマイクロサービス群が含まれています。

## アーキテクチャ概要

```
┌─────────────────┐
│     Browser      │
└────────┬────────┘
         │ HTTP
         ▼
    ┌────────┐
    │  Web   │ ← Next.js フロントエンド
    └───┬────┘
        │ HTTP/REST
        ▼
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

### [web/](web/)
**フロントエンド** - Next.js App Routerによるユーザーインターフェース

- 役割: フィード一覧・記事一覧の表示、フィード登録、同期操作
- 技術: Next.js, TypeScript, Tailwind CSS
- ポート: 3000

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
# Go サービス
mise tasks | grep go

# Web フロントエンド
mise tasks | grep web
```

### Protocol Buffers コード生成
```bash
mise run proto:generate feed
mise run proto:generate collector
mise run proto:generate bff
```
