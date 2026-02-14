---
name: import-jira-to-todo
description: Jira の課題をローカルの todo システムにインポートするスキル。「Jira の課題を todo にインポートして」「スプリントのタスクを取り込んで」「Jira から課題を同期して」など、atl jira の出力を todo datasource import に流し込む操作を依頼された場合に使用。
---

# Import Jira to Todo

`atl jira` の JSON 出力を `todo datasource import jira --stdin` にインポートするパイプライン。

## Usage

```bash
atl jira <command> --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin
```

### Examples

```bash
# プロジェクトの課題をインポート
atl jira issue list --project PROJ --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin

# スプリントの課題をインポート
atl jira sprint issues --sprint 100 --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin

# 単一課題をインポート（url, description 付き）
atl jira issue view --key PROJ-123 --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin
```

## Status Mapping

| Jira | todo |
|------|------|
| To Do | pending |
| In Progress | in_progress |
| In Review | in_progress |
| Done | done |
| その他 | pending |

## Resources

- `scripts/import-helper.sh` - jq で JSON を JSONL に変換
- `references/import.md` - JSONL フォーマット仕様
