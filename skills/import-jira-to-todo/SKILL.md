---
name: import-jira-to-todo
description: Jira の課題をローカルの todo システムにインポートするスキル。「Jira の課題を todo にインポートして」「スプリントのタスクを取り込んで」「Jira から課題を同期して」など、atl jira の出力を todo datasource import に流し込む操作を依頼された場合に使用。
---

# Import Jira to Todo

`atl jira` の課題を取得・変換し `todo datasource import` にインポートするパイプライン。

## Default Behavior

常に現在のユーザ (`assignee = currentUser()`) の未完了タスク (`statusCategory not in (Done)`) をインポートする。

```bash
bash skills/import-jira-to-todo/scripts/import-helper.sh <SITE> <DATASOURCE>
```

## Examples

```bash
# サイト指定でインポート
bash skills/import-jira-to-todo/scripts/import-helper.sh mysite my-datasource
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

- `scripts/import-helper.sh` - fetch → transform → import パイプライン
- `references/import.md` - JSONL フォーマット仕様
