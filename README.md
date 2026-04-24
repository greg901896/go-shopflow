# go-shopflow

Go 電商 API，練習分層架構、PostgreSQL、JWT 認證，並與 [go-task-queue](https://github.com/greg901896/go-task-queue) 串接展示微服務生態系。

## Tech Stack

Go 1.26 / Gin / PostgreSQL 16 / jackc/pgx v5 / Docker Compose

## Architecture

分層架構（Rails-inspired）：

```
cmd/server/main.go       # 進入點、DI 裝配
internal/
  ├── model/             # 資料結構
  ├── repository/        # DB 存取（SQL）
  ├── service/           # 業務邏輯
  └── handler/           # HTTP handler（Gin）
migrations/              # SQL migrations
```

## Quick Start

```bash
# 1. 啟動 Postgres
docker compose up -d

# 2. 跑 migration
migrate -path migrations \
  -database "postgres://shopflow:shopflow_dev@localhost:5433/shopflow?sslmode=disable" \
  up

# 3. 啟動 server（port 8081）
go run cmd/server/main.go
```

## API

### Products

| Method | Path | 說明 |
|---|---|---|
| `POST` | `/products` | 新增商品 |
| `GET` | `/products?page=1&limit=20` | 商品列表（分頁） |
| `GET` | `/products/:id` | 查單一商品 |
| `PUT` | `/products/:id` | 更新商品（partial update） |

## Design Notes

- 金額用 `NUMERIC(10,2)` 儲存並以 string 在 API 傳遞
- 外鍵的 `ON DELETE` 依資料性質分流：
  - 歷史紀錄（`orders`、`order_items.product_id`）→ `RESTRICT`
  - 附屬資料（`carts`、`cart_items`）→ `CASCADE`
- `PUT` 採 partial update 語意：用指標區分「未傳」與「零值」，SQL 以 `COALESCE` 保留原值
- Postgres 外鍵手動加索引

## Development Status

- ✅ Phase 1：專案初始化、Docker Compose、DB migrations
- ✅ Phase 2：商品模組（CRUD）
- ⬜ Phase 3：會員認證（JWT + bcrypt）
- ⬜ Phase 4：購物車
- ⬜ Phase 5：訂單 + 串接 go-task-queue
- ⬜ Phase 6：收尾（logging、tests、env）
