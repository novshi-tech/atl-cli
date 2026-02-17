# atl

Atlassian Cloud (Jira / Bitbucket) 向け CLI ツール。

Jira の課題管理や Bitbucket のリポジトリ・PR 操作をターミナルから実行できます。`--json` フラグによる機械可読出力に対応しており、AI エディタ（Claude Code / Codex / Cursor）向けのスキル定義も同梱しています。

## 機能

**Jira Cloud**
- 課題の作成・閲覧・更新・検索（JQL）
- スプリント一覧・スプリント内課題の表示
- ユーザー検索・課題へのアサイン
- コメントの追加

**Bitbucket Cloud**
- リポジトリ一覧・詳細の取得
- プルリクエストの作成・コメント取得

**その他**
- `--json` フラグで全コマンドの出力を JSON 形式に変換
- AI エディタ向けスキル定義の同梱と `install-skills` コマンドによる展開

## インストール

```bash
go install github.com/novshi-tech/atl-cli/cmd/atl@latest
```

## 初期セットアップ

### 1. サイトの設定

```bash
atl configure --site mysite --default
```

対話形式で以下のフィールドを入力します（既存値がある場合は Enter でスキップ可能）。

| フィールド | 説明 |
|---|---|
| Site URL | Atlassian サイトの URL（例: `https://myteam.atlassian.net`） |
| Email | Atlassian アカウントのメールアドレス |
| API Token | Jira 用 API トークン |
| Bitbucket API Token | Bitbucket 用 API トークン（省略可） |

### 2. API トークンの取得

- **Jira**: https://id.atlassian.com/manage-profile/security/api-tokens
- **Bitbucket**: https://bitbucket.org/account/settings/app-passwords/

### 3. 認証情報の保存先

認証情報は OS のキーリング（macOS Keychain / GNOME Keyring 等）に保存されます。キーリングが利用できない環境では `pass(1)` にフォールバックします。

## 複数サイトの管理

設定済みサイトの一覧を確認:

```bash
atl sites
```

コマンド実行時に `--site` フラグでサイトを切り替えられます:

```bash
atl jira issue list --site other-site --jql "project = DEMO"
```

## スキルのインストール

AI エディタ向けのスキル定義をローカルに展開します:

```bash
atl install-skills
```

| オプション | 説明 |
|---|---|
| `--list` | 利用可能なスキルの一覧を表示 |
| `--dry-run` | 書き込みせずにプレビュー |
| `--path <dir>` | 展開先ディレクトリを指定 |

## 基本的な使い方

### Jira

```bash
# JQL で課題を検索
atl jira issue list --jql "project = MYPROJ AND status = 'In Progress'"

# 課題の詳細を表示
atl jira issue view --key MYPROJ-123

# 課題を作成
atl jira issue create -p MYPROJ -s "バグ修正" -t Bug -d "ログイン画面のエラーを修正"

# 課題を更新
atl jira issue update --key MYPROJ-123 --status Done --assignee none

# コメントを追加
atl jira issue comment --key MYPROJ-123 -b "対応完了しました"

# スプリント一覧
atl jira sprint list --board 1 --state active

# スプリント内の課題
atl jira sprint issues --sprint 10

# ユーザー検索
atl jira user search -q "tanaka"
```

### Bitbucket

```bash
# リポジトリ一覧
atl bitbucket repo list --workspace myws

# リポジトリ詳細
atl bitbucket repo get --workspace myws --repo myrepo

# プルリクエスト作成
atl bitbucket pr create --workspace myws --repo myrepo --title "機能追加" --source feature-branch

# PR のコメント取得
atl bitbucket pr comment --workspace myws --repo myrepo --pr 42
```

### JSON 出力

全コマンドで `--json` フラグが使用可能です:

```bash
atl jira issue list --jql "project = MYPROJ" --json
```

## ライセンス

MIT License — 詳細は [LICENSE](LICENSE) を参照。
