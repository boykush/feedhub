# Claude Code Development Guide

このファイルには、Claude Codeを使用してこのプロジェクトを開発する際のガイドラインが記載されています。

## プロジェクト概要

@README.md を参照してください。

## 言語設定

- **コミットメッセージ**: 英語
- **コミュニケーション**: 日本語

## mise ツールの実行

Claude CodeのBashツールでmiseツール（buf、kubectl、helm等）を実行する場合は、`mise exec` 経由で実行してください。

```bash
# ❌ 直接実行（失敗する可能性あり）
buf generate

# ✅ mise exec経由で実行
mise exec -- buf generate
mise exec -- kubectl get pods
```

## mise File Tasks

タスク定義は @.mise-tasks/ に配置されています。

```bash
# タスク一覧
mise tasks

# タスク実行
mise run proto:generate
mise run go:build transaction
```
