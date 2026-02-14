# Kubernetes Manifests

このディレクトリには、Foreseeアプリケーションのkubernetesマニフェストが含まれています。
Kustomizeを使用した階層構造で管理しており、環境ごとのオーバーレイで設定を上書きできます。

## ディレクトリ構成

```
k8s/
├── base/                      # ベースとなるマニフェスト
│   ├── platform/              # プラットフォーム層（ArgoCD, Istio等）
│   │   ├── argocd/
│   │   ├── istio-base/
│   │   ├── istio-ingressgateway/
│   │   └── istiod/
│   └── workloads/             # アプリケーション層
│       ├── bff/
│       ├── feed/
│       ├── collector/
│       ├── web/
│       ├── postgres/
│       └── istio/             # Istio設定（Gateway, VirtualService等）
└── overlays/                  # 環境別のオーバーレイ
    └── local/                 # ローカル開発環境
        ├── platform/
        └── workloads/
```

## 構成の考え方

### Platform vs Workloads

- **platform/**: インフラストラクチャ層
  - ArgoCD: GitOpsデプロイメントツール
  - Istio: サービスメッシュ（istio-base, istiod, ingressgateway）

- **workloads/**: アプリケーション層
  - マイクロサービス（bff, feed, collector, web）
  - データベース（postgres - CloudNativePG）
  - Istio設定（Gateway, VirtualService, AuthorizationPolicy等）

この分離により、プラットフォーム層とアプリケーション層を独立して管理できます。

### Base vs Overlays

- **base/**: 環境に依存しない共通のマニフェスト
- **overlays/**: 環境ごとの設定（現在はlocalのみ）

## ローカル開発

### クラスタの作成
```bash
mise run k8s:local:cluster:create
```

### プラットフォームリソースのデプロイ
```bash
mise run k8s:local:deploy-platform
```

### アプリケーションのデプロイ
```bash
# Dockerイメージのビルドとロード
mise run k8s:local:cluster:load-image

# 全リソースのデプロイ（platform + workloads）
mise run k8s:local:deploy
```

### クラスタの削除
```bash
mise run k8s:local:cluster:delete
```

## 利用可能なコマンド
```bash
mise tasks | grep k8s
```
