# basicauth-proxy

HTTP proxy with basicauth

# Images

```
guoyk/basicauth-proxy
ghcr.io/guoyk93/basicauth-proxy
```

# Metrics and Readiness

```
GET /metrics
GET /ready
```

# Environment Variables

- `PORT`, default to `80` listening port
- `PROXY_TARGET`, proxy target
- `PROXY_TARGET_INSECURE`, ignore TLS certificate errors if proxy target is https
- `BASICAUTH_USERNAME`, username
- `BASICAUTH_PASSWORD`, password
- `BASICAUTH_REALM`, default to `BasicAuth Proxy`, realm name of HTTP basic auth

# Donation

View https://guoyk.xyz/donation

# Credits

GUO YANKE, MIT License
