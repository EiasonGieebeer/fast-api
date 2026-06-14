import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Modal,
  Button,
  Space,
  Typography,
  Input,
  Banner,
} from '@douyinfe/semi-ui';
import { API, copy, showError, showSuccess } from '../../../../helpers';

const { Text } = Typography;

const CodexOAuthModal = ({ visible, onCancel, onSuccess }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [authorizeUrl, setAuthorizeUrl] = useState('');
  const [input, setInput] = useState('');

  useEffect(() => {
    if (!visible) return;
    setAuthorizeUrl('');
    setInput('');
  }, [visible]);

  const startOAuth = async () => {
    setLoading(true);
    try {
      const res = await API.post('/api/channel/codex/oauth/start', {});
      const url = res?.data?.data?.authorize_url;
      if (!res?.data?.success || !url) {
        throw new Error(res?.data?.message || t('启动授权失败'));
      }
      setAuthorizeUrl(url);
      window.open(url, '_blank', 'noopener,noreferrer');
    } catch (error) {
      showError(error?.message || t('启动授权失败'));
    } finally {
      setLoading(false);
    }
  };

  const completeOAuth = async () => {
    setLoading(true);
    try {
      const res = await API.post('/api/channel/codex/oauth/complete', {
        input,
      });
      const key = res?.data?.data?.key;
      if (!res?.data?.success || !key) {
        throw new Error(res?.data?.message || t('授权失败'));
      }
      onSuccess?.(key);
      showSuccess(t('已生成授权凭据'));
      onCancel?.();
    } catch (error) {
      showError(error?.message || t('授权失败'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      title={t('Codex 授权')}
      visible={visible}
      onCancel={onCancel}
      maskClosable={false}
      width={720}
      footer={
        <Space>
          <Button onClick={onCancel} disabled={loading}>
            {t('取消')}
          </Button>
          <Button
            theme='solid'
            type='primary'
            onClick={completeOAuth}
            loading={loading}
            disabled={!input.trim()}
          >
            {t('生成并填入')}
          </Button>
        </Space>
      }
    >
      <Space vertical spacing='tight' style={{ width: '100%' }}>
        <Banner
          type='info'
          description={t(
            '完成登录后，复制浏览器地址栏中的完整 localhost 回调 URL 并粘贴到下方。',
          )}
        />
        <Space wrap>
          <Button type='primary' onClick={startOAuth} loading={loading}>
            {t('打开授权页面')}
          </Button>
          <Button
            theme='outline'
            disabled={!authorizeUrl || loading}
            onClick={() => copy(authorizeUrl)}
          >
            {t('复制授权链接')}
          </Button>
        </Space>
        <Input
          value={input}
          onChange={setInput}
          placeholder={t('请粘贴完整回调 URL（包含 code 与 state）')}
          showClear
        />
        <Text type='tertiary' size='small'>
          {t('生成的 JSON 凭据会自动填入渠道密钥。')}
        </Text>
      </Space>
    </Modal>
  );
};

export default CodexOAuthModal;
