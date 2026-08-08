<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import PlugsApi from '#/api/core/plugs';

import type { PlugCandidate, PlugCategory } from './types';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  category: null | PlugCategory;
}>();

const emit = defineEmits<{
  success: [];
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

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen && props.category) {
      void loadCandidates();
    }
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      drawerApi.setState({ title: '新增插件' }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(
  () => props.category?.plus_category_id,
  () => {
    if (open.value && props.category) {
      void loadCandidates();
    }
  },
);

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
      ElMessage.success('新增成功');
      open.value = false;
      emit('success');
    }
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :footer="false"
    title="新增插件"
  >
    <Grid>
      <template #action="{ row }">
        <ElButton
          :loading="submitting"
          link
          type="primary"
          @click="handleAdd(row)"
        >
          新增
        </ElButton>
      </template>
    </Grid>
  </Drawer>
</template>
