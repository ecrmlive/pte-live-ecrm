<script setup lang="ts">
import { computed, ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElPagination,
  ElRadio,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import {
  listPlatformProductStoresApi,
  type PlatformProductStoreOption,
} from '#/api/core/platform-catalog';

const modelValue = defineModel<number | undefined>({ default: undefined });

const props = withDefaults(
  defineProps<{ placeholder?: string }>(),
  { placeholder: '请选择店铺' },
);

const keyword = ref('');
const loading = ref(false);
const options = ref<PlatformProductStoreOption[]>([]);
const page = ref(1);
const pageSize = 10;
const selected = ref<PlatformProductStoreOption>();
const selectedLabel = ref('');

const filteredRows = computed(() => {
  const value = keyword.value.trim().toLowerCase();
  if (!value) return options.value;
  return options.value.filter((store) =>
    `${store.store_name} ${store.merchant_name}`.toLowerCase().includes(value),
  );
});
const rows = computed(() => {
  const start = (page.value - 1) * pageSize;
  return filteredRows.value.slice(start, start + pageSize);
});

const [Modal, modalApi] = useVbenModal({
  title: '请选择店铺：',
  class: 'h-[min(76dvh,820px)] w-[min(94vw,1200px)] max-w-[94vw]',
  cancelText: '关闭',
  confirmText: '确定',
  onConfirm: () => {
    if (!selected.value) {
      ElMessage.warning('请选择店铺');
      return;
    }
    modelValue.value = selected.value.store_id;
    selectedLabel.value = `${selected.value.store_name}（${selected.value.merchant_name}）`;
    modalApi.close();
  },
});

async function load() {
  loading.value = true;
  try {
    const result = await listPlatformProductStoresApi();
    options.value = result.list || [];
  } finally {
    loading.value = false;
  }
}

async function openPicker() {
  page.value = 1;
  keyword.value = '';
  selected.value = options.value.find((store) => store.store_id === modelValue.value);
  modalApi.open();
  if (!options.value.length) await load();
}

function search() {
  page.value = 1;
}

function reset() {
  keyword.value = '';
  page.value = 1;
}

function clear() {
  modelValue.value = undefined;
  selected.value = undefined;
  selectedLabel.value = '';
}
</script>

<template>
  <div class="store-relation-select">
    <ElButton plain type="primary" @click="openPicker">选择店铺</ElButton>
    <span v-if="selectedLabel" class="store-relation-select__summary">{{ selectedLabel }}</span>
    <ElTag v-if="modelValue && !selectedLabel" closable @close="clear">已关联店铺</ElTag>
    <ElButton v-if="modelValue" link type="primary" @click="clear">清除</ElButton>
  </div>

  <Modal>
    <div class="relation-picker">
      <ElForm inline class="relation-picker__search" @submit.prevent="search">
        <ElFormItem label="店铺搜索：">
          <ElInput
            v-model="keyword"
            class="relation-picker__keyword"
            clearable
            :placeholder="props.placeholder"
            @keyup.enter="search"
          />
        </ElFormItem>
        <ElFormItem>
          <ElButton @click="reset">重置</ElButton>
          <ElButton type="primary" @click="search">搜索</ElButton>
        </ElFormItem>
      </ElForm>
      <div class="relation-picker__table">
        <ElTable v-loading="loading" :data="rows" height="100%" @row-click="(row) => (selected = row)">
          <ElTableColumn align="center" width="64">
            <template #default="{ row }">
              <ElRadio :model-value="selected?.store_id" :value="row.store_id" @change="selected = row" @click.stop>&nbsp;</ElRadio>
            </template>
          </ElTableColumn>
          <ElTableColumn label="店铺名称" min-width="280" prop="store_name" />
          <ElTableColumn label="所属商户" min-width="280" prop="merchant_name" />
        </ElTable>
      </div>
      <div class="relation-picker__pager">
        <ElPagination background layout="total, prev, pager, next" :current-page="page" :page-size="pageSize" :total="filteredRows.length" @current-change="(next) => (page = next)" />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.store-relation-select { display: flex; flex-wrap: wrap; gap: 8px; align-items: center; width: 100%; }
.store-relation-select__summary { color: var(--el-text-color-secondary); font-size: 13px; }
.relation-picker { display: flex; flex-direction: column; height: 100%; min-height: 0; overflow: hidden; }
.relation-picker__search, .relation-picker__pager { flex: 0 0 auto; }
.relation-picker__table { flex: 1 1 auto; min-height: 0; overflow: hidden; }
.relation-picker__keyword { width: 320px; }
.relation-picker__pager { display: flex; justify-content: flex-end; padding: 16px 0; }
</style>
