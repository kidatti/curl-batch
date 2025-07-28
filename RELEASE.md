# Release Guide

curl-batchプロジェクトのリリース手順について説明します。

## リリース手順

### 1. 準備

リリース前に以下を確認してください：

- すべてのテストが通ること
- コードの品質チェック（lint、format）が通ること
- ドキュメントが最新であること

```bash
# テスト実行
make test

# コード品質チェック
make fmt
make lint
```

### 2. CHANGELOG.mdの更新

リリース内容をCHANGELOG.mdに記録します：

1. `[Unreleased]`セクションの内容を確認
2. 新しいバージョンセクションを作成
3. 変更内容を適切なカテゴリに分類：
   - `Added` - 新機能
   - `Changed` - 既存機能の変更
   - `Deprecated` - 非推奨機能
   - `Removed` - 削除された機能
   - `Fixed` - バグ修正
   - `Security` - セキュリティ関連

例：
```markdown
## [v1.0.0] - 2024-01-15

### Added
- 新機能の説明

### Fixed
- 修正されたバグの説明
```

### 3. バージョンタグの作成とプッシュ

#### 方法1: Makefileを使用（推奨）

```bash
# バージョンタグを作成してプッシュ
make tag VERSION=v1.0.0
```

#### 方法2: 手動でgitコマンド実行

```bash
# バージョンタグを作成
git tag v1.0.0

# タグをリモートにプッシュ
git push origin v1.0.0
```

### 4. 自動リリース処理

タグがプッシュされると、GitHub Actionsが自動的に以下を実行します：

1. **マルチプラットフォームビルド**
   - Linux (AMD64, ARM64)
   - macOS (AMD64, ARM64)
   - Windows (AMD64)

2. **リリースの作成**
   - GitHub Releasesページに新しいリリースを作成
   - ビルド済みバイナリを自動アップロード

3. **アーティファクトの配布**
   - Linux/macOS: tar.gz形式
   - Windows: zip形式

### 5. リリース後の確認

リリース完了後に以下を確認してください：

1. **GitHub Releasesページ**
   - https://github.com/[ユーザー名]/curl-batch/releases
   - 新しいリリースが作成されていること
   - すべてのプラットフォーム向けバイナリが添付されていること

2. **ダウンロードテスト**
   - 各プラットフォームのバイナリをダウンロード
   - 基本動作確認

## 対応プラットフォーム

### リリースされるバイナリ

| プラットフォーム | アーキテクチャ | ファイル形式 | ファイル名例 |
|------------------|----------------|--------------|--------------|
| Linux | AMD64 | tar.gz | curl-batch-linux-amd64.tar.gz |
| Linux | ARM64 | tar.gz | curl-batch-linux-arm64.tar.gz |
| macOS | AMD64 | tar.gz | curl-batch-darwin-amd64.tar.gz |
| macOS | ARM64 | tar.gz | curl-batch-darwin-arm64.tar.gz |
| Windows | AMD64 | zip | curl-batch-windows-amd64.zip |

### 手動ビルド

ローカルで特定プラットフォーム向けにビルドする場合：

```bash
# 現在のプラットフォーム向け
make build

# 全プラットフォーム向け
make build-all

# 特定プラットフォーム向け
make build-linux    # Linux (AMD64, ARM64)
make build-darwin   # macOS (AMD64, ARM64)  
make build-windows  # Windows (AMD64)

# リリースパッケージ作成
make package
```

## バージョニング

このプロジェクトは[Semantic Versioning](https://semver.org/)に従います：

- **MAJOR** (例: v1.0.0 → v2.0.0): 互換性のない変更
- **MINOR** (例: v1.0.0 → v1.1.0): 後方互換性を保った機能追加
- **PATCH** (例: v1.0.0 → v1.0.1): 後方互換性を保ったバグ修正

## トラブルシューティング

### GitHub Actionsが失敗する場合

1. **ビルドエラー**
   - ローカルで`make build-all`を実行してエラーを確認
   - 依存関係やソースコードの問題を修正

2. **権限エラー**
   - リポジトリの設定でGitHub Actionsが有効になっていることを確認
   - `GITHUB_TOKEN`の権限設定を確認

3. **リリース作成エラー**
   - 同じタグ名のリリースが既に存在していないか確認
   - タグ名の形式が正しいか確認（例: v1.0.0）

### 緊急時のリリース削除

```bash
# ローカルタグ削除
git tag -d v1.0.0

# リモートタグ削除
git push origin --delete v1.0.0
```

その後、GitHub ReleasesページからリリースとAssetを手動削除してください。

## 参考資料

- [GitHub Actions Workflow](.github/workflows/release.yml)
- [Changelog](CHANGELOG.md)
- [Makefile](Makefile)