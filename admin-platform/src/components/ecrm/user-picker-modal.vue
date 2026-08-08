<script setup lang="ts">
import { ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElAvatar,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElPagination,
  ElRadio,
  ElTable,
  ElTableColumn,
} from 'element-plus';

import {
  fetchPlatformUsers,
  type PlatformUserRow,
} from '#/api/core/ecrm';

export type PickedPlatformUser = {
  id: number;
  nickname: string;
  avatar_url: string;
  mobile: string;
};

const SOURCE_CHANNEL_LABEL: Record<string, string> = {
  wechat: '微信公众号',
  mini_program: '小程序',
  h5: 'H5',
  pc: 'PC',
  ios: 'iOS',
  android: 'Android',
  harmony: '鸿蒙',
};

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  select: [user: PickedPlatformUser];
}>();

const keyword = ref('');
const loading = ref(false);
const rows = ref<PlatformUserRow[]>([]);
const total = ref(0);
const page = ref(1);
const limit = ref(10);
const selectedId = ref(0);
const selectedRow = ref<PlatformUserRow | null>(null);

const [Modal, modalApi] = useVbenModal({
  title: '请选择用户',
  class: 'w-[860px] max-w-[96vw]',
  contentClass: 'user-picker-modal__content !overflow-hidden !p-0',
  confirmText: '确定',
  cancelText: '关闭',
  onConfirm: () => {
    if (!selectedRow.value?.id) {
      ElMessage.warning('请选择用户');
      return;
    }
    emit('select', {
      id: selectedRow.value.id,
      nickname: selectedRow.value.nickname || `用户#${selectedRow.value.id}`,
      avatar_url: selectedRow.value.avatar_url || '',
      mobile: selectedRow.value.mobile || '',
    });
    open.value = false;
  },
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

function platformLabel(channel?: string) {
  const key = (channel || '').trim();
  if (!key) return '—';
  return SOURCE_CHANNEL_LABEL[key] || key;
}

async function loadUsers() {
  loading.value = true;
  try {
    const result = await fetchPlatformUsers({
      page: page.value,
      limit: limit.value,
      keyword: keyword.value.trim() || undefined,
      // 有效用户：启用；注销/停用均为 status=0
      status: 1,
    });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function onSearch() {
  page.value = 1;
  void loadUsers();
}

function onPageChange(next: number) {
  page.value = next;
  void loadUsers();
}

function onRadioChange(row: PlatformUserRow) {
  selectedId.value = row.id;
  selectedRow.value = row;
}

function onRowClick(row: PlatformUserRow) {
  onRadioChange(row);
}

watch(open, (visible) => {
  if (visible) {
    keyword.value = '';
    page.value = 1;
    selectedId.value = 0;
    selectedRow.value = null;
    modalApi.open();
    void loadUsers();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal>
    <div class="user-picker">
      <div class="user-picker__filter">
        <ElForm inline @submit.prevent="onSearch">
          <ElFormItem label="用户搜索">
            <ElInput
              v-model="keyword"
              clearable
              class="picker-field"
              placeholder="请输入昵称/ID/手机号"
              @keyup.enter="onSearch"
            />
          </ElFormItem>
          <ElFormItem>
            <ElButton type="primary" @click="onSearch">查询</ElButton>
          </ElFormItem>
        </ElForm>
      </div>
      <div class="user-picker__table">
        <ElTable
          v-loading="loading"
          :data="rows"
          row-key="id"
          border
          height="100%"
          highlight-current-row
          @row-click="onRowClick"
        >
          <ElTableColumn label="" width="52" align="center">
            <template #default="{ row }">
              <ElRadio
                :model-value="selectedId"
                :value="row.id"
                @change="() => onRadioChange(row)"
                @click.stop
              >
                &nbsp;
              </ElRadio>
            </template>
          </ElTableColumn>
          <ElTableColumn label="ID" prop="id" width="90" />
          <ElTableColumn label="头像" width="80" align="center">
            <template #default="{ row }">
              <ElAvatar
                v-if="row.avatar_url"
                :size="36"
                :src="row.avatar_url"
                shape="circle"
              />
              <ElAvatar v-else :size="36">
                {{ (row.nickname || '?').slice(0, 1) }}
              </ElAvatar>
            </template>
          </ElTableColumn>
          <ElTableColumn label="昵称" prop="nickname" min-width="140" />
          <ElTableColumn label="手机号" prop="mobile" width="140" />
          <ElTableColumn label="用户平台" min-width="110">
            <template #default="{ row }">
              {{ platformLabel(row.source_channel) }}
            </template>
          </ElTableColumn>
        </ElTable>
      </div>
      <div class="user-picker__pager">
        <ElPagination
          background
          layout="total, prev, pager, next"
          :current-page="page"
          :page-size="limit"
          :total="total"
          @current-change="onPageChange"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.user-picker {
  display: flex;
  flex-direction: column;
  height: min(60vh, 520px);
  min-height: 360px;
  overflow: hidden;
}

.user-picker__filter {
  flex-shrink: 0;
  padding: 12px 12px 0;
}

.user-picker__filter :deep(.el-form--inline) {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  row-gap: 8px;
  margin-bottom: 0;
}

.user-picker__filter :deep(.el-form-item) {
  margin-bottom: 12px;
}

.picker-field {
  width: 260px;
}

.user-picker__table {
  flex: 1;
  min-height: 0;
  padding: 0 12px;
  overflow: hidden;
}

.user-picker__pager {
  display: flex;
  flex-shrink: 0;
  justify-content: flex-end;
  padding: 12px;
  border-top: 1px solid hsl(var(--border));
}
</style>

<style>
.user-picker-modal__content {
  display: flex;
  flex-direction: column;
}
</style>
