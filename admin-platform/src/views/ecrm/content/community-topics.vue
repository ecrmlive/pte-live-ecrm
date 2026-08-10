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
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createCommunityTopicApi,
  deleteCommunityTopicApi,
  listCommunityCategoriesApi,
  listCommunityTopicsApi,
  updateCommunityTopicApi,
  updateCommunityTopicHotApi,
  updateCommunityTopicStatusApi,
  type CommunityCategory,
  type CommunityTopic,
} from '#/api/core/platform-community';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const canRead = ref(false);
const canManage = ref(false);
const categories = ref<CommunityCategory[]>([]);
const editing = ref<CommunityTopic>();
const form = reactive({
  category_id: undefined as number | undefined,
  is_hot: 0,
  pic: '',
  sort: 0,
  status: 1,
  topic_name: '',
});

const gridOptions: VxeGridProps<CommunityTopic> = {
  columns: [
    {
      field: 'topic_name',
      minWidth: 180,
      showOverflow: false,
      title: '话题名称',
    },
    {
      field: 'pic',
      slots: { default: 'pic' },
      title: '话题图标',
      width: 100,
    },
    {
      field: 'cate_name',
      formatter: ({ cellValue, row }) =>
        String(cellValue || '').trim() ||
        categoryName(row.category_id) ||
        '—',
      minWidth: 120,
      showOverflow: false,
      title: '上级分类',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'count_use',
      title: '文章数',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'is_hot',
      slots: { default: 'is_hot' },
      title: '是否推荐',
      width: 120,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const list = (await listCommunityTopicsApi()).list || [];
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'topic_id' },
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

function categoryName(id: number) {
  return (
    categories.value.find((c) => c.category_id === id)?.cate_name || ''
  );
}

function imageOf(row: CommunityTopic) {
  return resolveCosMediaUrl(String(row.pic || '').trim());
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    category_id: undefined,
    is_hot: 0,
    pic: '',
    sort: 0,
    status: 1,
    topic_name: '',
  });
}

async function loadCategories() {
  try {
    categories.value = (await listCommunityCategoriesApi()).list || [];
  } catch {
    categories.value = [];
  }
}

function openCreate() {
  resetForm();
  void loadCategories();
  formDrawerApi.setState({ title: '添加社区话题' }).open();
}

function openEdit(row: CommunityTopic) {
  editing.value = row;
  Object.assign(form, {
    category_id: row.category_id || undefined,
    is_hot: row.is_hot === 1 ? 1 : 0,
    pic: row.pic || '',
    sort: Number(row.sort || 0),
    status: row.status === 1 ? 1 : 0,
    topic_name: row.topic_name || '',
  });
  void loadCategories();
  formDrawerApi.setState({ title: '编辑社区话题' }).open();
}

async function save() {
  const topicName = form.topic_name.trim();
  if (!topicName) {
    ElMessage.warning('请填写话题名称');
    return;
  }
  if (!form.category_id) {
    ElMessage.warning('请选择上级分类');
    return;
  }
  const sort = Math.max(0, Math.min(99999, Number(form.sort) || 0));
  formDrawerApi.lock();
  try {
    const payload = {
      category_id: Number(form.category_id),
      is_hot: form.is_hot === 1 ? 1 : 0,
      pic: String(form.pic || '').trim(),
      sort,
      status: form.status === 1 ? 1 : 0,
      topic_name: topicName,
    };
    if (editing.value) {
      await updateCommunityTopicApi(editing.value.topic_id, payload);
    } else {
      await createCommunityTopicApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: CommunityTopic, enabled: boolean) {
  if (!canManage.value) return;
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateCommunityTopicStatusApi(row.topic_id, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function changeHot(row: CommunityTopic, enabled: boolean) {
  if (!canManage.value) return;
  const before = row.is_hot;
  row.is_hot = enabled ? 1 : 0;
  try {
    await updateCommunityTopicHotApi(row.topic_id, enabled ? 1 : 0);
  } catch {
    row.is_hot = before;
  }
}

async function remove(row: CommunityTopic) {
  try {
    await confirm({
      content: `删除话题“${row.topic_name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteCommunityTopicApi(row.topic_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('content.community_topic.manage');
  canRead.value =
    canManage.value ||
    permissions.includes('content.community_topic.read');
  if (canRead.value) {
    await loadCategories();
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
          添加社区话题
        </ElButton>
      </template>
      <template #pic="{ row }">
        <ElImage
          v-if="imageOf(row)"
          :src="imageOf(row)"
          fit="cover"
          class="topic-icon"
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
      <template #is_hot="{ row }">
        <ElSwitch
          :model-value="row.is_hot === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="推荐"
          inactive-text="否"
          @change="
            (enabled: string | number | boolean) =>
              changeHot(row, Boolean(enabled))
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
        <ElFormItem label="话题名称" required>
          <ElInput
            v-model="form.topic_name"
            maxlength="32"
            show-word-limit
            placeholder="请输入话题名称"
          />
        </ElFormItem>
        <ElFormItem label="话题图标">
          <ImageField v-model="form.pic" />
        </ElFormItem>
        <ElFormItem label="上级分类" required>
          <ElSelect
            v-model="form.category_id"
            class="w-full"
            clearable
            filterable
            placeholder="请选择社区分类"
          >
            <ElOption
              v-for="item in categories"
              :key="item.category_id"
              :label="item.cate_name"
              :value="item.category_id"
            />
          </ElSelect>
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
        <ElFormItem label="是否推荐">
          <ElSwitch
            v-model="form.is_hot"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="推荐"
            inactive-text="否"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.topic-icon {
  width: 40px;
  height: 40px;
  border-radius: 4px;
}
</style>
