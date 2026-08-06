# Frequently Asked Questions

## 401: Unauthorized

Check that the API key is complete, has not been disabled, and is sent using:

```http
Authorization: Bearer YOUR_API_KEY
```

## 403: Permission denied

The key may not have access to the model, or no available channel in its fixed group or Auto-group order may support it. Check the key's model restrictions, group settings, and account permissions.

## 404: Endpoint not found

Check the Base URL, request path, and protocol supported by the model. Common mistakes include adding `/v1` twice or calling the Responses API with a Chat-Completions-only model.

## 429: Too many requests or insufficient quota

Possible causes include insufficient balance, an exhausted key quota, excessive concurrency, or a rate limit. Check the usage logs for the exact reason.

## 5xx: Service error

This usually means an upstream model or network is temporarily unavailable. Retry later. If the issue continues, keep the request time, model name, and request ID from the logs.

## Why does the same model have different costs?

Cost can vary with input length, output length, caching, image count, audio duration, the group that actually handled the request, and tool-call surcharges. If a tool icon appears next to the charge, open the log details to review the tool fee.

## What is the difference between Auto and fixed groups?

A fixed group uses only the selected group. Auto tries groups in the global order or a key-specific custom order. With cross-group retry enabled, it continues to later groups when the current group is unavailable, so the final multiplier can depend on the group that eventually handles the request.

## What should I do if a key is exposed?

Disable or delete the old key on [API Keys](https://www.fastapi.cool/keys) immediately and create a replacement. Never continue using a key that has been made public.

## Where are the new console pages?

After signing in, use the sidebar. Common destinations are:

- Overview: `/dashboard/overview`
- API Keys: `/keys`
- Usage Logs: `/usage-logs/common`
- Task Logs: `/usage-logs/task`
- Wallet: `/wallet`
- Profile: `/profile`
- Playground: `/playground`

Old `/console/...` paths normally redirect, but bookmarks should be updated. See [New Console Workflow](/en/guide/console) for the complete map.

## Why is a menu item missing?

Menus can be hidden by account role, permissions, or administrator configuration. Channels, model management, users, and system settings require the corresponding administrative role. Top-ups, subscriptions, and rankings can also be disabled site-wide.
