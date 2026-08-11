<script setup lang="ts">
import { reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { changePlatformUserStatus } from '#/api/core/ecrm';
import UserRelationSelect from '#/components/ecrm/UserRelationSelect.vue';
const submitting = ref(false);
const form = reactive({ user_id: undefined as number | undefined, status: 0 as 0 | 1, reason: '' });
async function submit() { const reason = form.reason.trim(); if (!form.user_id || reason.length < 2 || reason.length > 500) { ElMessage.warning('请填写用户 ID 与 2 至 500 字的调整原因'); return; } submitting.value = true; try { await changePlatformUserStatus(form.user_id, { status: form.status, reason, idempotency_key: `user-status-${form.user_id}-${crypto.randomUUID()}` }); ElMessage.success(form.status ? '用户已启用，需重新登录后访问' : '用户已停用，既有 C 端会话已失效'); Object.assign(form,{user_id:undefined,status:0,reason:''}); } finally { submitting.value = false; } }
</script>
<template><Page title="用户启停" description="仅平台角色可执行。停用或重新启用都会递增身份版本，使既有 C 端令牌失效；不会取消订单、退款、资产、优惠券或佣金事实。"><el-alert class="mb-4" type="warning" :closable="false" title="停用只阻断新的受保护 C 端访问，不自动撤销已产生的交易或资金记录。"/><el-card shadow="never" class="max-w-3xl"><el-form label-width="116px" @submit.prevent="submit"><el-form-item label="关联用户" required><UserRelationSelect v-model="form.user_id" /></el-form-item><el-form-item label="目标状态" required><el-radio-group v-model="form.status"><el-radio :value="0">停用</el-radio><el-radio :value="1">启用</el-radio></el-radio-group></el-form-item><el-form-item label="调整原因" required><el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文风控工单演示，暂时停用账号。"/></el-form-item><el-form-item><el-button type="primary" :loading="submitting" @click="submit">确认调整</el-button></el-form-item></el-form></el-card></Page></template>
