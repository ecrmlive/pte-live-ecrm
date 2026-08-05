<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { createSvipPlan, listSvipPlans, updateSvipPlan, type SvipPlan, type SvipPlanInput } from '#/api/core/platform-svip-plan';
import { listSvipInterests, type SvipInterest } from '#/api/core/platform-svip-interest';

const canManage = ref(false);
const rows = ref<SvipPlan[]>([]);
const open = ref(false);
const saving = ref(false);
const editing = ref<number>();
const form = reactive<SvipPlanInput>({ name: '', price: 0, plan_type: 'period', duration_days: 30, benefits: [], status: 1, sort: 0 });
const interests = ref<SvipInterest[]>([]);

function planTypeName(type: SvipPlan['plan_type']) { return ({ trial: '体验会员', period: '期限会员', lifetime: '永久会员' }[type]); }
function parseBenefits(raw: string) { try { const parsed = JSON.parse(raw); return Array.isArray(parsed) ? parsed.filter((item): item is string => typeof item === 'string') : []; } catch { return []; } }
async function load() { if (canManage.value) rows.value = (await listSvipPlans()).list || []; }
function edit(row?: SvipPlan) {
  editing.value = row?.id;
  Object.assign(form, row ? { name: row.name, price: Number(row.price), plan_type: row.plan_type, duration_days: row.duration_days || 0, benefits: [], status: row.status, sort: row.sort } : { name: '', price: 0, plan_type: 'period', duration_days: 30, benefits: [], status: 1, sort: 0 });
  form.benefits = row ? parseBenefits(row.benefits) : [];
  open.value = true;
}
async function save() {
  if (!form.name.trim() || !form.benefits.length || (form.plan_type === 'trial' && (form.price !== 0 || form.duration_days < 1)) || (form.plan_type === 'period' && (form.price <= 0 || form.duration_days < 1)) || (form.plan_type === 'lifetime' && form.price <= 0)) { ElMessage.warning('请完整填写符合会员类型规则的套餐信息'); return; }
  if (form.plan_type === 'lifetime') form.duration_days = 0;
  saving.value = true;
  try { if (editing.value) await updateSvipPlan(editing.value, form); else await createSvipPlan(form); ElMessage.success('会员类型已保存'); open.value = false; await load(); } finally { saving.value = false; }
}
onMounted(async () => { const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]); canManage.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && codes.includes('user.svip.plan.manage'); if (canManage.value) interests.value = (await listSvipInterests()).list.filter((item) => item.status === 1); await load(); });
</script>

<template>
  <Page title="会员类型" description="维护 C 端可售 SVIP 会员类型。停用仅停止新购；付款回调、已付款订单和既有会员权益不在后台改写。">
    <el-alert v-if="!canManage" title="当前账号没有会员类型维护权限" type="warning" :closable="false" />
    <el-card v-else shadow="never"><template #header><div class="flex justify-between"><span>可售会员类型</span><el-button type="primary" @click="edit()">新增类型</el-button></div></template>
      <el-table :data="rows"><el-table-column prop="name" label="名称" min-width="140"/><el-table-column label="类型" width="110"><template #default="{ row }">{{ planTypeName(row.plan_type) }}</template></el-table-column><el-table-column label="售价" width="110"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column><el-table-column label="有效期" width="110"><template #default="{ row }">{{ row.plan_type === 'lifetime' ? '永久' : `${row.duration_days} 天` }}</template></el-table-column><el-table-column label="权益" min-width="220"><template #default="{ row }">{{ parseBenefits(row.benefits).join('、') }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column prop="sort" label="排序" width="80"/><el-table-column label="操作" width="80"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button></template></el-table-column></el-table>
    </el-card>
    <el-dialog v-model="open" :title="editing ? '编辑会员类型' : '新增会员类型'" width="520px" destroy-on-close><el-form label-width="98px"><el-form-item label="名称"><el-input v-model="form.name" maxlength="64"/></el-form-item><el-form-item label="会员类型"><el-select v-model="form.plan_type" class="w-full"><el-option label="体验会员（免费）" value="trial"/><el-option label="期限会员" value="period"/><el-option label="永久会员" value="lifetime"/></el-select></el-form-item><el-form-item label="售价"><el-input-number v-model="form.price" :disabled="form.plan_type === 'trial'" :min="0" :precision="2"/></el-form-item><el-form-item v-if="form.plan_type !== 'lifetime'" label="有效天数"><el-input-number v-model="form.duration_days" :min="1" :max="form.plan_type === 'trial' ? 31 : 3660"/></el-form-item><el-form-item label="会员权益" required><el-checkbox-group v-model="form.benefits"><el-checkbox v-for="item in interests" :key="item.id" :label="item.name">{{ item.name }}</el-checkbox></el-checkbox-group><div v-if="!interests.length" class="text-sm text-red-500">暂无启用权益，请先在“会员权益”页面维护。</div></el-form-item><el-form-item label="状态"><el-radio-group v-model="form.status"><el-radio :value="1">启用</el-radio><el-radio :value="0">停用</el-radio></el-radio-group></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0"/></el-form-item></el-form><template #footer><el-button @click="open = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
