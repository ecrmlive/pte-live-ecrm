<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

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
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createBusinessZoneAgent,
  fetchBusinessZoneAgentMerchants,
  fetchBusinessZoneAgents,
  revokeBusinessZoneAgent,
  resetBusinessZoneAgentPassword,
  updateBusinessZoneAgent,
  type BusinessZoneAgentRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { LIST_KEYWORD_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const merchantRows = ref<
  Array<{
    merchant_id: number;
    merchant_name: string;
    region_id: number;
    status: number;
  }>
>([]);
const resetTarget = ref<BusinessZoneAgentRow>();
const editingID = ref<number>();
const form = reactive({
  uid: 0,
  name: '',
  phone: '',
  qualification: '',
  remark: '',
  payment_method: 0,
  payment_name: '',
  payment_account: '',
  payment_bank: '',
  payment_qr_img: '',
  type: 0,
  business_name: '',
  business_store_category: 0,
  business_store_type: 0,
});
const passwordReset = reactive({
  password: '',
  confirmPassword: '',
  reason: '',
});

const statusText = (value: number) =>
  ({
    '-2': '已撤销',
    '-1': '已驳回',
    '0': '待审核',
    '1': '已通过',
  })[String(value)] || '未知';
const statusTag = (value: number) =>
  ({
    '-2': 'info',
    '-1': 'danger',
    '0': 'warning',
    '1': 'success',
  })[String(value)] || 'info';
const payText = (value: number) =>
  ['银行卡', '微信', '支付宝'][value] || '银行卡';
const dialogTitle = computed(() =>
  editingID.value ? '编辑代理申请' : '新增代理申请',
);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('姓名 / 手机号 / 商户名'),
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '已通过', value: 1 },
        { label: '已驳回', value: -1 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<BusinessZoneAgentRow> = {
  columns: [
    { field: 'circle_agent_id', title: 'ID', width: 72 },
    { field: 'name', minWidth: 110, title: '代理姓名' },
    { field: 'phone', title: '手机号', width: 140 },
    { field: 'business_name', minWidth: 140, title: '关联商户' },
    {
      field: 'payment_method',
      formatter: ({ cellValue }) => payText(Number(cellValue)),
      title: '结算方式',
      width: 100,
    },
    { field: 'balance', title: '佣金余额', width: 110 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '审核状态',
      width: 100,
    },
    platformListActionColumn({ minWidth: 270 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const statusRaw = formValues?.status;
        const result = await fetchBusinessZoneAgents({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            statusRaw === 0 || statusRaw === 1 || statusRaw === -1
              ? Number(statusRaw)
              : undefined,
        });
        return { items: result.list, total: result.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_agent_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [MerchantsDrawer, merchantsDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '代理关联商户',
});

const [ResetDrawer, resetDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确认重置',
  cancelText: '取消',
  placement: 'right',
  title: '重置区域代理后台密码',
  onConfirm: async () => submitPasswordReset(),
});

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  Object.assign(form, {
    uid: 0,
    name: '',
    phone: '',
    qualification: '',
    remark: '',
    payment_method: 0,
    payment_name: '',
    payment_account: '',
    payment_bank: '',
    payment_qr_img: '',
    type: 0,
    business_name: '',
    business_store_category: 0,
    business_store_type: 0,
  });
}

function openCreate() {
  editingID.value = undefined;
  resetForm();
  formDrawerApi.setState({ title: dialogTitle.value }).open();
}

function openEdit(row: BusinessZoneAgentRow) {
  if (row.status !== 0) {
    ElMessage.warning('已审核的代理资料不可编辑');
    return;
  }
  editingID.value = row.circle_agent_id;
  Object.assign(form, {
    uid: row.uid,
    name: row.name,
    phone: row.phone,
    qualification: row.qualification,
    remark: row.remark,
    payment_method: row.payment_method,
    payment_name: row.payment_name,
    payment_account: '',
    payment_bank: '',
    payment_qr_img: '',
    type: row.type,
    business_name: row.business_name,
    business_store_category: 0,
    business_store_type: 0,
  });
  formDrawerApi.setState({ title: dialogTitle.value }).open();
}

async function save() {
  if (!form.name.trim() || !form.phone.trim()) {
    ElMessage.warning('代理姓名和手机号必填');
    return;
  }
  formDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateBusinessZoneAgent(editingID.value, form);
    } else {
      await createBusinessZoneAgent(form);
    }
    formDrawerApi.close();
    ElMessage.success('已保存，新增申请需在代理审核中处理');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function revoke(row: BusinessZoneAgentRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      `撤销“${row.name}”的代理资格不会删除其历史审核、结算或关联事实。已关联区域或仍有佣金余额的代理不能撤销。`,
      '撤销代理资格',
      {
        inputPattern: /.{2,}/,
        inputErrorMessage: '撤销原因至少 2 个字符',
        confirmButtonText: '确认撤销',
        cancelButtonText: '取消',
        type: 'warning',
      },
    );
    await revokeBusinessZoneAgent(row.circle_agent_id, {
      reason: value.trim(),
      idempotency_key: `agent-revoke-${row.circle_agent_id}-${Date.now()}`,
    });
    ElMessage.success('代理资格已撤销');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

async function openMerchants(row: BusinessZoneAgentRow) {
  merchantsDrawerApi.open();
  merchantRows.value = [];
  try {
    merchantRows.value =
      (await fetchBusinessZoneAgentMerchants(row.circle_agent_id)).list || [];
  } catch {
    merchantsDrawerApi.close();
  }
}

function openPasswordReset(row: BusinessZoneAgentRow) {
  if (row.status !== 1) {
    ElMessage.warning('仅审核通过的代理可以重置后台密码');
    return;
  }
  resetTarget.value = row;
  Object.assign(passwordReset, { password: '', confirmPassword: '', reason: '' });
  resetDrawerApi.open();
}

async function submitPasswordReset() {
  const reason = passwordReset.reason.trim();
  const target = resetTarget.value;
  if (
    !target ||
    passwordReset.password.length < 12 ||
    passwordReset.password.length > 72 ||
    passwordReset.password !== passwordReset.confirmPassword ||
    reason.length < 2 ||
    reason.length > 500
  ) {
    ElMessage.warning(
      '请填写两次一致的 12 至 72 位新密码和 2 至 500 字的重置原因',
    );
    return;
  }
  await resetBusinessZoneAgentPassword(target.circle_agent_id, {
    password: passwordReset.password,
    reason,
    idempotency_key: `agent-password-${target.circle_agent_id}-${crypto.randomUUID()}`,
  });
  ElMessage.success('后台密码已重置，该代理旧后台会话已失效');
  resetDrawerApi.close();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增代理
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusTag(row.status)">{{ statusText(row.status) }}</ElTag>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="row.status === 0"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton link @click="openMerchants(row)">关联商户</ElButton>
        <ElButton
          v-if="row.status === 1"
          link
          type="warning"
          @click="openPasswordReset(row)"
        >
          重置密码
        </ElButton>
        <ElButton
          v-if="row.status === 1"
          link
          type="danger"
          @click="revoke(row)"
        >
          撤销
        </ElButton>
      </template>
    </Grid>

    <FormDrawer :title="dialogTitle">
      <ElForm label-width="110px">
        <ElFormItem label="关联用户ID">
          <ElInputNumber v-model="form.uid" :min="0" />
        </ElFormItem>
        <ElFormItem label="代理姓名" required>
          <ElInput v-model="form.name" />
        </ElFormItem>
        <ElFormItem label="联系电话" required>
          <ElInput v-model="form.phone" />
        </ElFormItem>
        <ElFormItem label="代理类型">
          <ElRadioGroup v-model="form.type">
            <ElRadio :value="0">区域代理</ElRadio>
            <ElRadio :value="1">商户型代理</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="关联商户">
          <ElInput v-model="form.business_name" placeholder="商户型代理填写" />
        </ElFormItem>
        <ElFormItem label="结算方式">
          <ElSelect v-model="form.payment_method">
            <ElOption label="银行卡" :value="0" />
            <ElOption label="微信" :value="1" />
            <ElOption label="支付宝" :value="2" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="结算名称">
          <ElInput v-model="form.payment_name" />
        </ElFormItem>
        <ElFormItem label="结算账号">
          <ElInput
            v-model="form.payment_account"
            autocomplete="off"
            :placeholder="
              editingID ? '留空保持原资料，系统不会回显' : '仅写入，不会在列表或详情回显'
            "
          />
        </ElFormItem>
        <ElFormItem label="开户行">
          <ElInput
            v-model="form.payment_bank"
            autocomplete="off"
            :placeholder="editingID ? '留空保持原资料' : '仅写入，不会回显'"
          />
        </ElFormItem>
        <ElFormItem label="资质材料">
          <ElInput v-model="form.qualification" type="textarea" />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model="form.remark" type="textarea" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <MerchantsDrawer>
      <ElTable :data="merchantRows" border>
        <ElTableColumn label="商户ID" prop="mer_id" width="100" />
        <ElTableColumn label="商户名称" min-width="180" prop="mer_name" />
        <ElTableColumn label="状态" width="100">
          <template #default="{ row }">
            <ElTag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '停用' }}
            </ElTag>
          </template>
        </ElTableColumn>
      </ElTable>
    </MerchantsDrawer>

    

    <ResetDrawer>
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="仅对已关联且启用的统一后台区域账号生效；不回显、不记录或回传密码，提交后旧后台会话立即失效。"
      />
      <ElForm label-width="108px" @submit.prevent="submitPasswordReset">
        <ElFormItem label="代理">
          <ElInput :model-value="resetTarget?.name || ''" disabled />
        </ElFormItem>
        <ElFormItem label="新密码" required>
          <ElInput
            v-model="passwordReset.password"
            autocomplete="new-password"
            show-password
            type="password"
          />
        </ElFormItem>
        <ElFormItem label="确认新密码" required>
          <ElInput
            v-model="passwordReset.confirmPassword"
            autocomplete="new-password"
            show-password
            type="password"
          />
        </ElFormItem>
        <ElFormItem label="重置原因" required>
          <ElInput
            v-model="passwordReset.reason"
            maxlength="500"
            show-word-limit
            type="textarea"
            :rows="3"
          />
        </ElFormItem>
      </ElForm>
    </ResetDrawer>
  </Page>
</template>
