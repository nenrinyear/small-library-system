# small-library-system

Go (Echo) + Next.js (Chakra UI v3) で構成した QR ベースの所蔵物管理システムです。

## 構成

- Frontend: `next/` (Next.js 14 + Chakra UI v3)
- Backend: `go-api/` (Go + Echo + GORM + MySQL)
- DB schema/init: `go-api/migrations/`, `mysql/init.d/`

## 開発起動

1. `.env.example` を `.env` にコピーして値を調整
2. `docker compose -f compose.dev.yml up --build`

公開ポート:

- Frontend: `http://localhost:3000`
- API: `http://localhost:8080`
- MySQL: `localhost:3306`

## 主要環境変数

- `NEXT_PUBLIC_API_BASE_URL`
- `ADMIN_API_TOKEN`
- `MYSQL_DATABASE`, `MYSQL_USER`, `MYSQL_PASSWORD`

