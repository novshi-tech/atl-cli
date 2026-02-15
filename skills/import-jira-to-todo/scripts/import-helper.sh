#!/usr/bin/env bash
# atl jira の JSON 出力を /todo の JSONL 形式に変換するヘルパースクリプト
#
# 対応入力:
#   - atl jira issue list --json   (JSON 配列)
#   - atl jira sprint issues --json (JSON 配列)
#   - atl jira issue view --json   (単一オブジェクト)
#
# 使用例:
#   atl jira issue list --jql "project = PROJ ORDER BY updated DESC" --json | bash import-helper.sh
#   atl jira issue view --key PROJ-123 --json | bash import-helper.sh

jq -c '
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
'
