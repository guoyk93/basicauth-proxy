# basicauth-proxy

为 HTTP 服务提供 BasicAuth 反向代理

# 容器镜像

```
yankeguo/basicauth-proxy
ghcr.io/yankeguo/basicauth-proxy
```

# 遥测指标和健康检查

```
GET /metrics
GET /ready
```

# 环境变量

- `PORT`, 默认为 `80`, 监听端口
- `PROXY_TARGET`, 转发目标
- `PROXY_TARGET_INSECURE`, 如果转发目标为 HTTPS, 则忽略证书错误
- `BASICAUTH_USERNAME`, 设置认证用户名
- `BASICAUTH_PASSWORD`, 设置认证密码
- `BASICAUTH_REALM`, 默认为 `BasicAuth Proxy`, HTTP 认证的 Realm 名称

# 许可证

GUO YANKE, MIT License
