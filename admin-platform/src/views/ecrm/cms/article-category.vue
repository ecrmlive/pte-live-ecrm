<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createArticleCategoryApi,
  deleteArticleCategoryApi,
  getArticleCategoryListApi,
  updateArticleCategoryApi,
  updateArticleCategoryStatusApi,
  type ArticleCategoryOption,
} from '#/api/core/plus-article';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const canRead = ref(false);
const canManage = ref(false);
const editing = ref<ArticleCategoryOption>();
const form = reactive({
  image: '',
  info: '',
  sort: 0,
  status: 1,
  title: '',
});

const gridOptions: VxeGridProps<ArticleCategoryOption> = {
  columns: [
    {
      field: 'title',
      formatter: ({ row }) => formatCateName(row),
      minWidth: 200,
      showOverflow: false,
      title: '分类名称',
    },
    {
      field: 'info',
      formatter: ({ cellValue }) => String(cellValue || '').trim() || '—',
      minWidth: 160,
      showOverflow: false,
      title: '配置分类说明',
    },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '分类图片',
      width: 120,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 180,
      showOverflow: false,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const list = (await getArticleCategoryListApi()).list || [];
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'cid' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function formatCateName(row: ArticleCategoryOption) {
  const name = String(row.title || '').trim() || '—';
  const id = Number(row.cid || 0);
  // CRMEB：分类名称后附 ID，如 生活  [ 1 ]
  return id > 0 ? `${name}  [ ${id}  ]` : name;
}

function imageOf(row: ArticleCategoryOption) {
  return resolveCosMediaUrl(String(row.image || '').trim());
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { image: '', info: '', sort: 0, status: 1, title: '' });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加文章分类' }).open();
}

function openEdit(row: ArticleCategoryOption) {
  editing.value = row;
  Object.assign(form, {
    image: row.image || '',
    info: row.info || '',
    sort: Number(row.sort || 0),
    status: row.status === 1 ? 1 : 0,
    title: row.title || '',
  });
  formDrawerApi.setState({ title: '编辑文章分类' }).open();
}

async function save() {
  const title = form.title.trim();
  if (!title) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const sort = Math.max(0, Math.min(99999, Number(form.sort) || 0));
  formDrawerApi.lock();
  try {
    const payload = {
      image: String(form.image || '').trim(),
      info: form.info.trim(),
      sort,
      status: form.status === 1 ? 1 : 0,
      title,
    };
    if (editing.value) {
      await updateArticleCategoryApi(editing.value.cid, payload);
    } else {
      await createArticleCategoryApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: ArticleCategoryOption, enabled: boolean) {
  if (!canManage.value) return;
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateArticleCategoryStatusApi(row.cid, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function remove(row: ArticleCategoryOption) {
  try {
    await confirm({
      content: `删除分类“${row.title}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteArticleCategoryApi(row.cid);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('content.article_category.manage');
  canRead.value =
    canManage.value ||
    permissions.includes('content.article_category.read') ||
    permissions.includes('content.article.read') ||
    permissions.includes('content.article.manage');
  if (canRead.value) {
    gridApi.reload();
  }
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
          添加文章分类
        </ElButton>
      </template>
      <template #image="{ row }">
        <ElImage
          v-if="imageOf(row)"
          :src="imageOf(row)"
          fit="cover"
          class="cate-cover"
          :preview-src-list="[imageOf(row)]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
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
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="primary" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="分类名称" required>
          <ElInput
            v-model="form.title"
            maxlength="32"
            show-word-limit
            placeholder="请输入分类名称"
          />
        </ElFormItem>
        <ElFormItem label="分类简介">
          <ElInput
            v-model="form.info"
            maxlength="255"
            show-word-limit
            placeholder="请输入分类简介"
          />
        </ElFormItem>
        <ElFormItem label="分类图片">
          <ImageField v-model="form.image" />
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
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="99999"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.cate-cover {
  width: 48px;
  height: 48px;
  border-radius: 4px;
}
</style>
