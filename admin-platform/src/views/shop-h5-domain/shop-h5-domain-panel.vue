<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';
import type { H5DomainRow } from '#/api/core/h5-domain';

import { useAccess } from '@vben/access';
import { useVbenDrawer } from '@vben/common-ui';
import { Plus } from '@element-plus/icons-vue';
import { ElButton, ElMessage, ElTag } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenForm, z } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import PlatformH5DomainApi from '#/api/core/platform-h5-domain';
import {
  certSyncTagType,
  deployDisplayTagType,
  enableStatusTagType,
} from '#/api/core/h5-domain';
import {
  PLATFORM_LIST_GRID_LAYOUT,
  platformListActionColumn,
} from '#/constants/platform-list-grid';
import { overlayConfirm } from '#/utils/overlay-confirm';

const props = defineProps<{
  appId: number;
  appName?: string;
}>();

const { hasAccessByCodes } = useAccess();

const DEFAULT_GATEWAY_CNAME = 'cname.qxkejiwl.top';
const gatewayCname = ref('');
const maxCount = ref(5);
const totalCount = ref(0);
const addOpen = ref(false);
const saving = ref(false);
const actionLoadingKey = ref('');

async function fetchDomainPage(pageSize: number, currentPage: number) {
  if (props.appId <= 0) {
    return { items: [], total: 0 };
  }
  const res = await PlatformH5DomainApi.index(
    {
      app_id: props.appId,
      page: currentPage,
      list_rows: pageSize,
    },
    true,
  );
  if (res.code === 1) {
    gatewayCname.value = res.data?.gateway_cname || DEFAULT_GATEWAY_CNAME;
    maxCount.value = res.data?.max_count ?? 5;
    totalCount.value = res.data?.list?.total ?? 0;
    return {
      items: res.data?.list?.data ?? [],
      total: totalCount.value,
    };
  }
  return { items: [], total: 0 };
}

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { placeholder: 'qxkejiwl.top' },
    description: `将域名 CNAME 解析到 ${gatewayCname.value || DEFAULT_GATEWAY_CNAME}`,
    fieldName: 'domain',
    label: '域名',
    rules: 'required',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '直播', value: 1 },
        { label: '商城', value: 2 },
      ],
    },
    defaultValue: 1,
    fieldName: 'category',
    label: '分类',
  },
  {
    component: 'RadioGroup',
    componentProps: {
      options: [
        { label: '自有证书', value: 1 },
        { label: '免费证书', value: 2 },
      ],
    },
    defaultValue: 1,
    fieldName: 'cert_type',
    label: '证书类型',
  },
  {
    component: 'Textarea',
    componentProps: { placeholder: '请输入密钥', rows: 6 },
    dependencies: {
      if: (values) => Number(values.cert_type) === 1,
      rules(values) {
        if (Number(values.cert_type) !== 1) {
          return z.string().optional();
        }
        return z.string().min(1, '请输入密钥');
      },
      triggerFields: ['cert_type'],
    },
    fieldName: 'cert_key',
    formItemClass: 'col-span-full',
    label: '密钥(KEY)',
  },
  {
    component: 'Textarea',
    componentProps: { placeholder: '请输入证书 PEM', rows: 6 },
    dependencies: {
      if: (values) => Number(values.cert_type) === 1,
      rules(values) {
        if (Number(values.cert_type) !== 1) {
          return z.string().optional();
        }
        return z.string().min(1, '请输入证书 PEM');
      },
      triggerFields: ['cert_type'],
    },
    fieldName: 'cert_pem',
    formItemClass: 'col-span-full',
    label: '证书(PEM)',
  },
]);

const [AddForm, addFormApi] = useVbenForm(
  reactive({
    commonConfig: {
      componentProps: { class: 'w-full' },
      labelWidth: 100,
    },
    handleSubmit: async (values) => {
      if (props.appId <= 0) {
        return;
      }
      saving.value = true;
      try {
        const certType = Number(values.cert_type ?? 1);
        const res = await PlatformH5DomainApi.add(
          {
            app_id: props.appId,
            category: Number(values.category ?? 1),
            cert_key: certType === 1 ? String(values.cert_key ?? '').trim() : '',
            cert_pem: certType === 1 ? String(values.cert_pem ?? '').trim() : '',
            cert_type: certType,
            domain: String(values.domain ?? '').trim(),
          },
          true,
        );
        if (res.code === 1) {
          ElMessage.success(res.msg || '添加成功');
          addOpen.value = false;
          reload();
        }
      } finally {
        saving.value = false;
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
    submitButtonOptions: { show: false },
  }),
);

const [AddDrawer, addDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    addOpen.value = isOpen;
  },
});

const gridOptions = {
  ...PLATFORM_LIST_GRID_LAYOUT,
  border: true,
  columns: [
    { field: 'domain_id', title: 'ID', width: 64 },
    {
      field: 'domain',
      minWidth: 160,
      showOverflow: true,
      title: '域名',
    },
    { field: 'category_text', title: '分类', width: 72 },
    { field: 'cert_type_text', title: '证书类型', width: 88 },
    {
      field: 'enable_status',
      slots: { default: 'enableStatus' },
      title: '证书启用',
      width: 88,
    },
    {
      field: 'cert_sync_status',
      slots: { default: 'certSyncStatus' },
      title: '证书更新',
      width: 88,
    },
    {
      field: 'deploy_display_status',
      slots: { default: 'deployDisplayStatus' },
      title: '部署状态',
      width: 88,
    },
    { field: 'create_time', minWidth: 150, title: '添加时间' },
    platformListActionColumn({ width: 200 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 30] },
  proxyConfig: {
    ajax: {
      query: async ({ page }) =>
        fetchDomainPage(page.pageSize, page.currentPage),
    },
  },
  rowConfig: { isHover: true, keyField: 'domain_id' },
  toolbarConfig: { refresh: true },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function reload() {
  gridApi.reload();
}

function openAdd() {
  if (totalCount.value >= maxCount.value) {
    ElMessage.warning(`最多添加 ${maxCount.value} 条域名`);
    return;
  }
  void addFormApi.setValues({
    category: 1,
    cert_key: '',
    cert_pem: '',
    cert_type: 1,
    domain: '',
  });
  addOpen.value = true;
}

async function submitAdd() {
  await addFormApi.validateAndSubmitForm();
}

watch(addOpen, (visible) => {
  if (visible) {
    addDrawerApi.open();
    return;
  }
  addDrawerApi.close();
});

function actionKey(prefix: string, domainId: number) {
  return `${prefix}-${domainId}`;
}

function isActionLoading(prefix: string, domainId: number) {
  return actionLoadingKey.value === actionKey(prefix, domainId);
}

const DEPLOY_FAIL_HINT = '部署失败，请联系技术部';

async function runDomainAction(
  key: string,
  task: () => Promise<{ code: number; msg?: string; data?: unknown } | void>,
  successMsg?: string,
  failHint?: string,
) {
  if (actionLoadingKey.value) {
    return;
  }
  actionLoadingKey.value = key;
  try {
    const res = await task();
    if (res && res.code === 1) {
      ElMessage.success(successMsg || res.msg || '操作成功');
      reload();
    }
  } catch {
    if (failHint) {
      ElMessage.error(failHint);
    }
    reload();
  } finally {
    actionLoadingKey.value = '';
  }
}

async function deployRow(row: H5DomainRow, isUpdate: boolean) {
  const tip = isUpdate
    ? `确认更新部署域名「${row.domain}」吗？将使用最新证书写入 nginx 并 reload。`
    : row.cert_type === 2
      ? `确认部署域名「${row.domain}」吗？将写入 nginx 配置并 reload。`
      : `确认部署域名「${row.domain}」吗？将使用商户证书写入 nginx 并 reload。`;
  try {
    await overlayConfirm(tip, isUpdate ? '更新部署' : '部署', {
      type: 'warning',
    });
  } catch {
    return;
  }
  await runDomainAction(
    actionKey(isUpdate ? 'update' : 'deploy', row.domain_id),
    () => PlatformH5DomainApi.deploy(props.appId, row.domain_id, true),
    isUpdate ? '更新部署成功' : '部署成功，域名已启用',
    DEPLOY_FAIL_HINT,
  );
}

async function disableRow(row: H5DomainRow) {
  try {
    await overlayConfirm(
      `确认停用域名「${row.domain}」吗？将从服务器移除该域名的 nginx 配置并 reload。`,
      '停用',
      { type: 'warning' },
    );
  } catch {
    return;
  }
  await runDomainAction(
    actionKey('disable', row.domain_id),
    () => PlatformH5DomainApi.disable(props.appId, row.domain_id, true),
    '已停用',
    '操作失败，请联系技术部',
  );
}

async function deleteRow(row: H5DomainRow) {
  try {
    await overlayConfirm(
      '删除仅移除域名数据记录，不会操作 nginx。请确认该域名已在平台停用。',
      '删除',
      { type: 'warning' },
    );
  } catch {
    return;
  }
  await runDomainAction(
    actionKey('delete', row.domain_id),
    () => PlatformH5DomainApi.delete(props.appId, row.domain_id, true),
    '删除成功',
  );
}

function needsDeploy(row: H5DomainRow) {
  return (
    row.deploy_display_status === 'not_deployed' ||
    row.deploy_display_status === 'failed' ||
    row.deploy_status === 'offline' ||
    row.deploy_status === 'pending'
  );
}

function needsUpdateDeploy(row: H5DomainRow) {
  return row.deploy_display_status === 'updated' || row.cert_sync_status === 'outdated';
}

function canDisable(row: H5DomainRow) {
  return row.deploy_status === 'deployed' || row.deploy_display_status === 'updated';
}

function canDelete(row: H5DomainRow) {
  if (
    row.deploy_display_status === 'deploying' ||
    row.deploy_status === 'deploying'
  ) {
    return false;
  }
  return row.can_delete === true;
}

watch(
  () => props.appId,
  (appId) => {
    if (appId > 0) {
      reload();
    }
  },
);
</script>

<template>
  <div class="shop-h5-domain-panel">
    <Grid class="shop-h5-domain-panel__grid">
      <template #toolbar-actions>
        <ElButton
          v-if="hasAccessByCodes(['platform:shopH5Domain:add'])"
          :icon="Plus"
          type="primary"
          @click="openAdd"
        >
          添加域名
        </ElButton>
        <span class="ml-3 text-xs text-muted-foreground">
          最多添加 {{ maxCount }} 条记录
        </span>
      </template>

      <template #toolbar-tools>
        <span class="text-xs text-muted-foreground">
          部署/更新部署=写 nginx 并启用；停用=移除 nginx；删除=仅删数据
        </span>
      </template>

      <template #enableStatus="{ row }">
        <ElTag :type="enableStatusTagType(row.enable_status)" size="small">
          {{ row.enable_status_text || '未启用' }}
        </ElTag>
        <ElTag
          v-if="row.biz_disabled_pending"
          class="ml-1"
          size="small"
          type="warning"
        >
          商户已停用
        </ElTag>
      </template>

      <template #certSyncStatus="{ row }">
        <ElTag :type="certSyncTagType(row.cert_sync_status)" size="small">
          {{ row.cert_sync_status_text || '—' }}
        </ElTag>
      </template>

      <template #deployDisplayStatus="{ row }">
        <div class="deploy-status-cell">
          <ElTag :type="deployDisplayTagType(row.deploy_display_status)" size="small">
            {{ row.deploy_display_text || '未部署' }}
          </ElTag>
          <span
            v-if="row.deploy_display_status === 'failed'"
            class="deploy-status-cell__hint"
          >
            {{ DEPLOY_FAIL_HINT }}
          </span>
        </div>
      </template>

      <template #action="{ row }">
        <ElButton
          v-if="hasAccessByCodes(['platform:shopH5Domain:deploy']) && needsDeploy(row)"
          :loading="isActionLoading('deploy', row.domain_id)"
          link
          size="small"
          type="primary"
          @click="deployRow(row, false)"
        >
          部署
        </ElButton>
        <ElButton
          v-if="hasAccessByCodes(['platform:shopH5Domain:deploy']) && needsUpdateDeploy(row)"
          :loading="isActionLoading('update', row.domain_id)"
          link
          size="small"
          type="primary"
          @click="deployRow(row, true)"
        >
          更新部署
        </ElButton>
        <ElButton
          v-if="hasAccessByCodes(['platform:shopH5Domain:disable']) && canDisable(row)"
          :loading="isActionLoading('disable', row.domain_id)"
          link
          size="small"
          type="warning"
          @click="disableRow(row)"
        >
          停用
        </ElButton>
        <ElButton
          v-if="hasAccessByCodes(['platform:shopH5Domain:delete']) && canDelete(row)"
          :loading="isActionLoading('delete', row.domain_id)"
          link
          size="small"
          type="danger"
          @click="deleteRow(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <AddDrawer
      :close-on-click-modal="false"
      :destroy-on-close="true"
      class="w-[640px]"
      title="添加域名"
    >
      <AddForm />
      <template #footer>
        <ElButton @click="addOpen = false">取消</ElButton>
        <ElButton :loading="saving" type="primary" @click="submitAdd">提交</ElButton>
      </template>
    </AddDrawer>
  </div>
</template>

<style scoped>
.shop-h5-domain-panel__grid :deep(.h-full.rounded-md.bg-card),
.shop-h5-domain-panel__grid :deep(.vxe-grid) {
  height: auto !important;
  min-height: 0 !important;
}

.shop-h5-domain-panel__grid :deep(.vxe-table--body-inner-wrapper),
.shop-h5-domain-panel__grid :deep(.vxe-table--body-wrapper),
.shop-h5-domain-panel__grid :deep(.vxe-grid--layout-body-content-wrapper) {
  height: auto !important;
  min-height: 0 !important;
  max-height: none !important;
}

.deploy-status-cell {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  max-width: 100%;
}

.deploy-status-cell__hint {
  font-size: 11px;
  line-height: 1.35;
  color: var(--el-color-danger);
}
</style>
