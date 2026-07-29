# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

Atlassian Cloud (Jira / Bitbucket) 向けの Go 製 CLI ツール。AI エディタ（Claude Code / Cursor）との統合を意識した設計で、`--json` フラグによる機械可読出力をサポート。

## よく使うコマンド

```bash
# ビルド
go build -o atl ./cmd/atl

# 型チェック
go vet ./...

# 依存関係の整理
go mod tidy

# インストール
go install ./cmd/atl
```

テストファイルは現在存在しない。

## アーキテクチャ

### レイヤー構成

```
cmd/               ← CLI インターフェース層（Cobra コマンド実装）
internal/
  auth/            ← 認証情報管理（OS Keyring / pass フォールバック）
  jira/            ← Jira REST API v3 クライアント
  bitbucket/       ← Bitbucket Cloud API クライアント
  adf/             ← Atlassian Document Format パーサー（Markdown → ADF 変換）
skills/            ← Claude Code / Cursor 向けスキル定義（go embed でバイナリに同梱）
```

### エントリーポイント

`cmd/atl/main.go` → `cmd.Execute()`（Cobra ルートコマンド）

### データフロー

コマンド実行 → `newJiraClient(cmd)` / `newBitbucketClient(cmd)` で認証情報取得 → `internal/` の API クライアントで HTTP リクエスト → テーブルまたは JSON 出力

### 認証システム

- `internal/auth/auth.go`: サイトエイリアス単位で認証情報を管理
- OS Keyring（macOS Keychain / GNOME Keyring / Windows Credential Manager）→ `pass` コマンドへフォールバック
- `ATL_CRED_BACKEND=boid-cli` を明示指定した場合のみ `BoidCLIStore`（`internal/auth/boid_cli.go`）を使う。`boid secret get/set/delete --namespace atl` を exec するだけの薄い実装で、boid の host_command として daemon container 内で atl を動かす用途向け（daemon 自身の UNIX socket は同一コンテナ内からは無認証の trusted transport なので、追加の認証情報は不要）。keyring も pass も使えない環境（GUI セッションも GPG 鍵も無い headless container）向けの第三の選択肢。auto-detect ではなく明示 opt-in にしているのは、`boid` が PATH にあるだけの普通のワークステーションで誤って選ばれないようにするため
- 保存項目: `base-url`, `email`, `api-token`, `bb-api-token`（Bitbucket 専用）
- デフォルトサイトは `default-site` キーで管理
- `atl configure --site <alias>` で初期設定

### JSON 出力型

`cmd/json_types.go` に全コマンドの JSON 出力型を一元定義。`--json` フラグ指定時はこれらの型にマーシャルして出力。

### ADF パーサー

`internal/adf/adf.go`: Jira の課題作成・更新・コメント時に使用するテキスト → ADF 変換ロジック。サポート構文: 見出し（`## `）、箇条書き（`*`/`-`）、テーブル（`|`/`||`）、太字（`**`）、リンク（`[text](url)`）、ユーザーメンション（`@[名前:accountId]`）。

## 新コマンドの追加手順

1. `cmd/` に `<service>_<noun>_<verb>.go` の命名でファイルを作成
2. コマンド関数を定義し、`cmd/<service>.go` の `init()` で `AddCommand()` に追加
3. JSON 出力が必要な場合は `cmd/json_types.go` に型を追加
4. スキルのリファレンスを更新: `skills/<service>/references/commands.md`

## スキル管理

`skills/` ディレクトリのスキル定義は `skills.go` で `go:embed` によりバイナリに埋め込まれ、`atl install-skills` コマンドで Claude Code / Cursor の設定ディレクトリにインストールされる。

スキルを追加・更新した場合は `skills/<service>/SKILL.md` と `references/commands.md` を合わせて更新すること。
