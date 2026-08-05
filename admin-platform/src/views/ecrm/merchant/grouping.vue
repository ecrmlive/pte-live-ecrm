<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  deleteStoreGroup,
  fetchPlatformDiyPages,
  fetchPlatformMerchants,
  fetchStoreGroup,
  fetchStoreGroupMerchants,
  fetchStoreGroups,
  saveStoreGroup,
  setStoreGroupStatus,
  setStoreGroupTemplate,
  type DiyPageOption,
  type PlatformMerchantRow,
  type StoreGroupMerchantRow,
  type StoreGroupRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';

const rows = ref<StoreGroupRow[]>([]);
const merchants = ref<PlatformMerchantRow[]>([]);
const templates = ref<DiyPageOption[]>([]);
const loading = ref(false);
const dialogOpen = ref(false);
const storesOpen = ref(false);
const canManage = ref(false);
const editing = ref<StoreGroupRow>();
const linkedStores = ref<StoreGroupMerchantRow[]>([]);
const form = reactive({
  parent_id: 0,
  name: '',
  sort: 0,
  status: true,
  diy_page_id: 0,
  positioning_status: false,
  longitude: undefined as number | undefined,
  latitude: undefined as number | undefined,
  address: '',
  merchant_ids: [] as number[],
});

async function load() {
  loading.value = true;
  try {
    rows.value = (await fetchStoreGroups()).list || [];
  } finally {
    loading.value = false;
  }
}

async function loadOptions() {
  const [merchantResult, diyResult] = await Promise.all([
    fetchPlatformMerchants({ page: 1, limit: 100 }),
    fetchPlatformDiyPages(),
  ]);
  merchants.value = merchantResult.list || [];
  templates.value = diyResult.list || [];
}

function reset(row?: StoreGroupRow) {
  editing.value = row;
  Object.assign(form, {
    parent_id: row?.parent_id || 0,
    name: row?.name || '',
    sort: row?.sort || 0,
    status: row ? row.status === 1 : true,
    diy_page_id: row?.diy_page_id || 0,
    positioning_status: row ? row.positioning_status === 1 : false,
    longitude: row?.longitude ?? undefined,
    latitude: row?.latitude ?? undefined,
    address: row?.address || '',
    merchant_ids: row?.merchant_ids || [],
  });
}

function add() {
  reset();
  dialogOpen.value = true;
}

async function edit(row: StoreGroupRow) {
  const detail = await fetchStoreGroup(row.id);
  reset(detail);
  dialogOpen.value = true;
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分组名称');
    return;
  }
  if ((form.longitude === undefined) !== (form.latitude === undefined)) {
    ElMessage.warning('经度和纬度需同时填写，或同时留空');
    return;
  }
  await saveStoreGroup(editing.value?.id, { ...form, name: form.name.trim(), address: form.address.trim() });
  dialogOpen.value = false;
  ElMessage.success('店铺分组已保存');
  await load();
}

async function toggleStatus(row: StoreGroupRow) {
  await setStoreGroupStatus(row.id, row.status !== 1);
  ElMessage.success('分组状态已更新，并已同步至子分组');
  await load();
}

async function updateTemplate(row: StoreGroupRow) {
  try {
    const { value } = await ElMessageBox.prompt('填写装修模板 ID；填 0 可清空绑定。', `设置“${row.name}”装修模板`, {
      inputValue: String(row.diy_page_id || 0),
      inputPattern: /^\d+$/,
      inputErrorMessage: '请输入非负整数 ID',
    });
    await setStoreGroupTemplate(row.id, Number(value));
    ElMessage.success('装修模板已更新');
    await load();
  } catch {
    // 取消操作或统一请求层已提示错误。
  }
}

async function remove(row: StoreGroupRow) {
  try {
    await ElMessageBox.confirm(`删除“${row.name}”后不可恢复；含子分组时系统会拒绝删除。是否继续？`, '删除店铺分组', { type: 'warning' });
    await deleteStoreGroup(row.id);
    ElMessage.success('店铺分组已删除');
    await load();
  } catch {
    // 取消操作或统一请求层已提示错误。
  }
}

async function showStores(row: StoreGroupRow) {
  linkedStores.value = (await fetchStoreGroupMerchants(row.id)).list || [];
  storesOpen.value = true;
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load(), loadOptions()]);
  canManage.value = permissions.includes('merchant.group.manage');
});
</script>

<template>
  <Page title="店铺分组" description="平台维护最多三级店铺分组、关联店铺及商城装修模板；区域角色和商户角色无此全局权限。">
    <template #extra><el-button v-if="canManage" type="primary" @click="add">新增分组</el-button></template>
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="id" default-expand-all>
        <el-table-column label="分组名称" min-width="260" prop="name" />
        <el-table-column label="层级" width="80"><template #default="{ row }">第 {{ row.level + 1 }} 级</template></el-table-column>
        <el-table-column label="排序" width="80" prop="sort" />
        <el-table-column label="关联店铺" width="110"><template #default="{ row }">{{ row.merchant_count }} 家</template></el-table-column>
        <el-table-column label="定位" width="90"><template #default="{ row }">{{ row.positioning_status ? '启用' : '关闭' }}</template></el-table-column>
        <el-table-column label="装修模板" width="120"><template #default="{ row }">{{ row.diy_page_id || '未绑定' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column label="操作" min-width="280" fixed="right"><template #default="{ row }">
          <el-button link type="primary" @click="showStores(row)">关联店铺</el-button>
          <template v-if="canManage"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="primary" @click="toggleStatus(row)">{{ row.status ? '停用' : '启用' }}</el-button><el-button link type="primary" @click="updateTemplate(row)">模板</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template>
        </template></el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogOpen" :title="editing ? '编辑店铺分组' : '新增店铺分组'" width="680px" destroy-on-close>
      <el-form label-width="112px">
        <el-form-item label="上级分组"><el-tree-select v-model="form.parent_id" :data="rows" node-key="id" :props="{ label: 'name', value: 'id', children: 'children' }" check-strictly clearable placeholder="不选则为一级分组" /></el-form-item>
        <el-form-item label="分组名称" required><el-input v-model="form.name" maxlength="128" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item label="初始状态"><el-switch v-model="form.status" /></el-form-item>
        <el-form-item label="关联店铺"><el-select v-model="form.merchant_ids" multiple filterable class="w-full" placeholder="可关联多个统一后台商户投影"><el-option v-for="item in merchants" :key="item.mer_id" :label="item.mer_name" :value="item.mer_id" /></el-select></el-form-item>
        <el-form-item label="装修模板"><el-select v-model="form.diy_page_id" clearable class="w-full" placeholder="不选则不绑定"><el-option :value="0" label="不绑定模板" /><el-option v-for="item in templates" :key="item.id" :label="item.name || `模板 #${item.id}`" :value="item.id" /></el-select></el-form-item>
        <el-form-item label="启用定位"><el-switch v-model="form.positioning_status" /></el-form-item>
        <el-form-item label="经纬度"><el-input-number v-model="form.longitude" :precision="7" :min="-180" :max="180" placeholder="经度" /><span class="mx-2">/</span><el-input-number v-model="form.latitude" :precision="7" :min="-90" :max="90" placeholder="纬度" /></el-form-item>
        <el-form-item label="地址"><el-input v-model="form.address" maxlength="255" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="storesOpen" title="关联店铺" width="600px"><el-table :data="linkedStores"><el-table-column label="店铺 ID" prop="merchant_id" width="100" /><el-table-column label="店铺名称" prop="merchant_name" /><el-table-column label="区域" prop="region_id" width="100" /><el-table-column label="状态" width="90"><template #default="{ row }">{{ row.status ? '启用' : '停用' }}</template></el-table-column></el-table></el-dialog>
  </Page>
</template>
