---
name: bitbucket
description: Bitbucket Cloud の操作を行うスキル。リポジトリの一覧・詳細表示、プルリクエストの一覧・作成・コメント取得をサポート。「Bitbucketのリポジトリを一覧して」「PRを一覧して」「PRを作成して」「PRのコメントを確認して」など Bitbucket 関連の操作を依頼された場合に使用。
---

# Bitbucket

## Overview

`atl bitbucket` を使って Bitbucket Cloud を操作するスキル。リポジトリの一覧・詳細取得、プルリクエストの一覧・作成・コメント取得をコマンドラインから実行できる。

## Quick Start

すべてのコマンドは `atl bitbucket` を経由して実行する。事前に `atl configure` でサイト設定が必要。

```bash
# リポジトリを一覧
atl bitbucket repo list --workspace myteam

# リポジトリの詳細を表示
atl bitbucket repo get --workspace myteam --repo my-app

# プルリクエストを一覧
atl bitbucket pr list --workspace myteam --repo my-app

# プルリクエストを作成
atl bitbucket pr create --workspace myteam --repo my-app --title "新機能追加" --source feature/new-feature

# PRのコメントを取得
atl bitbucket pr comment --workspace myteam --repo my-app --pr 42
```

## リポジトリ操作

### リポジトリを一覧する (`repo list`)

ワークスペース内のリポジトリを一覧表示する。

```bash
atl bitbucket repo list --workspace myteam
atl bitbucket repo list --workspace myteam --max 50
```

**フラグ:**
- `--workspace` - ワークスペースのスラッグ（必須）
- `--max` - 最大件数（デフォルト: 25）

### リポジトリの詳細を表示する (`repo get`)

リポジトリのスラッグ、名前、説明、言語、公開/非公開、メインブランチ、更新日時を表示する。

```bash
atl bitbucket repo get --workspace myteam --repo my-app
```

**フラグ:**
- `--workspace` - ワークスペースのスラッグ（必須）
- `--repo` - リポジトリのスラッグ（必須）

## プルリクエスト操作

### プルリクエストを一覧する (`pr list`)

リポジトリのプルリクエストを一覧表示する。

```bash
# オープンなPRを一覧（デフォルト）
atl bitbucket pr list --workspace myteam --repo my-app

# マージ済みPRを一覧
atl bitbucket pr list --workspace myteam --repo my-app --state MERGED
```

**フラグ:**
- `--workspace` - ワークスペースのスラッグ（必須）
- `--repo` - リポジトリのスラッグ（必須）
- `--state` - 状態でフィルタ: `OPEN` / `MERGED` / `DECLINED` / `SUPERSEDED`（デフォルト: OPEN）
- `--max` - 最大件数（デフォルト: 25）

### プルリクエストを作成する (`pr create`)

新しいプルリクエストを作成する。

```bash
# 基本的な作成（dest はリポジトリのメインブランチがデフォルト）
atl bitbucket pr create --workspace myteam --repo my-app --title "バグ修正" --source fix/login-bug

# ターゲットブランチと説明を指定
atl bitbucket pr create --workspace myteam --repo my-app --title "新機能" --source feature/auth --dest develop -d "OAuth2 認証を追加"
```

**フラグ:**
- `--workspace` - ワークスペースのスラッグ（必須）
- `--repo` - リポジトリのスラッグ（必須）
- `--title` - プルリクエストのタイトル（必須）
- `--source` - ソースブランチ名（必須）
- `--dest` - デスティネーションブランチ（省略時はリポジトリのメインブランチ）
- `--description` / `-d` - プルリクエストの説明

### PRのコメントを取得する (`pr comment`)

プルリクエストのコメントを一覧表示する（インラインコードレビューコメントは除外）。

```bash
atl bitbucket pr comment --workspace myteam --repo my-app --pr 42
```

**フラグ:**
- `--workspace` - ワークスペースのスラッグ（必須）
- `--repo` - リポジトリのスラッグ（必須）
- `--pr` - プルリクエスト ID（必須）

## 共通フラグ

すべてのサブコマンドで以下のフラグが使用可能:

- `--json` - 機械可読な JSON 形式で出力する（AI エージェント連携に推奨）
- `--site` - 使用するサイトエイリアス（未指定時はデフォルトサイト）

### `--json` フラグ

`--json` を付けるとすべてのコマンドが構造化 JSON を出力する。AI エージェントから呼び出す場合は常に `--json` を使用すること。

```bash
# リポジトリ一覧を JSON で取得
atl bitbucket repo list --workspace myteam --json

# リポジトリの詳細を JSON で取得
atl bitbucket repo get --workspace myteam --repo my-app --json

# PRを作成して結果を JSON で取得
atl bitbucket pr create --workspace myteam --repo my-app --title "修正" --source fix/bug --json
```

## Resources

### references/

- `commands.md` - 全コマンドの詳細リファレンス
- `workflows.md` - 典型的なワークフロー例
