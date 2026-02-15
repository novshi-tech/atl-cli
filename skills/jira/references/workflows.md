# 典型的なワークフロー

## 1. 課題を確認して作業を開始する

自分にアサインされた課題を確認し、作業を開始する。

```bash
# 自分の課題を確認
atl jira issue list --jql "project = PROJ AND assignee = currentUser() AND status = 'To Do' ORDER BY updated DESC"

# 課題の詳細を確認
atl jira issue view --key PROJ-123

# ステータスを "In Progress" に変更
atl jira issue update --key PROJ-123 --status "In Progress"
```

## 2. 課題を作成してブランチで作業する

新しい課題を作成し、対応するブランチで作業を進める。

```bash
# 課題を作成
atl jira issue create --project PROJ --summary "ユーザー認証の実装" --type Story --description "OAuth2 を使った認証機能を実装する"

# 出力された課題キーでブランチを作成
git checkout -b feature/PROJ-456-user-auth

# 作業中にステータスを更新
atl jira issue update --key PROJ-456 --status "In Progress"

# 作業完了後
atl jira issue comment --key PROJ-456 --body "PR #42 を作成しました"
atl jira issue update --key PROJ-456 --status "In Review"
```

## 3. スプリントの進捗を確認する

現在のスプリントの状況を確認する。

```bash
# アクティブなスプリントを確認
atl jira sprint list --board 42 --state active

# スプリント内の課題を確認
atl jira sprint issues --sprint 100
```

## 4. バグ報告と対応

バグを報告し、対応する。

```bash
# バグを作成
atl jira issue create --project PROJ --summary "ログイン画面で500エラー" --type Bug --description "メールアドレスに+を含む場合にサーバーエラーが発生する"

# 自分で対応する場合はステータスを変更
atl jira issue update --key PROJ-789 --status "In Progress"

# 修正後にコメントとステータス更新
atl jira issue comment --key PROJ-789 --body "入力バリデーションを修正。PR #55 参照"
atl jira issue update --key PROJ-789 --status "Done"
```

## 5. Jira 課題を todo にインポートする

Jira の課題を `/todo` にインポートして、ローカルでタスク管理する。

```bash
# プロジェクトの課題をインポート
atl jira issue list --jql "project = PROJ ORDER BY updated DESC" --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin

# スプリントの課題をインポート
atl jira sprint issues --sprint 100 --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin

# 自分の未完了課題だけをインポート
atl jira issue list --jql "assignee = currentUser() AND statusCategory not in (Done)" --json \
  | bash skills/import-jira-to-todo/scripts/import-helper.sh \
  | todo datasource import jira --stdin
```

詳細は `skills/import-jira-to-todo/references/import.md` を参照。

## 6. プロジェクトの状況を把握する

プロジェクト全体の状況を JQL で確認する。

```bash
# 未完了の課題を確認
atl jira issue list --jql "project = PROJ AND statusCategory not in (Done) ORDER BY priority DESC"

# 特定のスプリントで残っている課題
atl jira issue list --jql "sprint in openSprints() AND project = PROJ AND statusCategory not in (Done)"

# 最近更新された課題
atl jira issue list --jql "project = PROJ AND updated >= -7d ORDER BY updated DESC"
```
