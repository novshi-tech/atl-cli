---
name: update-check
description: atl CLI の新しいバージョンが利用可能かチェックし、ユーザーにアップデートを提案するスキル。atl コマンド実行時に「新しいバージョンが利用可能です」というメッセージが stderr に表示された場合に、ユーザーにアップデートするか確認する。
---

# Update Check

atl CLI が新しいバージョンの存在を通知した場合に、ユーザーにアップデートを提案する。

## When to Use

atl コマンドの実行結果の stderr に以下のようなメッセージが含まれている場合にのみ発動する:

```
atl: new version available (vX.Y.Z → vA.B.C)
atl: run `go install github.com/novshi-tech/atl-cli/cmd/atl@latest` to update.
```

**このスキルは proactive にバージョンチェックを行わない。** atl CLI 自体が組み込みのバージョンチェック機能（1日1回、Go module proxy に問い合わせ）を持っており、新しいバージョンがあれば上記メッセージを自動で出力する。

## Procedure

### Step 1: メッセージを検知

atl コマンドの出力に `new version available` が含まれていることを確認する。

### Step 2: ユーザーに確認

> atl CLI の新しいバージョンが利用可能です。アップデートしますか?

ユーザーが同意した場合のみ Step 3 に進む。

### Step 3: アップデートを実行

```bash
go install github.com/novshi-tech/atl-cli/cmd/atl@latest
```

### Step 4: アップデート結果を確認

```bash
go version -m $(which atl) 2>/dev/null | grep mod
```

バージョンが更新されたことを確認し、ユーザーに報告する。

## Notes

- ユーザーが拒否した場合は何もせず、メインの操作を続行する
- ネットワークエラーなどで `go install` が失敗した場合はエラーを報告し、メインの操作を続行する
