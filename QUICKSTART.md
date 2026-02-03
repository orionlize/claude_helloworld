# 快速启动

## 使用 Docker Compose

```bash
docker-compose up --build
```

- 前端地址：http://localhost:3000
- 后端 API：http://localhost:8080/api/health

## 本地开发（可选）

### 后端

```bash
cd backend
go run main.go
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

前端默认使用 `VITE_API_BASE` 环境变量指向 API 地址，未设置时会使用 `/api`。
