
Caddy reverse proxy


casbin plugin for caddy :
https://github.com/casbin/caddy-authz/blob/master/authz.go#L100

add oauth or jwt auth to plugin

### caddy setup for reverse proxy
```text
Get backends dynamically from SRV record DNS queries:

reverse_proxy {
dynamic srv _api._tcp.example.com
}
```
