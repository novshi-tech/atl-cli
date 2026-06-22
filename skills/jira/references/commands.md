# atl コマンドリファレンス

## jira me

認証済みの Jira ユーザー（自分自身）のアカウント情報を表示する。

```
atl jira me [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Account ID:  5b10ac8d14c052e1e6c2e251
Name:        John Doe
Email:       john@example.com
Active:      active
```

**JSON 出力例** (`--json`):
```json
{
  "accountId": "5b10ac8d14c052e1e6c2e251",
  "displayName": "John Doe",
  "emailAddress": "john@example.com",
  "active": true
}
```

## jira issue create

新しい Jira 課題を作成する。

```
atl jira issue create [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--project` | `-p` | Yes | - | プロジェクトキー |
| `--summary` | `-s` | Yes | - | 課題サマリー |
| `--type` | `-t` | Yes | - | 課題タイプ名（プロジェクトごとに異なるため、後述の `atl jira issuetype list --project <key>` で確認すること） |
| `--description` | `-d` | No | - | 課題の説明 |
| `--due` | - | No | - | 期日（YYYY-MM-DD） |
| `--epic` | - | No | - | 紐づけるエピックのキー（例: `PROJ-10`） |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Created issue: PROJ-456
URL: https://example.atlassian.net/browse/PROJ-456
```

**JSON 出力例** (`--json`):
```json
{
  "key": "PROJ-456",
  "url": "https://example.atlassian.net/browse/PROJ-456"
}
```

## jira issue list

JQL で課題を検索する。

```
atl jira issue list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--jql` | - | Yes | - | JQL クエリ文字列 |
| `--max` | - | No | `50` | 最大取得件数 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 3 issue(s):

PROJ-123      In Progress      John Doe              ログイン画面のバグ修正
PROJ-124      To Do            Unassigned            新機能の設計
PROJ-125      Done             Jane Smith            API ドキュメント更新
```

**JSON 出力例** (`--json`):
```json
[
  {
    "key": "PROJ-123",
    "summary": "ログイン画面のバグ修正",
    "status": "In Progress",
    "type": "Bug",
    "assignee": "John Doe"
  }
]
```

## jira issue view

課題の詳細情報（サマリー、ステータス、タイプ、アサイニー、説明、コメント）を表示する。

```
atl jira issue view [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--key` | `-k` | Yes | - | 課題キー |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Key:       PROJ-123
Summary:   ログイン画面のバグ修正
Status:    In Progress
Type:      Bug
Assignee:  John Doe
Epic:      PROJ-10
URL:       https://example.atlassian.net/browse/PROJ-123

--- Description ---
ログイン画面でエラーが発生する問題を修正する。

--- Comments (1) ---

[2024-01-15T10:30:00.000+0000] Jane Smith:
修正方針を確認しました。
```

**JSON 出力例** (`--json`):
```json
{
  "key": "PROJ-123",
  "summary": "ログイン画面のバグ修正",
  "status": "In Progress",
  "type": "Bug",
  "assignee": "John Doe",
  "url": "https://example.atlassian.net/browse/PROJ-123",
  "description": "ログイン画面でエラーが発生する問題を修正する。",
  "epic": "PROJ-10",
  "comments": [
    {
      "author": "Jane Smith",
      "created": "2024-01-15T10:30:00.000+0000",
      "body": "修正方針を確認しました。"
    }
  ]
}
```

## jira issue update

既存の課題を更新する。`--summary`、`--description`、`--status`、`--assignee`、`--epic` のいずれかを指定する。

```
atl jira issue update [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--key` | `-k` | Yes | - | 課題キー |
| `--summary` | `-s` | No | - | 新しいサマリー |
| `--description` | `-d` | No | - | 新しい説明 |
| `--status` | - | No | - | 遷移先ステータス |
| `--assignee` | - | No | - | 担当者の accountId（`none` で担当者解除） |
| `--due` | - | No | - | 期日（YYYY-MM-DD） |
| `--epic` | - | No | - | 紐づけるエピックのキー（例: `PROJ-10`） |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

## jira issue comment

課題にコメントを追加する。

```
atl jira issue comment [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--key` | `-k` | Yes | - | 課題キー |
| `--body` | `-b` | Yes | - | コメント本文 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

### メンション構文

`--body` 内で `@[表示名:accountId]` 構文を使うと、Jira のメンション通知が送信される。

```bash
atl jira issue comment --key PROJ-123 --body "@[山田太郎:5b10ac8d14c052e1e6c2e251] 確認をお願いします。"
```

accountId は `atl jira user search --query "名前" --json` で取得できる。

## jira issue attachment list

課題に添付されたファイルの一覧を表示する。

```
atl jira issue attachment list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--key` | `-k` | Yes | - | 課題キー |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 2 attachment(s) on PROJ-123:

10001         12.3KB      John Doe              screenshot.png
10002         4.2MB       Jane Smith            spec.pdf
```

**JSON 出力例** (`--json`):
```json
[
  {
    "id": "10001",
    "filename": "screenshot.png",
    "size": 12582,
    "mimeType": "image/png",
    "author": "John Doe",
    "created": "2024-01-15T10:30:00.000+0000",
    "content": "https://example.atlassian.net/rest/api/3/attachment/content/10001"
  }
]
```

## jira issue attachment download

課題の添付ファイルをダウンロードする。`--id` で直接指定するか、`--key` + `--filename` でファイル名から取得する。

```
atl jira issue attachment download [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--id` | - | No* | - | 添付ファイル ID |
| `--key` | `-k` | No* | - | 課題キー（`--filename` と併用） |
| `--filename` | - | No* | - | ファイル名（`--key` と併用） |
| `--output` | `-o` | No | サーバ提供のファイル名 | 出力先パス（ディレクトリ可、`-` で標準出力） |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

\* `--id` または (`--key` + `--filename`) のいずれかが必須。

**出力例:**
```
Downloaded attachment 10001 (12.3KB) to screenshot.png
```

**JSON 出力例** (`--json`):
```json
{
  "id": "10001",
  "filename": "screenshot.png",
  "path": "screenshot.png",
  "size": 12582
}
```

## jira sprint list

ボードのスプリント一覧を表示する。

```
atl jira sprint list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--board` | - | Yes | - | ボード ID |
| `--state` | - | No | - | 状態フィルタ（`active` / `closed` / `future`） |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
42      active      Sprint 10
43      future      Sprint 11
```

**JSON 出力例** (`--json`):
```json
[
  {
    "id": 42,
    "name": "Sprint 10",
    "state": "active"
  }
]
```

## jira sprint issues

スプリント内の課題一覧を表示する。

```
atl jira sprint issues [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--sprint` | - | Yes | - | スプリント ID |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

## jira user search

ユーザーを表示名またはメールアドレスで検索する。

```
atl jira user search [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--query` | `-q` | Yes | - | 検索文字列（表示名・メールアドレス） |
| `--max` | - | No | `50` | 最大取得件数 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 2 user(s):

5b10ac8d14c052e1e6c2e251          John Doe                  john@example.com                active
5b10a2844c20165700ede21g          Jane Smith                                                 active
```

**JSON 出力例** (`--json`):
```json
[
  {
    "accountId": "5b10ac8d14c052e1e6c2e251",
    "displayName": "John Doe",
    "emailAddress": "john@example.com",
    "active": true
  }
]
```

## jira issuetype list

利用可能な課題タイプを一覧表示する。`--project` を指定するとそのプロジェクトで実際に作成可能な課題タイプ（Jira の createmeta エンドポイントが返すもの）に絞り込まれる。

```
atl jira issuetype list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--project` | `-p` | No | - | プロジェクトキーまたは数値 ID。指定するとそのプロジェクトで作成可能なタイプのみ表示 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
タスク                   行うべきタスクを表します。
ストーリー                 ストーリー追跡機能、またはユーザー目標として表明された機能
バグ                    バグ
エピック                  大きなユーザーストーリー向けの課題タイプです。
サブタスク                 The sub-task of the issue [subtask]
```

**JSON 出力例** (`--json`):
```json
[
  {
    "id": "10001",
    "name": "ストーリー",
    "description": "ストーリー追跡機能、またはユーザー目標として表明された機能",
    "subtask": false
  }
]
```

> Jira プロジェクトの設定によっては、ここに表示されたタイプ名でも `jira issue create` が `選択した課題タイプは無効です` を返すことがある。これはプロジェクトの課題タイプスキーム／画面構成スキームの不整合に起因し、Jira 管理者側での修正が必要となる。

## configure

サイトの認証情報を設定する。

```
atl configure --site <alias>
```

対話的にサイト URL、メールアドレス、API トークンを入力する。

## sites

設定済みサイトの一覧を表示する。

```
atl sites [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--json` | - | No | `false` | JSON 形式で出力 |

**JSON 出力例** (`--json`):
```json
[
  {
    "name": "mysite",
    "default": true
  }
]
```
