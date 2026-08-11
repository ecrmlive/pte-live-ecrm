<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { fetchPlatformCouponTemplates, issuePlatformUserCoupon, revokePlatformUserCoupon, type PlatformCouponTemplate } from '#/api/core/ecrm';
import UserRelationSelect from '#/components/ecrm/UserRelationSelect.vue';

const templates = ref<PlatformCouponTemplate[]>([]);
const loading = ref(false);
const submitting = ref(false);
const form = reactive({ user_id: undefined as number | undefined, coupon_id: undefined as number | undefined, action: 'issue' as 'issue' | 'revoke', reason: '' });
const selected = computed(() => templates.value.find((item) => item.coupon_id === form.coupon_id));

async function loadTemplates() { loading.value = true; try { templates.value = (await fetchPlatformCouponTemplates()).list || []; } finally { loading.value = false; } }
function reset() { form.user_id = undefined; form.coupon_id = undefined; form.action = 'issue'; form.reason = ''; }

async function submit() {
  const reason = form.reason.trim();
  if (!form.user_id || !form.coupon_id || reason.length < 2 || reason.length > 500) { ElMessage.warning('请填写用户、优惠券和 2 至 500 字的操作原因'); return; }
  submitting.value = true;
  try {
    const input = { reason, idempotency_key: `coupon-${form.action}-${form.user_id}-${form.coupon_id}-${crypto.randomUUID()}` };
    if (form.action === 'issue') await issuePlatformUserCoupon(form.user_id, form.coupon_id, input); else await revokePlatformUserCoupon(form.user_id, form.coupon_id, input);
    ElMessage.success(form.action === 'issue' ? '优惠券已发放并写入审计记录' : '未锁定优惠券已撤销并保留审计记录'); reset();
  } finally { submitting.value = false; }
}
onMounted(() => void loadTemplates());
</script>

<template>
  <Page title="用户优惠券操作" description="仅平台角色可执行。发券只允许当前有效模板且同一用户同一模板仅一次；撤销只允许未锁定、未使用的用户券，绝不删除订单关联事实。">
    <el-alert class="mb-4" type="warning" :closable="false" title="锁定、已使用或已过期优惠券不能通过本页修改。订单锁券、支付核销和取消恢复仍由订单状态机负责。" />
    <el-card v-loading="loading" shadow="never" class="max-w-3xl"><el-form label-width="108px" @submit.prevent="submit"><el-form-item label="关联用户" required><UserRelationSelect v-model="form.user_id" /></el-form-item><el-form-item label="操作" required><el-radio-group v-model="form.action"><el-radio value="issue">发放</el-radio><el-radio value="revoke">撤销</el-radio></el-radio-group></el-form-item><el-form-item label="优惠券模板" required><el-select v-model="form.coupon_id" filterable class="w-96" placeholder="选择当前有效优惠券"><el-option v-for="item in templates" :key="item.coupon_id" :value="item.coupon_id" :label="`${item.name}（#${item.coupon_id}）`" /></el-select></el-form-item><el-form-item v-if="selected" label="规则"><span>{{ selected.discount_type === 'rate' ? `${selected.discount_value / 10} 折` : `减 ¥${Number(selected.discount_value).toFixed(2)}` }}，满 ¥{{ Number(selected.min_amount).toFixed(2) }} 可用</span></el-form-item><el-form-item label="操作原因" required><el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文工单演示，补发活动券。" /></el-form-item><el-form-item><el-button type="primary" :loading="submitting" @click="submit">确认{{ form.action === 'issue' ? '发放' : '撤销' }}</el-button><el-button :disabled="submitting" @click="reset">重置</el-button></el-form-item></el-form></el-card>
  </Page>
</template>
