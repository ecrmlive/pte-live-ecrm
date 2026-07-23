<script lang="ts" setup>
import { computed } from 'vue';

import { VxePager } from 'vxe-pc-ui';

import {
  PLATFORM_LIST_PAGE_SIZES,
  PLATFORM_LIST_PAGER_DEFAULTS,
  PLATFORM_LIST_PAGER_LAYOUTS,
} from '#/constants/platform-list-pager';

defineOptions({ name: 'PlatformListPager' });

const props = withDefaults(
  defineProps<{
    currentPage: number;
    pageSize: number;
    total: number;
    pageSizes?: number[];
  }>(),
  {
    pageSizes: () => [...PLATFORM_LIST_PAGE_SIZES],
  },
);

const emit = defineEmits<{
  change: [];
  sizeChange: [size: number];
  'update:currentPage': [page: number];
  'update:pageSize': [size: number];
}>();

const innerCurrentPage = computed({
  get: () => props.currentPage,
  set: (page: number) => emit('update:currentPage', page),
});

const innerPageSize = computed({
  get: () => props.pageSize,
  set: (size: number) => emit('update:pageSize', size),
});

function onPageChange({
  currentPage,
  pageSize,
}: {
  currentPage: number;
  pageSize: number;
}) {
  const pageChanged = currentPage !== props.currentPage;
  const sizeChanged = pageSize !== props.pageSize;

  if (pageChanged) {
    emit('update:currentPage', currentPage);
  }
  if (sizeChanged) {
    emit('update:pageSize', pageSize);
    emit('sizeChange', pageSize);
  }
  if (pageChanged || sizeChanged) {
    emit('change');
  }
}
</script>

<template>
  <!-- 独立 VxePager（非 ElPagination）；单页也展示 Total/Sizes -->
  <VxePager
    v-model:current-page="innerCurrentPage"
    v-model:page-size="innerPageSize"
    :auto-hidden="false"
    :background="PLATFORM_LIST_PAGER_DEFAULTS.background"
    :layouts="[...PLATFORM_LIST_PAGER_LAYOUTS]"
    :page-sizes="pageSizes"
    :size="PLATFORM_LIST_PAGER_DEFAULTS.size"
    :total="total"
    class="platform-list-pager"
    @page-change="onPageChange"
  />
</template>
