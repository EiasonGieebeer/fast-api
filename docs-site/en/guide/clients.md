# Client Setup

Most clients require only an API key, endpoint, and model name. Create an API key in the console first, then copy an available model name from the model catalog.

| Setting | Value |
| --- | --- |
| API Key | An API key created in the console |
| OpenAI Base URL | `https://www.fastapi.cool/v1` |
| Anthropic Base URL | `https://www.fastapi.cool` |
| Model | A model name from the model catalog |

The clients below are listed alphabetically.

## CC Switch

[CC Switch](https://github.com/farion1231/cc-switch) manages provider settings for Claude Code, Codex, OpenCode, OpenClaw, and other tools.

1. Open CC Switch and select the application you want to configure.
2. Click **Add Provider** and choose a custom configuration.
3. Enter `Fast API` as the name, then provide your API key and model name.
4. Use `https://www.fastapi.cool` for Claude Code. Use `https://www.fastapi.cool/v1` for Codex, OpenCode, and OpenClaw.
5. Save the provider, select `Fast API`, and click **Enable**.

Restart the relevant terminal tool after switching providers. Claude Code usually applies the change immediately.

## ChatBox / NextChat

1. Select OpenAI or OpenAI Compatible in the settings.
2. Set API Host / Base URL to:

```text
https://www.fastapi.cool/v1
```

3. Enter your API key and model name, save, and send a test message.

## Cherry Studio

1. Open **Settings → Model Providers**.
2. Add an OpenAI-compatible provider.
3. Set the API endpoint to `https://www.fastapi.cool/v1`.
4. Enter your API key and add the model name manually.
5. Run the connection test before using the provider.

## Claude Code

Set the following environment variables:

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

For persistent settings, add the variables to `~/.claude/settings.json`:

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://www.fastapi.cool",
    "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
    "ANTHROPIC_MODEL": "YOUR_MODEL"
  }
}
```

After startup, use `/status` to verify the active authentication method and `/model` to view or change the model.

## CodeBuddy

CodeBuddy currently uses the OpenAI API format and requires the complete Chat Completions endpoint.

1. Create or edit `~/.codebuddy/models.json`.
2. Add the following configuration and replace `YOUR_MODEL` with a name from the model catalog:

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

3. Set the API key and start CodeBuddy:

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

For project-specific settings, place the same configuration in `.codebuddy/models.json` inside the project.

## Codex CLI

Edit `~/.codex/config.toml`:

```toml
model = "YOUR_MODEL"
model_provider = "fast_api"

[model_providers.fast_api]
name = "Fast API"
base_url = "https://www.fastapi.cool/v1"
env_key = "FAST_API_KEY"
wire_api = "responses"
```

Set the API key and start Codex:

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

If the selected model supports only Chat Completions and not the Responses API, choose a model with Responses support.

## Cursor

1. Open **Cursor Settings → Models**.
2. Expand API Keys and enter your OpenAI API key.
3. Enable **Override OpenAI Base URL**.
4. Set the Base URL to `https://www.fastapi.cool/v1`.
5. Add or select the model, click **Verify**, and test a text-only conversation.

::: warning Current compatibility
Cursor controls how its custom Base URL is used. In some versions, Agent mode, image input, or built-in Cursor models may not fully use the custom endpoint. If you encounter an error, test a plain text chat first or use another client from this page.
:::

## OpenClaw

Edit `~/.openclaw/openclaw.json` and add an OpenAI-compatible provider:

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

Set the `FAST_API_KEY` environment variable and restart OpenClaw. The model reference must use the `fast-api/YOUR_MODEL` format.

## OpenCode

Create `opencode.json` in the project root, or edit the user configuration at `~/.config/opencode/opencode.json`:

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

Set `FAST_API_KEY`, restart OpenCode, and select `fast-api/YOUR_MODEL` from the model list.

::: tip Troubleshooting
For a 401 error, check the API key. For a 404 error, make sure the URL does not contain `/v1/v1`. For a model-not-found error, copy the model name from the model catalog again.
:::
