<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElPagination,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  auditCommunityPostApi,
  deleteCommunityPostApi,
  deleteCommunityReplyApi,
  getCommunityPostApi,
  listCommunityPostsApi,
  listCommunityRepliesApi,
  type CommunityPost,
  type CommunityReply,
} from '#/api/core/platform-community';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { LIST_KEYWORD_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canAudit = ref(false);
const canDelete = ref(false);
const saving = ref(false);
const current = ref<CommunityPost>();
const replies = ref<CommunityReply[]>([]);
const repliesTotal = ref(0);
const replyQuery = reactive({ limit: 20, page: 1 });
const rejectForm = reactive({ refusal: '' });

function statusInfo(status: number) {
  if (status === 1) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  return { label: '待审核', type: 'warning' as const };
}

function canAuditRow(row: CommunityPost) {
  return canAudit.value && row.status === 0;
}

function canUpdatePresentation(row: CommunityPost) {
  return canAudit.value && row.status === 1;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('标题或正文'),
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '已驳回', value: -1 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<CommunityPost> = {
  columns: [
    {
      field: 'title',
      minWidth: 200,
      showOverflow: false,
      title: '标题',
    },
    { field: 'nickname', minWidth: 110, title: '作者' },
    { field: 'mer_id', title: '商户 ID', width: 100 },
    { field: 'topic_name', minWidth: 120, title: '话题' },
    { field: 'count_reply', title: '评论', width: 76 },
    { field: 'pv', title: '浏览', width: 76 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '审核',
      width: 104,
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '展示',
      width: 76,
    },
    {
      field: 'is_hot',
      slots: { default: 'is_hot' },
      title: '置顶',
      width: 76,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '发布时间',
    },
    platformListActionColumn({ minWidth: 290 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const statusRaw = formValues?.status;
        const result = await listCommunityPostsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            statusRaw === 0 || statusRaw === 1 || statusRaw === -1
              ? Number(statusRaw)
              : undefined,
        });
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
  footer: false,
  title: '社区帖子详情',
});

const [RejectModal, rejectModalApi] = useVbenModal({
  onConfirm: async () => reject(),
});

async function loadReplies() {
  if (!current.value) return;
  const result = await listCommunityRepliesApi(
    current.value.community_id,
    replyQuery,
  );
  replies.value = result.list || [];
  repliesTotal.value = result.total || 0;
}

async function openDetail(row: CommunityPost) {
  replyQuery.page = 1;
  const [post, replyResult] = await Promise.all([
    getCommunityPostApi(row.community_id),
    listCommunityRepliesApi(row.community_id, replyQuery),
  ]);
  current.value = post;
  replies.value = replyResult.list || [];
  repliesTotal.value = replyResult.total || 0;
  detailDrawerApi.open();
}

async function approve(row: CommunityPost) {
  try {
    await ElMessageBox.confirm(
      '确认通过该社区帖子审核？通过后将按当前显示状态在 C 端生效。',
      '审核通过确认',
      {
        cancelButtonText: '取消',
        confirmButtonText: '确认通过',
        type: 'warning',
      },
    );
    await auditCommunityPostApi(row.community_id, {
      is_show: row.is_show === 1 ? 1 : 0,
      status: 1,
    });
    ElMessage.success('帖子已审核通过');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

function openReject(row: CommunityPost) {
  current.value = row;
  rejectForm.refusal = '';
  rejectModalApi.open();
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return false;
  }
  if (!current.value) return false;
  saving.value = true;
  try {
    await auditCommunityPostApi(current.value.community_id, {
      refusal,
      status: -1,
    });
    rejectModalApi.close();
    ElMessage.success('帖子已驳回');
    gridApi.reload();
    return true;
  } finally {
    saving.value = false;
  }
}

async function updatePresentation(
  row: CommunityPost,
  data: { is_hot?: 0 | 1; is_show?: 0 | 1 },
  action: string,
) {
  try {
    await ElMessageBox.confirm(
      `确认${action}该已审核通过的帖子？`,
      `${action}确认`,
      {
        cancelButtonText: '取消',
        confirmButtonText: `确认${action}`,
        type: 'warning',
      },
    );
    await auditCommunityPostApi(row.community_id, { ...data, status: 0 });
    ElMessage.success(`帖子已${action}`);
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function deletePost(row: CommunityPost) {
  try {
    await ElMessageBox.confirm(
      '确认删除该帖子？删除后会从 C 端隐藏且不可恢复。',
      '删除帖子',
      {
        cancelButtonText: '取消',
        confirmButtonText: '确认删除',
        type: 'error',
      },
    );
    await deleteCommunityPostApi(row.community_id);
    ElMessage.success('帖子已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function deleteReply(row: CommunityReply) {
  try {
    await ElMessageBox.confirm('确认删除该评论？删除后不可恢复。', '删除评论', {
      cancelButtonText: '取消',
      confirmButtonText: '确认删除',
      type: 'error',
    });
    await deleteCommunityReplyApi(row.reply_id);
    ElMessage.success('评论已删除');
    await loadReplies();
    if (current.value) {
      current.value = await getCommunityPostApi(current.value.community_id);
    }
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const isContentOperator =
    profile.roles.includes('platform') || profile.roles.includes('operations');
  canAudit.value =
    isContentOperator && permissions.includes('content.community.audit');
  canDelete.value =
    isContentOperator && permissions.includes('content.community.delete');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>
      <template #is_hot="{ row }">
        <ElTag :type="row.is_hot === 1 ? 'warning' : 'info'">
          {{ row.is_hot === 1 ? '是' : '否' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <template v-if="canAuditRow(row)">
          <ElButton link type="success" @click="approve(row)">通过</ElButton>
          <ElButton link type="danger" @click="openReject(row)">驳回</ElButton>
        </template>
        <template v-else-if="canUpdatePresentation(row)">
          <ElButton
            link
            type="warning"
            @click="
              updatePresentation(
                row,
                { is_hot: row.is_hot === 1 ? 0 : 1 },
                row.is_hot === 1 ? '取消置顶' : '置顶',
              )
            "
          >
            {{ row.is_hot === 1 ? '取消置顶' : '置顶' }}
          </ElButton>
          <ElButton
            link
            type="warning"
            @click="
              updatePresentation(
                row,
                { is_show: row.is_show === 1 ? 0 : 1 },
                row.is_show === 1 ? '隐藏' : '显示',
              )
            "
          >
            {{ row.is_show === 1 ? '隐藏' : '显示' }}
          </ElButton>
        </template>
        <ElButton v-if="canDelete" link type="danger" @click="deletePost(row)">
          删除
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer class="w-[680px]">
      <template v-if="current">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="标题" :span="2">
            {{ current.title }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="作者">
            {{ current.nickname || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="商户 ID">
            {{ current.mer_id || '平台内容' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="分类">
            {{ current.cate_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="话题">
            {{ current.topic_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="关联商品">
            {{
              current.product_name ||
              (current.product_id ? `商品 #${current.product_id}` : '—')
            }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="审核状态">
            <ElTag :type="statusInfo(current.status).type">
              {{ statusInfo(current.status).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="发布时间">
            {{ formatShanghaiDateTime(current.create_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="审核时间">
            {{ current.status_time ? formatShanghaiDateTime(current.status_time) : '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="正文">
            <div class="whitespace-pre-wrap">{{ current.content }}</div>
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.image" :span="2" label="图片">
            <ElImage
              :preview-src-list="[current.image]"
              :src="current.image"
              class="h-24 w-24"
              fit="cover"
              preview-teleported
            />
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.refusal" :span="2" label="驳回原因">
            {{ current.refusal }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div class="mb-3 mt-6 text-base font-medium">
          评论（{{ repliesTotal }}）
        </div>
        <ElTable :data="replies" border empty-text="暂无评论">
          <ElTableColumn label="用户" min-width="110" prop="nickname" />
          <ElTableColumn
            label="内容"
            min-width="220"
            prop="content"
            show-overflow-tooltip
          />
          <ElTableColumn label="时间" min-width="170" prop="create_time">
            <template #default="{ row }">
              {{ formatShanghaiDateTime(row.create_time) }}
            </template>
          </ElTableColumn>
          <ElTableColumn
            v-if="canDelete"
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
      </template>
    </DetailDrawer>

    <RejectModal title="驳回帖子">
      <ElForm label-width="84px">
        <ElFormItem label="驳回原因" required>
          <ElInput
            v-model="rejectForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请填写可供发布者查看的驳回原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </RejectModal>
  </Page>
</template>
