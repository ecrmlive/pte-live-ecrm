<script setup lang="ts">
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { computed, reactive, ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElImage,
  ElInput,
  ElMessage,
  ElOption,
  ElPagination,
  ElSelect,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { listPlatformProductsApi } from '#/api/core/platform-catalog';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  select: [PlatformProduct[]];
}>();

const props = defineProps<{
  excludeIds?: number[];
  categoryOptions?: { label: string; value: number }[];
}>();

const loading = ref(false);
const rows = ref<PlatformProduct[]>([]);
const total = ref(0);
const selected = ref<PlatformProduct[]>([]);
const tableRef = ref<{ clearSelection: () => void } | null>(null);

const filter = reactive({
  keyword: '',
  cate_id: undefined as number | undefined,
  mer_cate_name: '',
  brand_name: '',
  type: 1 as number,
  store_name: '',
  page: 1,
  limit: 10,
});

const canConfirm = computed(() => selected.value.length > 0);

function coverOf(row: PlatformProduct) {
  return resolveCosMediaUrl(String(row.image || '').trim());
}

function statusLabel(row: PlatformProduct) {
  if (Number(row.status) === 1) return '上架';
  if (Number(row.status) === 0) return '下架';
  return '—';
}

async function load() {
  loading.value = true;
  try {
    const res = await listPlatformProductsApi({
      page: filter.page,
      limit: filter.limit,
      keyword: filter.keyword.trim() || undefined,
      cate_id: filter.cate_id || undefined,
      store_name: filter.store_name.trim() || undefined,
      brand_name: filter.brand_name.trim() || undefined,
      type: filter.type || 1,
      is_trader: 1,
    });
    let list = res.list || [];
    const merCate = filter.mer_cate_name.trim();
    if (merCate) {
      list = list.filter((p) =>
        String(p.mer_cate_name || '').includes(merCate),
      );
    }
    const exclude = new Set(props.excludeIds || []);
    if (exclude.size) {
      list = list.filter((p) => !exclude.has(p.product_id));
    }
    rows.value = list;
    total.value = res.total || 0;
  } catch {
    ElMessage.error('加载商品失败');
    rows.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  filter.page = 1;
  void load();
}

function resetFilter() {
  filter.keyword = '';
  filter.cate_id = undefined;
  filter.mer_cate_name = '';
  filter.brand_name = '';
  filter.type = 1;
  filter.store_name = '';
  filter.page = 1;
  void load();
}

function onSelectionChange(list: PlatformProduct[]) {
  selected.value = list;
}

function confirm() {
  if (!selected.value.length) {
    ElMessage.warning('请选择商品');
    return;
  }
  emit('select', [...selected.value]);
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    selected.value = [];
    filter.page = 1;
    void load();
    modalApi.open();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[960px] max-w-[96vw]"
    title="商品信息"
  >
    <ElAlert
      type="error"
      :closable="false"
      show-icon
      class="picker-alert"
      title="平台添加秒杀商品，此处仅展示所有自营店铺商品"
    />

    <div class="picker-filter">
      <ElInput
        v-model="filter.keyword"
        clearable
        placeholder="商品搜索"
        class="picker-filter__item"
        @keyup.enter="search"
      />
      <ElSelect
        v-model="filter.cate_id"
        clearable
        filterable
        placeholder="平台分类"
        class="picker-filter__item"
      >
        <ElOption
          v-for="opt in categoryOptions || []"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </ElSelect>
      <ElInput
        v-model="filter.mer_cate_name"
        clearable
        placeholder="店铺分类"
        class="picker-filter__item"
        @keyup.enter="search"
      />
      <ElInput
        v-model="filter.brand_name"
        clearable
        placeholder="商品标签"
        class="picker-filter__item"
        @keyup.enter="search"
      />
      <ElSelect
        v-model="filter.type"
        clearable
        placeholder="商品状态"
        class="picker-filter__item"
      >
        <ElOption label="出售中" :value="1" />
        <ElOption label="仓库中" :value="2" />
      </ElSelect>
      <ElInput
        v-model="filter.store_name"
        clearable
        placeholder="店铺名称"
        class="picker-filter__item"
        @keyup.enter="search"
      />
      <ElButton @click="resetFilter">重置</ElButton>
      <ElButton type="primary" @click="search">搜索</ElButton>
    </div>

    <ElTable
      ref="tableRef"
      v-loading="loading"
      :data="rows"
      row-key="product_id"
      @selection-change="onSelectionChange"
    >
      <ElTableColumn type="selection" width="48" />
      <ElTableColumn prop="product_id" label="ID" width="80" />
      <ElTableColumn label="商品图" width="80">
        <template #default="{ row }">
          <ElImage
            v-if="coverOf(row)"
            :src="coverOf(row)"
            fit="cover"
            class="picker-cover"
          />
          <span v-else class="text-muted">—</span>
        </template>
      </ElTableColumn>
      <ElTableColumn
        prop="store_name"
        label="名称"
        min-width="180"
        show-overflow-tooltip
      />
      <ElTableColumn
        prop="cate_name"
        label="分类"
        min-width="100"
        show-overflow-tooltip
      />
      <ElTableColumn
        prop="mer_name"
        label="店铺"
        min-width="120"
        show-overflow-tooltip
      />
      <ElTableColumn prop="stock" label="库存" width="80" />
      <ElTableColumn label="状态" width="90">
        <template #default="{ row }">
          <ElTag size="small" :type="row.status === 1 ? 'success' : 'info'">
            {{ statusLabel(row) }}
          </ElTag>
        </template>
      </ElTableColumn>
    </ElTable>

    <div class="picker-pager">
      <ElPagination
        v-model:current-page="filter.page"
        v-model:page-size="filter.limit"
        background
        layout="total, prev, pager, next"
        :total="total"
        @current-change="load"
        @size-change="
          () => {
            filter.page = 1;
            load();
          }
        "
      />
    </div>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">
        确定{{ selected.length ? `（${selected.length}）` : '' }}
      </ElButton>
    </template>
  </Modal>
</template>

<style scoped>
.picker-alert {
  margin-bottom: 12px;
}

.picker-filter {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
  align-items: center;
}

.picker-filter__item {
  width: 150px;
}

.picker-cover {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  display: block;
}

.picker-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.text-muted {
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}
</style>
