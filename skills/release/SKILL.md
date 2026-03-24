---
name: release
description: atl CLI の新しいバージョンをリリースするスキル。セマンティックバージョニングに従い、Git タグを作成・プッシュして Go module proxy に公開する。「新しいバージョンをリリースして」「バージョンを上げて」「リリースして」など、バージョンリリースを依頼された場合に使用。
---

# Release

atl CLI の新しいバージョンをリリースする。Go module proxy (`proxy.golang.org`) はGit タグを元にバージョンを配信するため、リリース = タグの作成・プッシュとなる。

## Prerequisites

- `main` ブランチに最新の変更がマージされていること
- ワーキングツリーがクリーンであること

## Procedure

### Step 1: 現在の状態を確認

```bash
# 現在のブランチを確認（main であるべき）
git branch --show-current

# ワーキングツリーがクリーンか確認
git status --short

# リモートと同期しているか確認
git fetch origin && git log --oneline origin/main..HEAD
```

**main ブランチでない場合**、または未プッシュのコミットがある場合はユーザーに確認する。

### Step 2: 既存のタグを確認して次のバージョンを決定

```bash
# 既存タグを降順で表示
git tag --sort=-v:refname | head -10
```

セマンティックバージョニング (`vMAJOR.MINOR.PATCH`) に従い、次のバージョンを決定する:

- **patch** (v0.1.0 -> v0.1.1): バグ修正、ドキュメント修正
- **minor** (v0.1.0 -> v0.2.0): 新機能追加、後方互換性のある変更
- **major** (v0.1.0 -> v1.0.0): 破壊的変更

タグがまだ存在しない場合は `v0.1.0` から開始する。

前回タグからのコミットログを確認して適切なバージョンを判断する:

```bash
# 前回タグからの変更を確認（タグがない場合は全コミット）
git log --oneline $(git describe --tags --abbrev=0 2>/dev/null || git rev-list --max-parents=0 HEAD)..HEAD
```

### Step 3: ユーザーに確認

次のバージョン番号をユーザーに提案し、承認を得る。例:

> 前回リリース以降の変更内容:
> - feat: xxx
> - fix: yyy
>
> 次のバージョンとして `vX.Y.Z` を提案します。よろしいですか?

**ユーザーの承認なしにタグを作成・プッシュしてはならない。**

### Step 4: タグを作成してプッシュ

```bash
# アノテーションタグを作成（メッセージにはリリース内容のサマリーを含める）
git tag -a vX.Y.Z -m "Release vX.Y.Z: 変更内容のサマリー"

# タグをリモートにプッシュ
git push origin vX.Y.Z
```

### Step 5: Go module proxy への反映を確認

```bash
# proxy に反映されるまで少し待つ場合がある
curl -s https://proxy.golang.org/github.com/novshi-tech/atl-cli/@latest | jq .
```

### Step 6: ユーザーに完了を報告

> `vX.Y.Z` をリリースしました。
> ユーザーは以下のコマンドでアップデートできます:
>
> ```bash
> go install github.com/novshi-tech/atl-cli/cmd/atl@latest
> ```

## Important Notes

- **必ずユーザーの承認を得てからタグを作成・プッシュすること** (リモートへの変更は取り消しが難しい)
- `--force` でのタグ上書きは禁止。既存タグを修正する必要がある場合はユーザーに確認する
- v0.x.x の間は破壊的変更を minor で行ってもよい（Go の慣習）
