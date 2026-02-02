# Supabase 集成指南

本文档说明如何将项目配置为使用 Supabase 作为后端数据库。

## 什么是 Supabase?

Supabase 是一个开源的 Firebase 替代方案,提供:
- PostgreSQL 数据库
- 身份认证
- 实时订阅
- 存储功能
- Edge Functions

## 配置步骤

### 1. 创建 Supabase 项目

1. 访问 [https://supabase.com](https://supabase.com)
2. 注册账号并登录
3. 点击 "New Project"
4. 选择组织
5. 设置项目信息:
   - 项目名称
   - 数据库密码(请妥善保管)
   - 区域(选择靠近你的区域)
6. 等待项目创建完成(通常需要 1-2 分钟)

### 2. 获取凭证

1. 进入项目仪表板
2. 点击左侧菜单的 "Settings" -> "API"
3. 复制以下信息:
   - **Project URL**: 类似 `https://xxxxx.supabase.co`
   - **anon public**: 匿名公开密钥
   - **service_role**: 服务角色密钥(仅在服务端使用)

### 3. 执行数据库初始化脚本

1. 在 Supabase 仪表板中,点击 "SQL Editor"
2. 点击 "New Query"
3. 复制 `backend/supabase_setup.sql` 文件的内容
4. 粘贴到 SQL 编辑器中
5. 点击 "Run" 执行脚本

这将创建以下表:
- `users` - 用户表
- `projects` - 项目表
- `collections` - 集合表
- `endpoints` - 接口表
- `environments` - 环境变量表

### 4. 配置环境变量

#### 方式一:使用 `.env` 文件(本地开发)

复制环境变量模板:
```bash
cp .env.example .env
```

编辑 `.env` 文件,填入你的 Supabase 凭证:
```env
# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key

# Storage Mode
STORAGE_MODE=supabase
```

#### 方式二:使用 Docker Compose

创建 `.env` 文件:
```bash
cp .env.example .env
```

编辑 `.env` 文件,填入你的 Supabase 凭证

然后启动服务:
```bash
# 仅使用 Supabase(推荐)
docker-compose up -d

# 或使用传统 PostgreSQL(不推荐)
docker-compose --profile with-postgres up -d
```

### 5. 验证配置

启动后端服务:
```bash
cd backend
go run cmd/api/main.go
```

访问健康检查端点:
```bash
curl http://localhost:8080/health
```

如果返回 `{"status":"ok"}`,说明服务运行正常。

## 数据库结构

### 表结构

#### users
```sql
- id: UUID (主键)
- username: VARCHAR(255) (唯一)
- email: VARCHAR(255) (唯一)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

#### projects
```sql
- id: UUID (主键)
- name: VARCHAR(255)
- description: TEXT
- user_id: UUID (外键 -> users)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

#### collections
```sql
- id: UUID (主键)
- project_id: UUID (外键 -> projects)
- name: VARCHAR(255)
- description: TEXT
- parent_id: UUID (外键 -> collections, 可为空)
- sort_order: INTEGER
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

#### endpoints
```sql
- id: UUID (主键)
- collection_id: UUID (外键 -> collections)
- name: VARCHAR(255)
- method: VARCHAR(10)
- url: TEXT
- headers: JSONB
- body: TEXT
- description: TEXT
- request_params: JSONB
- request_body: JSONB
- response_params: JSONB
- response_body: JSONB
- sort_order: INTEGER
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

#### environments
```sql
- id: UUID (主键)
- project_id: UUID (外键 -> projects)
- name: VARCHAR(255)
- variables: JSONB
- is_default: BOOLEAN
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
```

## 安全性

项目使用 Supabase 的 Row Level Security (RLS) 来保护数据:

1. **数据隔离**: 每个用户只能访问自己的数据
2. **自动策略**: RLS 策略自动确保数据安全
3. **认证集成**: 使用 Supabase Auth 或自定义 JWT

## 从传统 PostgreSQL 迁移

如果你之前使用传统 PostgreSQL,迁移到 Supabase:

1. 导出现有数据:
```bash
pg_dump -h localhost -U postgres apihub > backup.sql
```

2. 在 Supabase 中执行 `supabase_setup.sql` 创建表结构

3. 导入数据(需要调整 UUID 格式)

4. 更新环境变量设置 `STORAGE_MODE=supabase`

## 常见问题

### Q: 忘记数据库密码怎么办?
A: 在 Supabase 项目设置中可以重置数据库密码

### Q: 如何查看数据库内容?
A: 使用 Supabase 的 "Table Editor" 功能

### Q: 如何执行自定义 SQL?
A: 使用 Supabase 的 "SQL Editor"

### Q: 连接失败怎么办?
A: 检查:
1. SUPABASE_URL 是否正确
2. 密钥是否正确
3. 网络连接是否正常
4. Supabase 项目是否处于活动状态

## 生产环境建议

1. **环境变量**: 永远不要将密钥提交到代码仓库
2. **备份**: 定期备份 Supabase 数据库
3. **监控**: 使用 Supabase Dashboard 监控性能
4. **限制**: 设置适当的 API 请求限制
5. **安全**: 使用 `service_role` 密钥仅在服务端使用

## 更多资源

- [Supabase 官方文档](https://supabase.com/docs)
- [Supabase Go 客户端](https://github.com/supabase/supabase-go)
- [项目 GitHub Issues](https://github.com/your-repo/issues)
