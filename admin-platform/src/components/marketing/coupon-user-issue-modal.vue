<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { ElAvatar, ElPagination } from 'element-plus';

import {
  listCouponReceiptRecords,
  type CouponReceiptRecord,
} from '#/api/core/platform-coupon-command';
import { formatShanghaiDateTime } from '#/utils/date-time';

export type CouponIssueMode = 'receive' | 'used';

const props = withDefaults(
  defineProps<{
    couponId?: number;
    couponScope?: 'platform' | 'store';
    mode?: CouponIssueMode;
  }>(),
  {
    couponId: 0,
    mode: 'used',
  },
);

const open = defineModel<boolean>('open', { default: false });

const loading = ref(false);
const rows = ref<CouponReceiptRecord[]>([]);
const pager = reactive({ page: 1, limit: 10, total: 0 });

const modalTitle = computed(() =>
  props.mode === 'used' ? '使用记录' : '领取记录',
);
const timeColumnTitle = computed(() =>
  props.mode === 'used' ? '使用时间' : '领取时间',
);

const [Modal, modalApi] = useVbenModal({
  class: 'w-[720px] max-w-[96vw]',
  contentClass: 'coupon-issue-modal__content',
  footer: false,
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

function formatTime(value?: string | null) {
  if (!value) return '—';
  const raw = String(value);
  if (
    raw.startsWith('0000-00-00') ||
    raw.startsWith('0001-01-01') ||
    raw === 'Invalid Date'
  ) {
    return '—';
  }
  return formatShanghaiDateTime(raw);
}

function rowTime(row: CouponReceiptRecord) {
  if (props.mode === 'used') {
    return formatTime(row.use_time || row.used_at || row.obtained_at);
  }
  return formatTime(row.obtained_at);
}

function displayName(row: CouponReceiptRecord) {
  return row.nickname || row.recipient || `用户#${row.user_id}`;
}

async function loadRows() {
  if (!props.couponId) return;
  loading.value = true;
  try {
    const data = await listCouponReceiptRecords({
      page: pager.page,
      limit: pager.limit,
      coupon_id: props.couponId,
      coupon_scope: props.couponScope,
      status: props.mode === 'used' ? 'used' : undefined,
    });
    rows.value = data.list || [];
    pager.total = data.total || 0;
  } finally {
    loading.value = false;
  }
}

function onPageChange(page: number) {
  pager.page = page;
  void loadRows();
}

watch(
  () => open.value,
  async (isOpen) => {
    if (!isOpen) {
      modalApi.close();
      return;
    }
    pager.page = 1;
    rows.value = [];
    modalApi.setState({ title: modalTitle.value }).open();
    await loadRows();
  },
);

watch(
  () => [props.couponId, props.mode, props.couponScope] as const,
  () => {
    if (!open.value) return;
    modalApi.setState({ title: modalTitle.value });
    pager.page = 1;
    void loadRows();
  },
);
</script>

<template>
  <Modal>
    <div v-loading="loading" class="issue-modal">
      <table class="issue-table">
        <thead>
          <tr>
            <th>用户名</th>
            <th>用户头像</th>
            <th>{{ timeColumnTitle }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id">
            <td>{{ displayName(row) }}</td>
            <td>
              <ElAvatar :size="36" :src="row.avatar_url || undefined">
                {{ displayName(row).slice(0, 1) }}
              </ElAvatar>
            </td>
            <td>{{ rowTime(row) }}</td>
          </tr>
          <tr v-if="!rows.length">
            <td colspan="3" class="issue-empty">暂无数据</td>
          </tr>
        </tbody>
      </table>
      <div class="issue-pager">
        <ElPagination
          background
          layout="prev, pager, next, jumper"
          :current-page="pager.page"
          :page-size="pager.limit"
          :total="pager.total"
          @current-change="onPageChange"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.issue-modal {
  min-height: 220px;
}

.issue-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.issue-table th,
.issue-table td {
  padding: 12px 10px;
  border-bottom: 1px solid hsl(var(--border));
  text-align: left;
  vertical-align: middle;
}

.issue-table th {
  color: hsl(var(--muted-foreground));
  font-weight: 600;
  background: hsl(var(--muted) / 0.35);
}

.issue-empty {
  padding: 36px 0 !important;
  color: hsl(var(--muted-foreground));
  text-align: center !important;
}

.issue-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
