<script lang="ts" setup>
import { computed } from 'vue';

import {
  MERCHANT_LIST_PAGE_SIZES,
  MERCHANT_LIST_PAGER_DEFAULTS,
  MERCHANT_LIST_PAGER_LAYOUTS,
} from '#/constants/merchant-list-pager';

defineOptions({ name: 'MerchantListPagination' });

const props = withDefaults(
  defineProps<{
    currentPage: number;
    pageSize: number;
    total: number;
    pageSizes?: number[];
  }>(),
  {
    pageSizes: () => [...MERCHANT_LIST_PAGE_SIZES],
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
  <!-- VxePager 由 #/adapter/vxe-table 全局注册 -->
  <vxe-pager
    v-model:current-page="innerCurrentPage"
    v-model:page-size="innerPageSize"
    :background="MERCHANT_LIST_PAGER_DEFAULTS.background"
    :layouts="[...MERCHANT_LIST_PAGER_LAYOUTS]"
    :page-sizes="pageSizes"
    :size="MERCHANT_LIST_PAGER_DEFAULTS.size"
    :total="total"
    class="merchant-list-pager"
    @page-change="onPageChange"
  />
</template>
