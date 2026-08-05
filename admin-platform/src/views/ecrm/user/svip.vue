<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';

import {
  listPlatformSvipUsersApi,
  setPlatformUserSvipApi,
  type PlatformSvipUser,
} from '#/api/core/platform-svip';

const loading = ref(false);
const saving = ref(false);
const rows = ref<PlatformSvipUser[]>([]);
const canManage = ref(false);
const total = ref(0);
const dialogOpen = ref(false);
const editing = ref<PlatformSvipUser>();
const query = reactive({ limit: 20, page: 1 });
const form = reactive({ is_svip: 0, svip_endtime: '' });

function svipText(value: number) {
  return ({ '-1': '已关闭', 0: '普通用户', 1: '体验会员', 2: '有效期会员', 3: '永久会员' }[value] || '未知');
}

function svipType(value: number) {
  return ({ '-1': 'info', 0: 'info', 1: 'warning', 2: 'success', 3: 'success' }[value] || 'info') as 'info' | 'success' | 'warning';
}

async function load() {
  if (!canManage.value) return;
  loading.value = true;
  try {
    const result = await listPlatformSvipUsersApi(query);
    rows.value = result.list;
    total.value = result.total;
  } finally { loading.value = false; }
}

function openEdit(row: PlatformSvipUser) {
  editing.value = row;
  Object.assign(form, { is_svip: row.is_svip, svip_endtime: row.svip_endtime?.slice(0, 19).replace('T', ' ') || '' });
  dialogOpen.value = true;
}

async function save() {
  if (!editing.value) return;
  if (form.is_svip === 2 && !form.svip_endtime) {
    ElMessage.warning('有效期会员必须填写到期时间');
    return;
  }
  saving.value = true;
  try {
    await setPlatformUserSvipApi(editing.value.uid, {
      is_svip: form.is_svip,
      ...(form.is_svip === 2 ? { svip_endtime: form.svip_endtime } : {}),
    });
    dialogOpen.value = false;
    ElMessage.success('会员状态已更新');
    await load();
  } finally { saving.value = false; }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value = profile.roles.includes('platform') && permissions.includes('user.svip.manage');
  await load();
});
</script>

<template>
  <Page title="付费会员" description="管理 C 端用户的 SVIP 状态；页面仅显示脱敏联系方式，会员折扣与优惠券叠加规则由订单计价服务统一执行。">
    <el-alert v-if="!canManage" class="mb-4" title="当前账号没有 SVIP 监管权限" type="warning" :closable="false" />
    <el-card v-else shadow="never"><el-table v-loading="loading" :data="rows" row-key="uid"><el-table-column label="用户 ID" prop="uid" width="92" /><el-table-column label="用户" min-width="160"><template #default="{ row }"><div>{{ row.nickname || '未设置昵称' }}</div><div class="text-xs text-muted-foreground">{{ row.phone_masked || '—' }}</div></template></el-table-column><el-table-column label="余额" width="110"><template #default="{ row }">¥{{ Number(row.now_money).toFixed(2) }}</template></el-table-column><el-table-column label="积分" prop="integral" width="90" /><el-table-column label="会员状态" width="120"><template #default="{ row }"><el-tag :type="svipType(row.is_svip)">{{ svipText(row.is_svip) }}</el-tag></template></el-table-column><el-table-column label="到期时间" min-width="170"><template #default="{ row }">{{ row.svip_endtime || '—' }}</template></el-table-column><el-table-column fixed="right" label="操作" width="76"><template #default="{ row }"><el-button link type="primary" @click="openEdit(row)">设置</el-button></template></el-table-column></el-table></el-card>
    <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div>
    <el-dialog v-model="dialogOpen" title="设置付费会员" width="480px" destroy-on-close><el-form label-width="96px"><el-form-item label="用户"><span>{{ editing?.nickname || `用户 #${editing?.uid || ''}` }}</span></el-form-item><el-form-item label="会员类型"><el-select v-model="form.is_svip" class="w-full"><el-option label="普通用户" :value="0" /><el-option label="体验会员" :value="1" /><el-option label="有效期会员" :value="2" /><el-option label="永久会员" :value="3" /><el-option label="关闭会员" :value="-1" /></el-select></el-form-item><el-form-item v-if="form.is_svip === 2" label="到期时间" required><el-date-picker v-model="form.svip_endtime" class="w-full" format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" type="datetime" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
