<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
  ElTreeSelect,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteStoreGroup,
  fetchPlatformDiyPages,
  fetchPlatformMerchants,
  fetchStoreGroup,
  fetchStoreGroupMerchants,
  fetchStoreGroups,
  saveStoreGroup,
  setStoreGroupStatus,
  setStoreGroupTemplate,
  type DiyPageOption,
  type PlatformMerchantRow,
  type StoreGroupMerchantRow,
  type StoreGroupRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const treeRows = ref<StoreGroupRow[]>([]);
const merchants = ref<PlatformMerchantRow[]>([]);
const templates = ref<DiyPageOption[]>([]);
const editing = ref<StoreGroupRow>();
const linkedStores = ref<StoreGroupMerchantRow[]>([]);
const form = reactive({
  parent_id: 0,
  name: '',
  sort: 0,
  status: true,
  diy_page_id: 0,
  positioning_status: false,
  longitude: undefined as number | undefined,
  latitude: undefined as number | undefined,
  address: '',
  merchant_ids: [] as number[],
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '分组名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

function filterTreeByStatus(
  nodes: StoreGroupRow[],
  status?: number,
): StoreGroupRow[] {
  if (status !== 0 && status !== 1) return nodes;
  return nodes
    .map((node) => {
      const children = node.children
        ? filterTreeByStatus(node.children, status)
        : undefined;
      const selfMatch = node.status === status;
      if (selfMatch || (children && children.length)) {
        return { ...node, children };
      }
      return null;
    })
    .filter((node): node is StoreGroupRow => node !== null);
}

const gridOptions: VxeGridProps<StoreGroupRow> = {
  columns: [
    {
      field: 'name',
      minWidth: 260,
      showOverflow: false,
      title: '分组名称',
      treeNode: true,
    },
    {
      field: 'level',
      formatter: ({ cellValue }) => `第 ${Number(cellValue) + 1} 级`,
      title: '层级',
      width: 80,
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'merchant_count',
      formatter: ({ cellValue }) => `${cellValue} 家`,
      title: '关联店铺',
      width: 110,
    },
    {
      field: 'positioning_status',
      formatter: ({ cellValue }) => (cellValue ? '启用' : '关闭'),
      title: '定位',
      width: 90,
    },
    {
      field: 'diy_page_id',
      formatter: ({ cellValue }) => cellValue || '未绑定',
      title: '装修模板',
      width: 120,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ minWidth: 280 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim() || undefined;
        const statusRaw = formValues?.status;
        let list = (await fetchStoreGroups(keyword)).list || [];
        treeRows.value = list;
        if (statusRaw === 0 || statusRaw === 1) {
          list = filterTreeByStatus(list, Number(statusRaw));
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  treeConfig: {
    childrenField: 'children',
    expandAll: true,
  },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

const [StoresModal, storesModalApi] = useVbenModal({
  showConfirmButton: false,
  cancelText: '关闭',
});

function resetForm(row?: StoreGroupRow) {
  editing.value = row;
  Object.assign(form, {
    parent_id: row?.parent_id || 0,
    name: row?.name || '',
    sort: row?.sort || 0,
    status: row ? row.status === 1 : true,
    diy_page_id: row?.diy_page_id || 0,
    positioning_status: row ? row.positioning_status === 1 : false,
    longitude: row?.longitude ?? undefined,
    latitude: row?.latitude ?? undefined,
    address: row?.address || '',
    merchant_ids: row?.merchant_ids || [],
  });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: '新增店铺分组' }).open();
}

async function openEdit(row: StoreGroupRow) {
  resetForm(await fetchStoreGroup(row.id));
  formModalApi.setState({ title: '编辑店铺分组' }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分组名称');
    return;
  }
  if (
    (form.longitude === undefined) !== (form.latitude === undefined)
  ) {
    ElMessage.warning('经度和纬度需同时填写，或同时留空');
    return;
  }
  formModalApi.lock();
  try {
    await saveStoreGroup(editing.value?.id, {
      ...form,
      name: form.name.trim(),
      address: form.address.trim(),
    });
    formModalApi.close();
    ElMessage.success('店铺分组已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function toggleStatus(row: StoreGroupRow) {
  await setStoreGroupStatus(row.id, row.status !== 1);
  ElMessage.success('分组状态已更新，并已同步至子分组');
  gridApi.reload();
}

async function updateTemplate(row: StoreGroupRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写装修模板 ID；填 0 可清空绑定。',
      `设置“${row.name}”装修模板`,
      {
        inputValue: String(row.diy_page_id || 0),
        inputPattern: /^\d+$/,
        inputErrorMessage: '请输入非负整数 ID',
      },
    );
    await setStoreGroupTemplate(row.id, Number(value));
    ElMessage.success('装修模板已更新');
    gridApi.reload();
  } catch {
    /* 取消操作或统一请求层已提示错误 */
  }
}

async function remove(row: StoreGroupRow) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”后不可恢复；含子分组时系统会拒绝删除。是否继续？`,
      '删除店铺分组',
      { type: 'warning' },
    );
    await deleteStoreGroup(row.id);
    ElMessage.success('店铺分组已删除');
    gridApi.reload();
  } catch {
    /* 取消操作或统一请求层已提示错误 */
  }
}

async function showStores(row: StoreGroupRow) {
  linkedStores.value = (await fetchStoreGroupMerchants(row.id)).list || [];
  storesModalApi.setState({ title: `关联店铺：${row.name}` }).open();
}

async function loadOptions() {
  const [merchantResult, diyResult] = await Promise.all([
    fetchPlatformMerchants({ page: 1, limit: 100 }),
    fetchPlatformDiyPages(),
  ]);
  merchants.value = merchantResult.list || [];
  templates.value = diyResult.list || [];
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), loadOptions()]);
  canManage.value = permissions.includes('merchant.group.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增分组
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="showStores(row)">
          关联店铺
        </ElButton>
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="primary" @click="toggleStatus(row)">
            {{ row.status ? '停用' : '启用' }}
          </ElButton>
          <ElButton link type="primary" @click="updateTemplate(row)">
            模板
          </ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="112px">
        <ElFormItem label="上级分组">
          <ElTreeSelect
            v-model="form.parent_id"
            :data="treeRows"
            node-key="id"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            check-strictly
            clearable
            placeholder="不选则为一级分组"
          />
        </ElFormItem>
        <ElFormItem label="分组名称" required>
          <ElInput v-model="form.name" maxlength="128" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
        <ElFormItem label="初始状态">
          <ElSwitch v-model="form.status" />
        </ElFormItem>
        <ElFormItem label="关联店铺">
          <ElSelect
            v-model="form.merchant_ids"
            multiple
            filterable
            class="w-full"
            placeholder="可关联多个统一后台商户投影"
          >
            <ElOption
              v-for="item in merchants"
              :key="item.mer_id"
              :label="item.mer_name"
              :value="item.mer_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="装修模板">
          <ElSelect
            v-model="form.diy_page_id"
            clearable
            class="w-full"
            placeholder="不选则不绑定"
          >
            <ElOption :value="0" label="不绑定模板" />
            <ElOption
              v-for="item in templates"
              :key="item.id"
              :label="item.name || `模板 #${item.id}`"
              :value="item.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="启用定位">
          <ElSwitch v-model="form.positioning_status" />
        </ElFormItem>
        <ElFormItem label="经纬度">
          <ElInputNumber
            v-model="form.longitude"
            :precision="7"
            :min="-180"
            :max="180"
            placeholder="经度"
          />
          <span class="mx-2">/</span>
          <ElInputNumber
            v-model="form.latitude"
            :precision="7"
            :min="-90"
            :max="90"
            placeholder="纬度"
          />
        </ElFormItem>
        <ElFormItem label="地址">
          <ElInput v-model="form.address" maxlength="255" />
        </ElFormItem>
      </ElForm>
    </FormModal>

    <StoresModal>
      <ElTable :data="linkedStores">
        <ElTableColumn label="店铺 ID" prop="merchant_id" width="100" />
        <ElTableColumn label="店铺名称" prop="merchant_name" />
        <ElTableColumn label="区域" prop="region_id" width="100" />
        <ElTableColumn label="状态" width="90">
          <template #default="{ row }">
            {{ row.status ? '启用' : '停用' }}
          </template>
        </ElTableColumn>
      </ElTable>
    </StoresModal>
  </Page>
</template>
