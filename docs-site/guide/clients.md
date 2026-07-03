# 客户端配置

多数客户端只需要 API 密钥、接口地址和模型名称。请先在控制台创建 API 密钥，再从模型广场复制实际可用的模型名称。

| 配置项 | 填写内容 |
| --- | --- |
| API Key | 控制台创建的 API 密钥 |
| OpenAI Base URL | `https://www.fastapi.cool/v1` |
| Anthropic Base URL | `https://www.fastapi.cool` |
| Model | 模型广场中的模型名称 |

以下教程按客户端英文名称排序。

## CC Switch

[CC Switch](https://github.com/farion1231/cc-switch) 可以集中管理 Claude Code、Codex、OpenCode 和 OpenClaw 等工具的服务商配置。

1. 打开 CC Switch，切换到需要配置的应用。
2. 点击「添加供应商」，选择自定义配置。
3. 名称填写 `Fast API`，填入 API 密钥和模型名称。
4. Claude Code 使用 `https://www.fastapi.cool`；Codex、OpenCode 和 OpenClaw 使用 `https://www.fastapi.cool/v1`。
5. 保存后选择 `Fast API`，点击「启用」。

切换后建议重启对应的终端工具；Claude Code 通常可以直接生效。

## ChatBox / NextChat

1. 在设置中选择 OpenAI 或 OpenAI Compatible。
2. API Host / Base URL 填写：

```text
https://www.fastapi.cool/v1
```

3. 填入 API 密钥和模型名称，保存后发送一条测试消息。

## Cherry Studio

1. 打开「设置 → 模型服务」。
2. 添加 OpenAI 兼容服务。
3. API 地址填写 `https://www.fastapi.cool/v1`。
4. 填入 API 密钥，手动添加模型名称。
5. 点击连通性测试，成功后即可使用。

## Claude Code

在终端中设置环境变量：

::: code-group

```bash [macOS / Linux]
export ANTHROPIC_BASE_URL="https://www.fastapi.cool"
export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"
export ANTHROPIC_MODEL="YOUR_MODEL"
claude
```

```powershell [Windows PowerShell]
$env:ANTHROPIC_BASE_URL="https://www.fastapi.cool"
$env:ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"
$env:ANTHROPIC_MODEL="YOUR_MODEL"
claude
```

:::

需要长期生效时，可将变量写入 `~/.claude/settings.json`：

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://www.fastapi.cool",
    "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
    "ANTHROPIC_MODEL": "YOUR_MODEL"
  }
}
```

启动后运行 `/status` 检查当前认证方式，运行 `/model` 查看或切换模型。

## CodeBuddy

CodeBuddy 目前使用 OpenAI 接口格式，并要求填写完整的 Chat Completions 地址。

1. 创建或编辑用户配置 `~/.codebuddy/models.json`。
2. 写入以下内容，并将 `YOUR_MODEL` 替换为模型广场中的名称：

```json
{
  "models": [
    {
      "id": "YOUR_MODEL",
      "name": "Fast API Model",
      "vendor": "OpenAI",
      "apiKey": "${FAST_API_KEY}",
      "url": "https://www.fastapi.cool/v1/chat/completions",
      "supportsToolCall": true
    }
  ],
  "availableModels": ["YOUR_MODEL"]
}
```

3. 设置密钥并启动：

::: code-group

```bash [macOS / Linux]
export FAST_API_KEY="YOUR_API_KEY"
codebuddy --model "YOUR_MODEL"
```

```powershell [Windows PowerShell]
$env:FAST_API_KEY="YOUR_API_KEY"
codebuddy --model "YOUR_MODEL"
```

:::

项目级配置也可以放在项目目录的 `.codebuddy/models.json`。

## Codex CLI

编辑用户配置 `~/.codex/config.toml`：

```toml
model = "YOUR_MODEL"
model_provider = "fast_api"

[model_providers.fast_api]
name = "Fast API"
base_url = "https://www.fastapi.cool/v1"
env_key = "FAST_API_KEY"
wire_api = "responses"
```

设置密钥后启动 Codex：

::: code-group

```bash [macOS / Linux]
export FAST_API_KEY="YOUR_API_KEY"
codex
```

```powershell [Windows PowerShell]
$env:FAST_API_KEY="YOUR_API_KEY"
codex
```

:::

如果所选模型只支持 Chat Completions、不支持 Responses API，请改用支持 Responses 的模型。

## Cursor

1. 打开「Cursor Settings → Models」。
2. 展开 API Keys，填入 OpenAI API Key。
3. 开启「Override OpenAI Base URL」。
4. Base URL 填写 `https://www.fastapi.cool/v1`。
5. 添加或选择对应模型，点击 Verify 后进行文本对话测试。

::: warning 当前兼容性
Cursor 的自定义 Base URL 由客户端自身控制。部分版本的 Agent、图片输入或 Cursor 内置模型可能不会完整使用自定义地址；遇到异常时请先使用普通文本聊天测试，或改用本页其他客户端。
:::

## OpenClaw

编辑 `~/.openclaw/openclaw.json`，添加一个 OpenAI 兼容供应商：

```json
{
  "agents": {
    "defaults": {
      "model": {
        "primary": "fast-api/YOUR_MODEL"
      }
    }
  },
  "models": {
    "mode": "merge",
    "providers": {
      "fast-api": {
        "baseUrl": "https://www.fastapi.cool/v1",
        "apiKey": "${FAST_API_KEY}",
        "api": "openai-completions",
        "models": [
          {
            "id": "YOUR_MODEL",
            "name": "Fast API Model"
          }
        ]
      }
    }
  }
}
```

设置 `FAST_API_KEY` 环境变量后重启 OpenClaw。模型引用必须使用 `fast-api/YOUR_MODEL` 格式。

## OpenCode

在项目根目录创建 `opencode.json`，或编辑用户级配置 `~/.config/opencode/opencode.json`：

```json
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "fast-api": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Fast API",
      "options": {
        "baseURL": "https://www.fastapi.cool/v1",
        "apiKey": "{env:FAST_API_KEY}"
      },
      "models": {
        "YOUR_MODEL": {
          "name": "Fast API Model"
        }
      }
    }
  },
  "model": "fast-api/YOUR_MODEL"
}
```

设置 `FAST_API_KEY` 后重启 OpenCode，并在模型列表中选择 `fast-api/YOUR_MODEL`。

::: tip 排查顺序
报 401 时检查 API 密钥；报 404 时检查地址是否出现 `/v1/v1`；报模型不存在时，从模型广场重新复制模型名称。
:::
