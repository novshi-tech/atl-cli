# atl bitbucket コマンドリファレンス

## bitbucket me

認証済みの Bitbucket ユーザー（自分自身）のアカウント情報を表示する。

```
atl bitbucket me [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Account ID:    5b10ac8d14c052e1e6c2e251
Display Name:  John Doe
Nickname:      johndoe
UUID:          {1234-5678-abcd}
Created On:    2020-01-15T10:30:00.000000+00:00
```

**JSON 出力例** (`--json`):
```json
{
  "accountId": "5b10ac8d14c052e1e6c2e251",
  "displayName": "John Doe",
  "nickname": "johndoe",
  "uuid": "{1234-5678-abcd}",
  "createdOn": "2020-01-15T10:30:00.000000+00:00"
}
```

## bitbucket repo list

ワークスペース内のリポジトリを一覧表示する。

```
atl bitbucket repo list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | No | サイト設定値 | ワークスペースのスラッグ（サイト設定で制限） |
| `--max` | - | No | `25` | 最大取得件数 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 3 repositor(ies):

my-app                          My App                Go          private
frontend                        Frontend              TypeScript  public
infra                           Infrastructure        HCL         private
```

**JSON 出力例** (`--json`):
```json
[
  {
    "slug": "my-app",
    "name": "My App",
    "language": "Go",
    "is_private": true
  }
]
```

## bitbucket repo get

リポジトリの詳細情報を表示する。

```
atl bitbucket repo get [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | No | サイト設定値 | ワークスペースのスラッグ（サイト設定で制限） |
| `--repo` | - | Yes | - | リポジトリのスラッグ |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Slug:         my-app
Name:         My App
Full Name:    myteam/my-app
Description:  メインのアプリケーション
Language:     Go
Private:      Yes
Main Branch:  main
Updated:      2024-06-15T10:30:00.000000+00:00
```

**JSON 出力例** (`--json`):
```json
{
  "slug": "my-app",
  "name": "My App",
  "full_name": "myteam/my-app",
  "description": "メインのアプリケーション",
  "language": "Go",
  "is_private": true,
  "mainbranch": "main",
  "updated_on": "2024-06-15T10:30:00.000000+00:00"
}
```

## bitbucket pr list

リポジトリのプルリクエストを一覧表示する。

```
atl bitbucket pr list [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | No | サイト設定値 | ワークスペースのスラッグ（サイト設定で制限） |
| `--repo` | - | Yes | - | リポジトリのスラッグ |
| `--state` | - | No | `OPEN` | 状態フィルタ: `OPEN` / `MERGED` / `DECLINED` / `SUPERSEDED` |
| `--max` | - | No | `25` | 最大取得件数 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 2 pull request(s):

#42      OPEN        Alice               feature/auth→main       認証機能を追加
#41      OPEN        Bob                 fix/login-bug→main      ログインバグを修正
```

**JSON 出力例** (`--json`):
```json
[
  {
    "id": 42,
    "title": "認証機能を追加",
    "state": "OPEN",
    "author": "Alice",
    "source": "feature/auth",
    "dest": "main"
  }
]
```

## bitbucket pr create

新しいプルリクエストを作成する。

```
atl bitbucket pr create [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | No | サイト設定値 | ワークスペースのスラッグ（サイト設定で制限） |
| `--repo` | - | Yes | - | リポジトリのスラッグ |
| `--title` | - | Yes | - | プルリクエストのタイトル |
| `--source` | - | Yes | - | ソースブランチ名 |
| `--dest` | - | No | リポジトリのメインブランチ | デスティネーションブランチ |
| `--description` | `-d` | No | - | プルリクエストの説明 |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Created pull request: #42
URL: https://bitbucket.org/myteam/my-app/pull-requests/42
```

**JSON 出力例** (`--json`):
```json
{
  "key": "42",
  "url": "https://bitbucket.org/myteam/my-app/pull-requests/42"
}
```

## bitbucket pr comment

プルリクエストのコメントを一覧表示する。インラインコードレビューコメントはデフォルトで含まれる。解決済みコメントはデフォルトで除外される。

```
atl bitbucket pr comment [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | No | サイト設定値 | ワークスペースのスラッグ（サイト設定で制限） |
| `--repo` | - | Yes | - | リポジトリのスラッグ |
| `--pr` | - | Yes | - | プルリクエスト ID |
| `--inline` | - | No | `true` | インラインコードレビューコメントも含める |
| `--include-resolved` | - | No | `false` | 解決済みコメントも含める（デフォルトでは除外） |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 2 comment(s):

[#101][2024-06-15T10:30:00.000000+00:00] John Doe:
LGTM! マージしてください。

[#102 reply to #101][2024-06-15T11:00:00.000000+00:00] Jane Smith:
修正を確認しました。

Found 1 inline comment(s):

[#201][2024-06-15T09:00:00.000000+00:00] Alice on src/auth.go (lines 42-45):
この変数名はもう少し分かりやすくした方が良いです。
```

コメントヘッダの `#<id>` はコメント ID、`reply to #<id>` は返信先コメント ID。`bitbucket pr comment create --parent` でこの ID を使う。

**出力例** (`--inline=false`):
```
Found 2 comment(s):

[#101][2024-06-15T10:30:00.000000+00:00] John Doe:
LGTM! マージしてください。

[#102 reply to #101][2024-06-15T11:00:00.000000+00:00] Jane Smith:
修正を確認しました。
```

**JSON 出力例** (`--json`):
```json
{
  "comments": [
    {
      "id": 101,
      "author": "John Doe",
      "created": "2024-06-15T10:30:00.000000+00:00",
      "body": "LGTM! マージしてください。"
    },
    {
      "id": 102,
      "parent_id": 101,
      "author": "Jane Smith",
      "created": "2024-06-15T11:00:00.000000+00:00",
      "body": "修正を確認しました。"
    }
  ],
  "inline_comments": []
}
```

**JSON 出力例** (`--json --inline`):
```json
{
  "comments": [
    {
      "id": 101,
      "author": "John Doe",
      "created": "2024-06-15T10:30:00.000000+00:00",
      "body": "LGTM! マージしてください。"
    }
  ],
  "inline_comments": [
    {
      "id": 201,
      "author": "Alice",
      "created": "2024-06-15T09:00:00.000000+00:00",
      "path": "src/auth.go",
      "from": 42,
      "to": 45,
      "body": "この変数名はもう少し分かりやすくした方が良いです。"
    },
    {
      "id": 202,
      "parent_id": 201,
      "author": "Bob",
      "created": "2024-06-15T09:30:00.000000+00:00",
      "path": "src/auth.go",
      "from": 42,
      "to": 45,
      "body": "指摘の通り修正します。"
    }
  ]
}
```

## bitbucket pr comment create

プルリクエストにコメントを投稿する。`--path` と `--line` を指定するとインラインコメントになる。

```
atl bitbucket pr comment create [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | No | サイト設定値 | ワークスペースのスラッグ |
| `--repo` | - | Yes | - | リポジトリのスラッグ |
| `--pr` | - | Yes | - | プルリクエスト ID |
| `--body` | `-b` | Yes | - | コメント本文 |
| `--path` | - | No | - | インラインコメントのファイルパス |
| `--line` | - | No | - | インラインコメントの行番号 |
| `--parent` | - | No | - | 返信先のコメント ID（コメントへの返信に使用） |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Comment added to pull request #42
URL: https://bitbucket.org/myteam/my-app/pull-requests/42
```

**JSON 出力例** (`--json`):
```json
{
  "key": "1",
  "url": "https://bitbucket.org/myteam/my-app/pull-requests/42"
}
```

**インラインコメントへの返信**:

返信先のコメント ID は `bitbucket pr comment` の出力（`#<id>` あるいは JSON の `id`）から取得する。インライン返信では返信先と同じ `--path` / `--line` を指定する。

```bash
# 1) まずコメント一覧から返信先インラインコメントの id を取得
atl bitbucket pr comment --repo my-app --pr 42 --json

# 2) 同じ path / line と --parent でインライン返信
atl bitbucket pr comment create --repo my-app --pr 42 \
    --path src/auth.go --line 45 --parent 201 \
    --body "指摘の通り修正しました"

# 通常コメント（非インライン）への返信は --path / --line を省略
atl bitbucket pr comment create --repo my-app --pr 42 \
    --parent 101 --body "確認しました"
```
