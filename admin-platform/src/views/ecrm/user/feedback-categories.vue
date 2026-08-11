<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createUserFeedbackCategory,
  deleteUserFeedbackCategory,
  fetchUserFeedbackCategories,
  setUserFeedbackCategoryStatus,
  updateUserFeedbackCategory,
  type UserFeedbackCategory,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const canManage = ref(false);
const treeRows = ref<UserFeedbackCategory[]>([]);
const editing = ref<UserFeedbackCategory>();
const form = reactive({
  id: 0,
  name: '',
  pid: 0,
  sort: 0,
  status: 1 as 0 | 1,
});

/** 两级树：仅顶级分类可作为上级 */
const topParents = computed(() =>
  treeRows.value
    .filter((node) => node.id !== editing.value?.id)
    .map((node) => ({ label: node.name, value: node.id })),
);

function idempotencyKey(action: string, id = 0) {
  return `feedback-category-${action}-${id}-${crypto.randomUUID()}`;
}

function filterTree(
  nodes: UserFeedbackCategory[],
  keyword: string,
  status?: number,
): UserFeedbackCategory[] {
  return nodes
    .map((node) => {
      const children = node.children
        ? filterTree(node.children, keyword, status)
        : undefined;
      const nameMatch =
        !keyword || node.name.toLowerCase().includes(keyword);
      const statusMatch =
        status !== 0 && status !== 1 ? true : node.status === status;
      if ((nameMatch && statusMatch) || (children && children.length)) {
        return { ...node, children };
      }
      return null;
    })
    .filter((node): node is UserFeedbackCategory => node !== null);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('分类名称'),
  LIST_ENABLE_STATUS_FIELD('显示状态'),
]);

const gridOptions: VxeGridProps<UserFeedbackCategory> = {
  columns: [
    {
      align: 'left',
      field: 'name',
      headerAlign: 'left',
      minWidth: 260,
      showOverflow: false,
      slots: { default: 'name' },
      title: '分类名称',
      treeNode: true,
    },
    { align: 'center', field: 'sort', title: '排序', width: 90 },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        const status =
          statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined;
        let list = (await fetchUserFeedbackCategories()).list || [];
        treeRows.value = list;
        if (keyword || status !== undefined) {
          list = filterTree(list, keyword, status);
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  treeConfig: {
    childrenField: 'children',
    expandAll: true,
    indent: 10,
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

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm(parentID = 0) {
  editing.value = undefined;
  Object.assign(form, {
    id: 0,
    name: '',
    pid: parentID,
    sort: 0,
    status: 1,
  });
}

function openCreate(parentID = 0) {
  resetForm(parentID);
  formDrawerApi.setState({ title: '新增反馈分类' }).open();
}

function openEdit(row: UserFeedbackCategory) {
  editing.value = row;
  Object.assign(form, {
    id: row.id,
    name: row.name,
    pid: row.pid ?? 0,
    sort: row.sort,
    status: row.status,
  });
  formDrawerApi.setState({ title: '编辑反馈分类' }).open();
}

async function save() {
  const name = form.name.trim();
  if (!name || [...name].length > 32) {
    ElMessage.warning('请填写不超过 32 字的分类名称');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      name,
      pid: Number(form.pid) || 0,
      sort: Math.max(0, Math.min(9999, Number(form.sort) || 0)),
      status: form.status,
      idempotency_key: idempotencyKey(form.id ? 'update' : 'create', form.id),
    };
    if (form.id) await updateUserFeedbackCategory(form.id, payload);
    else await createUserFeedbackCategory(payload);
    formDrawerApi.close();
    ElMessage.success('已保存');
    resetForm();
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: UserFeedbackCategory, enabled: boolean) {
  const before = row.status === 1;
  row.status = enabled ? 1 : 0;
  try {
    await setUserFeedbackCategoryStatus(row.id, {
      status: enabled ? 1 : 0,
      idempotency_key: idempotencyKey('status', row.id),
    });
  } catch {
    row.status = before ? 1 : 0;
  }
}

async function remove(row: UserFeedbackCategory) {
  try {
    await confirm({
      content: `删除“${row.name}”仅从可选分类中移除，既有反馈仍保留原分类文本。存在子分类时请先处理子集。`,
      icon: 'warning',
      title: '提示',
    });
    await deleteUserFeedbackCategory(row.id, {
      idempotency_key: idempotencyKey('delete', row.id),
    });
    ElMessage.success('已删除');
    if (form.id === row.id) resetForm();
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('user.feedback.category.read') ||
      codes.includes('user.feedback.category.manage'));
  canManage.value =
    roleOK && codes.includes('user.feedback.category.manage');
  if (canRead.value) gridApi.reload();
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
          @click="openCreate(0)"
        >
          新增反馈分类
        </ElButton>
      </template>

      <template #name="{ row }">
        <span>{{ row.name }} [ {{ row.id }} ]</span>
      </template>

      <template #status="{ row }">
        <ElSwitch
          :disabled="!canManage"
          :model-value="row.status === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeShow(row, Boolean(enabled))
          "
        />
      </template>

      <template #action="{ row }">
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          type="danger"
          @click="remove(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="88px">
        <ElFormItem label="分类名称" required>
          <ElInput
            v-model="form.name"
            maxlength="32"
            placeholder="请输入分类名称"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="上级分类">
          <ElSelect
            v-model="form.pid"
            clearable
            class="w-full"
            placeholder="顶级分类"
          >
            <ElOption label="顶级分类" :value="0" />
            <ElOption
              v-for="item in topParents"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :max="9999" class="w-full" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="显示"
            inactive-text="隐藏"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
:deep(.vxe-cell--tree-node) {
  text-align: left;
}

:deep(.vxe-cell--tree-btn) {
  width: 1em;
  height: 1em;
  margin: 0;
}

:deep(.vxe-tree-cell) {
  padding-left: 1em;
}
</style>
