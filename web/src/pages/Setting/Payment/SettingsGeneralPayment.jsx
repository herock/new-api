/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useEffect, useState, useRef } from 'react';
import { Button, Form, Spin } from '@douyinfe/semi-ui';
import {
  API,
  removeTrailingSlash,
  showError,
  showSuccess,
} from '../../../helpers';
import { useTranslation } from 'react-i18next';

export default function SettingsGeneralPayment(props) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({
    ServerAddress: '',
    EnableSubscriptionPurchase: true,
  });
  const formApiRef = useRef(null);

  useEffect(() => {
    if (props.options && formApiRef.current) {
      const currentInputs = {
        ServerAddress: props.options.ServerAddress || '',
        EnableSubscriptionPurchase:
          props.options.EnableSubscriptionPurchase !== false,
      };
      setInputs(currentInputs);
      formApiRef.current.setValues(currentInputs);
    }
  }, [props.options]);

  const handleFormChange = (values) => {
    setInputs(values);
  };

  const submitServerAddress = async () => {
    setLoading(true);
    try {
      let ServerAddress = removeTrailingSlash(inputs.ServerAddress || '');
      const [serverAddressRes, subscriptionPurchaseRes] = await Promise.all([
        API.put('/api/option/', {
          key: 'ServerAddress',
          value: ServerAddress,
        }),
        API.put('/api/option/', {
          key: 'payment_setting.enable_subscription_purchase',
          value: String(inputs.EnableSubscriptionPurchase !== false),
        }),
      ]);
      if (serverAddressRes.data.success && subscriptionPurchaseRes.data.success) {
        showSuccess(t('更新成功'));
        props.refresh && props.refresh();
      } else {
        showError(
          serverAddressRes.data.message || subscriptionPurchaseRes.data.message,
        );
      }
    } catch (error) {
      showError(t('更新失败'));
    }
    setLoading(false);
  };

  return (
    <Spin spinning={loading}>
      <Form
        initValues={inputs}
        onValueChange={handleFormChange}
        getFormApi={(api) => (formApiRef.current = api)}
      >
        <Form.Section text={t('通用设置')}>
          <Form.Input
            field='ServerAddress'
            label={t('服务器地址')}
            placeholder={'https://yourdomain.com'}
            style={{ width: '100%' }}
            extraText={t(
              '该服务器地址将影响支付回调地址以及默认首页展示的地址，请确保正确配置',
            )}
          />
          <Form.Switch
            field='EnableSubscriptionPurchase'
            label={t('开放订阅购买入口')}
            extraText={t(
              '关闭后，普通用户将看不到可购买订阅；已绑定订阅的用户仍可查看自己的订阅信息。',
            )}
            size='large'
          />
          <Button onClick={submitServerAddress}>
            {t('更新通用支付设置')}
          </Button>
        </Form.Section>
      </Form>
    </Spin>
  );
}
