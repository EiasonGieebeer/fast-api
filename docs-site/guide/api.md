# API 调用

## 先确认模型使用的协议

新版首页和模型服务支持多种兼容协议。请先在 [模型广场](https://www.fastapi.cool/pricing) 确认模型名称和能力，再选择对应接口。

| 协议 | Base URL | 请求路径 | 认证方式 |
| --- | --- | --- | --- |
| OpenAI Chat Completions | `https://www.fastapi.cool/v1` | `/chat/completions` | `Authorization: Bearer` |
| OpenAI Responses | `https://www.fastapi.cool/v1` | `/responses` | `Authorization: Bearer` |
| Claude Messages | `https://www.fastapi.cool` | `/v1/messages` | `x-api-key` |
| Gemini | `https://www.fastapi.cool` | `/v1beta/models/{model}:generateContent` | `x-goog-api-key` |

::: warning
并非每个模型都同时支持四种协议。模型名称正确但接口不兼容时，可能返回 404 或协议转换错误。
:::

## OpenAI SDK

### Python

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://www.fastapi.cool/v1",
)

response = client.chat.completions.create(
    model="MODEL_ID",
    messages=[{"role": "user", "content": "你好"}],
)

print(response.choices[0].message.content)
```

### Node.js

```javascript
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: "YOUR_API_KEY",
  baseURL: "https://www.fastapi.cool/v1",
});

const response = await client.chat.completions.create({
  model: "MODEL_ID",
  messages: [{ role: "user", content: "你好" }],
});

console.log(response.choices[0].message.content);
```

## Responses API

```bash
curl https://www.fastapi.cool/v1/responses \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "MODEL_ID",
    "input": "你好，请用一句话介绍你自己。"
  }'
```

## Claude Messages

```bash
curl https://www.fastapi.cool/v1/messages \
  -H "x-api-key: YOUR_API_KEY" \
  -H "anthropic-version: 2023-06-01" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "MODEL_ID",
    "max_tokens": 256,
    "messages": [
      {"role": "user", "content": "你好"}
    ]
  }'
```

## Gemini

```bash
curl "https://www.fastapi.cool/v1beta/models/MODEL_ID:generateContent" \
  -H "x-goog-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [
      {"parts": [{"text": "你好"}]}
    ]
  }'
```

## 常用接口

| 用途 | 路径 |
| --- | --- |
| OpenAI Chat Completions | `/v1/chat/completions` |
| OpenAI Responses | `/v1/responses` |
| Claude Messages | `/v1/messages` |
| Embeddings | `/v1/embeddings` |
| 图片生成 | `/v1/images/generations` |
| 语音合成 | `/v1/audio/speech` |
| 语音转文字 | `/v1/audio/transcriptions` |
| Gemini | `/v1beta/models/{model}:generateContent` |

不同模型支持的接口和参数可能不同，请以模型广场及对应模型官方 API 规范为准。首次接入可先在 [游乐场](https://www.fastapi.cool/playground) 测试，再到 [使用日志](https://www.fastapi.cool/usage-logs/common) 核对请求。

## 流式响应

兼容 OpenAI 流式输出的模型可在请求中加入：

```json
{
  "stream": true
}
```

## 认证格式

OpenAI 兼容接口使用 Bearer Token：

```http
Authorization: Bearer YOUR_API_KEY
```

Claude 与 Gemini 原生兼容接口分别使用 `x-api-key` 和 `x-goog-api-key`，如上方示例所示。
