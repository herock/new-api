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

import React from 'react';
import {
  Avatar,
  Typography,
  Card,
  Button,
  Input,
  Badge,
  Space,
} from '@douyinfe/semi-ui';
import { Copy, Users, BarChart2, TrendingUp, Gift, Zap } from 'lucide-react';

const { Text } = Typography;

const InvitationCard = ({
  t,
  userState,
  renderQuota,
  setOpenTransfer,
  affLink,
  handleAffLinkClick,
  quotaForInviter,
  quotaForInvitee,
}) => {
  // 根据奖励配置生成奖励说明文案
  const getRewardDescription = () => {
    const inviterReward = renderQuota(quotaForInviter || 0);
    const inviteeReward = renderQuota(quotaForInvitee || 0);

    if (quotaForInviter > 0 && quotaForInvitee > 0) {
      return t('邀请双方奖励说明', {
        inviterReward,
        inviteeReward,
      });
    } else if (quotaForInviter > 0) {
      return t('邀请仅邀请者奖励说明', {
        inviterReward,
      });
    } else if (quotaForInvitee > 0) {
      return t('邀请仅被邀请者奖励说明', {
        inviteeReward,
      });
    }
    return '';
  };

  const rewardDescription = getRewardDescription();
  const hasInviterReward = quotaForInviter > 0;

  return (
    <Card className='!rounded-2xl shadow-sm border-0'>
      {/* 卡片头部 */}
      <div className='flex items-center mb-4'>
        <Avatar size='small' color='green' className='mr-3 shadow-md'>
          <Gift size={16} />
        </Avatar>
        <div>
          <Typography.Text className='text-lg font-medium'>
            {t('邀请奖励')}
          </Typography.Text>
          <div className='text-xs'>{t('邀请好友获得额外奖励')}</div>
        </div>
      </div>

      {/* 收益展示区域 */}
      <Space vertical style={{ width: '100%' }}>
        {/* 统计数据统一卡片 */}
        <Card
          className='!rounded-xl w-full'
          cover={
            <div className='relative h-30 bg-semi-color-fill-0'>
              {/* 标题和按钮 */}
              <div className='relative z-10 h-full flex flex-col justify-between p-4'>
                <div className='flex justify-between items-center'>
                  <Text strong style={{ fontSize: '16px' }}>
                    {t('收益统计')}
                  </Text>
                  <Button
                    type='primary'
                    theme='solid'
                    size='small'
                    disabled={
                      !userState?.user?.aff_quota ||
                      userState?.user?.aff_quota <= 0
                    }
                    onClick={() => setOpenTransfer(true)}
                    className='!rounded-lg'
                  >
                    <Zap size={12} className='mr-1' />
                    {t('划转到余额')}
                  </Button>
                </div>

                {/* 统计数据 */}
                <div className='grid grid-cols-3 gap-6 mt-4'>
                  {/* 待使用收益 */}
                  <div className='text-center'>
                    <div
                      className='text-base sm:text-2xl font-bold mb-2'
                    >
                      {renderQuota(userState?.user?.aff_quota || 0)}
                    </div>
                    <div className='flex items-center justify-center text-sm'>
                      <TrendingUp
                        size={14}
                        className='mr-1'
                        style={{ color: 'var(--semi-color-text-2)' }}
                      />
                      <Text
                        style={{
                          color: 'var(--semi-color-text-2)',
                          fontSize: '12px',
                        }}
                      >
                        {t('待使用收益')}
                      </Text>
                    </div>
                  </div>

                  {/* 总收益 */}
                  <div className='text-center'>
                    <div
                      className='text-base sm:text-2xl font-bold mb-2'
                    >
                      {renderQuota(userState?.user?.aff_history_quota || 0)}
                    </div>
                    <div className='flex items-center justify-center text-sm'>
                      <BarChart2
                        size={14}
                        className='mr-1'
                        style={{ color: 'var(--semi-color-text-2)' }}
                      />
                      <Text
                        style={{
                          color: 'var(--semi-color-text-2)',
                          fontSize: '12px',
                        }}
                      >
                        {t('总收益')}
                      </Text>
                    </div>
                  </div>

                  {/* 邀请人数 */}
                  <div className='text-center'>
                    <div
                      className='text-base sm:text-2xl font-bold mb-2'
                    >
                      {userState?.user?.aff_count || 0}
                    </div>
                    <div className='flex items-center justify-center text-sm'>
                      <Users
                        size={14}
                        className='mr-1'
                        style={{ color: 'var(--semi-color-text-2)' }}
                      />
                      <Text
                        style={{
                          color: 'var(--semi-color-text-2)',
                          fontSize: '12px',
                        }}
                      >
                        {t('邀请人数')}
                      </Text>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          }
        >
          {/* 邀请链接部分 */}
          <Input
            value={affLink}
            readonly
            className='!rounded-lg'
            prefix={t('邀请链接')}
            suffix={
              <Button
                type='primary'
                theme='solid'
                onClick={handleAffLinkClick}
                icon={<Copy size={14} />}
                className='!rounded-lg'
              >
                {t('复制')}
              </Button>
            }
          />
        </Card>

        {/* 奖励说明 */}
        {rewardDescription && (
          <Card
            className='!rounded-xl w-full'
            title={<Text type='tertiary'>{t('奖励说明')}</Text>}
          >
            <div className='space-y-3'>
              <div className='flex items-start gap-2'>
                <Badge dot type='success' />
                <Text type='tertiary' className='text-sm'>
                  {rewardDescription}
                </Text>
              </div>

              {hasInviterReward && (
                <div className='flex items-start gap-2'>
                  <Badge dot type='success' />
                  <Text type='tertiary' className='text-sm'>
                    {t('邀请划转提示')}
                  </Text>
                </div>
              )}
            </div>
          </Card>
        )}
      </Space>
    </Card>
  );
};

export default InvitationCard;
