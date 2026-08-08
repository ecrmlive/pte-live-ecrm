<script setup lang="ts">
import { nextTick, ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElPagination,
  ElTable,
  ElTableColumn,
} from 'element-plus';

import {
  fetchPlatformMerchants,
  type PlatformMerchantRow,
} from '#/api/core/ecrm';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

export type PickedStore = {
  mer_id: number;
  mer_name: string;
  real_name: string;
  mer_phone: string;
  mer_avatar?: string;
};

const props = withDefaults(
  defineProps<{
    /** 打开时预勾选（跨页保留） */
    selected?: Array<
      Pick<PickedStore, 'mer_id' | 'mer_name'> &
        Partial<Omit<PickedStore, 'mer_id' | 'mer_name'>>
    >;
    /** 仅展示指定状态店铺；默认 1=开启 */
    status?: number;
  }>(),
  {
    selected: () => [],
    status: 1,
  },
);

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  confirm: [stores: PickedStore[]];
}>();

const keyword = ref('');
const loading = ref(false);
const rows = ref<PlatformMerchantRow[]>([]);
const total = ref(0);
const page = ref(1);
const limit = ref(10);
const picked = ref<PickedStore[]>([]);
const syncingSelection = ref(false);
const tableRef = ref<{
  clearSelection: () => void;
  toggleRowSelection: (row: PlatformMerchantRow, selected?: boolean) => void;
}>();

const [Modal, modalApi] = useVbenModal({
  title: '选择店铺',
  class: 'w-[920px] max-w-[96vw]',
  contentClass: 'store-picker-modal__content !overflow-hidden !p-0',
  confirmText: '确定',
  cancelText: '取消',
  onConfirm: () => {
    emit('confirm', [...picked.value]);
    open.value = false;
  },
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

function toPicked(row: PlatformMerchantRow): PickedStore {
  return {
    mer_id: row.mer_id,
    mer_name: row.mer_name || '',
    real_name: row.real_name || '',
    mer_phone: row.mer_phone || '',
    mer_avatar: row.mer_avatar || '',
  };
}

async function syncTableSelection() {
  const table = tableRef.value;
  if (!table) return;
  syncingSelection.value = true;
  try {
    table.clearSelection();
    const selectedIds = new Set(picked.value.map((item) => item.mer_id));
    for (const row of rows.value) {
      if (selectedIds.has(row.mer_id)) {
        table.toggleRowSelection(row, true);
      }
    }
  } finally {
    await nextTick();
    syncingSelection.value = false;
  }
}

async function loadStores() {
  loading.value = true;
  try {
    const result = await fetchPlatformMerchants({
      page: page.value,
      limit: limit.value,
      keyword: keyword.value.trim() || undefined,
      status: props.status,
    });
    rows.value = result.list || [];
    total.value = result.total || 0;
    await nextTick();
    await syncTableSelection();
  } finally {
    loading.value = false;
  }
}

function onSearch() {
  page.value = 1;
  void loadStores();
}

function onReset() {
  keyword.value = '';
  page.value = 1;
  void loadStores();
}

function onPageChange(next: number) {
  page.value = next;
  void loadStores();
}

function onSelectionChange(selectedRows: PlatformMerchantRow[]) {
  if (syncingSelection.value) return;
  const pageIds = new Set(rows.value.map((row) => row.mer_id));
  const kept = picked.value.filter((item) => !pageIds.has(item.mer_id));
  const added = selectedRows.map(toPicked);
  const map = new Map<number, PickedStore>();
  for (const item of [...kept, ...added]) {
    map.set(item.mer_id, item);
  }
  picked.value = [...map.values()];
}

watch(open, (visible) => {
  if (visible) {
    keyword.value = '';
    page.value = 1;
    picked.value = (props.selected || []).map((item) => ({
      mer_id: item.mer_id,
      mer_name: item.mer_name || '',
      real_name: item.real_name || '',
      mer_phone: item.mer_phone || '',
      mer_avatar: item.mer_avatar || '',
    }));
    modalApi.open();
    void loadStores();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal>
    <div class="store-picker">
      <div class="store-picker__filter">
        <ElForm inline @submit.prevent="onSearch">
          <ElFormItem label="关键字">
            <ElInput
              v-model="keyword"
              clearable
              class="picker-field"
              placeholder="店铺ID / 店铺名称 / 联系人"
              @keyup.enter="onSearch"
            />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="onSearch">搜索</ElButton>
            <ElButton @click="onReset">重置</ElButton>
          </ElFormItem>
        </ElForm>
      </div>
      <div class="store-picker__table">
        <ElTable
          ref="tableRef"
          v-loading="loading"
          :data="rows"
          row-key="mer_id"
          border
          height="100%"
          @selection-change="onSelectionChange"
        >
          <ElTableColumn type="selection" width="48" />
          <ElTableColumn label="店铺ID" prop="mer_id" width="90" />
          <ElTableColumn label="封面" width="80" align="center">
            <template #default="{ row }">
              <ElImage
                v-if="row.mer_avatar"
                :src="resolveCosMediaUrl(row.mer_avatar)"
                fit="cover"
                class="store-cover"
                :preview-src-list="[resolveCosMediaUrl(row.mer_avatar)]"
                preview-teleported
              >
                <template #error>
                  <span class="store-cover-empty">—</span>
                </template>
              </ElImage>
              <span v-else class="store-cover-empty">—</span>
            </template>
          </ElTableColumn>
          <ElTableColumn label="店铺名称" prop="mer_name" min-width="160" />
          <ElTableColumn label="联系人" prop="real_name" width="120" />
          <ElTableColumn label="联系电话" prop="mer_phone" width="140" />
        </ElTable>
      </div>
      <div class="store-picker__pager">
        <ElPagination
          background
          layout="total, prev, pager, next"
          :current-page="page"
          :page-size="limit"
          :total="total"
          @current-change="onPageChange"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.store-picker {
  display: flex;
  flex-direction: column;
  height: min(60vh, 520px);
  min-height: 360px;
  overflow: hidden;
}

.store-picker__filter {
  flex-shrink: 0;
  padding: 12px 12px 0;
}

.store-picker__filter :deep(.el-form--inline) {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  row-gap: 8px;
  margin-bottom: 0;
}

.store-picker__filter :deep(.el-form-item) {
  margin-bottom: 12px;
}

.picker-field {
  width: 280px;
}

.store-picker__table {
  flex: 1;
  min-height: 0;
  padding: 0 12px;
  overflow: hidden;
}

.store-cover {
  width: 40px;
  height: 40px;
  border-radius: 4px;
}

.store-cover-empty {
  color: var(--el-text-color-secondary);
}

.store-picker__pager {
  display: flex;
  flex-shrink: 0;
  justify-content: flex-end;
  padding: 12px;
  border-top: 1px solid hsl(var(--border));
}
</style>

<style>
.store-picker-modal__content {
  display: flex;
  flex-direction: column;
}
</style>
