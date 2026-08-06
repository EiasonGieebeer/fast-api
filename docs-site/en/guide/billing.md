# Billing and Quota

## How usage is charged

Each request is charged according to the model price, input and output usage, the group that actually handled the request, and any applicable tool-call surcharge. The final charge shown in the console usage log is authoritative.

## View pricing

The [model catalog](https://www.fastapi.cool/pricing) lists available models, capabilities, protocols, and prices. Models may be billed by tokens, requests, images, audio duration, or dynamic rules.

For token-priced models, tools such as Web Search, File Search, and image generation can add a fee based on actual calls. Per-request models normally do not receive an additional tool fee. Administrators can change tool prices, so use the current Usage Log as the source of truth.

## Wallet, top-ups, and subscriptions

The new console combines related features in [Wallet](https://www.fastapi.cool/wallet):

- current balance and account quota;
- supported top-up methods;
- subscription plans;
- redemption codes;
- billing and top-up history;
- referral rewards and eligible balance transfers.

Some cards may be hidden by site configuration. A subscription can have an expiration date, model restrictions, or a separate quota; the Wallet page shows the applicable terms.

## Quota display

The frontend may display quota as USD, CNY, or tokens depending on site settings. Internally, the platform uses smaller quota units for precise settlement, so users normally do not need to convert them manually.

## Control costs

When creating an API key, consider setting:

- A quota limit for the key
- An expiration date
- Allowed models
- An IP allowlist when required by your application

If you notice unexpected usage, disable or delete the affected key immediately and review the usage logs.

## Verify a charge

The final charge in [Usage Logs](https://www.fastapi.cool/usage-logs/common) is authoritative. Logs include model, request time, input and output usage, the selected group, stream status, latency, quota consumption, and errors. A tool icon next to the amount marks a tool-call surcharge. For image, video, and other asynchronous jobs, also review [Task Logs](https://www.fastapi.cool/usage-logs/task).
