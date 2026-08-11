<script setup lang="ts">
import { ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import {
  ElAvatar,
  ElButton,
  ElCheckbox,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElPagination,
  ElRadio,
  ElTable,
  ElTableColumn,
} from 'element-plus';
import { UserFilled } from '@element-plus/icons-vue';

import {
  fetchPlatformUsers,
  type PlatformUserRow,
} from '#/api/core/ecrm';

export type PickedCustomerUser = {
  avatar_url: string;
  id: number;
  mobile: string;
  nickname: string;
};

const open = defineModel<boolean>('open', { default: false });
const props = withDefaults(defineProps<{ multiple?: boolean }>(), {
  multiple: false,
});

const emit = defineEmits<{
  select: [user: PickedCustomerUser];
  selectMultiple: [users: PickedCustomerUser[]];
}>();

const keyword = ref('');
const loading = ref(false);
const rows = ref<PlatformUserRow[]>([]);
const total = ref(0);
const page = ref(1);
const selectedID = ref(0);
const selectedRow = ref<PlatformUserRow>();
const selectedRows = ref<PlatformUserRow[]>([]);

const [Modal, modalApi] = useVbenModal({
  title: '请选择用户：',
  class: 'h-[min(76dvh,820px)] w-[min(94vw,1200px)] max-w-[94vw]',
  confirmText: '确定',
  cancelText: '关闭',
  onConfirm: () => {
    if (props.multiple) {
      if (!selectedRows.value.length) {
        ElMessage.warning('请选择用户');
        return;
      }
      emit(
        'selectMultiple',
        selectedRows.value.map((user) => ({
          id: user.id,
          nickname: user.nickname || `用户${user.id}`,
          avatar_url: user.avatar_url || '',
          mobile: user.mobile || '',
        })),
      );
      open.value = false;
      return;
    }
    if (!selectedRow.value) {
      ElMessage.warning('请选择用户');
      return;
    }
    emit('select', {
      id: selectedRow.value.id,
      nickname: selectedRow.value.nickname || `用户${selectedRow.value.id}`,
      avatar_url: selectedRow.value.avatar_url || '',
      mobile: selectedRow.value.mobile || '',
    });
    open.value = false;
  },
  onOpenChange(visible) {
    open.value = visible;
  },
});

async function loadUsers() {
  loading.value = true;
  try {
    const result = await fetchPlatformUsers({
      page: page.value,
      limit: 10,
      keyword: keyword.value.trim() || undefined,
      status: 1,
    });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  page.value = 1;
  void loadUsers();
}

function resetSearch() {
  keyword.value = '';
  search();
}

function select(row: PlatformUserRow) {
  if (props.multiple) {
    const index = selectedRows.value.findIndex((item) => item.id === row.id);
    if (index >= 0) selectedRows.value.splice(index, 1);
    else selectedRows.value.push(row);
    return;
  }
  selectedID.value = row.id;
  selectedRow.value = row;
}

function selected(row: PlatformUserRow) {
  return selectedRows.value.some((item) => item.id === row.id);
}

watch(open, (visible) => {
  if (!visible) {
    modalApi.close();
    return;
  }
  keyword.value = '';
  page.value = 1;
  selectedID.value = 0;
  selectedRow.value = undefined;
  selectedRows.value = [];
  modalApi.open();
  void loadUsers();
});
</script>

<template>
  <Modal>
    <div class="customer-user-picker">
      <ElForm inline class="customer-user-picker__search" @submit.prevent="search">
        <ElFormItem label="用户搜索：">
          <ElInput
            v-model="keyword"
            class="customer-user-picker__keyword"
            clearable
            placeholder="请输入用户昵称/ID/手机号搜索"
            @keyup.enter="search"
          />
        </ElFormItem>
        <ElFormItem>
          <ElButton @click="resetSearch">重置</ElButton>
          <ElButton type="primary" @click="search">搜索</ElButton>
        </ElFormItem>
      </ElForm>

      <div class="customer-user-picker__table">
        <ElTable
          v-loading="loading"
          :data="rows"
          height="100%"
          highlight-current-row
          row-key="id"
          @row-click="select"
        >
          <ElTableColumn align="center" width="64">
            <template #default="{ row }">
              <ElCheckbox
                v-if="props.multiple"
                :model-value="selected(row)"
                @change="select(row)"
                @click.stop
              />
              <ElRadio
                v-else
                :model-value="selectedID"
                :value="row.id"
                @change="select(row)"
                @click.stop
              >&nbsp;</ElRadio>
            </template>
          </ElTableColumn>
          <ElTableColumn label="ID" prop="id" width="180" />
          <ElTableColumn label="微信用户名" min-width="250" prop="nickname" />
          <ElTableColumn align="center" label="客服头像" min-width="180">
            <template #default="{ row }">
              <ElAvatar :icon="UserFilled" :size="48" :src="row.avatar_url || undefined" shape="square" />
            </template>
          </ElTableColumn>
        </ElTable>
      </div>

      <div class="customer-user-picker__pager">
        <ElPagination
          background
          layout="total, prev, pager, next"
          :current-page="page"
          :page-size="10"
          :total="total"
          @current-change="(next) => { page = next; void loadUsers(); }"
        />
      </div>
    </div>
  </Modal>
</template>

<style scoped>
.customer-user-picker {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  padding: 8px 12px 0;
  overflow: hidden;
}

.customer-user-picker__search {
  flex: 0 0 auto;
}

.customer-user-picker__table {
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
}

.customer-user-picker__keyword {
  width: 320px;
}

.customer-user-picker__pager {
  display: flex;
  flex: 0 0 auto;
  justify-content: flex-end;
  padding: 16px 0;
}
</style>
