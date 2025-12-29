# Terraform Configuration

ForeseeプロジェクトのGCPインフラストラクチャをTerraformで管理します。

## ディレクトリ構成

```
terraform/
├── versions.tf           # Terraformバージョンとbackend設定
├── provider.tf           # GCPプロバイダー設定
├── variables.tf          # 変数定義
├── main.tf               # 共通リソース定義
└── artifact-registry.tf  # Artifact Registryリソース
```

## セットアップ

### 1. GCP認証

```bash
gcloud auth application-default login
gcloud config set project YOUR_PROJECT_ID
```

### 2. 環境変数の設定

`.mise.local.toml`を作成し、プロジェクトIDを設定します:

```toml
[env]
TF_VAR_project_id = "YOUR_PROJECT_ID"
```

**注意**: `.mise.local.toml`は`.gitignore`に含まれており、Gitにコミットされません。

### 3. Terraform の初期化

```bash
mise run tf:init
```

## 使い方

### mise タスク経由（推奨）

```bash
# 初期化
mise run tf:init

# Plan の実行
mise run tf:plan

# Apply の実行
mise run tf:apply
```

### 直接実行

```bash
cd terraform

# Plan
mise exec -- terraform plan

# Apply
mise exec -- terraform apply
```

**重要**: `mise exec`経由で実行することで、`.mise.local.toml`の環境変数が自動的に読み込まれます。
