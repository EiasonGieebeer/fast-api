# API Usage

## Choose the model protocol first

The new homepage and model gateway support several compatible protocols. Check the model name and capabilities in the [model catalog](https://www.fastapi.cool/pricing), then use the matching endpoint.

| Protocol | Base URL | Request path | Authentication |
| --- | --- | --- | --- |
| OpenAI Chat Completions | `https://www.fastapi.cool/v1` | `/chat/completions` | `Authorization: Bearer` |
| OpenAI Responses | `https://www.fastapi.cool/v1` | `/responses` | `Authorization: Bearer` |
| Claude Messages | `https://www.fastapi.cool` | `/v1/messages` | `x-api-key` |
| Gemini | `https://www.fastapi.cool` | `/v1beta/models/{model}:generateContent` | `x-goog-api-key` |

::: warning
Not every model supports all four protocols. A correct model name used with an incompatible endpoint may return a 404 or protocol-conversion error.
:::

The updated protocol conversion layer improves interoperability among OpenAI Chat, Responses, Claude, and Gemini, and adds DeepSeek Responses support. Availability still depends on the selected model and backend channel configuration; do not infer protocol support from the provider name alone.

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
    messages=[{"role": "user", "content": "Hello"}],
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
  messages: [{ role: "user", content: "Hello" }],
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
    "input": "Introduce yourself in one sentence."
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
      {"role": "user", "content": "Hello"}
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
      {"parts": [{"text": "Hello"}]}
    ]
  }'
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

Supported endpoints and parameters vary by model. Refer to the model catalog and the official API specification for that model. Test a model in [Playground](https://www.fastapi.cool/playground), then review the request in [Usage Logs](https://www.fastapi.cool/usage-logs/common).

## Streaming

For models that support OpenAI-compatible streaming, add:

```json
{
  "stream": true
}
```

After the request, [Usage Logs](https://www.fastapi.cool/usage-logs/common) show its stream status. If the client disconnects early, the upstream stream does not finish normally, or the request fails, the stream status and error details help identify where it stopped.

## Authentication

OpenAI-compatible endpoints use a Bearer token:

```http
Authorization: Bearer YOUR_API_KEY
```

The native Claude and Gemini compatible endpoints use `x-api-key` and `x-goog-api-key`, respectively, as shown above.
