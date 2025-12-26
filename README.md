# Foresee

マイクロサービスアーキテクチャで構築されたアプリケーション。

## プロジェクト構成

```
foresee/
├── services/          # マイクロサービス群
│   ├── bff/          # Backend for Frontend (gRPC Gateway)
│   └── transaction/  # 取引管理サービス
├── k8s/              # Kubernetesマニフェスト
│   ├── base/         # ベースマニフェスト
│   └── overlays/     # 環境別設定
└── .mise-tasks/      # タスク定義（mise用）
```

詳細は各ディレクトリのREADMEを参照してください：
- [services/README.md](services/README.md) - マイクロサービスアーキテクチャとサービス一覧
- [k8s/README.md](k8s/README.md) - Kubernetesマニフェストの構成と管理

## 開発環境のセットアップ

### mise のインストール

このプロジェクトでは、[mise](https://mise.jdx.dev/)を使用してツールのバージョン管理とタスク実行を行っています。

```bash
# mise のインストール（macOS）
brew install mise

# シェル設定の追加（例: zsh）
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc
source ~/.zshrc
```

### 依存ツールのインストール

プロジェクトルートで以下を実行すると、[.mise.toml](.mise.toml)で定義された必要なツールが自動的にインストールされます：

```bash
mise install

# インストール済みツールの確認
mise ls
```

## タスク管理

利用可能なタスクは[.mise-tasks/](.mise-tasks/)ディレクトリに定義されています。

```bash
# タスク一覧の確認
mise tasks
```

## クイックスタート

```bash
# 1. 依存ツールのインストール
mise install

# 2. ローカルk8sクラスタ作成
mise run k8s:local:cluster:create

# 3. Dockerイメージのビルドとロード
mise run docker:build
mise run k8s:local:cluster:load-image

# 4. Kubernetesリソースのデプロイ
mise run k8s:local:apply
```
