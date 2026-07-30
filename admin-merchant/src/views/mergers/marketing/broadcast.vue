<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  createMerchantBroadcastRoomApi,
  deleteMerchantBroadcastRoomApi,
  getMerchantBroadcastRoomApi,
  listMerchantBroadcastRoomsApi,
  setMerchantBroadcastGoodsApi,
  setMerchantBroadcastLiveApi,
  updateMerchantBroadcastRoomApi,
  type MerchantBroadcastRoom,
  type MerchantBroadcastRoomInput,
} from '#/api/core/merchant-broadcast';
import { listMerchantProductsApi, type MerchantProduct } from '#/api/core/merchant-catalog';
import ImageField from '#/components/shop/image-field.vue';

const loading = ref(false);
const saving = ref(false);
const rows = ref<MerchantBroadcastRoom[]>([]);
const products = ref<MerchantProduct[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const goodsOpen = ref(false);
const liveOpen = ref(false);
const editingID = ref<number>();
const goodsRoom = ref<MerchantBroadcastRoom>();
const liveRoom = ref<MerchantBroadcastRoom>();
const selectedGoods = ref<number[]>([]);
const selectedLiveStatus = ref(102);
const canCreate = ref(false);
const canDelete = ref(false);
const canGoods = ref(false);
const canLive = ref(false);
const query = reactive({ limit: 20, page: 1 });
const form = reactive<MerchantBroadcastRoomInput>({
  anchor_name: '', cover_img: '', end_time: '', feeds_img: '', is_show: 1, mark: '', name: '', phone: '', play_url: '', product_ids: [], push_url: '', sort: 0, star: 1, start_time: '',
});

function auditInfo(status: number) {
  if (status === 2) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  return { label: '待平台审核', type: 'warning' as const };
}

function liveText(status: number) {
  return status === 101 ? '直播中' : status === 102 ? '未开始' : '已结束';
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { anchor_name: '', cover_img: '', end_time: '', feeds_img: '', is_show: 1, mark: '', name: '', phone: '', play_url: '', product_ids: [], push_url: '', sort: 0, star: 1, start_time: '' });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
}

async function openEdit(row: MerchantBroadcastRoom) {
  const room = await getMerchantBroadcastRoomApi(row.broadcast_room_id);
  editingID.value = room.broadcast_room_id;
  Object.assign(form, {
    anchor_name: room.anchor_name || '', cover_img: room.cover_img || '', end_time: room.end_time || '', feeds_img: room.feeds_img || '', is_show: room.is_show, mark: room.mark || '', name: room.name, phone: room.phone || '', play_url: room.play_url || '', product_ids: (room.goods || []).map((item) => item.product_id), push_url: room.push_url || '', sort: room.sort || 0, star: room.star || 1, start_time: room.start_time || '',
  });
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantBroadcastRoomsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

async function loadProducts() {
  const result = await listMerchantProductsApi({ limit: 100, page: 1, status: 1 });
  products.value = result.list || [];
}

function dateRangeChanged(value: string[]) {
  form.start_time = value[0] || '';
  form.end_time = value[1] || '';
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写直播间名称');
    return;
  }
  saving.value = true;
  try {
    const body: MerchantBroadcastRoomInput = {
      ...form,
      name: form.name.trim(),
      product_ids: editingID.value && !canGoods.value ? undefined : (form.product_ids || []),
    };
    if (editingID.value) await updateMerchantBroadcastRoomApi(editingID.value, body);
    else await createMerchantBroadcastRoomApi(body);
    dialogOpen.value = false;
    ElMessage.success(editingID.value ? '直播间已更新，已重新提交平台审核' : '直播间已创建，等待平台审核');
    await load();
  } finally {
    saving.value = false;
  }
}

function openGoods(row: MerchantBroadcastRoom) {
  goodsRoom.value = row;
  selectedGoods.value = (row.goods || []).map((item) => item.product_id);
  goodsOpen.value = true;
}

async function saveGoods() {
  if (!goodsRoom.value) return;
  saving.value = true;
  try {
    await setMerchantBroadcastGoodsApi(goodsRoom.value.broadcast_room_id, selectedGoods.value);
    goodsOpen.value = false;
    ElMessage.success('挂货已更新，直播间已重新提交平台审核');
    await load();
  } finally {
    saving.value = false;
  }
}

function openLive(row: MerchantBroadcastRoom) {
  liveRoom.value = row;
  selectedLiveStatus.value = row.live_status;
  liveOpen.value = true;
}

async function saveLive() {
  if (!liveRoom.value) return;
  saving.value = true;
  try {
    await setMerchantBroadcastLiveApi(liveRoom.value.broadcast_room_id, selectedLiveStatus.value);
    liveOpen.value = false;
    ElMessage.success('直播状态已更新，直播间已重新提交平台审核');
    await load();
  } finally {
    saving.value = false;
  }
}

async function remove(row: MerchantBroadcastRoom) {
  try {
    await ElMessageBox.confirm(`确认删除直播间“${row.name}”？删除后不可恢复。`, '删除确认', { cancelButtonText: '取消', confirmButtonText: '确认删除', type: 'warning' });
    await deleteMerchantBroadcastRoomApi(row.broadcast_room_id);
    ElMessage.success('直播间已删除');
    if (rows.value.length === 1 && query.page > 1) query.page -= 1;
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canCreate.value = permissions.includes('broadcast/create');
  canDelete.value = permissions.includes('broadcast/delete');
  canGoods.value = permissions.includes('broadcast/goods');
  canLive.value = permissions.includes('broadcast/live');
  await Promise.all([load(), loadProducts()]);
});
</script>

<template>
  <Page title="直播间管理" description="维护本店直播间和挂货商品；创建、编辑、直播状态及挂货变更均需平台重新审核，审核通过且显示的直播间才会在 C 端展示。">
    <template #extra><el-button v-if="canCreate" type="primary" @click="openCreate">新建直播间</el-button></template>
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="broadcast_room_id">
        <el-table-column label="直播间" min-width="190" prop="name" show-overflow-tooltip />
        <el-table-column label="主播" width="120" prop="anchor_name" />
        <el-table-column label="挂货" width="76"><template #default="{ row }">{{ (row.goods || []).length }}</template></el-table-column>
        <el-table-column label="审核状态" width="118"><template #default="{ row }"><el-tag :type="auditInfo(row.status).type">{{ auditInfo(row.status).label }}</el-tag></template></el-table-column>
        <el-table-column label="直播状态" width="100"><template #default="{ row }">{{ liveText(row.live_status) }}</template></el-table-column>
        <el-table-column label="C 端显示" width="96"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag></template></el-table-column>
        <el-table-column label="驳回原因" min-width="160" prop="refusal" show-overflow-tooltip />
        <el-table-column label="创建时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="250">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-if="canGoods" link type="primary" @click="openGoods(row)">挂货</el-button>
            <el-button v-if="canLive" link type="warning" @click="openLive(row)">直播状态</el-button>
            <el-button v-if="canDelete" link type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>

    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑直播间' : '新建直播间'" destroy-on-close width="760px">
      <el-alert class="mb-4" :closable="false" show-icon title="直播间的编辑或挂货变更会重置为待平台审核；推流地址仅在本商户编辑表单中可见。" type="info" />
      <el-form class="grid grid-cols-2 gap-x-4" label-width="96px">
        <el-form-item class="col-span-2" label="直播间名称" required><el-input v-model="form.name" maxlength="100" /></el-form-item>
        <el-form-item label="主播名称"><el-input v-model="form.anchor_name" /></el-form-item>
        <el-form-item label="主播电话"><el-input v-model="form.phone" /></el-form-item>
        <el-form-item class="col-span-2" label="开播时段"><el-date-picker :model-value="form.start_time && form.end_time ? [form.start_time, form.end_time] : []" class="w-full" end-placeholder="结束时间" start-placeholder="开始时间" type="datetimerange" value-format="YYYY-MM-DD HH:mm:ss" @update:model-value="dateRangeChanged" /></el-form-item>
        <el-form-item label="封面"><ImageField v-model="form.cover_img" button-text="选择封面" /></el-form-item>
        <el-form-item label="分享图"><ImageField v-model="form.feeds_img" button-text="选择分享图" /></el-form-item>
        <el-form-item class="col-span-2" label="播放地址"><el-input v-model="form.play_url" placeholder="可选，供 C 端播放器使用" /></el-form-item>
        <el-form-item class="col-span-2" label="推流地址"><el-input v-model="form.push_url" placeholder="可选，仅商户可见" type="password" show-password /></el-form-item>
        <el-form-item v-if="canGoods" class="col-span-2" label="初始挂货"><el-select v-model="form.product_ids" collapse-tags collapse-tags-tooltip filterable multiple class="w-full" placeholder="仅可选择本店已审核商品"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" class="w-full" /></el-form-item>
        <el-form-item label="C 端显示"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item>
        <el-form-item class="col-span-2" label="备注"><el-input v-model="form.mark" :rows="3" maxlength="500" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">提交审核</el-button></template>
    </el-dialog>

    <el-dialog v-model="goodsOpen" destroy-on-close title="直播挂货" width="620px">
      <el-alert class="mb-4" :closable="false" title="挂货变更会将直播间重新提交平台审核。" type="warning" />
      <el-select v-model="selectedGoods" collapse-tags collapse-tags-tooltip filterable multiple class="w-full" placeholder="选择本店已审核商品"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select>
      <template #footer><el-button @click="goodsOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="saveGoods">保存挂货</el-button></template>
    </el-dialog>

    <el-dialog v-model="liveOpen" destroy-on-close title="更新直播状态" width="420px">
      <el-alert class="mb-4" :closable="false" title="更新直播状态会将直播间重新提交平台审核。" type="warning" />
      <el-radio-group v-model="selectedLiveStatus"><el-radio :value="102">未开始</el-radio><el-radio :value="101">直播中</el-radio><el-radio :value="103">已结束</el-radio></el-radio-group>
      <template #footer><el-button @click="liveOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="saveLive">保存状态</el-button></template>
    </el-dialog>
  </Page>
</template>
