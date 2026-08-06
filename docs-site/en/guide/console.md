# New Console Workflow

The new website replaces the old monolithic console with focused pages. After signing in, use the sidebar or bookmark the direct URLs below.

## Recommended workflow

1. Check account balance, keys, and recent usage on [Overview](https://www.fastapi.cool/dashboard/overview).
2. Confirm the model name, capabilities, protocol, and price in the [model catalog](https://www.fastapi.cool/pricing).
3. Add balance or purchase a plan through [Wallet](https://www.fastapi.cool/wallet) when necessary.
4. Create a credential on [API Keys](https://www.fastapi.cool/keys), choose a fixed or Auto group, and optionally configure quota, expiration, model, and IP restrictions.
5. Test inside [Playground](https://www.fastapi.cool/playground), or configure an SDK using [API Usage](/en/guide/api).
6. Review completed requests in [Usage Logs](https://www.fastapi.cool/usage-logs/common) or asynchronous jobs in [Task Logs](https://www.fastapi.cool/usage-logs/task).

## Page map

| Page | Path | Purpose |
| --- | --- | --- |
| Overview | `/dashboard/overview` | Account status, balance, keys, and recent usage |
| Dashboard | `/dashboard/models` | Model call analytics, trends, and traffic |
| API Keys | `/keys` | Create, reveal, copy, edit, enable, or disable keys |
| Usage Logs | `/usage-logs/common` | Standard API requests, billing, latency, and errors |
| Task Logs | `/usage-logs/task` | Image, video, and other asynchronous tasks |
| Wallet | `/wallet` | Balance, top-ups, subscriptions, redemption, and billing |
| Profile | `/profile` | Profile, password, sessions, passkeys, and two-factor authentication |
| Playground | `/playground` | Test models and request parameters in the browser |

::: info
Administrators can hide modules based on account permissions or system configuration. A missing menu item does not necessarily indicate an error.
:::

## API key workflow

The create drawer contains name, group, expiration, quantity, and quota settings. Advanced settings can restrict models and IP/CIDR ranges.

When the administrator enables **Auto groups**, a key can try multiple groups in priority order:

- **Inherit global Auto** follows the latest group order maintained by the administrator;
- a custom order lets you add, remove, and reorder groups for this key only;
- **Cross-group retry** continues with the next group when channels in the current group are unavailable;
- Usage Logs show the group that actually handled the request and the final charge.

After creation:

- use the copy button to copy the complete key;
- click the masked key when you need to reveal it;
- edit, disable, or delete a key at any time;
- if a key is exposed, disable or delete it immediately and create a replacement.

## Wallet and billing

Wallet combines balance, top-ups, subscriptions, redemption codes, referral rewards, and billing history. The exact cards shown depend on the payment and plan features enabled by the administrator.

The final charge recorded in [Usage Logs](https://www.fastapi.cool/usage-logs/common) is based on the model, input and output usage, account-group multiplier, and other applicable billing factors.

## New information in Usage Logs

Usage Logs now distinguish standard and streaming requests more clearly and record whether a stream completed normally. A tool icon next to the charge means the request includes a surcharge for Web Search, File Search, image generation, or another billable tool. Open the log details for the complete billing information.

## Migrating old links

Old console URLs may still redirect, but bookmarks and client instructions should use the new paths:

| Old path | New path |
| --- | --- |
| `/console` | `/dashboard/overview` |
| `/console/token` | `/keys` |
| `/console/log` | `/usage-logs/common` |
| `/console/task` | `/usage-logs/task` |
| `/console/topup` | `/wallet` |

Administrative features now live on separate Channels, Models, Users, Subscriptions, System Info, and System Settings pages. Standard users do not see those entries.
