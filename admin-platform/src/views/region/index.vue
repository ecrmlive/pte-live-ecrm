<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref, watch } from 'vue';

import { ElButton, ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import RegionApi from '#/api/core/region';
import { Page } from '@vben/common-ui';
import { parseApiListPage } from '#/utils/list-response';
import {
  PLATFORM_SEARCH_SELECT_PROPS
} from '#/utils/platform-list-search-form';

import RegionFormModal from './region-form-modal.vue';
import type { RegionAreaItem, RegionRow } from './types';
import { regionCityOptions } from './types';

const formOpen = ref(false);
const formMode = ref<'add' | 'edit'>('add');
const currentRegion = ref<RegionRow | undefined>();
const areaList = ref<Record<number, RegionAreaItem>>({});

const LEVEL_OPTIONS = [
  { label: '全部', value: 0 },
  { label: '省份', value: 1 },
  { label: '城市', value: 2 },
  { label: '地区', value: 3 },
];

function buildProvinceOptions() {
  return [
    { label: '全部', value: 0 },
    ...Object.values(areaList.value).map((p) => ({
      label: p.name,
      value: p.id,
    })),
  ];
}

function buildCityOptions(provinceId: number) {
  return [
    { label: '全部', value: 0 },
    ...regionCityOptions(areaList.value, provinceId).map((c) => ({
      label: c.name,
      value: c.id,
    })),
  ];
}

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '请输入名称' },
      defaultValue: '',
      fieldName: 'name',
      formItemClass: 'pb-0',
      label: '名称',
    },
    {
      component: 'Select',
      componentProps: {
        options: LEVEL_OPTIONS,
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      defaultValue: 0,
      fieldName: 'level',
      formItemClass: 'pb-0',
      label: '地区类型',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [{ label: '全部', value: 0 }],
        placeholder: '省',
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      defaultValue: 0,
      fieldName: 'province_id',
      formItemClass: 'pb-0',
      label: '选择省份',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [{ label: '全部', value: 0 }],
        placeholder: '市',
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      defaultValue: 0,
      dependencies: {
        componentProps(values) {
          const provinceId = Number(values.province_id || 0);
          return {
            clearable: true,
            options: buildCityOptions(provinceId),
            placeholder: '市',
            ...PLATFORM_SEARCH_SELECT_PROPS,
          };
        },
        show(values) {
          return Number(values.province_id) !== 0;
        },
        triggerFields: ['province_id'],
      },
      fieldName: 'city_id',
      formItemClass: 'pb-0',
      label: '选择城市',
    },
  ],
};

async function fetchRegionPage(
  pageSize: number,
  currentPage: number,
  formValues?: Record<string, unknown>,
) {
  const res = await RegionApi.regionList(
    {
      name: String(formValues?.name ?? ''),
      level: Number(formValues?.level ?? 0),
      province_id: Number(formValues?.province_id ?? 0),
      city_id: Number(formValues?.city_id ?? 0),
      list_rows: pageSize,
      page: currentPage,
    },
    true,
  );
  const page = parseApiListPage<RegionRow>(res.data);
  const regionData =
    (res.data as { regionData?: Record<number, RegionAreaItem> })?.regionData ??
    {};
  areaList.value = regionData;
  return { items: page.list, total: page.total };
}

const gridOptions: VxeGridProps = {
  border: true,
  height: 'auto',
  columns: [
    { field: 'id', title: 'ID', width: 88 },
    {
      field: 'ad_code',
      minWidth: 120,
      showOverflow: true,
      title: '区划编号',
    },
    { field: 'shortname', minWidth: 100, title: '简称' },
    { field: 'name', minWidth: 120, title: '名称' },
    {
      field: 'merger_name',
      minWidth: 200,
      showOverflow: true,
      title: '全称',
    },
    {
      field: 'level',
      slots: { default: 'level' },
      title: '级别',
      width: 88,
    },
    { field: 'lng', title: '经度', width: 100 },
    { field: 'lat', title: '纬度', width: 100 },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 140,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchRegionPage(page.pageSize, page.currentPage, formValues),
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    ...formOptions,
    handleValuesChange(values, fieldsChanged) {
      if (fieldsChanged.includes('province_id')) {
        gridApi.formApi?.setFieldValue('city_id', 0);
      }
    },
  },
  gridOptions,
});

watch(areaList, () => {
  gridApi.formApi?.updateSchema([
    {
      componentProps: {
        clearable: true,
        options: buildProvinceOptions(),
        placeholder: '省',
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      fieldName: 'province_id',
    },
  ]);
});

function reload() {
  gridApi.reload();
}

function openAdd() {
  formMode.value = 'add';
  currentRegion.value = undefined;
  formOpen.value = true;
}

function openEdit(item: RegionRow) {
  formMode.value = 'edit';
  currentRegion.value = item;
  formOpen.value = true;
}

async function deleteRegion(row: RegionRow) {
  try {
    await ElMessageBox.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  const res = await RegionApi.deleteRegion({ id: row.id }, true);
  ElMessage.success(res.msg || '删除成功');
  reload();
}
</script>

<template>
  <Page>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:region:add'"
          :icon="Plus"
          type="primary"
          @click="openAdd"
        >
          添加地区
        </ElButton>
      </template>

      <template #level="{ row }">
        <span v-if="row.level === 1">省份</span>
        <span v-else-if="row.level === 2">城市</span>
        <span v-else-if="row.level === 3">地区</span>
      </template>

      <template #action="{ row }">
        <ElButton
          v-access:code="'platform:region:edit'"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-access:code="'platform:region:delete'"
          link
          type="danger"
          @click="deleteRegion(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <RegionFormModal
      v-model:open="formOpen"
      :mode="formMode"
      :region="currentRegion"
      @success="reload"
    />
  </Page>
</template>
