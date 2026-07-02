# Quick Start

## 1. Register and sign in

Open [Fast API](https://newapi.fastapi.cool), create an account, and complete email verification.

## 2. Create an API key

Open the [API Keys](https://newapi.fastapi.cool/keys) page in the console:

1. Click **Create key**.
2. Enter a recognizable name.
3. Optionally set an expiration date, quota, and allowed models.
4. Store the key securely after it is created.

::: warning
An API key grants access to your account quota. Never publish it on a website, in a public repository, or in a chat message.
:::

## 3. Choose a model

Open the [model catalog](https://newapi.fastapi.cool/pricing) to review currently available models and pricing, then copy the model name.

## 4. Send your first request

Replace `YOUR_API_KEY` and `MODEL_ID` with your actual values:

```bash
curl https://newapi.fastapi.cool/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "MODEL_ID",
    "messages": [
      {"role": "user", "content": "Introduce yourself in one sentence."}
    ]
  }'
```

A JSON response means the connection is working. If the request fails, check the usage logs in the console for details.
