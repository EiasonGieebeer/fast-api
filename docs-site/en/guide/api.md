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

The updated protocol conversion layer improves interoperability among OpenAI Chat, Responses, Claude, and Gemini, and adds Responses support for DeepSeek and GLM channels. Ollama channels can pass Claude Messages to upstream `/v1/messages` and pass OpenAI Responses and Responses Compact through unchanged; vLLM-compatible requests preserve `thinking_token_budget`. Chat/Responses conversion preserves explicit `frequency_penalty` and `presence_penalty` values (including `0`) plus `prompt_cache_key`; an upstream may still reject unsupported fields, and Codex channels remove penalties that their backend does not accept. Claude conversion omits `tools` when no tool was supplied; when a tool exists without a parameter definition, it preserves the tool and supplies a valid empty-object input schema instead of silently dropping it. Availability still depends on the selected model and backend channel configuration; do not infer protocol support from the provider name alone.

The new conversion layer also preserves hosted tools, citations, Gemini grounding metadata, and authoritative upstream billing usage when requests cross OpenAI Chat, Responses, Claude, and Gemini protocols. Content that cannot be represented losslessly is recorded in conversion diagnostics instead of being silently presented as a complete conversion. Actual support still depends on the model and channel configuration.

## Model-name modifiers

You can append explicit `@key:value` modifiers to a model name to control reasoning and sampling. Supported keys are `thinking`, `effort`, `temperature`, and `topp`, for example `MODEL_ID@effort:high` or `MODEL_ID@thinking:on@temperature:0.2`. Modifiers override the corresponding request-body fields. For duplicate keys, the rightmost value wins and a diagnostic is recorded; unknown keys or invalid values return HTTP 400.

Legacy `-thinking`, `-nothinking`, and effort suffixes are now parsed only for positively matched GPT/o-series, Claude, Gemini, and DeepSeek families so real model IDs are not rewritten accidentally. DeepSeek V4.1 models `deepseek-flash` and `deepseek-v4-pro` support `-none`, `-low`, and `-max` variants. Explicit `@thinking` / `@effort` modifiers take precedence over those suffixes and map to the corresponding DeepSeek field in Chat, Claude, or Responses requests. OpenRouter no longer implicitly rewrites arbitrary `*-thinking` names. To preserve an old public alias, map it explicitly at the channel, for example `some-model-thinking` → `some-model@thinking:on`. Fully pass-through channels do not parse these modifiers.

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
| Model list | `/v1/models` |
| Claude Messages | `/v1/messages` |
| Embeddings | `/v1/embeddings` |
| Image generation | `/v1/images/generations` |
| OpenAI Video create | `/v1/videos` |
| OpenAI Video retrieve | `/v1/videos/{task_id}` |
| Generic plugin task create | `/v1/tasks/{plugin_key}` |
| Generic plugin task retrieve | `/v1/tasks/{task_id}` |
| Task artifact list | `/v1/tasks/{task_id}/artifacts` |
| Text to speech | `/v1/audio/speech` |
| Speech to text | `/v1/audio/transcriptions` |
| DeepSeek search | `/v1/tools/search` |
| DeepSeek page fetch | `/v1/tools/fetch` |
| Gemini | `/v1beta/models/{model}:generateContent` |

Supported endpoints and parameters vary by model. Refer to the model catalog and the official API specification for that model. Test a model in [Playground](https://www.fastapi.cool/playground), then review the request in [Usage Logs](https://www.fastapi.cool/usage-logs/common).

`GET /v1/models` returns a protocol-shaped list based on authentication: a Bearer token receives the OpenAI-style `data` response, while an `x-goog-api-key` header or `?key=` query parameter receives the Gemini-style `models` response.

A task plugin can claim provider-native routes, OpenAI Responses (stream, synchronous, or background modes), and OpenAI Video. A request is handled only when the currently enabled plugin explicitly claims the protocol and is bound to the selected model. Aliases and case variants in channel `model_mapping` that point to a plugin model also participate in routing, while billing identity retains the original requested name. The generic `/v1/tasks/{plugin_key}` create endpoint returns a public task ID, which can then be used with the task and artifact endpoints. The Hailuo plugin supports MiniMax-H3 `/v2/video_generation` and validates duration, resolution, aspect ratio, and input-media counts before submission. Request validation failures return HTTP 400 and should not be retried as server faults.

Protocol conversion preserves images in Responses tool outputs, reasoning and answer text that share one stream chunk, and tool JSON Schema constraints sent to Gemini. Responses image URLs converted for Chat upstreams remain structured objects. Chat/Claude-to-Responses conversion preserves `input_image.detail`, including `original`; assistant text is retained alongside tool calls; and a streamed reasoning item announced before its text still includes the protocol-required empty `summary`. Claude `thinking.block_binding`, not-yet-modeled thinking fields, and every `anthropic-beta` value are forwarded as well. Responses WebSocket sessions send a close frame when the session ends, which clients should treat as a normal close rather than a transport interruption.

`/v1/tools/search` and `/v1/tools/fetch` relay the DeepSeek Search and Fetch APIs under the fixed billing identities `deepseek-search` and `deepseek-fetch`. A call is charged only when the upstream returns search results or page content. DeepSeek server-side searches and fetches reported in usage for Claude Code or Codex are charged under the corresponding tool as well.

Body-read and JSON-parse failures retain their actual cause. A body larger than `MAX_REQUEST_BODY_MB` returns HTTP 413 across the relevant Chat, Responses, Images, Video, audio, and task-plugin entry points; reduce the request instead of retrying it as an upstream 5xx.

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

When **Allow querying user quota** is enabled in a key's advanced settings, the same Bearer key can call `GET /api/usage/token/`. In addition to the key's `total_usage`, `total_granted`, and `total_available`, the response includes remaining and used account quota under `user_quota`. Existing keys default this permission off, and the key's IP allowlist also applies to this endpoint and key-log queries.
