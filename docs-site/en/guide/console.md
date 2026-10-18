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
| Task Plugins (super administrator) | `/task-plugins` | Manage built-in and custom asynchronous task platforms |

::: info
Administrators can hide modules based on account permissions or system configuration. A missing menu item does not necessarily indicate an error.
:::

## API key workflow

The create drawer contains name, group, expiration, quantity, and quota settings. Advanced settings can restrict models and IP/CIDR ranges and decide whether the key may query its owner's remaining and used account quota. New keys allow account-quota queries by default; keys created before this upgrade remain off until explicitly enabled. When enabled, `GET /api/usage/token/` returns both key usage and account quota. Key-usage, quota, and key-log queries all enforce the key's IP allowlist.

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

The saved subscription billing preference is always shown as-is. With no active subscription, **Subscription First** falls back to the wallet, while **Subscription Only** rejects requests. The page explains the applicable behavior instead of displaying either saved setting as **Wallet First**.

Before opening checkout, the platform validates both the requested amount and the resulting wallet quota. An amount that cannot be credited safely is rejected before payment is created. Payment callbacks update the order and quota transactionally, so duplicate callbacks do not credit the wallet twice.

The final charge recorded in [Usage Logs](https://www.fastapi.cool/usage-logs/common) is based on the model, input and output usage, account-group multiplier, and other applicable billing factors.

## New information in Usage Logs

Usage Logs now distinguish standard and streaming requests more clearly and record whether a stream completed normally. Details show the recorded reasoning effort; for dynamic pricing, Conditional Multipliers identify the conditions and multipliers that actually matched this request. A tool icon next to the charge means the request includes a surcharge for Web Search, File Search, image generation, or another billable tool. Open the log details for the complete billing information.

The common Usage Logs page can now export a selected time range as CSV while retaining the view's other filters; column names use the current interface language. Users export only their own logs, while administrators in the all-users view export every user's logs. One export may run per account at a time. The administrator's **Longest log time range** also bounds regular-user log searches, personal dashboard charts, and every export; its default is 365 days. Administrators can disable user exports under **System Settings → Log Settings → Log search and export**, but administrator exports remain available.

Protocol-conversion warnings, rejection reasons, channel chains, quota saturation, and task-plugin diagnostics are grouped under the log's administrator information. Ordinary users receive only public fields; administrators can inspect privileged metadata such as rejection reasons in the details view.

Group, token-name, and username filters remain normal text filters, but their hidden state now uses an in-page mask with browser autocomplete disabled. This prevents password managers from inserting credentials into log criteria while preserving the show/hide control for reviewing a filter value.

Playground briefly fades in newly streamed prose; inline code, code blocks, and settled content do not replay the animation. When a history message has unsaved edits, cancelling the edit or leaving the page now asks for confirmation before discarding them.

Task Logs now present unified task details, billing usage facts, and available artifacts for plugin-backed jobs. Video, image, and other artifacts open through access URLs bound to the exact task and artifact key. Standard users do not receive channel keys, upstream task IDs, or other private fields, while diagnostic details are separated by administrator and super-administrator role.

The **Notifications** profile tab independently controls low-balance reminders, channel notifications, and upstream-model update notices. Low-balance reminders default on and use the user's threshold or the site default; another reminder is armed only after a top-up raises the balance above the threshold. Channel and upstream-model notices go only to administrators who enable them, except that the root user receives channel notices until explicitly disabling them. Verification codes, password resets, email binding, and account-security mail depend on SMTP alone and ignore these notification switches.

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

The channel list now has an optional drag-ordering mode in both table and card views. A channel can be placed at the same priority, before a tier, or after it, with a live preview of the priorities to be written. Drag mode is off by default and includes a pausable, replayable first-use guide. Concurrent edits are detected instead of overwritten, and a successful move can be undone.

Asynchronous image, video, and similar platforms now run through sandboxed JavaScript task plugins. Super administrators can review factory plugins, install from the marketplace, or upload custom plugins on **Task Plugins**. A custom version with the same key overrides the factory version; disabling or deleting the override restores the factory version. Disabling the master switch stops both factory and custom plugins and affects in-flight work, so review the listed channel and in-flight task counts first.

A directly served factory plugin updates only with a system release, so its marketplace card shows **Updates with the system** instead of an install action. Existing custom overrides and third-party plugins can still be installed or upgraded. When an override differs from the current built-in version, the table explains which built-in version will be restored if the custom copy is deleted.

Task-plugin models can use aliases from channel `model_mapping`, including case variants of the same declared name. If a requested model is claimed by a plugin but the current group has no enabled channel for it, the 503 response identifies the owning plugin; enable the correct channel or intentionally disable/override that factory plugin.

Task polling distinguishes terminal 404/410 responses from authentication failures, 429/5xx responses, transport errors, and unrecognized bodies. The latter categories increment a consecutive-failure counter that resets after a valid non-terminal 2xx response. At `TASK_POLL_MAX_FAILURES` (default 20), the task fails and enters the existing settlement/refund path; `TASK_TIMEOUT_MINUTES` (default 1440) remains the outer timeout.

The relay HTTP client waits up to 1,800 seconds for upstream response headers by default; streaming after headers arrive is unaffected. Administrators can tune `RELAY_RESPONSE_HEADER_TIMEOUT`, or set it to `0` for no limit. Non-streaming providers may not send headers until generation finishes, so avoid an unrealistically short value.

A task-plugin channel must select a registered plugin and provide a Base URL, and the administrator needs the `task_plugin.bind` permission. The plugin's models, public routes, OpenAI Responses / Video protocol claims, and billing usage fields determine the exposed API and pricing editor. Marketplace installation does not force past plugin or public-route conflicts. **Task Public Address** in system settings is used to build artifact URLs and falls back to the server address only when it is empty.

The upload review now shows the request domains declared by a task plugin so administrators can inspect its external access before installation. Driver hooks can read upstream data through the controlled `utils.fetch` API, and plugins that depend on JSON member order can declare `json-order@1`. Alibaba Bailian and Doubao image models now use task plugins for the OpenAI Images protocol; use a model name declared by the plugin or align an alias through the channel model mapping.

After an administrator fetches upstream models for a channel, new and existing models are grouped by provider for searching, bulk selection, and review. Qwen TTS models are classified under Qwen. Review model names and mappings before saving; unrecognized names remain in the `Other` group.

Channel connectivity tests now use the protocol selected for the channel. Claude and Gemini tests send their native request formats, and Gemini streaming tests use the `:streamGenerateContent` path. A successful test therefore confirms that channel's matching native endpoint; it does not imply that the same model supports every other protocol.

Advanced channel settings include a per-channel **Skip upstream TLS certificate verification** switch. It applies only to that channel's HTTP transport and weakens connection security, so reserve it for controlled self-signed or internal-network deployments rather than treating it as a general fix for public certificate errors.

AWS channels can explicitly select the Mantle Messages API; existing unconfigured channels continue to use InvokeModel. Mantle supports native Claude responses, API-key or SigV4 authentication, and top-level `cache_control`. Vertex API-key mode has an optional Project ID: when present, request URLs include the project and region, including the `global`, `us`, and `eu` multi-region hosts; when absent, Express/legacy URL behavior is preserved.

Under **System Settings → Models and Channels → Routing Reliability**, scheduled channel testing has three scopes: all channels except manually disabled ones, only channels with auto-disable enabled, or only auto-disabled channels awaiting recovery. The narrower modes avoid unnecessary probes of healthy channels; the separate re-enable setting still controls whether a successful check restores a channel. Administrators can also set test concurrency from 1 to 32; the task runs up to that many checks in parallel while preserving progress reporting and cancellation.

Advanced Custom channels now have a visual route editor. Start from All protocols, OpenAI only, Claude only, or Gemini only, then keep only the routes the upstream supports. Each forwarding route can define its incoming path, upstream path, converter, authentication, and exact model scope; routes may share an incoming path when their model scopes do not overlap, with at most one final catch-all. Model-list and balance-query routes are separate management routes and remain in place when a forwarding template is replaced. If a balance response is valid JSON but not OpenAI `credit_summary`, the raw upstream JSON is shown without incorrectly updating the stored channel balance.

Compatible and gateway channel types—including OpenAI, Anthropic, Codex, Advanced Custom, Sub2API, and New API—show protocol-appropriate **Field passthrough controls**. Enable only fields supported by the upstream: `service_tier`, `inference_geo`, `speed`, `store`, and obfuscation controls are not available on every protocol. Parameter override rules can also read `user_id`, `user_group`, `token_group`, and the currently selected `using_group` to tailor upstream parameters by user or group.

When configuring a custom OAuth provider, **Access Policy** offers ready-to-fill “level and active” and “organization or role” templates with nested `and` / `or` conditions. Denial messages can use variables such as `provider`, `field`, `op`, `required`, `current`, and `current.roles`. Leaving the policy empty adds no extra user restriction.

Administrators can unbind built-in email, GitHub, Discord, WeChat, OIDC, Telegram, and LinuxDO identities from the user binding dialog. When the server sets `PASSWORD_LOGIN_ENCRYPTION_ENABLED=true`, the sign-in page fetches a temporary public key and submits an encrypted password field. This option is disabled by default; when disabled, the existing login request format remains in use.

The password sign-in page shows the site's privacy-policy and terms-of-service consent notice. This notice explains the legal effect of signing in; it does not replace server-side authentication or session validation.

When editing a redemption code, the drawer loads the latest server record and prevents submission until loading finishes. If you change only the name or expiration and leave quota untouched, the original internal quota is preserved so display-currency conversion and decimal precision do not alter it. The quota input step follows the active currency or token display settings.

Profile does not reveal an existing personal access token, and regeneration requires confirmation. The new token is shown once and cleared from the page when the dialog closes, while the previous token becomes invalid immediately. Store the new value securely before closing, then update every client that used the old token. This protected operation is also subject to the critical endpoint rate limit, so avoid repeated regeneration.

The former dashboard token is now a named access token with editable permissions. Super-administrator user-management operations also require step-up verification; ordinary administrators cannot create non-standard roles and do not see system settings reserved for the root role.
