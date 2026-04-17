import { getCurrencyConfig } from './render';

export const getQuotaPerUnit = () => {
  const raw = parseFloat(localStorage.getItem('quota_per_unit') || '1');
  return Number.isFinite(raw) && raw > 0 ? raw : 1;
};

export const quotaToDisplayAmount = (quota) => {
  const q = Number(quota || 0);
  if (!Number.isFinite(q) || q <= 0) return 0;
  const { type } = getCurrencyConfig();
  if (type === 'TOKENS') return q;
  // 前端显示强制使用 USD
  const usd = q / getQuotaPerUnit();
  return usd;
};

export const displayAmountToQuota = (amount) => {
  const val = Number(amount || 0);
  if (!Number.isFinite(val) || val <= 0) return 0;
  const { type } = getCurrencyConfig();
  if (type === 'TOKENS') return Math.round(val);
  // 前端显示强制使用 USD
  const usd = val;
  return Math.round(usd * getQuotaPerUnit());
};
