# tomeit-api

tomeit-api は tomeit の Web API です.

## セットアップ

1. `.env` ファイルの作成

*注意点：tomeit-api のユーザ認証には Firebase Authentication を利用しています. そのため

```text
PORT=8080
DSN=minguu:password@tcp(tomeit-api-db-dev:3306)/tomeit_db_dev?parseTime=true
GOOGLE_APPLICATION_CREDENTIALS=./build/tomeit-dev-firebase-adminsdk.json
ALLOW_ORIGINS=http://localhost:3000,https://tomeit.vercel.app

MYSQL_ROOT_PASSWORD=root_password
MYSQL_DATABASE=tomeit_db_dev
MYSQL_USER=minguu
MYSQL_PASSWORD=password
```

2. コマンドの実行

## ドキュメント

この API のエンドポイントの詳細は[こちらのドキュメント](GitHub PagesのURL)に載っています.

## 開発用コマンド

### 自動整形

### 静的解析

### テスト

### ローカル実行
