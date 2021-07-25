# tomeit-api

tomeit-api は tomeit の REST API です.

## セットアップ

1. `.env` ファイルの作成

注意点：ユーザ認証には Firebase Authentication を利用しています. そのため動かすには Firebase Admin SDK が必要です.

```bash
PORT=8080
DSN=<user>:<password>@tcp(tomeit-api-db-dev:3306)/<database>?parseTime=true
GOOGLE_APPLICATION_CREDENTIALS=<Firebase Admin SDK filepath>
ALLOW_ORIGINS=http://localhost:3000

MYSQL_ROOT_PASSWORD=<root_password>
MYSQL_DATABASE=<database>
MYSQL_USER=<user>
MYSQL_PASSWORD=<password>
```

2. コマンドの実行

```bash
make dev
```

[TODO: GitHub Pages で Swagger UI をホスティングしたら追加する]: <> (## ドキュメント)

[comment]: <> (この API のエンドポイントは[こちら]&#40;GitHub PagesのURL&#41;に載っています.)

## 開発環境

- プログラミング言語：Go
- ライブラリ・フレームワーク：chi
- デプロイ：Heroku
- フォーマッタ・リンタ：gofmt, goimports, govet, staticcheck
- テスト：testing, net/http/httptest

### 自動整形

```bash
make fmt
```

### 静的解析

```bash
make lint
```

### テスト

```bash
make test
```

### ローカル実行

```bash
make dev 
# 終了時
make down
```
