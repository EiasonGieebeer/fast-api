# Epay Gateway

桥接 [epay（易支付）协议](https://github.com/Calcium-Ion/go-epay) 与支付宝 / 微信支付官方 API 的轻量网关服务。

## 架构

```
fast-api（epay client）
  → form POST /submit.php（epay 签名）
    → 支付宝：302 跳转支付宝收银台
    → 微信支付：返回扫码页面
      → 用户完成支付
        → 支付宝/微信回调 /alipay/notify 或 /wechat/notify
          → 验签 → 转发 epay 格式通知到 fast-api
            → fast-api 余额到账
```

## 功能

- 支付宝电脑网站支付（TradePagePay）
- 微信 Native 支付（扫码支付）
- Epay MD5 签名验证与生成
- SQLite 订单持久化，幂等防重
- 支付宝沙箱 / 正式环境一键切换
- 微信支付可选配置，未配置时优雅降级

## 快速开始

### 1. 编译

```bash
cd cmd/epay-gateway
go build -o epay-gateway .
```

### 2. 配置

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml`，填写支付宝应用密钥（和可选的微信支付商户密钥）：

```yaml
server:
  host: "0.0.0.0"
  port: 8080

epay:
  partner_id: "10001"              # 商户ID，对应 fast-api 后台的 EpayId
  key: "your-md5-secret-key"       # 签名密钥，对应 fast-api 后台的 EpayKey

alipay:
  app_id: "2021XXXXXXXXXXXX"       # 支付宝应用 APPID
  is_production: false             # false=沙箱, true=正式
  private_key: |                   # 商户私钥 PKCS1 格式
    -----BEGIN RSA PRIVATE KEY-----
    ...
    -----END RSA PRIVATE KEY-----
  alipay_public_key: |             # 支付宝公钥
    -----BEGIN PUBLIC KEY-----
    ...
    -----END PUBLIC KEY-----
  notify_url: "https://pay.example.com/alipay/notify"
  return_url: "https://pay.example.com/alipay/return"

# 微信支付（可选，不填则不启用）
wechat:
  mch_id: ""
  api_v3_key: ""
  serial_no: ""
  private_key: ""
  notify_url: "https://pay.example.com/wechat/notify"
```

### 3. 运行

```bash
./epay-gateway -config config.yaml
```

### 4. 配置 fast-api 后台

登录管理后台 → 系统设置 → 支付网关，填写：

- **Epay 端点**：`https://pay.example.com`
- **Epay merchant ID**：与 `config.yaml` 中 `epay.partner_id` 一致
- **Epay secret key**：与 `config.yaml` 中 `epay.key` 一致

### 5. 支付宝开放平台配置

在网页应用设置中：
- **应用网关**：`https://pay.example.com/`
- **授权回调地址**：`https://pay.example.com/alipay/return`

## 部署

推荐使用 systemd 管理（参考下面示例），配合 nginx 反向代理 + Let's Encrypt SSL。

```ini
[Unit]
Description=Epay Gateway
After=network.target

[Service]
Type=simple
ExecStart=/opt/epay-gateway/epay-gateway -config /opt/epay-gateway/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## API 端点

| 路径 | 方法 | 说明 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/submit.php` | POST/GET | Epay 支付提交 |
| `/alipay/return` | GET/POST | 支付宝同步跳转 |
| `/alipay/notify` | POST | 支付宝异步通知 |
| `/wechat/notify` | POST | 微信支付异步通知 |
| `/wechat/check` | GET | 微信支付状态轮询 |

## 测试流程

1. 沙箱环境：`is_production: false`，使用沙箱账号支付
2. 正式环境：切换为 `true`，替换正式密钥
3. 建议先小额测试，确认到账后再正式上线
