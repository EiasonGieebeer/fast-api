# Quick Start

## 1. Register and sign in

Open [Fast API](https://www.fastapi.cool), click **Get Started** or **Sign In**, and create an account with a username, password, email address, and email verification code.

After signing in, the site opens [Overview](https://www.fastapi.cool/dashboard/overview). Continue through the sidebar.

## 2. Choose a model and check quota

1. Open the [model catalog](https://www.fastapi.cool/pricing) and filter by provider, capability, or pricing type.
2. Open the model details and confirm its name, supported protocol, and price.
3. Use [Wallet](https://www.fastapi.cool/wallet) to review balance, top up, purchase a subscription, redeem a code, or inspect billing history. Visible features depend on site configuration.

## 3. Create an API key

Open **API Keys** in the sidebar or go directly to the [API Keys page](https://www.fastapi.cool/keys):

1. Click **Create API Key**.
2. Enter a name and select an available group.
3. Optionally set expiration, quantity, and per-key quota. Leave **Unlimited Quota** enabled when no per-key limit is required.
4. Expand **Advanced Settings** to restrict models or IP/CIDR ranges.
5. Save, then use the copy button in the key table to obtain the full key. Click the masked value when you need to reveal it.

::: warning
An API key grants access to your account quota. Never publish it on a website, in a public repository, or in a chat message.
:::

## 4. Verify in Playground

Open [Playground](https://www.fastapi.cool/playground), select the model, and send a simple message. This verifies your account, model, and platform routing before you configure an external client.

## 5. Send your first API request

Replace `YOUR_API_KEY` and `MODEL_ID` with your actual values:

```bash
curl https://www.fastapi.cool/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "MODEL_ID",
    "messages": [
      {"role": "user", "content": "Introduce yourself in one sentence."}
    ]
  }'
```

A JSON response means the connection is working. For Responses, Claude, or Gemini compatible models, use the matching example in [API Usage](/en/guide/api).

## 6. Review the result

- [Usage Logs](https://www.fastapi.cool/usage-logs/common): standard API requests, charges, latency, request IDs, and errors.
- [Task Logs](https://www.fastapi.cool/usage-logs/task): image, video, and other asynchronous tasks.
- [Overview](https://www.fastapi.cool/dashboard/overview): balance, keys, and recent usage.
- [Dashboard](https://www.fastapi.cool/dashboard/models): model call analytics and trends.

See [New Console Workflow](/en/guide/console) for the complete page map.
