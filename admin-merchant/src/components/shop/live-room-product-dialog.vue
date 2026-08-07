<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { LiveRoomProductItem } from '#/api/core/live';

import { useVbenDrawer } from '@vben/common-ui';
import { Plus } from '@element-plus/icons-vue';
import { ElButton, ElMessage, ElMessageBox, ElSwitch } from 'element-plus';
import { reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  addLiveRoomProductsApi,
  deleteLiveRoomProductApi,
  getLiveRoomProductListApi,
  pushLiveRoomProductApi,
  setLiveRoomProductOnSaleApi,
} from '#/api/core/live';

import LiveWxProductImportDialog from '#/components/shop/live-wx-product-import-dialog.vue';

defineOptions({ name: 'LiveRoomProductDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  liveId: number;
  roomId: number | string;
}>();

const emit = defineEmits<{
  success: [];
}>();

const importOpen = ref(false);

const gridOptions = reactive<VxeGridProps<LiveRoomProductItem>>({
  columns: [
    {
      field: 'cover_img',
      slots: { default: 'image' },
      title: '商品图片',
      width: 70,
    },
    {
      field: 'name',
      minWidth: 180,
      slots: { default: 'name' },
      title: '商品名称',
    },
    {
      field: 'on_sale',
      slots: { default: 'on_sale' },
      title: '上架状态',
      width: 140,
    },
    {
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      title: '推送状态',
      width: 180,
    },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.roomId) {
          return { items: [], total: 0 };
        }
        const res = await getLiveRoomProductListApi({
          list_rows: page.pageSize,
          live_id: props.liveId,
          page: page.currentPage,
          room_id: props.roomId,
        });
        const list = Array.isArray(res.list) ? res.list : [];
        return {
          items: list,
          total: Number(res.total) || 0,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'live_product_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

async function onImport(goodsIds: Array<number | string>) {
  await addLiveRoomProductsApi({
    live_id: props.liveId,
    productIds: goodsIds,
    room_id: props.roomId,
  });
  ElMessage.success('导入成功');
  await gridApi.reload();
  emit('success');
}

async function onSaleChange(value: boolean | number | string, row: LiveRoomProductItem) {
  await setLiveRoomProductOnSaleApi({
    live_product_id: row.live_product_id,
    on_sale: Number(value),
    room_id: props.roomId,
  });
  ElMessage.success('操作成功');
  await gridApi.reload();
}

async function pushProduct(row: LiveRoomProductItem) {
  await pushLiveRoomProductApi({
    live_product_id: row.live_product_id,
    room_id: props.roomId,
  });
  ElMessage.success('推送成功');
  await gridApi.reload();
}

async function deleteRow(row: LiveRoomProductItem) {
  await ElMessageBox.confirm('删除后不可恢复，确认删除该记录吗?', '提示', { type: 'warning' });
  await deleteLiveRoomProductApi(row.live_product_id);
  ElMessage.success('删除成功');
  await gridApi.reload();
  emit('success');
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    void gridApi.reload();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[900px]"
    title="选择商品"
  >
    <div class="mb-3">
      <ElButton :icon="Plus" size="small" type="primary" @click="importOpen = true">
        导入商品
      </ElButton>
    </div>

    <Grid>
      <template #image="{ row }">
        <img
          v-if="row.product?.cover_img"
          :src="row.product.cover_img"
          alt=""
          class="h-[30px] w-[30px] object-cover"
        />
      </template>
      <template #name="{ row }">
        {{ row.product?.name || row.name }}
      </template>
      <template #on_sale="{ row }">
        <ElSwitch
          :active-value="1"
          :inactive-value="0"
          :model-value="row.on_sale"
          @change="(val) => onSaleChange(val, row)"
        />
      </template>
      <template #action="{ row }">
        <ElButton
          :disabled="row.isPush === 0"
          size="small"
          type="primary"
          @click="pushProduct(row)"
        >
          推送
        </ElButton>
        <ElButton size="small" type="primary" @click="deleteRow(row)">删除</ElButton>
      </template>
    </Grid>

    <LiveWxProductImportDialog
      v-model:open="importOpen"
      :room-id="roomId"
      @confirm="onImport"
    />
  </Modal>
</template>
