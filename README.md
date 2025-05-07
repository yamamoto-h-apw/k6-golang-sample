# 🚀 Go + k6 負荷テスト環境（Docker Compose）

このプロジェクトは、Goで作成された HTTP API に対して [k6](https://k6.io/) を用いた負荷テストを行う、シンプルで再利用可能な環境を Docker Compose で提供します。

---

## 📁 プロジェクト構成

```
k6-golang-sample/
├── compose.yml
├── server/
│   ├── main.go               # Go HTTP サーバーコード（/health, /login, /private）
│   ├── go.mod
│   └── Dockerfile.server
└── k6/
    └── loadtest.js           # k6 スクリプト（JWT認証付き）
```

---

## 🧱 初回セットアップとテスト実行

以下のコマンドでサーバーのビルド・起動と、k6による1回目の負荷テストを行います：

```bash
docker compose up --build
```

- Goサーバーが http://localhost:8080 で起動します。
- `k6` コンテナが `/login → /private` の認証付きAPIテストを実行し、結果を表示します。

---

## 🔁 k6 のテストを再実行する方法

Goサーバーが起動している状態で、`k6` テストだけを再実行するには以下のコマンドを使います：

```bash
docker compose run --rm k6
```

- `--rm` により、実行後コンテナは自動で削除されます。
- `k6/loadtest.js` に定義されたテストが実行されます。

---

## 📡 API エンドポイント

| メソッド | パス     | 説明             |
|----------|----------|------------------|
| GET      | /health    | ヘルスチェック                       |
| POST     | /login     | ユーザー名とパスワードでJWTを取得     |
| GET      | /private   | 認証トークン必須の保護されたエンドポイント |

---

## 🔐 JWT認証付きAPIテストの流れ

`k6/loadtest.js` では以下の処理を自動実行しています：

1. `POST /login` にユーザー情報（例: `testuser`/`testpass`）を送信し、JWTを取得
2. `Authorization: Bearer <token>` ヘッダーを付けて `GET /private` を実行

すべてのレスポンスに対してステータスチェックを行い、テスト結果に反映されます。

---

## 📊 k6 出力の見方（実行結果の読み方）

以下は、k6 実行時に表示される出力の各項目の意味です。

| カテゴリ     | 指標                     | 意味・評価                                     |
|--------------|--------------------------|-----------------------------------------------|
| ✅ チェック    | `checks_succeeded`        | 200回中200回成功（100%） → テスト全通過         |
| 🔄 応答時間   | `http_req_duration`       | 平均 1.73ms（高速）／最大 12.69ms               |
| 📉 エラー率   | `http_req_failed`         | 0%（すべてのリクエストが成功）                 |
| 🔁 実行回数   | `iterations`              | 100回（10VU × 10秒）                           |
| 📦 通信量     | `data_sent / data_received` | 約 40KB（全体）                                |
| 👥 同時ユーザー | `vus`, `vus_max`           | 常時10VUで安定実行                              |

> 🟢 総合評価：全テスト成功・応答高速・エラーなし → 非常に良好なパフォーマンス

---

## 🧹 クリーンアップ

すべてのコンテナを停止・削除するには：

```bash
docker compose down
```

---

## 🛠 使用技術

- Golang 1.20
- k6（Grafana 製の負荷テストツール）
- JWT（認証処理）
- Docker / Docker Compose