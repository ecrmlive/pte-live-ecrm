<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref, watch } from 'vue';

import { ElButton, ElDialog, ElMessage } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import PlugsApi from '#/api/core/plugs';

import type { PlugCandidate, PlugCategory } from './types';

const props = defineProps<{
  category: null | PlugCategory;
  open: boolean;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const submitting = ref(false);
const candidates = ref<PlugCandidate[]>([]);

const gridOptions: VxeGridProps<PlugCandidate> = {
  border: true,
  columns: [
    { field: 'name', minWidth: 200, showOverflow: true, title: '名称' },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 88,
    },
  ],
  height: 'auto',
  maxHeight: 360,
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async () => ({
        items: candidates.value,
        total: candidates.value.length,
      }),
    },
  },
  rowConfig: { isHover: true, keyField: 'access_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

watch(
  () => [props.open, props.category?.plus_category_id] as const,
  ([open]) => {
    if (open && props.category) {
      void loadCandidates();
    }
  },
  { immediate: true },
);

async function loadCandidates() {
  if (!props.category) return;
  gridApi.setLoading(true);
  try {
    const res = await PlugsApi.getplugs(
      { plus_category_id: props.category.plus_category_id },
      true,
    );
    candidates.value =
      (res.data as { accessList?: PlugCandidate[] })?.accessList ?? [];
    await gridApi.reload();
  } catch {
    candidates.value = [];
    await gridApi.reload();
  } finally {
    gridApi.setLoading(false);
  }
}

async function handleAdd(row: PlugCandidate) {
  if (!props.category || submitting.value) return;
  submitting.value = true;
  try {
    const res = await PlugsApi.addplugs(
      {
        access_id: row.access_id,
        plus_category_id: props.category.plus_category_id,
      },
      true,
    );
    if (res.code === 1) {
      ElMessage.success('添加成功');
      emit('update:open', false);
      emit('success');
    }
  } finally {
    submitting.value = false;
  }
}

function handleClose() {
  emit('update:open', false);
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :model-value="open"
    title="添加插件"
    width="520px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
  >
    <Grid>
      <template #action="{ row }">
        <ElButton
          :loading="submitting"
          link
          type="primary"
          @click="handleAdd(row)"
        >
          添加
        </ElButton>
      </template>
    </Grid>
  </ElDialog>
</template>
