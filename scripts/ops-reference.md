# 日常运维命令清单

## fast-api（Docker 容器）

```bash
# 启动
docker start fast-api

# 停止
docker stop fast-api

# 重启
docker restart fast-api

# 查看状态
docker ps --filter name=fast-api --format '{{.Names}}  {{.Status}}  {{.Ports}}'

# 查看日志（实时）
docker logs -f fast-api

# 查看日志（最近100行）
docker logs --tail 100 fast-api

# 设置开机自启
docker update --restart=always fast-api

# 关闭开机自启
docker update --restart=no fast-api

# 进入容器
docker exec -it fast-api sh
```

## jetact（systemd 服务）

```bash
# 启动
systemctl start jetact

# 停止
systemctl stop jetact

# 重启
systemctl restart jetact

# 查看状态
systemctl status jetact

# 查看日志（实时）
journalctl -u jetact -f

# 查看日志（最近100行）
journalctl -u jetact --no-pager -n 100

# 设置开机自启
systemctl enable jetact

# 关闭开机自启
systemctl disable jetact

# 手动启动（不走 systemd）
cd /home/jetact/jetact-admin && sh startup.sh -m standalone
```

## Nginx

```bash
# 测试配置
nginx -t

# 重载配置（不中断服务）
nginx -s reload

# 重启
systemctl restart nginx
```

## 快速健康检查

```bash
# 所有服务状态一览
docker ps --filter name=fast-api --format 'fast-api: {{.Status}}'
systemctl is-active jetact nginx mysqld

# 端口监听检查
ss -tlnp | grep -E '3001|8083|80|443|3306'

# 内存使用
free -h

# HTTP 响应检查
curl -s -o /dev/null -w '%{http_code}' https://www.jetbrains.top/
curl -s -o /dev/null -w '%{http_code}' https://www.jetbrains.top/jetact/login
```
