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

#### 方式一：使用 Supabase（推荐）

```
ENVIRONMENT=production
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres
JWT_SECRET=your-random-secret-key-here
```

**获取 Supabase DATABASE_URL:**

1. 登录 [Supabase](https://supabase.com/)
2. 创建新项目或选择现有项目
3. 点击 Settings → Database
4. 复制 Connection string → URI
5. 格式: `postgresql://postgres:[YOUR-PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres`

**运行数据库迁移：**

在 Supabase SQL Editor 中执行 `backend/migrations/001_init.sql` 中的 SQL 语句。

#### 方式二：使用 Railway PostgreSQL（备选）

1. 在 Railway 项目中添加 PostgreSQL 插件
2. Railway 会自动添加 `DATABASE_URL` 环境变量
3. 运行数据库迁移（需要通过 Railway Console 或本地连接执行）

#### 方式三：内存模式（演示/测试）

```
ENVIRONMENT=production
DB_MODE=memory
JWT_SECRET=your-random-secret-key-here
```

**生成 JWT_SECRET:**
```bash
openssl rand -base64 32
```

### 4. 配置服务

Railway 会自动检测 `Dockerfile`。确保：

- **Builder**: Dockerfile
- **Dockerfile Path**: `Dockerfile` (根目录)
- **Context**: `/`

### 5. 获取部署 URL

部署完成后，Railway 会提供一个公网 URL，例如：
```
https://your-app.railway.app
```

## 本地开发

### 使用 Supabase

```bash
# 复制 .env.example 到 .env
cp .env.example .env

# 编辑 .env，填入 Supabase DATABASE_URL
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres

# 启动服务
docker-compose up --build
```

### 使用内存模式

```bash
# .env 文件
DB_MODE=memory

# 启动服务
docker-compose up --build
```

## Supabase 设置步骤

### 1. 创建 Supabase 项目

1. 访问 [supabase.com](https://supabase.com/)
2. 点击 "New Project"
3. 设置项目名称和密码
4. 选择区域（推荐选择靠近用户的区域）

### 2. 获取数据库连接信息

1. 进入项目 → Settings → Database
2. 找到 "Connection string" → 选择 "URI"
3. 复制连接字符串，替换 `[YOUR-PASSWORD]`

### 3. 运行数据库迁移

在 Supabase SQL Editor 中：

1. 点击 "SQL Editor" → "New Query"
2. 复制 `backend/migrations/001_init.sql` 的内容
3. 点击 "Run" 执行

或者使用本地 psql：

```bash
psql "$DATABASE_URL" < backend/migrations/001_init.sql
```

## 注意事项

1. **Supabase 连接**: Supabase 默认启用 SSL，连接字符串已包含 `sslmode=require`
2. **健康检查**: Railway 会定期访问 `/health` 端点
3. **日志**: 在 Railway 控制台查看应用日志
4. **环境变量**: 确保在 Railway 中正确设置 `DATABASE_URL`

## 故障排查

### 数据库连接失败

1. 检查 `DATABASE_URL` 格式是否正确
2. 确认 Supabase 项目状态为 Active
3. 验证密码是否正确
4. 检查 Supabase 项目是否暂停

### 部署失败

1. 检查构建日志
2. 确认 `Dockerfile` 存在且格式正确
3. 查看运行时日志
4. 验证环境变量设置

## 成本估算

- **Railway**: 免费套餐 $5/月（包含 512MB RAM）
- **Supabase**: 免费套餐 500MB 数据库
- **总计**: 完全免费（小型项目）
