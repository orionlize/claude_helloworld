# 登录说明

## 问题: "Invalid credentials" 错误

### 原因
后端默认配置为 PostgreSQL 模式,但数据库中没有用户数据。

### 解决方案

#### 方案 1: 使用内存模式(推荐用于快速测试)

1. 编辑 `backend/.env` 文件,设置:
   ```env
   DB_MODE=memory
   ```

2. 重启后端服务

3. 使用以下测试账号登录:
   - Email: `demo@example.com`
   - Password: `demo123`

#### 方案 2: 注册新用户

1. 访问登录页面
2. 点击"注册"按钮
3. 填写:
   - 用户名
   - 邮箱地址
   - 密码
4. 注册成功后使用该账号登录

#### 方案 3: 使用 PostgreSQL 数据库

1. 确保 PostgreSQL 正在运行
2. 配置 `backend/.env`:
   ```env
   DB_MODE=postgres
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your-password
   DB_NAME=apihub
   ```

3. 运行数据库迁移(如果有)
4. 注册新用户并登录

### 环境变量说明

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_MODE` | 数据库模式: `memory` 或 `postgres` | `postgres` |
| `JWT_SECRET` | JWT 密钥 | `your-secret-key-change-this` |

### 当前配置

项目已配置为使用**内存模式**,可以直接使用测试账号登录。

### 常见问题

**Q: 为什么会有这个错误?**
A: PostgreSQL 模式下需要先注册用户,内存模式有预设的 demo 用户。

**Q: 如何切换到生产环境?**
A: 配置 Supabase 或 PostgreSQL,设置 `DB_MODE=postgres`。

**Q: 忘记密码怎么办?**
A: 内存模式可以重启服务重置,数据库模式需要修改数据库。
