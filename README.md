# Emby Extension Server

一个用于扩展 Emby 服务器功能的 Go 应用，提供服务器域名信息的 API 端点。

## 项目概述

本项目是一个轻量级的 HTTP 服务器，基于 Gin 框架构建。它提供了一个 API 端点 `/emby/System/Ext/ServerDomains`，用于返回配置的服务器域名列表。所有请求都需要通过 Emby 服务器的 Token 验证。

## 主要功能

- **Token 验证**：通过 Emby 服务器验证请求的合法性
- **域名管理**：支持配置多个服务器域名
- **灵活的 Token 提取**：支持多种 Token 传递方式
- **生产级别配置**：使用 YAML 配置文件管理设置

## API 端点

### GET /emby/System/Ext/ServerDomains

返回配置的服务器域名列表。


**成功响应（200 OK）：**

```json
{
  "data": [
    {
      "name": "Server 1",
      "url": "https://server1.example.com"
    },
    {
      "name": "Server 2",
      "url": "https://server2.example.com"
    }
  ],
  "ok": true
}
```

**错误响应（401 Unauthorized）：**

```json
{
  "error": "Token not found",
  "ok": false
}
```

或

```json
{
  "error": "Invalid token",
  "ok": false
}
```

## Token 验证机制

### 验证 URL

项目使用以下 URL 对 Token 进行验证：

```
{Emby.ServerURL}/emby/System/Info?X-Emby-Token={token}
```

**示例：**

```
https://your-emby-server.com/emby/System/Info?X-Emby-Token=abc123def456
```

### 验证流程

1. 从请求中提取 Token
2. 构造验证 URL，将 Token 作为查询参数传递
3. 向 Emby 服务器发送 GET 请求（超时时间：3 秒）
4. 如果 Emby 服务器返回 HTTP 200 状态码，则 Token 有效
5. 其他状态码或请求失败则 Token 无效

### 验证请求头

验证请求包含以下 HTTP 请求头：

- `User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36`
- `Accept: */*`

## 配置说明

项目使用 `config.yaml` 文件进行配置。

### 配置文件示例

```yaml
emby:
  server_url: "https://your-emby-server.com" # 替换为你的 Emby 服务器地址

server:
  port: 52143 # 服务器监听端口

domains:
  - name: "Server 1"
    url: "https://server1.example.com"
  - name: "Server 2"
    url: "https://server2.example.com"
  # 可以添加更多服务器
```

### 配置参数说明

| 参数              | 类型   | 说明                  | 示例                          |
| ----------------- | ------ | --------------------- | ----------------------------- |
| `emby.server_url` | string | Emby 服务器的基础 URL | `https://emby.example.com`    |
| `server.port`     | int    | 本服务器监听的端口    | `52143`                       |
| `domains[].name`  | string | 服务器域名的显示名称  | `Server 1`                    |
| `domains[].url`   | string | 服务器的完整 URL      | `https://server1.example.com` |

## 安装与运行

### 前置要求

- Go 1.16 或更高版本
- 有效的 Emby 服务器实例

### 本地运行

1. 克隆或下载项目
2. 修改 `config.yaml`，填入你的 Emby 服务器地址和域名信息
3. 运行服务器：

```bash
go run main.go
```

服务器将在配置的端口启动（默认 52143）。

### Docker 运行

项目包含 Dockerfile 和 docker-compose.yml，可以使用 Docker 运行：

```bash
docker-compose up -d
```

## 重要注意事项

### 1. Emby 服务器配置

- **必须配置**：`emby.server_url` 必须指向有效的 Emby 服务器地址
- **HTTPS 支持**：确保 Emby 服务器地址使用正确的协议（http 或 https）
- **网络连接**：本服务器必须能够访问配置的 Emby 服务器

### 2. Token 验证

- **验证超时**：Token 验证请求的超时时间为 **3 秒**。如果 Emby 服务器响应缓慢，验证可能失败
- **验证端点**：验证使用的 Emby 端点是 `/emby/System/Info`，这是 Emby 的标准系统信息端点
- **Token 有效性**：只有有效的 Emby Token 才能通过验证

### 3. 安全性建议

- **HTTPS 使用**：在生产环境中，建议使用 HTTPS 协议
- **Token 保护**：不要在日志或错误消息中暴露完整的 Token
- **防火墙配置**：限制对本服务器的访问，仅允许授权的客户端连接
- **定期更新**：保持 Go 依赖库的最新版本，以获得安全补丁

### 4. 性能考虑

- **并发连接**：Gin 框架支持高并发，但要确保 Emby 服务器能够处理验证请求
- **缓存机制**：当前版本不缓存 Token 验证结果，每次请求都会验证。如果需要性能优化，可考虑添加 Token 缓存
- **域名列表**：域名列表在服务器启动时从配置文件加载，修改后需要重启服务器

### 5. 故障排除

- **Token 验证失败**：检查 Emby 服务器是否在线，Token 是否有效
- **连接超时**：检查网络连接和 Emby 服务器的响应时间
- **配置错误**：确保 `config.yaml` 格式正确，所有必需字段都已填写

## 依赖项

- `github.com/gin-gonic/gin`：Web 框架
- `github.com/spf13/viper`：配置管理

## 许可证

MIT

## 支持

如有问题或建议，请提交 Issue 或 Pull Request。
