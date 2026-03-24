---
name: bitbucket
description: Bitbucket Cloud の操作を行うスキル。リポジトリの一覧・詳細表示、プルリクエストの一覧・作成・コメント取得をサポート。「Bitbucketのリポジトリを一覧して」「PRを一覧して」「PRを作成して」「PRのコメントを確認して」など Bitbucket 関連の操作を依頼された場合に使用。
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

# PRのコメントを取得（インラインコメントもデフォルトで含まれる）
atl bitbucket pr comment --repo my-app --pr 42

# インラインコメントを除外して取得
atl bitbucket pr comment --repo my-app --pr 42 --inline=false
```

## ワークスペースの解決

`--workspace` はサイト設定にワークスペースが保存されていれば省略可能。両方指定した場合は一致チェックが行われる。サイト設定にもフラグにもない場合はエラーになる。

## PRにコメントを投稿する（API直接呼び出し）

`atl bitbucket` にはPRコメント投稿コマンドがないため、Bitbucket REST API を直接呼び出す。

```bash
# 認証情報を取得
BASE_URL=$(atl auth get --field base-url --json | jq -r '.value')
EMAIL=$(atl auth get --field email --json | jq -r '.value')
BB_TOKEN=$(atl auth get --field bb-api-token --json | jq -r '.value')

# PRにコメントを投稿
curl -s -u "${EMAIL}:${BB_TOKEN}" \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"content": {"raw": "レビューコメントです"}}' \
  "https://api.bitbucket.org/2.0/repositories/{workspace}/{repo_slug}/pullrequests/{pr_id}/comments"
```

## 共通フラグ

- `--json` - 機械可読な JSON 形式で出力する（AI エージェント連携に推奨、常に付与すること）
- `--site` - 使用するサイトエイリアス（未指定時はデフォルトサイト）
- `--workspace` - ワークスペースのスラッグ（省略時はサイト設定値を使用）

## Resources

- [references/commands.md](references/commands.md) - 全コマンドの詳細リファレンス（フラグ・出力例）
- [references/workflows.md](references/workflows.md) - 典型的なワークフロー例
