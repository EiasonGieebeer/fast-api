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

Before opening checkout, the platform validates both the requested amount and the resulting wallet quota. An amount that cannot be credited safely is rejected before payment is created. Payment callbacks update the order and quota transactionally, so duplicate callbacks do not credit the wallet twice.

The final charge recorded in [Usage Logs](https://www.fastapi.cool/usage-logs/common) is based on the model, input and output usage, account-group multiplier, and other applicable billing factors.

## New information in Usage Logs

Usage Logs now distinguish standard and streaming requests more clearly and record whether a stream completed normally. Details show the recorded reasoning effort; for dynamic pricing, Conditional Multipliers identify the conditions and multipliers that actually matched this request. A tool icon next to the charge means the request includes a surcharge for Web Search, File Search, image generation, or another billable tool. Open the log details for the complete billing information.

Group, token-name, and username filters remain normal text filters, but their hidden state now uses an in-page mask with browser autocomplete disabled. This prevents password managers from inserting credentials into log criteria while preserving the show/hide control for reviewing a filter value.

Playground briefly fades in newly streamed prose; inline code, code blocks, and settled content do not replay the animation. When a history message has unsaved edits, cancelling the edit or leaving the page now asks for confirmation before discarding them.

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

## Administrative channel and redemption workflows

After an administrator fetches upstream models for a channel, new and existing models are grouped by provider for searching, bulk selection, and review. Qwen TTS models are classified under Qwen. Review model names and mappings before saving; unrecognized names remain in the `Other` group.

Channel connectivity tests now use the protocol selected for the channel. Claude and Gemini tests send their native request formats, and Gemini streaming tests use the `:streamGenerateContent` path. A successful test therefore confirms that channel's matching native endpoint; it does not imply that the same model supports every other protocol.

Under **System Settings → Models and Channels → Routing Reliability**, scheduled channel testing has three scopes: all channels except manually disabled ones, only channels with auto-disable enabled, or only auto-disabled channels awaiting recovery. The narrower modes avoid unnecessary probes of healthy channels; the separate re-enable setting still controls whether a successful check restores a channel. Administrators can also set test concurrency from 1 to 32; the task runs up to that many checks in parallel while preserving progress reporting and cancellation.

Advanced Custom channels now have a visual route editor. Start from All protocols, OpenAI only, Claude only, or Gemini only, then keep only the routes the upstream supports. Each forwarding route can define its incoming path, upstream path, converter, authentication, and exact model scope; routes may share an incoming path when their model scopes do not overlap, with at most one final catch-all. Model-list and balance-query routes are separate management routes and remain in place when a forwarding template is replaced. If a balance response is valid JSON but not OpenAI `credit_summary`, the raw upstream JSON is shown without incorrectly updating the stored channel balance.

Compatible and gateway channel types—including OpenAI, Anthropic, Codex, Advanced Custom, Sub2API, and New API—show protocol-appropriate **Field passthrough controls**. Enable only fields supported by the upstream: `service_tier`, `inference_geo`, `speed`, `store`, and obfuscation controls are not available on every protocol. Parameter override rules can also read `user_id`, `user_group`, `token_group`, and the currently selected `using_group` to tailor upstream parameters by user or group.

When configuring a custom OAuth provider, **Access Policy** offers ready-to-fill “level and active” and “organization or role” templates with nested `and` / `or` conditions. Denial messages can use variables such as `provider`, `field`, `op`, `required`, `current`, and `current.roles`. Leaving the policy empty adds no extra user restriction.

When editing a redemption code, the drawer loads the latest server record and prevents submission until loading finishes. If you change only the name or expiration and leave quota untouched, the original internal quota is preserved so display-currency conversion and decimal precision do not alter it. The quota input step follows the active currency or token display settings.

Profile does not reveal an existing personal access token, and regeneration requires confirmation. The new token is shown once and cleared from the page when the dialog closes, while the previous token becomes invalid immediately. Store the new value securely before closing, then update every client that used the old token. This protected operation is also subject to the critical endpoint rate limit, so avoid repeated regeneration.
