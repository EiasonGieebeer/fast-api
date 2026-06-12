# 客户端配置

多数支持自定义 OpenAI 服务的客户端只需要填写三项：

| 配置项 | 填写内容 |
| --- | --- |
| API Key | 控制台创建的 API 密钥 |
| API Base URL | `https://www.jetbrains.top/v1` |
| Model | 从模型广场复制模型名称 |

## Cherry Studio

1. 打开「设置 → 模型服务」。
2. 添加 OpenAI 兼容服务。
3. 填写 API 密钥和 API 地址。
4. 添加或同步模型后进行连通性测试。

## ChatBox、NextChat 等

选择 OpenAI 或 OpenAI Compatible 类型，将 API Host/Base URL 改为：

```text
https://www.jetbrains.top/v1
```

## Claude Code 类工具

支持 Anthropic 自定义地址的工具可使用：

```text
ANTHROPIC_BASE_URL=https://www.jetbrains.top
ANTHROPIC_AUTH_TOKEN=YOUR_API_KEY
```

具体环境变量名称取决于客户端版本。若客户端自动在地址末尾追加 `/v1`，请避免重复填写。

::: tip
客户端报 404 时，优先检查地址是否出现 `/v1/v1`；报 401 时检查 API 密钥是否完整。
:::
