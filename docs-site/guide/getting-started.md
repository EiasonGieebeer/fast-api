# 快速开始

## 1. 注册并登录

打开 [Fast API](https://newapi.fastapi.cool)，注册账号并完成邮箱验证。

## 2. 创建 API 密钥

进入控制台的 [API 密钥](https://newapi.fastapi.cool/keys) 页面：

1. 点击「创建密钥」。
2. 填写便于识别的名称。
3. 按需设置有效期、额度和可用模型。
4. 创建后妥善保存密钥。

::: warning
API 密钥等同于账号调用凭证。不要发布到网页、公开仓库或聊天记录中。
:::

## 3. 选择模型

前往 [模型广场](https://newapi.fastapi.cool/pricing) 查看当前可用模型及价格，并复制模型名称。

## 4. 发起第一次请求

将下面的 `YOUR_API_KEY` 和 `MODEL_ID` 替换为你的实际信息：

```bash
curl https://newapi.fastapi.cool/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "MODEL_ID",
    "messages": [
      {"role": "user", "content": "你好，请用一句话介绍你自己。"}
    ]
  }'
```

返回 JSON 内容即表示接入成功。若请求失败，可在控制台的「使用日志」中查看错误详情。
