---
name: import-jira-to-todo
description: Jira の課題をローカルの todo システムにインポートするスキル。「Jira の課題を todo にインポートして」「スプリントのタスクを取り込んで」「Jira から課題を同期して」など、atl jira の出力を todo datasource import に流し込む操作を依頼された場合に使用。
---

# Import Jira to Todo

`atl jira` の JSON 出力を `todo datasource import jira --stdin` にインポートするパイプライン。

## Default Behavior

常に現在のユーザ (`assignee = currentUser()`) の未完了タスク (`statusCategory not in (Done)`) をインポートする。

```bash
atl jira issue list --site <SITE> \
  --jql 'statusCategory not in (Done) AND assignee = currentUser()' \
  --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import <DATASOURCE> --stdin
```

## JQL Notes

- `statusCategory not in (Done)` を使うこと。`status not in (Done)` ではカスタム完了ステータス（採用・不採用など）が除外されない
- Claude Code の Bash ツールは `!` を `\!` にエスケープするため、`!=` 演算子は使えない。代わりに `not in ()` を使う（通常のターミナルでは `!=` も動作する）
## Examples

```bash
# サイト指定でインポート
atl jira issue list --site mysite \
  --jql 'statusCategory not in (Done) AND assignee = currentUser()' \
  --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import my-datasource --stdin

# プロジェクトを追加で絞り込む場合
atl jira issue list --site mysite \
  --jql 'statusCategory not in (Done) AND assignee = currentUser() AND project = PROJ' \
  --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import my-datasource --stdin
```

## Status Mapping

| Jira statusCategory | todo |
|------|------|
| To Do | pending |
| In Progress | in_progress |
| Done | done |

Individual status mapping in `import-helper.sh`:

| Jira status | todo |
|------|------|
| To Do | pending |
| In Progress / In Review | in_progress |
| Done | done |
| その他 | pending |

## Resources

- `scripts/import-helper.sh` - jq で JSON を JSONL に変換
- `references/import.md` - JSONL フォーマット仕様
