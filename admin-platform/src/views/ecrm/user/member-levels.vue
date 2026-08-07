<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRadio,
  ElRadioGroup,
  ElTag,
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
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const saving = ref(false);
const editing = ref<MemberLevel>();
const form = reactive<MemberLevelInput>({
  name: '',
  rank: 1,
  rules: '{\n  "description": "满足成长规则后自动升级"\n}',
  benefits: '[\n  "会员专享活动"\n]',
  status: 1,
});

function prettyJSON(value: string) {
  try {
    return JSON.stringify(JSON.parse(value), null, 2);
  } catch {
    return value;
  }
}

function benefitText(value: string) {
  try {
    const list = JSON.parse(value);
    return Array.isArray(list) ? list.join('、') : '—';
  } catch {
    return '—';
  }
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    rank: 1,
    rules: '{\n  "description": "满足成长规则后自动升级"\n}',
    benefits: '[\n  "会员专享活动"\n]',
    status: 1,
    version: undefined,
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '等级名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MemberLevel> = {
  columns: [
    { field: 'name', minWidth: 140, showOverflow: false, title: '等级名称' },
    { field: 'rank', title: '等级排序', width: 100 },
    {
      field: 'benefits',
      formatter: ({ cellValue }) => benefitText(String(cellValue ?? '')),
      minWidth: 240,
      showOverflow: false,
      title: '权益',
    },
    { field: 'assigned_count', title: '当前用户数', width: 120 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await listMemberLevels()).list || [];
        if (keyword) {
          list = list.filter((row) => row.name.toLowerCase().includes(keyword));
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return { items: list, total: list.length };
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
  formDrawerApi.setState({ title: '新增会员等级' }).open();
}

function openEdit(row: MemberLevel) {
  editing.value = row;
  Object.assign(form, {
    name: row.name,
    rank: row.rank,
    rules: prettyJSON(row.rules),
    benefits: prettyJSON(row.benefits),
    status: row.status,
    version: row.version,
  });
  formDrawerApi.setState({ title: '编辑会员等级' }).open();
}

async function save() {
  if (!form.name.trim() || form.rank < 1) {
    ElMessage.warning('请填写等级名称与排序等级');
    return;
  }
  try {
    JSON.parse(form.rules);
    const benefits = JSON.parse(form.benefits);
    if (!Array.isArray(benefits) || !benefits.length) throw new Error('benefits');
  } catch {
    ElMessage.warning(
      '等级规则必须是 JSON 对象，会员权益必须是非空 JSON 字符串数组',
    );
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const data = {
      ...form,
      name: form.name.trim(),
      rules: prettyJSON(form.rules),
      benefits: prettyJSON(form.benefits),
    };
    if (editing.value) await updateMemberLevel(editing.value.id, data);
    else await createMemberLevel(data);
    formDrawerApi.close();
    ElMessage.success(
      '会员等级已保存；不会修改现有用户等级或历史变更记录',
    );
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function remove(row: MemberLevel) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”只会逻辑隐藏配置。若有用户正在使用或已有历史变更记录，服务端将拒绝删除。`,
      '删除会员等级',
      { type: 'warning' },
    );
    await deleteMemberLevel(row.id);
    ElMessage.success('会员等级已逻辑删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManage.value =
    profile.roles.includes('platform') &&
    codes.includes('user.member.level.manage');
  if (canManage.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canManage"
      class="mb-4"
      title="当前账号没有会员等级管理权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增等级
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton
          link
          type="danger"
          :disabled="row.assigned_count > 0"
          @click="remove(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="98px">
        <ElFormItem label="等级名称" required>
          <ElInput v-model="form.name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="等级排序" required>
          <ElInputNumber v-model="form.rank" :min="1" :max="10000" />
        </ElFormItem>
        <ElFormItem label="成长规则 JSON" required>
          <ElInput
            v-model="form.rules"
            type="textarea"
            :rows="6"
            class="font-mono"
          />
        </ElFormItem>
        <ElFormItem label="会员权益 JSON" required>
          <ElInput
            v-model="form.benefits"
            type="textarea"
            :rows="6"
            class="font-mono"
            placeholder='["权益一", "权益二"]'
          />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">启用</ElRadio>
            <ElRadio :value="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
