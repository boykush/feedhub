# Services

このディレクトリには、Foreseeアプリケーションを構成するマイクロサービス群が含まれています。

## アーキテクチャ概要

```
┌─────────────────┐
│   Client App    │
└────────┬────────┘
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
   │ • transaction        │
   │ • (future services)  │
   └──────────────────────┘
```

## サービス一覧

### [bff/](bff/)
**Backend for Frontend** - クライアントアプリケーションのエントリポイント

- 役割: REST APIエンドポイントを提供し、内部のgRPCマイクロサービスへリクエストをルーティング
- 技術: grpc-gateway を使用してgRPCをHTTP/RESTに変換
- 設定ファイル: [rest_api.yaml](bff/rest_api.yaml) - HTTPマッピング定義

### [transaction/](transaction/)
取引管理サービス

- 役割: 取引データの管理と処理
- プロトコル: gRPC
- エンドポイント例: `/api/v1/transactions/*`

## 開発

### 利用可能なコマンド
```bash
# 利用可能なコマンド一覧
mise tasks | grep go
```

### Protocol Buffers コード生成
```bash
mise run proto:generate
```

### ローカル実行
```bash
# 各サービスディレクトリで直接実行
cd services/bff
go run cmd/server/main.go

cd services/transaction
go run cmd/server/main.go
```

## 新しいマイクロサービスの追加

1. `services/`配下に新しいディレクトリを作成
2. Protocol Buffersで定義ファイルを作成 (`proto/*.proto`)
3. BFFの[rest_api.yaml](bff/rest_api.yaml)にHTTPマッピングを追加
4. BFFのprotoディレクトリにシンボリックリンクを作成（必要に応じて）

```bash
# 例: user サービスの追加
mkdir -p services/user/proto
# proto定義を作成...
cd services/bff/proto
ln -s ../../user/proto user
```

各サービスディレクトリには個別の`go.mod`があり、独立してビルド・デプロイ可能です。
