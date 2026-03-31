---
name: bitbucket
description: Bitbucket Cloud の操作を行うスキル。リポジトリの一覧・詳細表示、プルリクエストの一覧・作成・コメント取得・コメント投稿をサポート。「Bitbucketのリポジトリを一覧して」「PRを一覧して」「PRを作成して」「PRのコメントを確認して」「PRにコメントして」など Bitbucket 関連の操作を依頼された場合に使用。
---

# Bitbucket

`atl bitbucket` を使って Bitbucket Cloud を操作する。事前に `atl configure` でサイト設定が必要。

## Quick Start

```bash
# リポジトリを一覧
atl bitbucket repo list

# リポジトリの詳細を表示
atl bitbucket repo get --repo my-app

# プルリクエストを一覧
atl bitbucket pr list --repo my-app

# プルリクエストを作成
atl bitbucket pr create --repo my-app --title "新機能追加" --source feature/new-feature

# PRのコメントを取得（インラインコメントもデフォルトで含まれる、解決済みは除外）
atl bitbucket pr comment --repo my-app --pr 42

# インラインコメントを除外して取得
atl bitbucket pr comment --repo my-app --pr 42 --inline=false

# 解決済みコメントも含めて取得
atl bitbucket pr comment --repo my-app --pr 42 --include-resolved

# PRにコメントを投稿
atl bitbucket pr comment create --repo my-app --pr 42 --body "LGTM!"

# インラインコメントを投稿
atl bitbucket pr comment create --repo my-app --pr 42 --body "この変数名を変更してください" --path src/auth.go --line 15
```

## ワークスペースの解決

`--workspace` はサイト設定にワークスペースが保存されていれば省略可能。両方指定した場合は一致チェックが行われる。サイト設定にもフラグにもない場合はエラーになる。

## 共通フラグ

- `--json` - 機械可読な JSON 形式で出力する（AI エージェント連携に推奨、常に付与すること）
- `--site` - 使用するサイトエイリアス（未指定時はデフォルトサイト）
- `--workspace` - ワークスペースのスラッグ（省略時はサイト設定値を使用）

## Resources

- [references/commands.md](references/commands.md) - 全コマンドの詳細リファレンス（フラグ・出力例）
- [references/workflows.md](references/workflows.md) - 典型的なワークフロー例
