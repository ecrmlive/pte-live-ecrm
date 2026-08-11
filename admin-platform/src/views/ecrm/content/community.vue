<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { ArrowDown } from '@element-plus/icons-vue';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElForm,
  ElFormItem,
  ElIcon,
  ElImage,
  ElInput,
  ElMessage,
  ElPagination,
  ElRadio,
  ElRadioGroup,
  ElRate,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  auditCommunityPostApi,
  deleteCommunityPostApi,
  deleteCommunityReplyApi,
  getCommunityPostApi,
  listCommunityCategoriesApi,
  listCommunityPostsApi,
  listCommunityRepliesApi,
  listCommunityTopicsApi,
  updateCommunityPostShowApi,
  updateCommunityPostStarApi,
  type CommunityCategory,
  type CommunityPost,
  type CommunityReply,
  type CommunityTopic,
} from '#/api/core/platform-community';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import {
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const canRead = ref(false);
const canManage = ref(false);
const categories = ref<CommunityCategory[]>([]);
const topics = ref<CommunityTopic[]>([]);
const typeTab = ref<1 | 2>(1);
const tabCounts = reactive({ image: 0, video: 0 });
const current = ref<CommunityPost>();
const detailTab = ref('content');
const replies = ref<CommunityReply[]>([]);
const repliesTotal = ref(0);
const replyQuery = reactive({ limit: 10, page: 1 });
const starForm = reactive({ start: 1 });
const forceForm = reactive({ refusal: '' });
const auditForm = reactive({ status: 1 as 1 | -1, refusal: '' });
const starTarget = ref<CommunityPost>();
const forceTarget = ref<CommunityPost>();
const auditTarget = ref<CommunityPost>();

function statusInfo(status: number) {
  if (status === 1) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  if (status === -2) return { label: '强制下架', type: 'info' as const };
  return { label: '待审核', type: 'warning' as const };
}

function mediaUrl(url?: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function splitImages(image?: string) {
  return String(image || '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s) => mediaUrl(s))
    .filter(Boolean);
}

function authorText(row: CommunityPost) {
  const nick = row.nickname || '用户';
  return row.uid ? `${nick} | ${row.uid}` : nick;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '已驳回', value: -1 },
        { label: '强制下架', value: -2 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '审核状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'category_id',
    label: '分类名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'topic_id',
    label: '话题名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '显示', value: 1 },
        { label: '隐藏', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'is_show',
    label: '是否显示',
  },
  listUserSearchFormField({
    fieldName: 'author_search',
    label: '作者搜索',
  }),
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入标题关键字',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
]);

const gridOptions: VxeGridProps<CommunityPost> = {
  columns: [
    { field: 'community_id', title: 'ID', width: 80 },
    {
      field: 'title',
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '标题',
    },
    {
      field: 'nickname',
      formatter: ({ row }) => authorText(row),
      minWidth: 130,
      showOverflow: 'tooltip',
      title: '作者',
    },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '图文封面',
      width: 140,
    },
    {
      field: 'start',
      slots: { default: 'start' },
      title: '推荐级别',
      width: 140,
    },
    { field: 'pv', title: '浏览量', width: 80 },
    { field: 'count_start', title: '点赞数', width: 80 },
    { field: 'count_reply', title: '评论数', width: 80 },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 170,
      title: '发布时间',
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '是否显示',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const statusRaw = formValues?.status;
        const showRaw = formValues?.is_show;
        const authorSearch = parseUserSearch(formValues, 'author_search');
        const result = await listCommunityPostsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            statusRaw === 0 ||
            statusRaw === 1 ||
            statusRaw === -1 ||
            statusRaw === -2
              ? Number(statusRaw)
              : undefined,
          category_id: formValues?.category_id
            ? Number(formValues.category_id)
            : undefined,
          topic_id: formValues?.topic_id
            ? Number(formValues.topic_id)
            : undefined,
          is_show:
            showRaw === 0 || showRaw === 1 ? Number(showRaw) : undefined,
          is_type: typeTab.value,
          author: authorSearch.keyword || undefined,
          author_type: authorSearch.type || 'nickname',
        });
        tabCounts.image = result.image_count || 0;
        tabCounts.video = result.video_count || 0;
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'community_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

const [StarDrawer, starDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => saveStar(),
});

const [ForceDrawer, forceDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => saveForceOff(),
});

const [AuditDrawer, auditDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => submitAudit(),
});

function setTypeTab(tab: 1 | 2) {
  if (typeTab.value === tab) return;
  typeTab.value = tab;
  gridApi.reload();
}

async function syncFilterOptions() {
  await gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'category_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: categories.value.map((c) => ({
          label: c.cate_name,
          value: c.category_id,
        })),
        placeholder: '请选择',
      },
    },
    {
      fieldName: 'topic_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: topics.value.map((t) => ({
          label: t.topic_name,
          value: t.topic_id,
        })),
        placeholder: '请选择',
      },
    },
  ]);
}

async function loadReplies() {
  if (!current.value) return;
  const result = await listCommunityRepliesApi(
    current.value.community_id,
    replyQuery,
  );
  replies.value = result.list || [];
  repliesTotal.value = result.total || 0;
}

async function openDetail(row: CommunityPost, tab: 'content' | 'comments' = 'content') {
  replyQuery.page = 1;
  detailTab.value = tab;
  const [post, replyResult] = await Promise.all([
    getCommunityPostApi(row.community_id),
    listCommunityRepliesApi(row.community_id, replyQuery),
  ]);
  current.value = post;
  replies.value = replyResult.list || [];
  repliesTotal.value = replyResult.total || 0;
  detailDrawerApi.setState({ title: '内容详情' }).open();
}

function openStar(row: CommunityPost) {
  starTarget.value = row;
  starForm.start = Number(row.start || 1);
  starDrawerApi.setState({ title: '编辑星级' }).open();
}

async function saveStar() {
  if (!starTarget.value) return false;
  if (starForm.start < 1 || starForm.start > 5) {
    ElMessage.warning('请选择 1-5 星');
    return false;
  }
  await updateCommunityPostStarApi(starTarget.value.community_id, starForm.start);
  ElMessage.success('星级已更新');
  starDrawerApi.close();
  gridApi.reload();
  return true;
}

function openForceOff(row: CommunityPost) {
  forceTarget.value = row;
  forceForm.refusal = '';
  forceDrawerApi.setState({ title: '强制下架' }).open();
}

async function saveForceOff() {
  const refusal = forceForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写下架理由');
    return false;
  }
  if (!forceTarget.value) return false;
  await auditCommunityPostApi(forceTarget.value.community_id, {
    status: -2,
    refusal,
  });
  ElMessage.success('已强制下架');
  forceDrawerApi.close();
  gridApi.reload();
  return true;
}

async function toggleShow(row: CommunityPost, enabled: boolean) {
  try {
    await updateCommunityPostShowApi(row.community_id, enabled ? 1 : 0);
    ElMessage.success(enabled ? '已显示' : '已隐藏');
    gridApi.reload();
  } catch {
    /* requestClient 已提示 */
  }
}

function openAudit(row: CommunityPost) {
  auditTarget.value = row;
  auditForm.status = 1;
  auditForm.refusal = '';
  auditDrawerApi.setState({ title: '审核内容' }).open();
}

async function submitAudit() {
  if (!auditTarget.value) return false;
  if (auditForm.status === -1 && !auditForm.refusal.trim()) {
    ElMessage.warning('请填写拒绝理由');
    return false;
  }
  await auditCommunityPostApi(auditTarget.value.community_id, {
    status: auditForm.status,
    is_show: auditForm.status === 1 ? 1 : 0,
    refusal:
      auditForm.status === -1 ? auditForm.refusal.trim() : undefined,
  });
  ElMessage.success(auditForm.status === 1 ? '已审核通过' : '已驳回');
  auditDrawerApi.close();
  gridApi.reload();
  return true;
}

function onMoreCommand(command: string, row: CommunityPost) {
  switch (command) {
    case 'audit':
      if (canAudit(row)) openAudit(row);
      break;
    case 'star':
      if (canEditStar(row)) openStar(row);
      break;
    case 'comments':
      openDetail(row, 'comments');
      break;
    case 'forceOff':
      if (canForceOff(row)) openForceOff(row);
      break;
    case 'delete':
      if (canManage.value) deletePost(row);
      break;
    default:
      break;
  }
}

async function deletePost(row: CommunityPost) {
  try {
    await confirm({
      title: '提示',
      content: '确认删除该内容？删除后不可恢复。',
      icon: 'warning',
    });
    await deleteCommunityPostApi(row.community_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

async function deleteReply(row: CommunityReply) {
  try {
    await confirm({
      title: '提示',
      content: '确认删除该评论？',
      icon: 'warning',
    });
    await deleteCommunityReplyApi(row.reply_id);
    ElMessage.success('评论已删除');
    await loadReplies();
    if (current.value) {
      current.value = await getCommunityPostApi(current.value.community_id);
    }
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

function canAudit(row: CommunityPost) {
  return canManage.value && row.status === 0;
}

function canForceOff(row: CommunityPost) {
  return canManage.value && row.status === 1;
}

function canEditStar(row: CommunityPost) {
  return canManage.value && (row.status === 1 || row.status === -2);
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
    (codes.includes('content.community_list.read') ||
      codes.includes('content.community_list.manage') ||
      codes.includes('content.community.audit') ||
      codes.includes('content.community.delete'));
  canManage.value =
    roleOK &&
    (codes.includes('content.community_list.manage') ||
      codes.includes('content.community.audit') ||
      codes.includes('content.community.delete'));
  if (!canRead.value) return;
  const [cateRes, topicRes] = await Promise.all([
    listCommunityCategoriesApi().catch(() => ({ list: [] as CommunityCategory[] })),
    listCommunityTopicsApi().catch(() => ({ list: [] as CommunityTopic[] })),
  ]);
  categories.value = cateRes.list || [];
  topics.value = topicRes.list || [];
  await syncFilterOptions();
  gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="community-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="community-tabs__item"
            :aria-selected="typeTab === 1"
            :class="{ 'is-active': typeTab === 1 }"
            @click="setTypeTab(1)"
          >
            图文列表({{ tabCounts.image }})
          </button>
          <button
            type="button"
            role="tab"
            class="community-tabs__item"
            :aria-selected="typeTab === 2"
            :class="{ 'is-active': typeTab === 2 }"
            @click="setTypeTab(2)"
          >
            短视频列表({{ tabCounts.video }})
          </button>
        </div>
      </template>

      <template #image="{ row }">
        <div class="cover-list">
          <ElImage
            v-for="(src, idx) in splitImages(row.image).slice(0, 3)"
            :key="`${row.community_id}-${idx}`"
            :src="src"
            :preview-src-list="splitImages(row.image)"
            :initial-index="idx"
            fit="cover"
            class="cover-thumb"
            preview-teleported
          >
            <template #error>
              <div class="cover-thumb cover-thumb--empty">无图</div>
            </template>
          </ElImage>
          <span v-if="!splitImages(row.image).length" class="text-xs text-gray-400">
            —
          </span>
        </div>
      </template>

      <template #start="{ row }">
        <ElRate :model-value="Number(row.start || 0)" :max="5" disabled />
      </template>

      <template #is_show="{ row }">
        <ElSwitch
          :model-value="row.is_show === 1"
          :disabled="!canManage || row.status !== 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              toggleShow(row, Boolean(enabled))
          "
        />
      </template>

      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row, 'content')">
          详情
        </ElButton>
        <ElDropdown
          trigger="click"
          @command="(cmd: string) => onMoreCommand(cmd, row)"
        >
          <ElButton link type="primary">
            更多
            <ElIcon class="el-icon--right"><ArrowDown /></ElIcon>
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem v-if="canAudit(row)" command="audit">
                审核
              </ElDropdownItem>
              <ElDropdownItem v-if="canEditStar(row)" command="star">
                编辑星级
              </ElDropdownItem>
              <ElDropdownItem command="comments">查看评论</ElDropdownItem>
              <ElDropdownItem
                v-if="canForceOff(row)"
                command="forceOff"
                divided
              >
                强制下架
              </ElDropdownItem>
              <ElDropdownItem
                v-if="canManage"
                command="delete"
                :divided="!canForceOff(row)"
              >
                删除
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="current">
        <div class="detail-head">
          <div class="detail-head__meta">
            <span>内容ID：{{ current.community_id }}</span>
            <ElTag :type="statusInfo(current.status).type" class="ml-2">
              {{ statusInfo(current.status).label }}
            </ElTag>
          </div>
          <div class="detail-head__sub">
            发布时间：{{ formatShanghaiDateTime(current.create_time) || '—' }}
            <template v-if="current.refusal">
              ；拒绝原因：{{ current.refusal }}
            </template>
          </div>
        </div>

        <ElTabs v-model="detailTab">
          <ElTabPane label="发布内容" name="content">
            <ElDescriptions :column="2" border>
              <ElDescriptionsItem label="作者">
                {{ authorText(current) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="参与话题">
                {{ current.topic_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem :span="2" label="内容图片">
                <div class="cover-list">
                  <ElImage
                    v-for="(src, idx) in splitImages(current.image)"
                    :key="`detail-${idx}`"
                    :src="src"
                    :preview-src-list="splitImages(current.image)"
                    :initial-index="idx"
                    fit="cover"
                    class="cover-thumb cover-thumb--lg"
                    preview-teleported
                  />
                  <span
                    v-if="!splitImages(current.image).length"
                    class="text-xs text-gray-400"
                  >
                    —
                  </span>
                </div>
              </ElDescriptionsItem>
              <ElDescriptionsItem
                v-if="current.is_type === 2"
                :span="2"
                label="视频链接"
              >
                {{ current.video_link || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem :span="2" label="文章内容">
                <div class="whitespace-pre-wrap">{{ current.content }}</div>
              </ElDescriptionsItem>
            </ElDescriptions>

            <div class="mb-3 mt-6 text-base font-medium">关联商品</div>
            <ElTable
              :data="
                current.product_id
                  ? [
                      {
                        product_id: current.product_id,
                        product_name: current.product_name,
                        product_price: current.product_price,
                      },
                    ]
                  : []
              "
              border
              empty-text="暂无关联商品"
            >
              <ElTableColumn label="商品ID" prop="product_id" width="100" />
              <ElTableColumn
                label="商品名称"
                min-width="180"
                prop="product_name"
                show-overflow-tooltip
              />
              <ElTableColumn label="价格" width="110">
                <template #default="{ row }">
                  ¥{{ Number(row.product_price || 0).toFixed(2) }}
                </template>
              </ElTableColumn>
            </ElTable>
          </ElTabPane>

          <ElTabPane label="评论内容" name="comments">
            <ElTable :data="replies" border empty-text="暂无评论">
              <ElTableColumn label="ID" prop="reply_id" width="80" />
              <ElTableColumn label="用户名|ID" min-width="140">
                <template #default="{ row }">
                  {{ row.nickname || '用户' }} | {{ row.uid }}
                </template>
              </ElTableColumn>
              <ElTableColumn
                label="评论内容"
                min-width="220"
                prop="content"
                show-overflow-tooltip
              />
              <ElTableColumn label="评论时间" min-width="170">
                <template #default="{ row }">
                  {{ formatShanghaiDateTime(row.create_time) || '—' }}
                </template>
              </ElTableColumn>
              <ElTableColumn
                v-if="canManage"
                fixed="right"
                label="操作"
                width="72"
              >
                <template #default="{ row }">
                  <ElButton link type="danger" @click="deleteReply(row)">
                    删除
                  </ElButton>
                </template>
              </ElTableColumn>
            </ElTable>
            <div class="mt-3 flex justify-end">
              <ElPagination
                small
                :current-page="replyQuery.page"
                :page-size="replyQuery.limit"
                :total="repliesTotal"
                layout="prev, pager, next"
                @current-change="
                  (page: number) => {
                    replyQuery.page = page;
                    loadReplies();
                  }
                "
              />
            </div>
          </ElTabPane>
        </ElTabs>
      </template>
    </DetailDrawer>

    <StarDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="推荐星级" required>
          <ElRate v-model="starForm.start" :max="5" />
        </ElFormItem>
      </ElForm>
    </StarDrawer>

    <ForceDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="下架理由" required>
          <ElInput
            v-model="forceForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请输入下架理由"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </ForceDrawer>

    <AuditDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="审核状态" required>
          <ElRadioGroup v-model="auditForm.status">
            <ElRadio :label="1">通过</ElRadio>
            <ElRadio :label="-1">拒绝</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="auditForm.status === -1" label="拒绝理由" required>
          <ElInput
            v-model="auditForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请填写拒绝理由"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </AuditDrawer>
  </Page>
</template>

<style scoped>
.community-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.community-tabs__item {
  padding: 6px 14px;
  color: hsl(var(--foreground) / 70%);
  cursor: pointer;
  background: transparent;
  border: 1px solid hsl(var(--border));
  border-radius: 4px;
}

.community-tabs__item.is-active {
  color: hsl(var(--primary));
  background: hsl(var(--primary) / 8%);
  border-color: hsl(var(--primary));
}

.cover-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
}

.cover-thumb {
  width: 44px;
  height: 44px;
  border-radius: 4px;
}

.cover-thumb--lg {
  width: 72px;
  height: 72px;
}

.cover-thumb--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #999;
  background: #f5f5f5;
}

.detail-head {
  margin-bottom: 16px;
}

.detail-head__meta {
  display: flex;
  align-items: center;
  margin-bottom: 6px;
  font-weight: 600;
}

.detail-head__sub {
  font-size: 13px;
  color: hsl(var(--foreground) / 65%);
}
</style>
