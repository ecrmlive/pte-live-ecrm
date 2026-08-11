<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createMemberLevel,
  deleteMemberLevel,
  listMemberLevels,
  updateMemberLevel,
  type MemberLevel,
  type MemberLevelInput,
} from '#/api/core/platform-member-level';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const canRead = ref(false);
const canManage = ref(false);
const showTips = ref(true);
const saving = ref(false);
const editing = ref<MemberLevel>();
const form = reactive<MemberLevelInput>({
  name: '',
  rank: 1,
  icon_url: '',
  growth_value: 0,
  bg_image: '',
  status: 1,
});

function iconSrc(url: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    rank: 1,
    icon_url: '',
    growth_value: 0,
    bg_image: '',
    status: 1,
    version: undefined,
  });
}

const gridOptions: VxeGridProps<MemberLevel> = {
  columns: [
    { field: 'id', title: 'ID', width: 88 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      title: '名称',
    },
    {
      field: 'icon_url',
      slots: { default: 'icon' },
      title: '等级图标',
      width: 100,
    },
    {
      field: 'user_count',
      title: '人数',
      width: 100,
    },
    {
      field: 'growth_value',
      title: '所需成长值',
      width: 120,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 130 }),
  ],
  emptyText: '暂无数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const list = (await listMemberLevels()).list || [];
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
        };
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增会员等级' }).open();
}

function openEdit(row: MemberLevel) {
  editing.value = row;
  Object.assign(form, {
    name: row.name,
    rank: row.rank,
    icon_url: row.icon_url || '',
    growth_value: Number(row.growth_value || 0),
    bg_image: row.bg_image || '',
    status: row.status,
    version: row.version,
  });
  formDrawerApi.setState({ title: '编辑会员等级' }).open();
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请填写会员名称');
    return;
  }
  if (!form.rank || form.rank < 1) {
    ElMessage.warning('请填写会员等级');
    return;
  }
  if (!String(form.icon_url || '').trim()) {
    ElMessage.warning('请选择会员图标');
    return;
  }
  if (form.growth_value == null || form.growth_value < 0) {
    ElMessage.warning('请输入所需成长值');
    return;
  }
  if (!String(form.bg_image || '').trim()) {
    ElMessage.warning('请上传背景图');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const data: MemberLevelInput = {
      name,
      rank: form.rank,
      icon_url: String(form.icon_url).trim(),
      growth_value: Number(form.growth_value),
      bg_image: String(form.bg_image).trim(),
      status: 1,
      version: form.version,
    };
    if (editing.value) await updateMemberLevel(editing.value.id, data);
    else await createMemberLevel(data);
    formDrawerApi.close();
    ElMessage.success(editing.value ? '会员等级已更新' : '会员等级已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function remove(row: MemberLevel) {
  if (Number(row.user_count || 0) > 0) {
    ElMessage.warning('该等级下有用户，不能删除');
    return;
  }
  try {
    await confirm({
      content: `确定删除等级「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteMemberLevel(row.id);
    ElMessage.success('会员等级已删除');
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
    (codes.includes('user.member.level.read') ||
      codes.includes('user.member.level.manage'));
  canManage.value = roleOK && codes.includes('user.member.level.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有会员等级查看权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <ElAlert
        v-if="showTips"
        class="mb-4"
        type="warning"
        show-icon
        closable
        @close="showTips = false"
      >
        用户经验值仅可累加、不做扣减，累计数值达到对应等级要求时，将自动完成等级提升
      </ElAlert>

      <Grid>
        <template #toolbar-actions>
          <ElButton
            v-if="canManage"
            :icon="Plus"
            type="primary"
            @click="openCreate"
          >
            新增会员等级
          </ElButton>
        </template>
        <template #icon="{ row }">
          <ElImage
            v-if="row.icon_url"
            class="level-list-icon"
            :src="iconSrc(row.icon_url)"
            fit="contain"
            alt="等级图标"
          >
            <template #error>
              <span>—</span>
            </template>
          </ElImage>
          <span v-else>—</span>
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
    </template>

    <FormDrawer>
      <ElForm label-width="120px" require-asterisk-position="left">
        <ElFormItem label="会员等级" required>
          <ElInputNumber
            v-model="form.rank"
            :min="1"
            :max="10000"
            :controls="false"
            class="w-full"
            placeholder="请输入会员等级"
          />
        </ElFormItem>
        <ElFormItem label="会员名称" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            show-word-limit
            placeholder="请输入会员名称"
          />
        </ElFormItem>
        <ElFormItem label="会员图标" required>
          <ImageField
            v-model="form.icon_url"
            :preview-size="72"
            default-library="system"
          />
        </ElFormItem>
        <ElFormItem label="所需成长值" required>
          <ElInputNumber
            v-model="form.growth_value"
            :min="0"
            :max="100000000"
            :controls="false"
            class="w-full"
            placeholder="请输入所需成长值"
          />
        </ElFormItem>
        <ElFormItem label="背景图" required>
          <ImageField
            v-model="form.bg_image"
            :preview-size="120"
            default-library="system"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.level-list-icon {
  width: 40px;
  height: 40px;
  border-radius: 6px;
}
</style>
