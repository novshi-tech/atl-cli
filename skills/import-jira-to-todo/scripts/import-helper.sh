#!/usr/bin/env bash
# Jira の課題を取得し、todo datasource import にインポートするヘルパースクリプト
#
# Usage:
#   bash import-helper.sh <SITE> <DATASOURCE>
#
# 使用例:
#   bash import-helper.sh mysite my-datasource

set -euo pipefail

SITE="${1:?Usage: import-helper.sh <SITE> <DATASOURCE>}"
DATASOURCE="${2:?Usage: import-helper.sh <SITE> <DATASOURCE>}"

atl jira issue list --site "${SITE}" \
  --jql 'statusCategory not in (Done) AND assignee = currentUser()' \
  --json \
  | jq -c '
    # ステータスマッピング
    def map_status:
      if . == "To Do" then "pending"
      elif . == "In Progress" then "in_progress"
      elif . == "In Review" then "in_progress"
      elif . == "Done" then "done"
      else "pending"
      end;

    # 配列・単一オブジェクト両対応
    (if type == "array" then .[] else . end) |
    {
      remote_id: .key,
      title: .summary,
      status: (.status | map_status),
      tags: [.type]
    } + (if .url then {url: .url} else {} end)
      + (if .description then {description: .description} else {} end)
  ' \
  | todo datasource import "${DATASOURCE}" --stdin
