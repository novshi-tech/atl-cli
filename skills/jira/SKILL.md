---
name: jira
description: Jira Cloud の操作を行うスキル。課題の作成・検索・表示・更新・ステータス遷移・コメント追加、スプリント管理をサポート。「Jiraで課題を作成して」「Jiraの課題を検索して」「スプリントの一覧を見せて」など Jira 関連の操作を依頼された場合に使用。
---

# Jira

## Overview

`jira-cli` を使って Jira Cloud を操作するスキル。課題の CRUD 操作、JQL 検索、スプリント管理をコマンドラインから実行できる。

## Quick Start

すべてのコマンドは `jira-cli` を経由して実行する。事前に `jira-cli configure` でサイト設定が必要。

```bash
# 課題を検索
jira-cli issue list --project PROJ --status "In Progress"

# 課題の詳細を表示
jira-cli issue view --key PROJ-123

# 課題を作成
jira-cli issue create --project PROJ --summary "バグ修正" --type Bug --description "詳細説明"

# 課題を更新
jira-cli issue update --key PROJ-123 --status "Done"

# コメントを追加
jira-cli issue comment --key PROJ-123 --body "対応完了しました"
```

## 課題操作

### 課題を検索する (`issue list`)

JQL クエリまたはフィルタフラグで課題を検索する。

```bash
# JQL で直接検索
jira-cli issue list --jql "project = PROJ AND status = 'In Progress' ORDER BY updated DESC"

# フィルタフラグで検索（JQL を自動生成）
jira-cli issue list --project PROJ --status "To Do" --assignee me --max 20
```

**フラグ:**
- `--jql` - JQL クエリ文字列（指定時は他のフィルタフラグは無視される）
- `--project` / `-p` - プロジェクトキーでフィルタ
- `--status` - ステータスでフィルタ
- `--assignee` - アサイニーでフィルタ（`me` で現在のユーザー）
- `--max` - 最大件数（デフォルト: 50）

### 課題の詳細を表示する (`issue view`)

課題のサマリー、ステータス、タイプ、アサイニー、説明、コメントを表示する。

```bash
jira-cli issue view --key PROJ-123
```

### 課題を作成する (`issue create`)

新しい課題を作成する。

```bash
jira-cli issue create --project PROJ --summary "新機能の実装" --type Story --description "機能の詳細説明"
```

**フラグ:**
- `--project` / `-p` - プロジェクトキー（必須）
- `--summary` / `-s` - 課題のサマリー（必須）
- `--type` / `-t` - 課題タイプ（デフォルト: Task）
- `--description` / `-d` - 説明

### 課題を更新する (`issue update`)

既存の課題のサマリー、説明、ステータスを更新する。

```bash
# サマリーを更新
jira-cli issue update --key PROJ-123 --summary "更新後のサマリー"

# ステータスを遷移
jira-cli issue update --key PROJ-123 --status "In Progress"

# 複数のフィールドを同時に更新
jira-cli issue update --key PROJ-123 --summary "新サマリー" --status "Done"
```

### コメントを追加する (`issue comment`)

課題にコメントを追加する。

```bash
jira-cli issue comment --key PROJ-123 --body "PR をレビューしてください"
```

## スプリント操作

### スプリント一覧 (`sprint list`)

ボードのスプリント一覧を表示する。

```bash
# すべてのスプリント
jira-cli sprint list --board 42

# アクティブなスプリントのみ
jira-cli sprint list --board 42 --state active
```

**フラグ:**
- `--board` - ボード ID（必須）
- `--state` - 状態フィルタ（active / closed / future）

### スプリント内の課題一覧 (`sprint issues`)

特定のスプリントに含まれる課題を一覧表示する。

```bash
jira-cli sprint issues --sprint 100
```

## 共通フラグ

すべてのコマンドで以下のフラグが使用可能:

- `--site` - 使用するサイトエイリアス（未指定時はデフォルトサイト）

## Resources

### references/

- `commands.md` - 全コマンドの詳細リファレンス
- `workflows.md` - 典型的なワークフロー例
