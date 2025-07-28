# curl-batch

CSVデータとcurlテンプレートを使用して、複数のcurlリクエストを一括実行するツールです。

## 機能

- CSVファイルからパラメータを読み込み
- curlテンプレート内の変数を置換
- HTTPリクエストを一括実行
- 結果をテキストファイルに保存
- 様々なHTTPメソッドとヘッダーをサポート
- リクエスト間のスリープ時間を指定可能（ミリ秒単位）

## インストール

```bash
git clone https://github.com/kidatti/curl-batch.git
cd curl-batch
go build -o curl-batch main.go
```

## 使用方法

```bash
./curl-batch -curl <curlテンプレートファイル> -csv <CSVファイル> -output <出力ファイル> [オプション]
```

### 引数とオプション

| フラグ | 説明 | 必須 | デフォルト値 |
|--------|------|------|-------------|
| `-curl` | curlテンプレートファイル | Yes | - |
| `-csv` | CSVデータファイル | Yes | - |
| `-output` | 出力ファイル | Yes | - |
| `-sleep` | リクエスト間のスリープ時間（ミリ秒） | No | 0 |

### 使用例

1. curlテンプレートファイル (`curl.txt`) を作成:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "${NAME}", "email": "${EMAIL}", "age": ${AGE}}' https://hogehoge.com/api/users
```

2. CSVファイル (`users.csv`) を作成:
```csv
NAME,EMAIL,AGE
田中太郎,tanaka@example.com,30
佐藤花子,sato@example.com,25
鈴木一郎,suzuki@example.com,35
```

3. バッチ実行:
```bash
# 基本実行
./curl-batch -curl sample/curl.txt -csv sample/users.csv -output results.txt

# リクエスト間に1秒のスリープを挿入
./curl-batch -curl sample/curl.txt -csv sample/users.csv -output results.txt -sleep 1000
```

## テンプレート変数

curlテンプレート内の変数は `${変数名}` の形式で記述し、CSVファイルの列ヘッダーと一致させる必要があります。

## 出力形式

このツールは以下の内容を含む詳細な出力ファイルを生成します:
- リクエスト番号
- 実行されたcurlコマンド
- そのリクエストで使用されたCSVデータ
- HTTPレスポンスのステータス、ヘッダー、ボディ
- 発生したエラー（もしあれば）

## ビルドとインストール

### Makefileを使用する場合
```bash
# ビルド
make build

# サンプルデータで実行
make run-sample

# 全プラットフォーム向けビルド
make build-all

# ヘルプを表示
make help
```

### 手動ビルド
```bash
go build -o curl-batch main.go
```

## リリース

リリース手順については [RELEASE.md](RELEASE.md) を参照してください。

### 対応プラットフォーム

GitHub Releasesで以下のプラットフォーム向けビルド済みバイナリをダウンロードできます：

- **Linux**: AMD64, ARM64 (tar.gz形式)
- **macOS**: AMD64, ARM64 (tar.gz形式)  
- **Windows**: AMD64 (zip形式)

## 必要要件

- Go 1.21以上
