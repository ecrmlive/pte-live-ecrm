<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  assignMerchantIntentionRegion,
  auditMerchantIntention,
  fetchMerchantIntentions,
  type MerchantIntentionRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { formatShanghaiDateTime } from '#/utils/date-time';

const loading = ref(false);
const rows = ref<MerchantIntentionRow[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const regionDialogOpen = ref(false);
const submitting = ref(false);
const selected = ref<MerchantIntentionRow>();
const canAudit = ref(false);
const canAssignRegion = ref(false);
const query = reactive({ keyword: '', page: 1, limit: 20, status: 0 as number | undefined });
const form = reactive({ account: '', fail_msg: '', mark: '', password: '', region_id: 0, status: 1 });
const assignment = reactive({ region_id: 0 });

const statusText = (status: number) => ({ 0: '待审核', 1: '已通过', 2: '已驳回' }[status] || '未知');
const statusType = (status: number) => ({ 0: 'warning', 1: 'success', 2: 'danger' }[status] || 'info');
const formatTime = (value?: string | null) => {
  if (!value) return '—';
  return formatShanghaiDateTime(value);
};

async function load() {
  loading.value = true;
  try {
    const result = await fetchMerchantIntentions({ ...query, keyword: query.keyword.trim() || undefined });
    rows.value = result.list;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.keyword = '';
  query.status = 0;
  query.page = 1;
  void load();
}

function openAudit(row: MerchantIntentionRow, status: number) {
  selected.value = row;
  form.status = status;
  form.mark = '';
  form.fail_msg = '';
  form.account = '';
  form.password = '';
  form.region_id = row.circle_id || 0;
  dialogOpen.value = true;
}

function openRegionAssignment(row: MerchantIntentionRow) {
  selected.value = row;
  assignment.region_id = row.circle_id || 0;
  regionDialogOpen.value = true;
}

async function assignRegion() {
  if (!selected.value || assignment.region_id <= 0) { ElMessage.warning('请填写有效区域 ID'); return; }
  submitting.value = true;
  try {
    await assignMerchantIntentionRegion(selected.value.mer_intention_id, assignment.region_id);
    ElMessage.success('入驻申请已分配区域');
    regionDialogOpen.value = false;
    await load();
  } finally { submitting.value = false; }
}

async function submit() {
  if (!selected.value) return;
  if (form.status === 2 && !form.fail_msg.trim()) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (form.status === 1 && (!form.account.trim() || !form.password)) {
    ElMessage.warning('通过入驻时必须设置商户管理账号和初始密码');
    return;
  }
  submitting.value = true;
  try {
    await auditMerchantIntention(selected.value.mer_intention_id, {
      account: form.account.trim(),
      fail_msg: form.fail_msg.trim(),
      mark: form.mark.trim(),
      password: form.password,
      region_id: form.region_id,
      status: form.status,
    });
    ElMessage.success(form.status === 1 ? '入驻审核已通过' : '入驻申请已驳回');
    dialogOpen.value = false;
    await load();
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canAudit.value = permissions.includes('merchant.intention.audit');
  canAssignRegion.value = permissions.includes('merchant.intention.assign_region');
});
</script>

<template>
  <Page title="店铺入驻审核" description="审核通过后创建商户及其管理账号；驳回必须保留原因。">
    <el-card shadow="never">
      <el-form class="grid gap-x-4 md:grid-cols-3 xl:grid-cols-4" label-width="76px" @submit.prevent="search">
        <el-form-item label="申请搜索"><el-input v-model="query.keyword" clearable placeholder="店铺名称 / 联系人 / 手机号" @keyup.enter="search" /></el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="query.status" class="w-full">
            <el-option label="待审核" :value="0" /><el-option label="已通过" :value="1" /><el-option label="已驳回" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="search">搜索</el-button><el-button @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>
    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="mer_intention_id">
        <el-table-column label="申请 ID" prop="mer_intention_id" width="88" />
        <el-table-column label="店铺名称" min-width="150" prop="mer_name" show-overflow-tooltip />
        <el-table-column label="联系人" min-width="110" prop="name" />
        <el-table-column label="联系电话" min-width="130" prop="phone" />
        <el-table-column label="申请时间" min-width="170"><template #default="{ row }">{{ formatTime(row.create_time) }}</template></el-table-column>
        <el-table-column label="状态" width="96"><template #default="{ row }"><el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="分配区域" width="108"><template #default="{ row }">{{ row.circle_id || '未分配' }}</template></el-table-column>
        <el-table-column label="驳回原因" min-width="160" prop="fail_msg" show-overflow-tooltip />
        <el-table-column fixed="right" label="操作" width="205">
          <template #default="{ row }"><template v-if="row.status === 0"><el-button v-if="canAssignRegion" link type="primary" @click="openRegionAssignment(row)">分配区域</el-button><el-button v-if="canAudit" link type="primary" @click="openAudit(row, 1)">通过</el-button><el-button v-if="canAudit" link type="danger" @click="openAudit(row, 2)">驳回</el-button><span v-if="!canAssignRegion && !canAudit">—</span></template><span v-else>—</span></template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="form.status === 1 ? '通过店铺入驻' : '驳回店铺入驻'" width="520px" destroy-on-close>
      <el-form label-width="105px">
        <el-form-item label="申请店铺"><span>{{ selected?.mer_name }}</span></el-form-item>
        <template v-if="form.status === 1">
          <el-form-item label="商户管理账号" required><el-input v-model="form.account" autocomplete="off" /></el-form-item>
          <el-form-item label="初始密码" required><el-input v-model="form.password" autocomplete="new-password" show-password type="password" /></el-form-item>
          <el-form-item label="所属区域ID"><el-input-number v-model="form.region_id" :min="0" /><div class="text-xs text-gray-500">填写区域商圈对应的商户 region_id，区域管理员将按此范围管理商户。</div></el-form-item>
        </template>
        <el-form-item v-else label="驳回原因" required><el-input v-model="form.fail_msg" :rows="3" maxlength="300" show-word-limit type="textarea" /></el-form-item>
        <el-form-item label="审核备注"><el-input v-model="form.mark" :rows="3" maxlength="300" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="submitting" type="primary" @click="submit">确定</el-button></template>
    </el-dialog>
    <el-dialog v-model="regionDialogOpen" title="分配入驻审核区域" width="420px" destroy-on-close>
      <el-form label-width="100px"><el-form-item label="申请店铺"><span>{{ selected?.mer_name }}</span></el-form-item><el-form-item label="区域 ID" required><el-input-number v-model="assignment.region_id" :min="1" /></el-form-item><p class="text-xs text-gray-500">分配后仅对应区域管理员可查看和审核该申请。</p></el-form>
      <template #footer><el-button @click="regionDialogOpen = false">取消</el-button><el-button :loading="submitting" type="primary" @click="assignRegion">确定</el-button></template>
    </el-dialog>
  </Page>
</template>
