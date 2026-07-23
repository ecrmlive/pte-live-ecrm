<script lang="ts" setup>
import { nextTick, onMounted, reactive, ref } from 'vue';

import {
  ElButton,
  ElCheckbox,
  ElForm,
  ElFormItem,
  ElMessage,
  ElMessageBox,
  ElSwitch,
  ElTable,
  ElTableColumn,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import MerchantAccessApi from '#/api/core/merchant-access';
import MenuIconPreview from '#/components/platform-menu/MenuIconPreview.vue';
import { Page } from '@vben/common-ui';
import { deepClone } from '#/utils/base';

import AccessFormModal from './access-form-modal.vue';
import type { AccessAddType, AccessNode } from '#/views/access/types';

const loading = ref(true);
const rawData = ref<AccessNode[]>([]);
const tableData = ref<AccessNode[]>([]);
const tableRef = ref<InstanceType<typeof ElTable>>();

const formSearch = reactive({
  is_menu: false,
  pack_up: true,
});

const formModalOpen = ref(false);
const formMode = ref<'add' | 'edit'>('add');
const addType = ref<AccessAddType>('');
const selectModel = ref<AccessNode | undefined>();

function routeTypeLabel(value: number) {
  if (value === 1) return '页面';
  if (value === 0) return '按钮';
  if (value === 2) return '独立单页面';
  return '—';
}

function showScreen(list: AccessNode[], type: number) {
  for (let i = 0; i < list.length; i++) {
    const item = list[i];
    if (item && typeof item.is_menu !== 'undefined' && item.is_menu !== type) {
      list.splice(i, 1);
      i--;
    } else if (item?.children?.length) {
      showScreen(item.children, type);
    }
  }
}

function forArr(arr: AccessNode[], isExpand: boolean) {
  arr.forEach((row) => {
    tableRef.value?.toggleRowExpansion(row, isExpand);
    if (row.children?.length) {
      forArr(row.children, isExpand);
    }
  });
}

function changeIsMenuFunc(checked: boolean) {
  const list = deepClone(rawData.value) as AccessNode[];
  if (checked) {
    showScreen(list, 1);
    tableData.value = list;
  } else {
    tableData.value = list;
  }
}

async function changePackUpFunc(checked: boolean) {
  await nextTick();
  forArr(tableData.value, !checked);
}

async function loadTableList() {
  loading.value = true;
  try {
    const res = await MerchantAccessApi.accessList({}, true);
    rawData.value = (res.data as AccessNode[]) ?? [];
    tableData.value = deepClone(rawData.value) as AccessNode[];
    if (formSearch.is_menu) {
      changeIsMenuFunc(true);
    }
    await nextTick();
    forArr(tableData.value, !formSearch.pack_up);
  } catch {
    rawData.value = [];
    tableData.value = [];
  } finally {
    loading.value = false;
  }
}

async function isShowFunc(row: AccessNode) {
  try {
    const res = await MerchantAccessApi.status(
      { access_id: row.access_id, status: row.is_show },
      true,
    );
    if (res.code === 1) {
      ElMessage.success(res.msg || '更新成功');
      await loadTableList();
    }
  } catch {
    row.is_show = row.is_show === 1 ? 0 : 1;
  }
}

function openAdd(row?: AccessNode, type: AccessAddType = '') {
  formMode.value = 'add';
  addType.value = type;
  selectModel.value = row ? deepClone(row) : undefined;
  formModalOpen.value = true;
}

function openEdit(row: AccessNode) {
  formMode.value = 'edit';
  addType.value = '';
  selectModel.value = row;
  formModalOpen.value = true;
}

async function handleDelete(row: AccessNode) {
  try {
    await ElMessageBox.confirm('删除后不可恢复，确认删除该记录吗?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  loading.value = true;
  try {
    const res = await MerchantAccessApi.delAccess({ access_id: row.access_id }, true);
    if (res.code === 1) {
      ElMessage.success(res.msg || '删除成功');
      await loadTableList();
    }
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadTableList();
});
</script>

<template>
  <Page v-loading="loading">
    <div class="mb-3 flex w-full flex-wrap items-center justify-between gap-3">
        <ElButton
          v-access:code="'platform:merchantAccess:add'"
          :icon="Plus"
          type="primary"
          @click="openAdd()"
        >
          添加菜单&权限
        </ElButton>
        <ElForm :inline="true" :model="formSearch" size="small">
          <ElFormItem>
            <ElCheckbox
              v-model="formSearch.is_menu"
              @change="changeIsMenuFunc(!!$event)"
            >
              只显示菜单
            </ElCheckbox>
            <ElCheckbox
              v-model="formSearch.pack_up"
              @change="changePackUpFunc(!!$event)"
            >
              收起
            </ElCheckbox>
          </ElFormItem>
        </ElForm>
    </div>

    <ElTable
      ref="tableRef"
      :data="tableData"
      :default-expand-all="false"
      :tree-props="{ children: 'children' }"
      border
      class="w-full"
      row-key="access_id"
      size="small"
      style="width: 100%"
    >
      <ElTableColumn label="菜单名称" min-width="180" prop="name">
        <template #default="{ row }">
          <span
            :class="row.path === '/plus' ? 'font-semibold text-destructive' : ''"
            class="inline-flex items-center gap-2"
          >
            <MenuIconPreview v-if="row.is_menu === 1" :icon="row.icon" />
            {{ row.name }}
          </span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="图标" min-width="140" prop="icon">
        <template #default="{ row }">
          <span v-if="row.is_menu === 1" class="inline-flex items-center gap-2">
            <MenuIconPreview :icon="row.icon" />
            <span class="text-xs text-muted-foreground">{{ row.icon || '—' }}</span>
          </span>
          <span v-else class="text-muted-foreground">—</span>
        </template>
      </ElTableColumn>
      <ElTableColumn label="路径" min-width="160" prop="path" />
      <ElTableColumn label="权限码" min-width="160" prop="permission_code" />
      <ElTableColumn label="类别" prop="is_route" width="110">
        <template #default="{ row }">
          {{ routeTypeLabel(row.is_route) }}
        </template>
      </ElTableColumn>
      <ElTableColumn label="是否显示" prop="is_show" width="90">
        <template #default="{ row }">
          <ElSwitch
            v-model="row.is_show"
            :active-value="1"
            :inactive-value="0"
            @change="isShowFunc(row)"
          />
        </template>
      </ElTableColumn>
      <ElTableColumn label="排序" prop="sort" width="70" />
      <ElTableColumn label="添加时间" prop="create_time" width="150" />
      <ElTableColumn fixed="right" label="操作" width="280">
        <template #default="{ row }">
          <ElButton
            v-access:code="'platform:merchantAccess:add'"
            link
            type="primary"
            @click="openAdd(row, 'copy')"
          >
            一键复制
          </ElButton>
          <ElButton
            v-access:code="'platform:merchantAccess:add'"
            link
            type="primary"
            @click="openAdd(row, 'child')"
          >
            添加子菜单
          </ElButton>
          <ElButton
            v-access:code="'platform:merchantAccess:edit'"
            link
            type="primary"
            @click="openEdit(row)"
          >
            编辑
          </ElButton>
          <ElButton
            v-access:code="'platform:merchantAccess:delete'"
            link
            type="primary"
            @click="handleDelete(row)"
          >
            删除
          </ElButton>
        </template>
      </ElTableColumn>
    </ElTable>

    <AccessFormModal
      v-model:open="formModalOpen"
      :add-type="addType"
      :mode="formMode"
      :raw-data="rawData"
      :select-model="selectModel"
      @success="loadTableList"
    />
  </Page>
</template>
