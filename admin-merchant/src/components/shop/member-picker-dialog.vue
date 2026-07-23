<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { MemberListItem } from '#/api/core/member';

import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { computed, reactive, ref, shallowRef, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getMemberListApi } from '#/api/core/member';

type MemberRow = MemberListItem & { noChoose?: boolean };

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
    multiple?: boolean;
  }>(),
  { excludeIds: () => [], multiple: false },
);

const emit = defineEmits<{
  confirm: [MemberListItem[]];
  select: [MemberListItem];
}>();

const gradeOptions = shallowRef<Array<{ label: string; value: string }>>([
  { label: '全部', value: '0' },
]);

const selectedRows = ref<MemberRow[]>([]);

const canConfirm = computed(() =>
  props.multiple ? selectedRows.value.length > 0 : selectedRows.value.length === 1,
);

function isSelectable(row: MemberRow) {
  if (typeof row.noChoose === 'boolean') {
    return row.noChoose;
  }
  return true;
}

const formOptions = computed((): VbenFormProps => ({
  actionLayout: 'inline',
  collapsed: false,
  schema: [
    {
      component: 'Select',
      componentProps: {
        options: gradeOptions.value,
        style: 'width: 120px',
      },
      defaultValue: '0',
      fieldName: 'grade_id',
      label: '等级',
    },
    {
      component: 'Input',
      componentProps: { placeholder: '昵称|手机号|ID' },
      fieldName: 'nick_name',
      label: '关键词',
    },
  ],
  showCollapseButton: false,
  submitButtonOptions: {
    content: '查询',
  },
  submitOnChange: false,
  submitOnEnter: true,
  wrapperClass: 'grid-cols-1 md:grid-cols-3',
}));

const gridOptions = reactive<VxeGridProps<MemberRow>>({
  checkboxConfig: props.multiple
    ? {
        checkMethod: ({ row }) => isSelectable(row),
        highlight: true,
        reserve: true,
      }
    : undefined,
  columns: [
    ...(props.multiple
      ? [{ type: 'checkbox' as const, width: 45 }]
      : [{ type: 'radio' as const, width: 45 }]),
    { field: 'user_id', title: 'ID', width: 70 },
    {
      field: 'avatarUrl',
      slots: { default: 'avatar' },
      title: '头像',
      width: 70,
    },
    { field: 'nickName', minWidth: 120, showOverflow: true, title: '昵称' },
    { field: 'mobile', title: '手机号', width: 120 },
    {
      field: 'balance',
      slots: { default: 'balance' },
      title: '用户余额',
      width: 100,
    },
    {
      field: 'grade.name',
      slots: { default: 'grade' },
      title: '会员等级',
      width: 100,
    },
    {
      field: 'pay_money',
      slots: { default: 'payMoney' },
      title: '消费金额',
      width: 100,
    },
    { field: 'create_time', title: '注册时间', width: 140 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 15,
    pageSizes: [15, 30, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getMemberListApi({
          grade_id: String(formValues?.grade_id ?? '0'),
          list_rows: page.pageSize,
          nick_name: String(formValues?.nick_name ?? ''),
          page: page.currentPage,
        });
        gradeOptions.value = [
          { label: '全部', value: '0' },
          ...(res.grade ?? []).map((item) => ({
            label: item.name,
            value: String(item.grade_id),
          })),
        ];
        const rows = (res.list.data ?? []).map((row) => ({
          ...row,
          noChoose: props.excludeIds.length
            ? !props.excludeIds.includes(row.user_id)
            : true,
        }));
        return {
          items: rows,
          total: res.list.total ?? 0,
        };
      },
    },
  },
  radioConfig: props.multiple
    ? undefined
    : {
        checkMethod: ({ row }) => isSelectable(row),
        highlight: true,
        strict: true,
        trigger: 'row',
      },
  rowConfig: {
    keyField: 'user_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents: {
    checkboxChange: onCheckboxChange,
    radioChange: onRadioChange,
  },
  gridOptions,
});

function onCheckboxChange() {
  selectedRows.value = (gridApi.grid?.getCheckboxRecords?.() ?? []) as MemberRow[];
}

function onRadioChange({ row }: { row: MemberRow }) {
  selectedRows.value = row ? [row] : [];
}

function syncSelectionFromGrid() {
  if (props.multiple) {
    selectedRows.value = (gridApi.grid?.getCheckboxRecords?.() ?? []) as MemberRow[];
    return;
  }
  const row = gridApi.grid?.getRadioRecord?.() as MemberRow | undefined;
  selectedRows.value = row ? [row] : [];
}

function confirm() {
  syncSelectionFromGrid();
  if (props.multiple) {
    if (!selectedRows.value.length) return;
    emit('confirm', [...selectedRows.value]);
  } else if (selectedRows.value[0]) {
    emit('select', selectedRows.value[0]);
  } else {
    return;
  }
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    selectedRows.value = [];
    gradeOptions.value = [{ label: '全部', value: '0' }];
    void gridApi.reload();
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
    class="w-[860px]"
    title="选择用户"
  >
    <Grid>
      <template #avatar="{ row }">
        <img
          v-if="row.avatarUrl"
          :alt="row.nickName"
          class="h-8 w-8 rounded-full object-cover"
          :src="row.avatarUrl"
        />
      </template>
      <template #balance="{ row }">
        <span class="text-orange-500">￥{{ row.balance }}</span>
      </template>
      <template #grade="{ row }">
        {{ row.grade?.name ?? '-' }}
      </template>
      <template #payMoney="{ row }">
        {{ row.pay_money ?? '0' }}
      </template>
    </Grid>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>
