# Jira 課題インポートガイド

`atl jira` の JSON 出力を `/todo` にインポートするためのリファレンス。

## ヘルパースクリプト

`skills/import-jira-to-todo/scripts/import-helper.sh` は `atl jira` の JSON 出力を `/todo` の JSONL 形式に変換する。

### 対応入力

| コマンド | 出力形式 |
|---------|---------|
| `atl jira issue list --json` | JSON 配列 |
| `atl jira sprint issues --json` | JSON 配列 |
| `atl jira issue view --json` | 単一オブジェクト（url, description 含む） |

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
atl jira issue list --site mysite \
  --jql 'statusCategory not in (Done) AND assignee = currentUser()' \
  --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import my-datasource --stdin
```

### プロジェクトを追加で絞り込む場合

```bash
atl jira issue list --site mysite \
  --jql 'statusCategory not in (Done) AND assignee = currentUser() AND project = PROJ' \
  --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import my-datasource --stdin
```
