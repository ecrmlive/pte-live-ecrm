<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, confirm } from '@vben/common-ui';
import { ElButton, ElImage, ElMessage } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  clearUserSearchRecords,
  exportUserSearchRecords,
  listUserSearchRecords,
  type UserSearchChannel,
  type UserSearchRecord,
} from '#/api/core/platform-user-search';
import { platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

type ChannelTab = UserSearchChannel;

const CHANNEL_TABS: Array<{ key: ChannelTab; label: string }> = [
  { key: '', label: '全部用户' },
  { key: 'wechat', label: '微信用户' },
  { key: 'mini_program', label: '小程序用户' },
  { key: 'h5', label: 'H5用户' },
  { key: 'app', label: 'APP用户' },
  { key: 'pc', label: 'PC用户' },
];

const CHANNEL_LABELS: Record<string, string> = {
  wechat: '微信用户',
  mini_program: '小程序用户',
  h5: 'H5用户',
  pc: 'PC用户',
  ios: 'APP用户',
  android: 'APP用户',
  harmony: 'APP用户',
};

const canRead = ref(false);
const canClear = ref(false);
const canExport = ref(false);
const channelTab = ref<ChannelTab>('');
const lastFilters = ref<Record<string, unknown>>({});

function buildFilters(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  return {
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    nickname: String(formValues?.nickname ?? '').trim() || undefined,
    user_type: (channelTab.value || undefined) as
      | Exclude<ChannelTab, ''>
      | undefined,
    start_date: range[0] as string | undefined,
    end_date: range[1] as string | undefined,
  };
}

function setChannelTab(key: ChannelTab) {
  if (channelTab.value === key) return;
  channelTab.value = key;
  gridApi.reload();
}

function channelLabel(value: string) {
  return CHANNEL_LABELS[String(value)] || String(value || '—');
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '搜索时间',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入搜索词' },
    fieldName: 'keyword',
    label: '搜索词',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入昵称' },
    fieldName: 'nickname',
    label: '用户昵称',
  },
]);

const gridOptions: VxeGridProps<UserSearchRecord> = {
  columns: [
    { field: 'user_id', title: '用户ID', width: 100 },
    {
      field: 'avatar_url',
      slots: { default: 'avatar' },
      title: '头像',
      width: 80,
    },
    {
      field: 'nickname',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      showOverflow: false,
      title: '昵称',
    },
    {
      field: 'user_type',
      formatter: ({ cellValue }) => channelLabel(String(cellValue || '')),
      minWidth: 110,
      title: '用户类型',
    },
    {
      field: 'keyword',
      minWidth: 180,
      showOverflow: false,
      title: '搜索词',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '搜索时间',
    },
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const filters = buildFilters(formValues);
        lastFilters.value = filters;
        const result = await listUserSearchRecords({
          page: page.currentPage,
          limit: page.pageSize,
          ...filters,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function clearRecords() {
  try {
    await confirm({
      title: '提示',
      content: '确认一键清空全部搜索记录吗？清空后列表将不再展示这些记录。',
      icon: 'warning',
    });
  } catch {
    return;
  }
  try {
    const result = await clearUserSearchRecords();
    ElMessage.success(`已清空 ${result.cleared_count} 条搜索记录`);
    gridApi.reload();
  } catch {
    /* 错误由 request 层提示 */
  }
}

async function exportRows() {
  try {
    const filters = lastFilters.value;
    const result = await exportUserSearchRecords({
      keyword: filters.keyword as string | undefined,
      nickname: filters.nickname as string | undefined,
      user_type: filters.user_type as Exclude<ChannelTab, ''> | undefined,
      start_date: filters.start_date as string | undefined,
      end_date: filters.end_date as string | undefined,
    });
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条搜索记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`,
    );
  } catch {
    /* 错误由 request 层提示 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const platform = profile.roles.includes('platform');
  canRead.value = platform && codes.includes('user.search_record.read');
  canClear.value = platform && codes.includes('user.search_record.clear');
  canExport.value = platform && codes.includes('user.search_record.export');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="search-record-toolbar">
          <div class="search-record-tabs" role="tablist">
            <button
              v-for="tab in CHANNEL_TABS"
              :key="tab.key || 'all'"
              type="button"
              role="tab"
              class="search-record-tabs__item"
              :aria-selected="channelTab === tab.key"
              :class="{ 'is-active': channelTab === tab.key }"
              @click="setChannelTab(tab.key)"
            >
              {{ tab.label }}
            </button>
          </div>
          <div class="search-record-toolbar__actions">
            <ElButton v-if="canExport" type="primary" @click="exportRows">
              导出搜索记录
            </ElButton>
            <ElButton v-if="canClear" type="danger" plain @click="clearRecords">
              一键清空
            </ElButton>
          </div>
        </div>
      </template>

      <template #avatar="{ row }">
        <ElImage
          v-if="row.avatar_url"
          :src="resolveCosMediaUrl(row.avatar_url)"
          fit="cover"
          class="user-avatar"
        >
          <template #error>
            <div class="user-avatar user-avatar--empty">无</div>
          </template>
        </ElImage>
        <div v-else class="user-avatar user-avatar--empty">—</div>
      </template>
    </Grid>
  </Page>
</template>

<style scoped>
.search-record-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.search-record-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 16px;
  min-width: 0;
}

.search-record-tabs__item {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 8px 2px;
  cursor: pointer;
  color: hsl(var(--foreground) / 0.75);
  font-size: 14px;
  line-height: 1.2;
}

.search-record-tabs__item.is-active {
  color: hsl(var(--primary));
  font-weight: 600;
  box-shadow: inset 0 -2px 0 hsl(var(--primary));
}

.search-record-toolbar__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  overflow: hidden;
}

.user-avatar--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: hsl(var(--muted) / 0.4);
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}
</style>
