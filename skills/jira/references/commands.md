# atl コマンドリファレンス

## jira issue create

新しい Jira 課題を作成する。

```
atl jira issue create [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--project` | `-p` | Yes | - | プロジェクトキー |
| `--summary` | `-s` | Yes | - | 課題サマリー |
| `--type` | `-t` | No | `Task` | 課題タイプ（Task, Bug, Story, Epic 等） |
| `--description` | `-d` | No | - | 課題の説明 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |

**出力例:**
```
Created issue: PROJ-456
URL: https://example.atlassian.net/browse/PROJ-456
```

## jira issue list

JQL で課題を検索する。`--jql` を指定した場合、他のフィルタフラグは無視される。`--jql` を指定しない場合、`--project`、`--status`、`--assignee` から JQL を自動生成する。

```
atl jira issue list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--jql` | - | No | - | JQL クエリ文字列 |
| `--project` | `-p` | No | - | プロジェクトキーでフィルタ |
| `--status` | - | No | - | ステータスでフィルタ |
| `--assignee` | - | No | - | アサイニーでフィルタ（`me` で自分） |
| `--max` | - | No | `50` | 最大取得件数 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |

**出力例:**
```
Found 3 issue(s):

PROJ-123      In Progress      John Doe              ログイン画面のバグ修正
PROJ-124      To Do            Unassigned            新機能の設計
PROJ-125      Done             Jane Smith            API ドキュメント更新
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

**出力例:**
```
Key:       PROJ-123
Summary:   ログイン画面のバグ修正
Status:    In Progress
Type:      Bug
Assignee:  John Doe
URL:       https://example.atlassian.net/browse/PROJ-123

--- Description ---
ログイン画面でエラーが発生する問題を修正する。

--- Comments (1) ---

[2024-01-15T10:30:00.000+0000] Jane Smith:
修正方針を確認しました。
```

## jira issue update

既存の課題を更新する。`--summary`、`--description`、`--status` のいずれかを指定する。

```
atl jira issue update [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--key` | `-k` | Yes | - | 課題キー |
| `--summary` | `-s` | No | - | 新しいサマリー |
| `--description` | `-d` | No | - | 新しい説明 |
| `--status` | - | No | - | 遷移先ステータス |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |

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

**出力例:**
```
42      active      Sprint 10
43      future      Sprint 11
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

## configure

サイトの認証情報を設定する。

```
atl configure --site <alias>
```

対話的にサイト URL、メールアドレス、API トークンを入力する。

## sites

設定済みサイトの一覧を表示する。

```
atl sites
```
