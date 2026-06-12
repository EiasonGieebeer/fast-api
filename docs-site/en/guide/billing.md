# Billing and Quota

## How usage is charged

Each request is charged according to the model price, input and output usage, and the multiplier for your account group. The final charge shown in the console usage log is authoritative.

## View pricing

The [model catalog](https://www.jetbrains.top/pricing) lists currently available models and prices. Models may be billed by tokens, requests, images, or audio duration.

## Top-ups and subscriptions

- Top-up balance can be used for pay-as-you-go requests.
- Subscription plans may have expiration dates, model restrictions, or separate quotas.
- The purchase page shows the scope that applies to each product.

## Quota display

The frontend displays balance in US dollars by default. Internally, the platform uses smaller quota units for precise settlement, so users normally do not need to convert them manually.

## Control costs

When creating an API key, consider setting:

- A quota limit for the key
- An expiration date
- Allowed models
- An IP allowlist when required by your application

If you notice unexpected usage, disable or delete the affected key immediately and review the usage logs.
