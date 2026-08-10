<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
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
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createSvipInterest,
  deleteSvipInterest,
  listSvipInterests,
  updateSvipInterest,
  updateSvipInterestStatus,
  type SvipInterest,
  type SvipInterestInput,
} from '#/api/core/platform-svip-interest';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const canRead = ref(false);
const canManage = ref(false);
const saving = ref(false);
const editing = ref<SvipInterest>();
const form = reactive<SvipInterestInput>({
  name: '',
  display_name: '',
  description: '',
  icon_url: '',
  on_icon_url: '',
  link: '',
  status: 1,
  sort: 0,
});

function iconSrc(url: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    display_name: '',
    description: '',
    icon_url: '',
    on_icon_url: '',
    link: '',
    status: 1,
    sort: 0,
    version: undefined,
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '权益名称 / 展示名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启', value: 1 },
        { label: '关闭', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '权益状态',
  },
]);

const gridOptions: VxeGridProps<SvipInterest> = {
  columns: [
    { field: 'name', minWidth: 120, showOverflow: false, title: '权益名称' },
    {
      field: 'display_name',
      minWidth: 120,
      showOverflow: false,
      title: '展示名称',
    },
    {
      field: 'icon_url',
      slots: { default: 'icon_off' },
      title: '未开通权益图标',
      width: 120,
    },
    {
      field: 'on_icon_url',
      slots: { default: 'icon_on' },
      title: '已开通权益图标',
      width: 120,
    },
    {
      className: 'col--remark',
      field: 'description',
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '权益简介',
      width: 220,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '权益状态',
      width: 110,
    },
    platformListActionColumn({ width: 130 }),
  ],
  emptyText: '暂无数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await listSvipInterests()).list || [];
        if (keyword) {
          list = list.filter(
            (row) =>
              row.name.toLowerCase().includes(keyword) ||
              String(row.display_name || '')
                .toLowerCase()
                .includes(keyword),
          );
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加会员权益' }).open();
}

function openEdit(row: SvipInterest) {
  editing.value = row;
  Object.assign(form, {
    name: row.name,
    display_name: row.display_name || '',
    description: row.description || '',
    icon_url: row.icon_url || '',
    on_icon_url: row.on_icon_url || '',
    link: row.link || '',
    status: row.status,
    sort: row.sort,
    version: row.version,
  });
  formDrawerApi.setState({ title: '编辑会员权益' }).open();
}

async function save() {
  const name = form.name.trim();
  const displayName = form.display_name.trim();
  const description = form.description.trim();
  const icon = String(form.icon_url || '').trim();
  const onIcon = String(form.on_icon_url || '').trim();
  if (!name) {
    ElMessage.warning('请填写权益名称');
    return;
  }
  if (!displayName) {
    ElMessage.warning('请填写展示名称');
    return;
  }
  if (!description) {
    ElMessage.warning('请填写权益简介');
    return;
  }
  if (!icon) {
    ElMessage.warning('请选择未开通权益图标');
    return;
  }
  if (!onIcon) {
    ElMessage.warning('请选择已开通权益图标');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const data: SvipInterestInput = {
      name,
      display_name: displayName,
      description,
      icon_url: icon,
      on_icon_url: onIcon,
      link: String(form.link || '').trim(),
      status: form.status ? 1 : 0,
      sort: Number(form.sort || 0),
      version: form.version,
    };
    if (editing.value) await updateSvipInterest(editing.value.id, data);
    else await createSvipInterest(data);
    formDrawerApi.close();
    ElMessage.success(editing.value ? '会员权益已更新' : '会员权益已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: SvipInterest, enabled: boolean) {
  if (!canManage.value) return;
  const before = row.status;
  const next = enabled ? 1 : 0;
  row.status = next;
  try {
    await updateSvipInterestStatus(row.id, {
      status: next,
      version: row.version,
    });
    row.version = Number(row.version || 0) + 1;
    ElMessage.success(enabled ? '已开启' : '已关闭');
  } catch {
    row.status = before;
  }
}

async function remove(row: SvipInterest) {
  try {
    await confirm({
      content: `确定删除权益「${row.name}」吗？若仍被启用会员类型使用将被拒绝。`,
      icon: 'warning',
      title: '提示',
    });
    await deleteSvipInterest(row.id);
    ElMessage.success('会员权益已删除');
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
    (codes.includes('user.svip.interest.read') ||
      codes.includes('user.svip.interest.manage'));
  canManage.value = roleOK && codes.includes('user.svip.interest.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有会员权益查看权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          添加会员权益
        </ElButton>
      </template>
      <template #icon_off="{ row }">
        <ElImage
          v-if="row.icon_url"
          class="equity-list-icon"
          :src="iconSrc(row.icon_url)"
          fit="contain"
          alt="未开通权益图标"
        >
          <template #error>
            <span>—</span>
          </template>
        </ElImage>
        <span v-else>—</span>
      </template>
      <template #icon_on="{ row }">
        <ElImage
          v-if="row.on_icon_url"
          class="equity-list-icon"
          :src="iconSrc(row.on_icon_url)"
          fit="contain"
          alt="已开通权益图标"
        >
          <template #error>
            <span>—</span>
          </template>
        </ElImage>
        <span v-else>—</span>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="开启"
          inactive-text="关闭"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
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

    <FormDrawer>
      <ElForm label-width="120px" require-asterisk-position="left">
        <ElFormItem label="权益名称" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            show-word-limit
            placeholder="请输入权益名称"
          />
        </ElFormItem>
        <ElFormItem label="展示名称" required>
          <ElInput
            v-model="form.display_name"
            maxlength="64"
            show-word-limit
            placeholder="请输入展示名称"
          />
        </ElFormItem>
        <ElFormItem label="权益简介" required>
          <ElInput
            v-model="form.description"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
            placeholder="请输入权益简介"
          />
        </ElFormItem>
        <ElFormItem label="未开通权益图标" required>
          <ImageField
            v-model="form.icon_url"
            :preview-size="80"
            default-library="system"
          />
        </ElFormItem>
        <ElFormItem label="已开通权益图标" required>
          <ImageField
            v-model="form.on_icon_url"
            :preview-size="80"
            default-library="system"
          />
        </ElFormItem>
        <ElFormItem label="跳转内部链接">
          <ElInput
            v-model="form.link"
            maxlength="500"
            placeholder="可选，C 端跳转路径或链接"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="999999"
            :controls="false"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="权益状态">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="开启"
            inactive-text="关闭"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.equity-list-icon {
  width: 48px;
  height: 48px;
  border-radius: 6px;
}
</style>
