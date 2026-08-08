<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';
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
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import ImageField from '#/components/shop/image-field.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const products = ref<MerchantProduct[]>([]);
const goodsRoom = ref<MerchantBroadcastRoom>();
const liveRoom = ref<MerchantBroadcastRoom>();
const selectedGoods = ref<number[]>([]);
const selectedLiveStatus = ref(102);
const canCreate = ref(false);
const canDelete = ref(false);
const canGoods = ref(false);
const canLive = ref(false);
const form = reactive<MerchantBroadcastRoomInput>({
  anchor_name: '',
  cover_img: '',
  end_time: '',
  feeds_img: '',
  is_show: 1,
  mark: '',
  name: '',
  phone: '',
  play_url: '',
  product_ids: [],
  push_url: '',
  sort: 0,
  star: 1,
  start_time: '',
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
  Object.assign(form, {
    anchor_name: '',
    cover_img: '',
    end_time: '',
    feeds_img: '',
    is_show: 1,
    mark: '',
    name: '',
    phone: '',
    play_url: '',
    product_ids: [],
    push_url: '',
    sort: 0,
    star: 1,
    start_time: '',
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '直播间名称 / 主播' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待平台审核', value: 0 },
        { label: '审核通过', value: 2 },
        { label: '已驳回', value: -1 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<MerchantBroadcastRoom> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'broadcast_room_id', title: 'ID', width: 80 },
    { field: 'name', minWidth: 190, showOverflow: false, title: '直播间' },
    { field: 'anchor_name', title: '主播', width: 120 },
    {
      field: 'goods',
      title: '挂货',
      width: 76,
      formatter: ({ cellValue }) => (cellValue || []).length,
    },
    {
      field: 'status',
      slots: { default: 'audit' },
      title: '审核状态',
      width: 118,
    },
    {
      field: 'live_status',
      title: '直播状态',
      width: 100,
      formatter: ({ cellValue }) => liveText(Number(cellValue)),
    },
    {
      field: 'is_show',
      slots: { default: 'show' },
      title: 'C 端显示',
      width: 96,
    },
    { field: 'refusal', minWidth: 160, showOverflow: false, title: '驳回原因' },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 250 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const data = await listMerchantBroadcastRoomsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            status === 0 || status === 2 || status === -1
              ? Number(status)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'broadcast_room_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: saveRoom,
});

const [GoodsDrawer, goodsDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: saveGoods,
});

const [LiveDrawer, liveDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: saveLive,
});

function openCreate() {
  resetForm();
  editDrawerApi.setState({ title: '新建直播间' }).open();
}

async function openEdit(row: MerchantBroadcastRoom) {
  const room = await getMerchantBroadcastRoomApi(row.broadcast_room_id);
  editingID.value = room.broadcast_room_id;
  Object.assign(form, {
    anchor_name: room.anchor_name || '',
    cover_img: room.cover_img || '',
    end_time: room.end_time || '',
    feeds_img: room.feeds_img || '',
    is_show: room.is_show,
    mark: room.mark || '',
    name: room.name,
    phone: room.phone || '',
    play_url: room.play_url || '',
    product_ids: (room.goods || []).map((item) => item.product_id),
    push_url: room.push_url || '',
    sort: room.sort || 0,
    star: room.star || 1,
    start_time: room.start_time || '',
  });
  editDrawerApi.setState({ title: '编辑直播间' }).open();
}

function dateRangeChanged(value: string[]) {
  form.start_time = value[0] || '';
  form.end_time = value[1] || '';
}

async function saveRoom() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写直播间名称');
    return;
  }
  saving.value = true;
  editDrawerApi.lock();
  try {
    const body: MerchantBroadcastRoomInput = {
      ...form,
      name: form.name.trim(),
      product_ids:
        editingID.value && !canGoods.value ? undefined : form.product_ids || [],
    };
    if (editingID.value) {
      await updateMerchantBroadcastRoomApi(editingID.value, body);
    } else {
      await createMerchantBroadcastRoomApi(body);
    }
    editDrawerApi.close();
    ElMessage.success(
      editingID.value
        ? '直播间已更新，已重新提交平台审核'
        : '直播间已创建，等待平台审核',
    );
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

function openGoods(row: MerchantBroadcastRoom) {
  goodsRoom.value = row;
  selectedGoods.value = (row.goods || []).map((item) => item.product_id);
  goodsDrawerApi.setState({ title: '直播挂货' }).open();
}

async function saveGoods() {
  if (!goodsRoom.value) return;
  saving.value = true;
  goodsDrawerApi.lock();
  try {
    await setMerchantBroadcastGoodsApi(
      goodsRoom.value.broadcast_room_id,
      selectedGoods.value,
    );
    goodsDrawerApi.close();
    ElMessage.success('挂货已更新，直播间已重新提交平台审核');
    gridApi.reload();
  } finally {
    saving.value = false;
    goodsDrawerApi.unlock();
  }
}

function openLive(row: MerchantBroadcastRoom) {
  liveRoom.value = row;
  selectedLiveStatus.value = row.live_status;
  liveDrawerApi.setState({ title: '更新直播状态' }).open();
}

async function saveLive() {
  if (!liveRoom.value) return;
  saving.value = true;
  liveDrawerApi.lock();
  try {
    await setMerchantBroadcastLiveApi(
      liveRoom.value.broadcast_room_id,
      selectedLiveStatus.value,
    );
    liveDrawerApi.close();
    ElMessage.success('直播状态已更新，直播间已重新提交平台审核');
    gridApi.reload();
  } finally {
    saving.value = false;
    liveDrawerApi.unlock();
  }
}

async function remove(row: MerchantBroadcastRoom) {
  try {
    await confirm({
      content: `确认删除直播间“${row.name}”？删除后不可恢复。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantBroadcastRoomApi(row.broadcast_room_id);
    ElMessage.success('直播间已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canCreate.value = permissions.includes('broadcast/create');
  canDelete.value = permissions.includes('broadcast/delete');
  canGoods.value = permissions.includes('broadcast/goods');
  canLive.value = permissions.includes('broadcast/live');
  const result = await listMerchantProductsApi({ limit: 100, page: 1, status: 1 });
  products.value = result.list || [];
});
</script>

<template>
  <Page auto-content-height>
    <template v-if="canCreate" #extra>
      <ElButton type="primary" @click="openCreate">新建直播间</ElButton>
    </template>

    <Grid>
      <template #audit="{ row }">
        <ElTag :type="auditInfo(row.status).type">
          {{ auditInfo(row.status).label }}
        </ElTag>
      </template>
      <template #show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton v-if="canGoods" link type="primary" @click="openGoods(row)">
          挂货
        </ElButton>
        <ElButton v-if="canLive" link type="warning" @click="openLive(row)">
          直播状态
        </ElButton>
        <ElButton v-if="canDelete" link type="danger" @click="remove(row)">
          删除
        </ElButton>
      </template>
    </Grid>

    <EditDrawer class="w-[760px] max-w-[96vw]">
      <ElAlert
        class="mb-4"
        :closable="false"
        show-icon
        title="直播间的编辑或挂货变更会重置为待平台审核；推流地址仅在本商户编辑表单中可见。"
        type="info"
      />
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="96px">
        <ElFormItem class="col-span-2" label="直播间名称" required>
          <ElInput v-model="form.name" maxlength="100" />
        </ElFormItem>
        <ElFormItem label="主播名称">
          <ElInput v-model="form.anchor_name" />
        </ElFormItem>
        <ElFormItem label="主播电话">
          <ElInput v-model="form.phone" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="开播时段">
          <ElDatePicker
            :model-value="
              form.start_time && form.end_time
                ? [form.start_time, form.end_time]
                : []
            "
            class="w-full"
            end-placeholder="结束时间"
            start-placeholder="开始时间"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
            @update:model-value="dateRangeChanged"
          />
        </ElFormItem>
        <ElFormItem label="封面">
          <ImageField v-model="form.cover_img" :preview-size="96" />
        </ElFormItem>
        <ElFormItem label="分享图">
          <ImageField v-model="form.feeds_img" :preview-size="96" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="播放地址">
          <ElInput v-model="form.play_url" placeholder="可选，供 C 端播放器使用" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="推流地址">
          <ElInput
            v-model="form.push_url"
            placeholder="可选，仅商户可见"
            type="password"
            show-password
          />
        </ElFormItem>
        <ElFormItem v-if="canGoods" class="col-span-2" label="初始挂货">
          <ElSelect
            v-model="form.product_ids"
            collapse-tags
            collapse-tags-tooltip
            filterable
            multiple
            class="w-full"
            placeholder="仅可选择本店已审核商品"
          >
            <ElOption
              v-for="item in products"
              :key="item.product_id"
              :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`"
              :value="item.product_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="C 端显示">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="备注">
          <ElInput
            v-model="form.mark"
            :rows="3"
            maxlength="500"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </EditDrawer>

    <GoodsDrawer class="w-[620px] max-w-[96vw]">
      <ElAlert
        class="mb-4"
        :closable="false"
        title="挂货变更会将直播间重新提交平台审核。"
        type="warning"
      />
      <ElSelect
        v-model="selectedGoods"
        collapse-tags
        collapse-tags-tooltip
        filterable
        multiple
        class="w-full"
        placeholder="选择本店已审核商品"
      >
        <ElOption
          v-for="item in products"
          :key="item.product_id"
          :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`"
          :value="item.product_id"
        />
      </ElSelect>
    </GoodsDrawer>

    <LiveDrawer class="w-[420px] max-w-[96vw]">
      <ElAlert
        class="mb-4"
        :closable="false"
        title="更新直播状态会将直播间重新提交平台审核。"
        type="warning"
      />
      <ElRadioGroup v-model="selectedLiveStatus">
        <ElRadio :value="102">未开始</ElRadio>
        <ElRadio :value="101">直播中</ElRadio>
        <ElRadio :value="103">已结束</ElRadio>
      </ElRadioGroup>
    </LiveDrawer>
  </Page>
</template>
