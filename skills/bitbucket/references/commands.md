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
| `--workspace` | - | Yes | - | ワークスペースのスラッグ |
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
| `--workspace` | - | Yes | - | ワークスペースのスラッグ |
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
| `--workspace` | - | Yes | - | ワークスペースのスラッグ |
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
| `--workspace` | - | Yes | - | ワークスペースのスラッグ |
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

プルリクエストのコメントを一覧表示する。インラインコードレビューコメントは除外される。

```
atl bitbucket pr comment [flags]
```

| フラグ | 短縮 | 必須 | デフォルト | 説明 |
|--------|------|------|-----------|------|
| `--workspace` | - | Yes | - | ワークスペースのスラッグ |
| `--repo` | - | Yes | - | リポジトリのスラッグ |
| `--pr` | - | Yes | - | プルリクエスト ID |
| `--site` | - | No | デフォルトサイト | サイトエイリアス |
| `--json` | - | No | `false` | JSON 形式で出力 |

**出力例:**
```
Found 2 comment(s):

[2024-06-15T10:30:00.000000+00:00] John Doe:
LGTM! マージしてください。

[2024-06-15T11:00:00.000000+00:00] Jane Smith:
修正を確認しました。
```

**JSON 出力例** (`--json`):
```json
[
  {
    "author": "John Doe",
    "created": "2024-06-15T10:30:00.000000+00:00",
    "body": "LGTM! マージしてください。"
  }
]
```
