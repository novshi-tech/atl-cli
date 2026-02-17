# 典型的なワークフロー

## 1. リポジトリの状況を把握する

ワークスペースのリポジトリを確認する。

```bash
# ワークスペース内のリポジトリを一覧
atl bitbucket repo list --workspace myteam

# 特定のリポジトリの詳細を確認
atl bitbucket repo get --workspace myteam --repo my-app
```

## 2. 機能ブランチからPRを作成する

ローカルで作業したブランチからプルリクエストを作成する。

```bash
# ブランチをプッシュ
git push -u origin feature/user-auth

# PRを作成
atl bitbucket pr create --workspace myteam --repo my-app --title "ユーザー認証の実装" --source feature/user-auth -d "OAuth2 を使った認証機能を追加"
```

## 3. PRのレビューコメントを確認する

プルリクエストに付いたコメントを確認する。

```bash
# PRのコメントを取得
atl bitbucket pr comment --workspace myteam --repo my-app --pr 42
```

## 4. Jira 課題と連携して PR を作成する

Jira で課題を確認し、対応するブランチで作業してPRを作成する。

```bash
# Jira 課題の詳細を確認
atl jira issue view --key PROJ-123

# ブランチで作業後、PRを作成
atl bitbucket pr create --workspace myteam --repo my-app --title "PROJ-123: ログイン画面のバグ修正" --source fix/PROJ-123-login-bug

# Jira にコメントでPR URLを記録
atl jira issue comment --key PROJ-123 --body "PR を作成しました: https://bitbucket.org/myteam/my-app/pull-requests/55"
```

## 5. 特定ブランチへの PR を作成する

develop ブランチなど、メインブランチ以外をターゲットにPRを作成する。

```bash
atl bitbucket pr create --workspace myteam --repo my-app --title "ホットフィックス" --source hotfix/critical-fix --dest develop -d "本番環境の緊急修正"
```
