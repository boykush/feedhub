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
│       ├── transaction/
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
  - マイクロサービス（bff, transaction等）
  - Istio設定（Gateway, VirtualService, AuthorizationPolicy等）

この分離により、プラットフォーム層とアプリケーション層を独立して管理できます。

### Base vs Overlays

- **base/**: 環境に依存しない共通のマニフェスト
- **overlays/**: 環境ごとの設定（local, staging, production等）

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
mise run docker:build
mise run k8s:local:cluster:load-image

# 全リソースのデプロイ（platform + workloads）
mise run k8s:local:deploy
```

### ArgoCD GUIの起動
```bash
mise run k8s:local:argocd-gui
```

### クラスタの削除
```bash
mise run k8s:local:cluster:delete
```

## 利用可能なコマンド
```bash
mise tasks | grep k8s
```

## GitOps with ArgoCD

ArgoCDを使用してGitリポジトリからデプロイを自動化します。

### ローカル環境

- **リポジトリ**: 公開GitHubリポジトリ（認証不要）
- **同期モード**: Auto sync有効
- **デプロイ方法**: ArgoCDが自動的にGitリポジトリと同期

### 本番環境

- **リポジトリ**: 公開GitHubリポジトリまたはプライベートリポジトリ（認証が必要な場合）
- **同期モード**: Auto sync有効
- **デプロイ方法**: ArgoCDが自動的にGitリポジトリと同期
