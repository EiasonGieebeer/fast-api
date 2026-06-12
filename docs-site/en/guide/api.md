# API Usage

## OpenAI SDK

### Python

```python
from openai import OpenAI

client = OpenAI(
    api_key="YOUR_API_KEY",
    base_url="https://www.jetbrains.top/v1",
)

response = client.chat.completions.create(
    model="MODEL_ID",
    messages=[{"role": "user", "content": "Hello"}],
)

print(response.choices[0].message.content)
```

### Node.js

```javascript
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: "YOUR_API_KEY",
  baseURL: "https://www.jetbrains.top/v1",
});

const response = await client.chat.completions.create({
  model: "MODEL_ID",
  messages: [{ role: "user", content: "Hello" }],
});

console.log(response.choices[0].message.content);
```

## Common endpoints

| Purpose | Path |
| --- | --- |
| OpenAI Chat Completions | `/v1/chat/completions` |
| OpenAI Responses | `/v1/responses` |
| Claude Messages | `/v1/messages` |
| Embeddings | `/v1/embeddings` |
| Image generation | `/v1/images/generations` |
| Text to speech | `/v1/audio/speech` |
| Speech to text | `/v1/audio/transcriptions` |
| Gemini | `/v1beta/models/{model}:generateContent` |

Supported endpoints and parameters vary by model. Refer to the model catalog and the official API specification for that model.

## Streaming

For models that support OpenAI-compatible streaming, add:

```json
{
  "stream": true
}
```

## Authentication

OpenAI-compatible endpoints use a Bearer token:

```http
Authorization: Bearer YOUR_API_KEY
```
