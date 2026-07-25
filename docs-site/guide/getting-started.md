# 快速开始

## 1. 注册并登录

打开 [Fast API](https://www.fastapi.cool)，点击「开始使用」或右上角「登录」。新用户需要填写用户名、密码和邮箱，并完成邮箱验证码验证。

登录后默认进入 [概览](https://www.fastapi.cool/dashboard/overview)，可从左侧导航继续后续操作。

## 2. 选择模型并确认额度

1. 前往 [模型广场](https://www.fastapi.cool/pricing)，按供应商、能力或计费类型筛选模型。
2. 打开模型详情，确认模型名称、支持的协议和价格。
3. 在 [钱包](https://www.fastapi.cool/wallet) 查看余额、充值、订阅套餐、兑换码和账单记录。是否显示某项功能取决于当前站点配置。

## 3. 创建 API 密钥

打开左侧导航的「API 密钥」，或直接进入 [API 密钥页面](https://www.fastapi.cool/keys)：

1. 点击「创建 API 密钥」。
2. 填写名称并选择可用分组。
3. 按需设置有效期、创建数量和密钥额度；不限制密钥额度时可保持「无限额度」。
4. 展开「高级设置」后，可限制允许调用的模型和 IP/CIDR。
5. 保存后，在密钥列表点击复制按钮获取完整密钥；点击脱敏密钥可按需显示完整内容。

::: warning
API 密钥等同于账号调用凭证。不要发布到网页、公开仓库或聊天记录中，也不要把真实密钥粘贴到截图或工单。
:::

## 4. 在游乐场验证

建议先打开 [游乐场](https://www.fastapi.cool/playground)，选择模型并发送一条简单消息。这样可以先确认账号、模型和站内路由正常，再配置外部客户端。

## 5. 发起第一次 API 请求

将下面的 `YOUR_API_KEY` 和 `MODEL_ID` 替换为你的实际信息：

```bash
curl https://www.fastapi.cool/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "MODEL_ID",
    "messages": [
      {"role": "user", "content": "你好，请用一句话介绍你自己。"}
    ]
  }'
```

返回 JSON 内容即表示接入成功。若所选模型使用 Responses、Claude 或 Gemini 协议，请参阅 [API 调用](/guide/api) 中的对应示例。

## 6. 查看调用结果

- [使用日志](https://www.fastapi.cool/usage-logs/common)：查看普通 API 请求、消耗、耗时、请求 ID 和错误信息。
- [任务日志](https://www.fastapi.cool/usage-logs/task)：查看图片、视频等异步任务；绘图日志也归在任务日志导航下。
- [概览](https://www.fastapi.cool/dashboard/overview)：查看余额、密钥和近期使用概况。
- [数据看板](https://www.fastapi.cool/dashboard/models)：查看模型调用分析和趋势。

完整页面说明见 [新版控制台流程](/guide/console)。
