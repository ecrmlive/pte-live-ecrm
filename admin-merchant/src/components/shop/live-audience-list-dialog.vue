<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type {
  LiveAudienceInviter,
  LiveAudienceListItem,
  LiveInviteStatsItem,
  LiveRoomListItem,
} from '#/api/core/live';
import type { VbenFormSchema } from '#/adapter/form';

import { computed, reactive, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getLiveAudienceListApi, getLiveInviteStatsApi } from '#/api/core/live';

defineOptions({ name: 'LiveAudienceListDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveRoomListItem | Record<string, never>;
}>();

const activeTab = ref<'audience' | 'stats'>('audience');
const statsLoaded = ref(false);

const searchSchema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入用户昵称' },
    fieldName: 'nick_name',
    label: '用户昵称',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入用户ID' },
    fieldName: 'user_id',
    label: '用户ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '本场邀请人' },
    fieldName: 'inviter_nick_name',
    label: '邀请人昵称',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '本场邀请人ID' },
    fieldName: 'inviter_id',
    label: '邀请人ID',
  },
  {
    component: 'DatePicker',
    componentProps: {
      endPlaceholder: '结束时间',
      rangeSeparator: '至',
      startPlaceholder: '开始时间',
      style: 'width: 100%',
      type: 'datetimerange',
      valueFormat: 'YYYY-MM-DD HH:mm:ss',
    },
    fieldName: 'enter_time_range',
    formItemClass: 'col-span-1 md:col-span-2',
    label: '进入时间',
  },
  {
    component: 'DatePicker',
    componentProps: {
      endPlaceholder: '结束时间',
      rangeSeparator: '至',
      startPlaceholder: '开始时间',
      style: 'width: 100%',
      type: 'datetimerange',
      valueFormat: 'YYYY-MM-DD HH:mm:ss',
    },
    fieldName: 'last_active_range',
    formItemClass: 'col-span-1 md:col-span-2',
    label: '最后活跃时间',
  },
]);

const [SearchForm, searchFormApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { size: 'small' },
      labelWidth: 96,
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema: searchSchema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4',
  }),
);

const summary = reactive({
  anonymous_count: 0,
  logged_in_count: 0,
});

const dialogTitle = computed(() => {
  const name = props.row?.name || '直播间';
  return `${name} - 观众列表`;
});

function formatLiveTimeText(text?: string) {
  const value = text?.trim() ?? '';
  if (!value || value.startsWith('0000-00-00')) return '-';
  return value;
}

function inviterUserId(inviter?: LiveAudienceInviter) {
  return inviter?.user_id ? inviter.user_id : 0;
}

function inviterNickName(inviter?: LiveAudienceInviter) {
  return inviter?.nick_name || '用户';
}

function inviterAvatar(inviter?: LiveAudienceInviter) {
  return inviter?.avatar || '';
}

async function buildAudienceParams(page: { currentPage: number; pageSize: number }) {
  const values = await searchFormApi.getValues();
  const enterRange = (values.enter_time_range as string[] | undefined) ?? [];
  const lastActiveRange = (values.last_active_range as string[] | undefined) ?? [];
  const params: Parameters<typeof getLiveAudienceListApi>[0] = {
    live_id: Number(props.row?.live_id ?? 0),
    list_rows: page.pageSize,
    nick_name: String(values.nick_name ?? ''),
    page: page.currentPage,
  };
  if (values.user_id) params.user_id = String(values.user_id);
  if (values.inviter_id) params.inviter_id = String(values.inviter_id);
  if (values.inviter_nick_name) params.inviter_nick_name = String(values.inviter_nick_name);
  if (enterRange.length === 2) {
    params.enter_time_start = enterRange[0];
    params.enter_time_end = enterRange[1];
  }
  if (lastActiveRange.length === 2) {
    params.last_active_start = lastActiveRange[0];
    params.last_active_end = lastActiveRange[1];
  }
  return params;
}

const audienceGridOptions = reactive<VxeGridProps<LiveAudienceListItem>>({
  columns: [
    {
      field: 'user',
      minWidth: 180,
      slots: { default: 'user' },
      title: '用户信息',
    },
    { field: 'enter_method_text', title: '进入直播间方式', width: 120 },
    {
      field: 'session_inviter',
      minWidth: 160,
      slots: { default: 'session_inviter' },
      title: '本场邀请人',
    },
    {
      field: 'first_inviter',
      minWidth: 160,
      slots: { default: 'first_inviter' },
      title: '首次邀请人',
    },
    {
      field: 'first_enter_time_text',
      minWidth: 170,
      slots: { default: 'first_enter_time' },
      title: '首次进入直播间时间',
    },
    {
      field: 'last_active_time_text',
      minWidth: 170,
      slots: { default: 'last_active_time' },
      title: '最后活跃时间',
    },
    {
      align: 'right',
      field: 'watch_duration_text',
      slots: { default: 'watch_duration' },
      title: '累计观看时长(分)',
      width: 140,
    },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50, 100],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.row?.live_id) {
          summary.logged_in_count = 0;
          summary.anonymous_count = 0;
          return { items: [], total: 0 };
        }
        try {
          const res = await getLiveAudienceListApi(await buildAudienceParams(page));
          summary.logged_in_count = res.summary?.logged_in_count ?? 0;
          summary.anonymous_count = res.summary?.anonymous_count ?? 0;
          return {
            items: res.list.data ?? [],
            total: res.list.total ?? 0,
          };
        } catch {
          summary.logged_in_count = 0;
          summary.anonymous_count = 0;
          return { items: [], total: 0 };
        }
      },
    },
  },
  toolbarConfig: {
    enabled: false,
  },
});

const statsGridOptions = reactive<VxeGridProps<LiveInviteStatsItem>>({
  columns: [
    {
      field: 'inviter',
      minWidth: 180,
      slots: { default: 'inviter' },
      title: '邀请人',
    },
    {
      align: 'right',
      field: 'invite_user_count',
      title: '邀请人数',
      width: 100,
    },
    {
      align: 'right',
      field: 'invite_watch_duration_text',
      slots: { default: 'invite_watch_duration' },
      title: '被邀请观看时长(分)',
      width: 160,
    },
    {
      align: 'right',
      field: 'order_count',
      title: '成交订单数',
      width: 110,
    },
    {
      align: 'right',
      field: 'order_pay_amount_text',
      slots: { default: 'order_pay_amount' },
      title: '支付金额(元)',
      width: 120,
    },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50, 100],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.row?.live_id) {
          return { items: [], total: 0 };
        }
        try {
          const res = await getLiveInviteStatsApi({
            live_id: props.row.live_id,
            list_rows: page.pageSize,
            page: page.currentPage,
          });
          return {
            items: res.list.data ?? [],
            total: res.list.total ?? 0,
          };
        } catch {
          return { items: [], total: 0 };
        }
      },
    },
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [AudienceGrid, audienceGridApi] = useVbenVxeGrid({ gridOptions: audienceGridOptions });
const [StatsGrid, statsGridApi] = useVbenVxeGrid({ gridOptions: statsGridOptions });

function resetSearchForm() {
  void searchFormApi.resetForm();
  statsLoaded.value = false;
}

function onTabChange(name: string | number) {
  if (name === 'stats' && !statsLoaded.value) {
    statsLoaded.value = true;
    void statsGridApi.reload();
  }
}

function onSearch() {
  void audienceGridApi.reload();
}

async function resetSearch() {
  resetSearchForm();
  if (activeTab.value === 'stats') {
    await statsGridApi.reload();
  } else {
    await audienceGridApi.reload();
  }
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    activeTab.value = 'audience';
    resetSearchForm();
    void audienceGridApi.reload();
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
    class="audience-list-dialog w-[1180px]"
    :title="dialogTitle"
  >
    <el-tabs v-model="activeTab" @tab-change="onTabChange">
      <el-tab-pane label="观众列表" name="audience">
        <div class="audience-panel">
          <SearchForm class="search-form" />
          <div class="search-actions">
            <el-button @click="resetSearch">重置</el-button>
            <el-button type="primary" @click="onSearch">查询</el-button>
          </div>

          <div class="summary-bar">
            <span>已登录观众人数：{{ summary.logged_in_count || 0 }}</span>
            <span v-if="summary.anonymous_count > 0" class="summary-bar__tip">
              还有 {{ summary.anonymous_count }} 位匿名用户未显示
            </span>
          </div>

          <AudienceGrid>
            <template #user="{ row }">
              <div class="user-cell">
                <img
                  v-if="row.avatar"
                  v-img-url="row.avatar"
                  alt=""
                  class="user-cell__avatar"
                />
                <div v-else class="user-cell__avatar user-cell__avatar--empty">-</div>
                <div class="user-cell__meta">
                  <div class="user-cell__name">{{ row.nick_name || '匿名用户' }}</div>
                  <div class="user-cell__id text-gray-400">ID：{{ row.user_id || '-' }}</div>
                </div>
              </div>
            </template>
            <template #session_inviter="{ row }">
              <div class="user-cell user-cell--compact">
                <template v-if="inviterUserId(row.session_inviter)">
                  <img
                    v-if="inviterAvatar(row.session_inviter)"
                    v-img-url="inviterAvatar(row.session_inviter)"
                    alt=""
                    class="user-cell__avatar user-cell__avatar--sm"
                  />
                  <div class="user-cell__meta">
                    <div class="user-cell__name">{{ inviterNickName(row.session_inviter) }}</div>
                    <div class="user-cell__id text-gray-400">
                      ID：{{ inviterUserId(row.session_inviter) }}
                    </div>
                  </div>
                </template>
                <span v-else class="text-gray-400">-</span>
              </div>
            </template>
            <template #first_inviter="{ row }">
              <div class="user-cell user-cell--compact">
                <template v-if="inviterUserId(row.first_inviter)">
                  <img
                    v-if="inviterAvatar(row.first_inviter)"
                    v-img-url="inviterAvatar(row.first_inviter)"
                    alt=""
                    class="user-cell__avatar user-cell__avatar--sm"
                  />
                  <div class="user-cell__meta">
                    <div class="user-cell__name">{{ inviterNickName(row.first_inviter) }}</div>
                    <div class="user-cell__id text-gray-400">
                      ID：{{ inviterUserId(row.first_inviter) }}
                    </div>
                  </div>
                </template>
                <span v-else class="text-gray-400">-</span>
              </div>
            </template>
            <template #first_enter_time="{ row }">
              {{ formatLiveTimeText(row.first_enter_time_text) }}
            </template>
            <template #last_active_time="{ row }">
              {{ formatLiveTimeText(row.last_active_time_text) }}
            </template>
            <template #watch_duration="{ row }">
              {{ row.watch_duration_text || '0' }}
            </template>
          </AudienceGrid>
        </div>
      </el-tab-pane>

      <el-tab-pane label="邀请统计" name="stats">
        <div class="audience-panel">
          <StatsGrid>
            <template #inviter="{ row }">
              <div class="user-cell">
                <img
                  v-if="row.inviter_avatar"
                  v-img-url="row.inviter_avatar"
                  alt=""
                  class="user-cell__avatar"
                />
                <div v-else class="user-cell__avatar user-cell__avatar--empty">-</div>
                <div class="user-cell__meta">
                  <div class="user-cell__name">{{ row.inviter_nick_name || '用户' }}</div>
                  <div class="user-cell__id text-gray-400">ID：{{ row.inviter_id || '-' }}</div>
                </div>
              </div>
            </template>
            <template #invite_watch_duration="{ row }">
              {{ row.invite_watch_duration_text || '0' }}
            </template>
            <template #order_pay_amount="{ row }">
              {{ row.order_pay_amount_text || '0.00' }}
            </template>
          </StatsGrid>
        </div>
      </el-tab-pane>
    </el-tabs>
  </Modal>
</template>

<style scoped lang="scss">
.audience-panel {
  min-height: 420px;
}

.search-form {
  margin-bottom: 8px;
}

.search-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-bottom: 12px;
}

.summary-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
  font-size: 13px;
  color: #303133;
}

.summary-bar__tip {
  color: #909399;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-cell--compact {
  gap: 8px;
}

.user-cell__avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.user-cell__avatar--sm {
  width: 32px;
  height: 32px;
}

.user-cell__avatar--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f2f3f5;
  color: #909399;
  font-size: 12px;
}

.user-cell__name {
  font-size: 13px;
  line-height: 1.4;
}

.user-cell__id {
  font-size: 12px;
  line-height: 1.4;
}
</style>
