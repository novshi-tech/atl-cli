# 典型的なワークフロー

## 1. リポジトリの状況を把握する

```bash
# ワークスペース内のリポジトリを一覧
atl bitbucket repo list

# 特定のリポジトリの詳細を確認
atl bitbucket repo get --repo my-app
```

## 2. 機能ブランチからPRを作成する

```bash
# ブランチをプッシュ
git push -u origin feature/user-auth

# PRを作成
atl bitbucket pr create --repo my-app --title "ユーザー認証の実装" --source feature/user-auth -d "OAuth2 を使った認証機能を追加"
```

## 3. PRのレビューコメントを確認する

```bash
# コメントを取得（インラインコメントもデフォルトで含まれる、解決済みは除外）
atl bitbucket pr comment --repo my-app --pr 42

# インラインコメントを除外して取得
atl bitbucket pr comment --repo my-app --pr 42 --inline=false

# 解決済みコメントも含めて取得
atl bitbucket pr comment --repo my-app --pr 42 --include-resolved
```

## 4. PRにレビューコメントを投稿する

```bash
# 一般コメントを投稿
atl bitbucket pr comment create --repo my-app --pr 42 --body "LGTM! マージしてください。"

# インラインコメントを投稿（特定ファイルの特定行に対して）
atl bitbucket pr comment create --repo my-app --pr 42 --body "この変数名はもう少し分かりやすくした方が良いです。" --path src/auth.go --line 42
```

## 5. Jira 課題と連携して PR を作成する

```bash
# Jira 課題の詳細を確認
atl jira issue view --key PROJ-123

# ブランチで作業後、PRを作成
atl bitbucket pr create --repo my-app --title "PROJ-123: ログイン画面のバグ修正" --source fix/PROJ-123-login-bug

# Jira にコメントでPR URLを記録
atl jira issue comment --key PROJ-123 --body "PR を作成しました: https://bitbucket.org/myteam/my-app/pull-requests/55"
```

## 6. 特定ブランチへの PR を作成する

```bash
atl bitbucket pr create --repo my-app --title "ホットフィックス" --source hotfix/critical-fix --dest develop -d "本番環境の緊急修正"
```
