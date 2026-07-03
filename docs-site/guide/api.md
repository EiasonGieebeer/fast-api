# API 调用

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

不同模型支持的接口和参数可能不同，请以模型广场及对应模型官方 API 规范为准。

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
