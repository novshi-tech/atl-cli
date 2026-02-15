# Jira 課題インポートガイド

`atl jira` の課題を `/todo` にインポートするためのリファレンス。

## ヘルパースクリプト

`skills/import-jira-to-todo/scripts/import-helper.sh` は Jira から課題を取得・変換・インポートを実行する。

### 引数

| 引数 | 説明 |
|------|------|
| `SITE` | Jira サイト名 |
| `DATASOURCE` | todo datasource 名 |

## JSONL フォーマット仕様

各行が 1 つの JSON オブジェクト。

### 必須フィールド

| フィールド | 型 | 説明 |
|-----------|------|------|
| `remote_id` | string | リモート側の一意識別子（Jira の課題キー） |
| `title` | string | タスクのタイトル |

### 任意フィールド

| フィールド | 型 | 説明 |
|-----------|------|------|
| `status` | string | `pending` / `in_progress` / `done` |
| `tags` | string[] | タグの配列 |
| `url` | string | 課題の URL |
| `description` | string | 課題の説明 |

### ステータスマッピング

| Jira ステータス | todo ステータス |
|----------------|----------------|
| `To Do` | `pending` |
| `In Progress` | `in_progress` |
| `In Review` | `in_progress` |
| `Done` | `done` |
| その他 | `pending` |

## JQL に関する注意事項

- 未完了タスクのフィルタには `statusCategory not in (Done)` を使う。`status not in (Done)` ではカスタム完了ステータス（採用・不採用など）が除外されない
- Claude Code の Bash ツールは `!` を `\!` にエスケープするため、`!=` 演算子は使えない。代わりに `not in ()` を使う（通常のターミナルでは `!=` も動作する）

## 使用例

### 現在のユーザの未完了課題をインポート（デフォルト）

```bash
bash skills/import-jira-to-todo/scripts/import-helper.sh mysite my-datasource
```
