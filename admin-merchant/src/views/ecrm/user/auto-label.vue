<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listStoreAutoLabelRulesApi,
  saveStoreAutoLabelRulesApi,
  type StoreAutoLabelRule,
} from '#/api/core/merchant-user-label';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const rows = ref<StoreAutoLabelRule[]>([]);
const loaded = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '规则名称' },
    fieldName: 'keyword',
    label: '规则搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已启用', value: 1 },
        { label: '已停用', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<StoreAutoLabelRule> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    {
      field: 'name',
      minWidth: 180,
      showOverflow: false,
      slots: { default: 'name' },
      title: '规则名称',
    },
    {
      field: 'rule_type',
      title: '规则类型',
      width: 180,
      slots: { default: 'type' },
    },
    {
      field: 'status',
      title: '启用',
      width: 100,
      slots: { default: 'status' },
    },
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_params, formValues) => {
        if (!loaded.value) {
          const data = await listStoreAutoLabelRulesApi();
          rows.value = data.list ?? [];
          loaded.value = true;
        }
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const status = formValues?.status;
        let list = rows.value;
        if (keyword) {
          list = list.filter((item) =>
            item.name.toLowerCase().includes(keyword),
          );
        }
        if (status === 0 || status === 1) {
          list = list.filter((item) => item.status === status);
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'rule_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function addRow() {
  rows.value.push({
    rule_id: rows.value.length + 1,
    name: '新规则',
    rule_type: 'order_count',
    status: 0,
  });
  gridApi.reload();
}

async function save() {
  saving.value = true;
  try {
    const result = await saveStoreAutoLabelRulesApi({ list: rows.value });
    rows.value = result.list ?? rows.value;
    loaded.value = true;
    ElMessage.success('自动标签规则已保存（规则引擎待后续接入）');
    gridApi.reload();
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton @click="addRow">新增规则</ElButton>
      <ElButton type="primary" :loading="saving" @click="save">保存</ElButton>
    </template>

    <ElAlert
      class="mb-4"
      type="info"
      :closable="false"
      title="当前仅保存规则配置，不会自动执行打标。"
    />

    <Grid>
      <template #name="{ row }">
        <ElInput v-model="row.name" />
      </template>
      <template #type="{ row }">
        <ElSelect v-model="row.rule_type">
          <ElOption label="下单次数" value="order_count" />
          <ElOption label="累计消费" value="total_pay" />
        </ElSelect>
      </template>
      <template #status="{ row }">
        <ElSwitch v-model="row.status" :active-value="1" :inactive-value="0" />
      </template>
    </Grid>
  </Page>
</template>
