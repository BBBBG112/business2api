# 配置说明

## 代理池配置 (`proxy_pool`)

**内置 xray-core**，支持 vmess、vless、shadowsocks、trojan 等协议，自动转换为本地 socks5 代理。

```json
"proxy_pool": {
  "proxy": "",                    // 备用单个代理 (http/socks5 格式)
  "subscribes": [                 // 订阅链接列表 (支持 base64 编码)
    "https://example.com/sub1",
    "https://example.com/sub2"
  ],
  "files": [                      // 本地代理文件列表
    "./proxies.txt"
  ],
  "health_check": true,           // 是否启用健康检查
  "check_on_startup": false       // 启动时是否检查所有节点
}
```

### 支持的代理格式

**代理文件/订阅内容格式** (每行一个):

```
# VMess
vmess://eyJ2IjoiMiIsInBzIjoi5ZCN56ewIiwiYWRkIjoic2VydmVyLmNvbSIsInBvcnQiOiI0NDMiLCJpZCI6InV1aWQiLCJhaWQiOiIwIiwic2N5IjoiYXV0byIsIm5ldCI6IndzIiwicGF0aCI6Ii9wYXRoIiwiaG9zdCI6Imhvc3QuY29tIiwidGxzIjoidGxzIn0=

# VLESS
vless://uuid@server.com:443?type=ws&security=tls&path=/path&host=host.com&sni=sni.com#名称

# Shadowsocks
ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@server.com:8388#名称

# Trojan
trojan://password@server.com:443?sni=sni.com#名称

# 直接代理
http://proxy.com:8080
socks5://127.0.0.1:1080
```

---

## 号池配置 (`pool`)

```json
"pool": {
  "target_count": 50,              // 目标账号数量
  "min_count": 10,                 // 最小账号数，低于此值触发注册
  "check_interval_minutes": 30,    // 检查间隔(分钟)
  "register_threads": 1,           // 注册线程数
  "register_headless": false,      // 注册时是否无头模式
  "refresh_on_startup": true,      // 启动时是否刷新账号
  "refresh_cooldown_sec": 240,     // 刷新冷却时间(秒)
  "use_cooldown_sec": 15,          // 使用冷却时间(秒)
  "max_fail_count": 3,             // 最大失败次数
  "enable_browser_refresh": true,  // 启用浏览器刷新
  "browser_refresh_headless": false, // 浏览器刷新无头模式
  "browser_refresh_max_retry": 1   // 浏览器刷新最大重试次数
}
```

---

## 号池服务器配置 (`pool_server`)

```json
"pool_server": {
  "enable": false,                 // 是否启用
  "mode": "local",                 // 模式: local/server/client
  "server_addr": "",               // 服务器地址 (客户端模式)
  "listen_addr": ":8000",          // 监听地址 (服务器模式)
  "secret": "",                    // 认证密钥
  "target_count": 50,              // 目标账号数
  "client_threads": 2,             // 客户端并发线程数
  "data_dir": "./data",            // 数据目录
  "expired_action": "delete"       // 过期账号处理方式
}
```

**模式说明**:
- `local`: 本地模式，独立运行
- `server`: 服务器模式，提供号池服务和API
- `client`: 客户端模式，连接服务器接收注册/续期任务

**expired_action 说明**:
- `delete`: 删除过期/失败账号
- `refresh`: 尝试浏览器刷新Cookie
- `queue`: 保留在队列等待重试

---

## Flow 配置 (`flow`)

```json
"flow": {
  "enable": false,                 // 是否启用 Flow 视频生成
  "tokens": [],                    // Flow ST Tokens
  "proxy": "",                     // Flow 专用代理
  "timeout": 120,                  // 超时时间(秒)
  "poll_interval": 3,              // 轮询间隔(秒)
  "max_poll_attempts": 500         // 最大轮询次数
}
```

---

## 其他配置

```json
{
  "api_keys": ["key1", "key2"],    // API 密钥列表
  "listen_addr": ":8000",          // 监听地址
  "data_dir": "./data",            // 数据目录
  "default_config": "",            // 默认 configId
  "debug": false,                  // 调试模式
  "proxy": "http://127.0.0.1:10808" // 全局代理 (兼容旧配置)
}
```
