# TODOアプリ（フェーズ1: 一覧表示のみ）

React + TypeScript（フロントエンド）と Go（API）で構成した TODO アプリの初期実装です。

## このフェーズのスコープ
- TODO 一覧取得 API（`GET /api/todos`）の実装
- DB 接続と一覧画面向けの最小スキーマ・データ準備
- フロントエンドでの一覧取得・表示
- ユニットテスト追加（Go + フロントエンド）
- 一覧表示確認の Playwright E2E 追加

## このフェーズで扱わないこと
- TODO の追加
- TODO 完了状態の切り替え
- TODO の削除

## DB選定
このフェーズでは SQLite を採用しています。

理由:
- 外部サービスが不要で、ローカル環境のセットアップが簡単
- 現時点の要件（読み取り専用の一覧表示）には十分
- データモデルをシンプルに保ちながら次フェーズへ拡張しやすい

## プロジェクト構成
- `backend/`: Go API サーバー
- `frontend/`: React + TypeScript アプリと Playwright テスト

## バックエンドのセットアップと起動
```bash
cd backend
go mod tidy
go run ./cmd/server -addr :8080 -db ./todo.db
```

API エンドポイント:
- `GET http://localhost:8080/api/todos`

## フロントエンドのセットアップと起動
```bash
cd frontend
npm install
VITE_API_BASE_URL=http://localhost:8080 npm run dev
```

`http://localhost:5173` を開いてください。

## ユニットテスト
バックエンド:
```bash
cd backend
go test ./...
```

フロントエンド:
```bash
cd frontend
npm test
```

## E2E（Playwright）
```bash
cd frontend
npx playwright install
npm run test:e2e
```

`test:e2e` 実行時は、Playwright の `webServer` 設定によりバックエンド/フロントエンドが自動起動します。
