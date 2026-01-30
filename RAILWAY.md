# Railway 部署指南

## 部署步骤

### 1. 准备工作

确保项目已推送到 GitHub 仓库。

### 2. 在 Railway 创建新项目

1. 访问 [railway.app](https://railway.app/)
2. 点击 "New Project" → "Deploy from GitHub repo"
3. 选择你的仓库

### 3. 配置环境变量

在 Railway 项目设置中添加以下环境变量：

```
ENVIRONMENT=production
DB_MODE=memory
SERVER_PORT=8081
JWT_SECRET=your-random-secret-key-here
```

**生成 JWT_SECRET:**
```bash
openssl rand -base64 32
```

### 4. 配置服务

Railway 会自动检测 `nixpacks.toml` 配置文件。确保：

- **Root Directory**: `/`
- **Build Command**: `cd backend && go build -o bin/api ./cmd/api`
- **Start Command**: `backend/bin/api`

### 5. 部署前端（可选）

如果需要单独部署前端：

1. 创建新的 Railway 服务
2. 选择 "Dockerfile" 构建
3. 设置 Root Directory 为 `frontend`
4. 端口设置为 `5173`

### 6. 获取部署 URL

部署完成后，Railway 会提供一个公网 URL，例如：
```
https://your-app.railway.app
```

## 本地测试

部署前可以本地测试：

```bash
# 使用 Railway 相同的环境变量
export ENVIRONMENT=production
export DB_MODE=memory
export SERVER_PORT=8081
export JWT_SECRET=test-secret

cd backend
go run cmd/api/main.go
```

## 注意事项

1. **内存模式**: 当前使用内存存储，重启后数据会丢失
2. **健康检查**: Railway 会定期访问 `/health` 端点
3. **日志**: 在 Railway 控制台查看应用日志
4. **域名**: 可在 Railway 设置中配置自定义域名

## 故障排查

如果部署失败：

1. 检查构建日志
2. 确认 `go.mod` 文件存在
3. 检查端口配置（Railway 默认使用 `$PORT` 环境变量）
4. 查看运行时日志

## 下一步

要持久化数据，可以：
1. 添加 Railway PostgreSQL 插件
2. 设置 `DB_MODE=postgres`
3. 配置数据库连接环境变量
