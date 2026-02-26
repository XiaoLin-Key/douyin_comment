# 仿抖音评论系统

基于 go-zero 框架开发的评论系统，支持评论树结构（评论和回复）。

## 项目结构

```
comment-system/
├── api/                    # API 服务
│   ├── comment.api         # API 定义文件
│   ├── etc/               # 配置文件
│   ├── internal/          # 内部代码
│   │   ├── handler/      # HTTP 处理器
│   │   ├── logic/        # 业务逻辑
│   │   ├── svc/          # 服务上下文
│   │   └── types/        # 类型定义
│   └── comment.go        # 入口文件
├── rpc/                   # RPC 服务
│   ├── comment.proto     # Protobuf 定义
│   ├── etc/              # 配置文件
│   ├── internal/         # 内部代码
│   │   ├── logic/       # RPC 业务逻辑
│   │   ├── server/      # RPC 服务
│   │   └── svc/         # 服务上下文
│   └── types/           # Protobuf 生成代码
└── model/                # 数据模型
    ├── comment.go       # Go 结构体定义
    └── comment.sql      # 数据库表结构
```

## 数据库设计

### comment 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint(20) | 评论ID（主键） |
| video_id | bigint(20) | 视频ID |
| user_id | bigint(20) | 用户ID |
| content | varchar(1000) | 评论内容 |
| parent_id | bigint(20) | 父评论ID，0表示顶级评论 |
| root_id | bigint(20) | 根评论ID，用于回复链 |
| like_count | bigint(20) | 点赞数 |
| reply_count | bigint(20) | 回复数 |
| is_deleted | tinyint(1) | 是否删除：0=否, 1=是 |
| created_at | datetime | 创建时间 |

### comment_like 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint(20) | 点赞ID（主键） |
| comment_id | bigint(20) | 评论ID |
| user_id | bigint(20) | 用户ID |
| created_at | datetime | 点赞时间 |

> **约束**：通过 `UNIQUE KEY uk_comment_user (comment_id, user_id)` 确保每个用户对同一条评论只能点赞一次。

### 评论树结构说明

- **顶级评论**：`parent_id = 0`，`root_id = 0`
- **一级回复**：`parent_id` 指向顶级评论ID，`root_id` 指向顶级评论ID
- **二级回复**：`parent_id` 指向一级回复ID，`root_id` 指向顶级评论ID

## API 接口

### 1. 发布评论
```
POST /api/v1/comment
```

请求体：
```json
{
  "video_id": 1,
  "user_id": 1001,
  "content": "这个视频太棒了！",
  "parent_id": 0
}
```

### 2. 获取评论列表
```
GET /api/v1/comments?video_id=1&root_id=0&page=1&page_size=10
```

#### 查询模式说明：
- **顶级评论查询 (`root_id = 0`)**：
  - 用于视频下方的评论区主列表。
  - **排序策略**：综合权重 = `点赞数*0.3 + 子评论数*0.3 + 时间权重*0.4`。
  - **分页逻辑**：支持无限滚动，每次通过 `page` 和 `page_size` 获取下一批顶级评论。
- **子评论查询 (`root_id > 0`)**：
  - 用于点击顶级评论下方的“查看更多”时调用。
  - **排序策略**：按 `created_at` 正序（最早的在上面）。
  - **分页逻辑**：每次固定加载 5 条回复。

### 3. 点赞/取消点赞评论
```
POST /api/v1/comment/like
```

请求体：
```json
{
  "comment_id": 1,
  "user_id": 1001,
  "action": 1
}
```
- action: 1=点赞, 0=取消点赞

### 4. 删除评论
```
DELETE /api/v1/comment
```

请求体：
```json
{
  "comment_id": 1,
  "user_id": 1001
}
```

#### 删除策略说明：
- **顶级评论 (parent_id = 0)**：执行 **物理删除**。会同时删除该评论及其下方所有的子评论（关联 `root_id`）。
- **子评论 (parent_id > 0)**：执行 **软删除**。为了保持评论树结构的连贯性，系统不会物理删除该记录，而是将 `content` 修改为 `"该评论已删除"`，并将 `is_deleted` 标记为 `1`。
}
```

## RPC 服务

### CommentService

- `CreateComment`: 创建评论
- `GetComments`: 获取评论列表（树形结构）
- `LikeComment`: 点赞/取消点赞
- `DeleteComment`: 删除评论

## 使用说明

### 1. 初始化数据库

执行 SQL 文件创建表：
```bash
mysql -u root -p your_database < model/comment.sql
```

### 2. 配置服务

修改 API 和 RPC 的配置文件：
- `api/etc/comment-api.yaml`
- `rpc/etc/comment.yaml`

### 3. 启动 RPC 服务

```bash
cd rpc
go run comment.go -f etc/comment.yaml
```

### 4. 启动 API 服务

```bash
cd api
go run comment.go -f etc/comment-api.yaml
```

## 注意事项

1. 需要先安装 go-zero 工具：
   ```bash
   go install github.com/zeromicro/go-zero/tools/goctl@latest
   ```

2. 需要安装 protoc 和相关插件：
   ```bash
   brew install protobuf
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. 数据库连接信息需要在配置文件中修改

4. 业务逻辑实现需要根据实际需求完善
