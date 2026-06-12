# Frequently Asked Questions

## 401: Unauthorized

Check that the API key is complete, has not been disabled, and is sent using:

```http
Authorization: Bearer YOUR_API_KEY
```

## 403: Permission denied

The key may not be allowed to use the selected model, or your account group may not include it. Check the key's model restrictions and account permissions.

## 404: Endpoint not found

Check the Base URL and request path. A common mistake is adding `/v1` twice.

## 429: Too many requests or insufficient quota

Possible causes include insufficient balance, an exhausted key quota, excessive concurrency, or a rate limit. Check the usage logs for the exact reason.

## 5xx: Service error

This usually means an upstream model or network is temporarily unavailable. Retry later. If the issue continues, keep the request time, model name, and request ID from the logs.

## Why does the same model have different costs?

Cost can vary with input length, output length, caching, image count, audio duration, and account-group multipliers.

## What should I do if a key is exposed?

Disable or delete the old key immediately and create a replacement. Never continue using a key that has been made public.
