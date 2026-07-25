# 新版控制台流程

新版网站将旧控制台拆分为独立页面。登录后可通过左侧导航进入各项功能，页面地址也可以直接收藏。

## 推荐使用顺序

1. 在 [概览](https://www.fastapi.cool/dashboard/overview) 检查账户余额、密钥和近期使用情况。
2. 在 [模型广场](https://www.fastapi.cool/pricing) 确认模型名称、能力、支持协议和价格。
3. 如需补充余额或购买套餐，进入 [钱包](https://www.fastapi.cool/wallet)。
4. 在 [API 密钥](https://www.fastapi.cool/keys) 创建调用凭证，并按需限制分组、额度、有效期、模型和 IP。
5. 在 [游乐场](https://www.fastapi.cool/playground) 完成站内测试，或按 [API 调用](/guide/api) 配置 SDK。
6. 请求完成后到 [使用日志](https://www.fastapi.cool/usage-logs/common) 或 [任务日志](https://www.fastapi.cool/usage-logs/task) 核对结果。

## 页面导航

| 新版页面 | 地址 | 主要用途 |
| --- | --- | --- |
| 概览 | `/dashboard/overview` | 账户状态、余额、密钥和近期用量 |
| 数据看板 | `/dashboard/models` | 模型调用分析、趋势和流量 |
| API 密钥 | `/keys` | 创建、显示、复制、编辑、启用或停用密钥 |
| 使用日志 | `/usage-logs/common` | 普通 API 请求、计费、耗时和错误 |
| 任务日志 | `/usage-logs/task` | 图片、视频等异步任务记录 |
| 钱包 | `/wallet` | 余额、充值、订阅、兑换码和账单 |
| 个人资料 | `/profile` | 资料、密码、登录会话、Passkey 和两步验证 |
| 游乐场 | `/playground` | 在浏览器内测试模型和请求参数 |

::: info
站点管理员可以按账号权限或系统配置隐藏部分菜单。某个入口未显示时，不代表页面故障，请以当前账户实际可见功能为准。
:::

## API 密钥的新操作方式

创建密钥时，基础设置包括名称、分组、有效期、数量和额度。高级设置可以限制模型与 IP/CIDR。

创建完成后：

- 点击密钥行旁的复制按钮，可复制完整密钥；
- 点击脱敏显示的密钥，可按需查看完整内容；
- 可随时编辑、停用或删除密钥；
- 密钥泄露时应立即停用或删除，再创建新密钥。

## 钱包与计费

钱包页将余额、充值、订阅、兑换码、邀请奖励和账单记录集中在同一页面。具体显示内容取决于管理员启用的支付与套餐功能。

最终扣费以 [使用日志](https://www.fastapi.cool/usage-logs/common) 中记录的模型、输入输出用量、分组倍率和其他计费项为准。

## 从旧地址迁移

旧控制台链接仍可能被自动重定向，但建议更新书签和客户端说明：

| 旧入口 | 新入口 |
| --- | --- |
| `/console` | `/dashboard/overview` |
| `/console/token` | `/keys` |
| `/console/log` | `/usage-logs/common` |
| `/console/task` | `/usage-logs/task` |
| `/console/topup` | `/wallet` |

管理员功能已分别移动到「渠道」「模型」「用户」「订阅」「系统信息」和「系统设置」页面，普通用户不会看到这些入口。
