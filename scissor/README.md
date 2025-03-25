# Scissor - 文章剪藏应用

Scissor是一个用于保存和整理微信公众号文章的应用。它可以自动提取文章内容，使用Kimi AI生成摘要和标签，并将所有信息保存到MySQL数据库中。

## 功能特点

- 保存微信公众号文章
- 使用Kimi AI自动生成文章摘要
- 自动提取关键词标签
- 数据持久化存储

## 环境要求

- Go 1.21或更高版本
- MySQL 5.7或更高版本
- Kimi AI API密钥

## 配置

在运行应用之前，需要设置以下环境变量：

```bash
export MYSQL_DSN="user:password@tcp(localhost:3306)/scissor?charset=utf8mb4&parseTime=True&loc=Local"
export KIMI_API_KEY="your-kimi-api-key"
```

## 安装和运行

1. 克隆仓库：
```bash
git clone https://github.com/gorexlv/cabinet.git
cd cabinet/scissor
```

2. 安装依赖：
```bash
go mod tidy
```

3. 运行应用：
```bash
go run main.go
```

应用将在 http://localhost:8080 启动。

## API接口

### 添加文章
POST /articles
```json
{
    "title": "文章标题",
    "content": "文章内容",
    "sourceURL": "原文链接"
}
```

### 获取文章列表
GET /articles

## 数据库结构

文章表包含以下字段：
- ID: 主键
- Title: 文章标题
- Content: 文章内容
- Summary: AI生成的摘要
- Tags: 关键词标签（逗号分隔）
- SourceURL: 原文链接
- CreatedAt: 创建时间
- UpdatedAt: 更新时间 