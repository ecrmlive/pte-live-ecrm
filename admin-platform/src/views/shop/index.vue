<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { useAccess } from '@vben/access';

import {
  ElButton,
  ElCheckbox,
  ElLink,
  ElMessage,
  ElMessageBox,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import ShopApi from '#/api/core/shop';
import LiveTrafficApi from '#/api/core/live-traffic';
import { Page } from '@vben/common-ui';
import { deepClone } from '#/utils/base';
import { formatUnixTime } from '#/utils/live-token-refresh';
import { parseApiListPage } from '#/utils/list-response';

import ShopAddModal from './shop-add-modal.vue';
import ShopEditModal from './shop-edit-modal.vue';
import ShopH5DomainManageModal from './shop-h5-domain-manage-modal.vue';
import ShopTrafficManageModal from './shop-traffic-manage-modal.vue';
import ShopUserManageModal from './shop-user-manage-modal.vue';
import type { ShopRow } from './types';

const { hasAccessByCodes } = useAccess();

const addModalOpen = ref(false);
const editModalOpen = ref(false);
const shopUserModalOpen = ref(false);
const h5DomainModalOpen = ref(false);
const trafficModalOpen = ref(false);
const curModel = ref<ShopRow | undefined>();
const manageShop = ref<ShopRow | undefined>();

function filterRows(rows: ShopRow[], formValues?: Record<string, unknown>) {
  const idKw = String(formValues?.appId ?? '').trim();
  const nameKw = String(formValues?.keyword ?? '')
    .trim()
    .toLowerCase();
  return rows.filter((row) => {
    if (idKw && !String(row.app_id).includes(idKw)) {
      return false;
    }
    if (!nameKw) {
      return true;
    }
    return (
      row.app_name?.toLowerCase().includes(nameKw) ||
      row.user_name?.toLowerCase().includes(nameKw) ||
      row.platform_operator_name?.toLowerCase().includes(nameKw)
    );
  });
}

function targetLinkAddress() {
  return `https://${window.location.hostname}/shop/#/login`;
}

function fmtRemainGB(v: null | number | string | undefined) {
  if (v == null || v === '') return '—';
  const n = Number(v || 0);
  return n.toFixed(n >= 100 ? 1 : 2);
}

function shopAdminNick(row: Pick<ShopRow, 'real_name' | 'user_name'>) {
  const nick = row.real_name?.trim();
  return nick || row.user_name || '—';
}

function formatCreateTime(ts?: null | number | string) {
  const n = Number(ts);
  if (!Number.isFinite(n) || n <= 0) {
    return '—';
  }
  return formatUnixTime(n);
}

function formatExpireTime(row: Pick<ShopRow, 'expire_time' | 'expire_time_text'>) {
  const text = row.expire_time_text?.trim();
  if (text) {
    return text;
  }
  if (!row.expire_time) {
    return '永不过期';
  }
  return formatUnixTime(row.expire_time).slice(0, 10);
}

async function fetchShopRows(
  pageSize: number,
  currentPage: number,
  formValues?: Record<string, unknown>,
) {
  const res = await ShopApi.shopList(
    { list_rows: pageSize, page: currentPage },
    true,
  );
  const page = parseApiListPage<ShopRow>(res.data);
  const shops = page.list;

  const rows = await Promise.all(
    shops.map(async (shop) => {
      try {
        const accRes = await LiveTrafficApi.account({ app_id: shop.app_id });
        return {
          ...shop,
          remain_gb: accRes.data?.remain_gb ?? null,
        };
      } catch {
        return { ...shop, remain_gb: null };
      }
    }),
  );

  const items = rows.map((row) => ({
    ...row,
    is_recycle: (row as ShopRow & { is_recycle: number }).is_recycle === 0,
    weixin_service:
      (row as ShopRow & { weixin_service: number }).weixin_service === 1,
  })) as ShopRow[];

  return { items: filterRows(items, formValues), total: page.total };
}

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '筛选 ID' },
      fieldName: 'appId',
      formItemClass: 'pb-0',
      label: '商城ID',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '筛选名称' },
      fieldName: 'keyword',
      formItemClass: 'pb-0',
      label: '商城名称',
    },
  ],
};

const gridOptions = {
  border: true,
  columns: [
    { field: 'app_id', minWidth: 90, title: '商城ID' },
    { field: 'app_name', minWidth: 120, showOverflow: true, title: '商城名称' },
    { field: 'user_name', minWidth: 110, showOverflow: true, title: '超管账号' },
    {
      field: 'real_name',
      minWidth: 120,
      slots: { default: 'adminNick' },
      title: '超管昵称',
    },
    {
      field: 'is_recycle',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'weixin_service',
      slots: { default: 'wxPay' },
      title: '微信服务商支付',
      width: 140,
    },
    {
      field: 'expire_time',
      minWidth: 120,
      slots: { default: 'expire' },
      title: '过期时间',
    },
    {
      field: 'remain_gb',
      slots: { default: 'remainGb' },
      title: '剩余流量(GB)',
      width: 120,
    },
    {
      field: 'create_time',
      minWidth: 160,
      slots: { default: 'createTime' },
      title: '添加时间',
    },
    {
      field: 'platform_operator_name',
      formatter: ({ cellValue }) => String(cellValue ?? '') || '—',
      title: '操作员',
      width: 100,
    },
    {
      className: 'shop-action-cell',
      fixed: 'right',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'action' },
      title: '操作',
    },
  ],
  cellConfig: { verticalAlign: 'top' },
  pagerConfig: { enabled: true, pageSize: 15, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchShopRows(page.pageSize, page.currentPage, formValues),
    },
  },
  rowConfig: { isHover: true, keyField: 'app_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function reload() {
  gridApi.reload();
}

function openAdd() {
  addModalOpen.value = true;
}

function openEdit(row: ShopRow) {
  const model = deepClone(row) as ShopRow;
  if (model.expire_time === 0) {
    model.expire_time_text = '';
    model.no_expire = true;
  } else {
    model.no_expire = false;
  }
  curModel.value = model;
  editModalOpen.value = true;
}

function openTrafficManage(row: ShopRow) {
  manageShop.value = row;
  trafficModalOpen.value = true;
}

function openShopUsers(row: ShopRow) {
  manageShop.value = row;
  shopUserModalOpen.value = true;
}

function openShopH5Domains(row: ShopRow) {
  manageShop.value = row;
  h5DomainModalOpen.value = true;
}

async function statusChange(checked: boolean, row: ShopRow) {
  try {
    await ShopApi.updateStatus({ app_id: row.app_id }, true);
    row.is_recycle = checked;
  } catch {
    row.is_recycle = !checked;
  }
}

async function wxStatusChange(checked: boolean, row: ShopRow) {
  try {
    await ShopApi.updateWxStatus({ app_id: row.app_id }, true);
    row.weixin_service = checked;
  } catch {
    row.weixin_service = !checked;
  }
}

async function deleteShop(row: ShopRow) {
  try {
    await ElMessageBox.confirm('删除后不可恢复，确认删除该记录吗?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  const res = await ShopApi.deleteShop({ app_id: row.app_id }, true);
  if (res.code === 1) {
    ElMessage.success(res.msg || '删除成功');
    reload();
  }
}
</script>

<template>
  <Page>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:shop:add'"
          :icon="Plus"
          type="primary"
          @click="openAdd"
        >
          添加商城
        </ElButton>
      </template>

      <template #adminNick="{ row }">
        {{ shopAdminNick(row) }}
      </template>

      <template #status="{ row }">
        <ElCheckbox
          v-if="hasAccessByCodes(['platform:shop:status'])"
          v-model="row.is_recycle"
          @change="(checked) => statusChange(Boolean(checked), row)"
        >
          启用
        </ElCheckbox>
        <span v-else>{{ row.is_recycle ? '启用' : '禁用' }}</span>
      </template>

      <template #wxPay="{ row }">
        <ElCheckbox
          v-if="hasAccessByCodes(['platform:shop:wxStatus'])"
          v-model="row.weixin_service"
          @change="(checked) => wxStatusChange(Boolean(checked), row)"
        >
          启用
        </ElCheckbox>
        <span v-else>{{ row.weixin_service ? '启用' : '禁用' }}</span>
      </template>

      <template #expire="{ row }">
        {{ formatExpireTime(row) }}
      </template>

      <template #remainGb="{ row }">
        <span :class="Number(row.remain_gb) < 0 ? 'text-[#e6a23c]' : ''">
          {{ fmtRemainGB(row.remain_gb) }}
        </span>
      </template>

      <template #createTime="{ row }">
        {{ formatCreateTime(row.create_time) }}
      </template>

      <template #action="{ row }">
        <div class="shop-op-col">
          <div class="shop-op-row">
            <ElLink
              :href="targetLinkAddress()"
              size="small"
              target="_blank"
              type="primary"
            >
              进入商城
            </ElLink>
            <ElLink
              v-access:code="'platform:shopUser:list'"
              class="shop-op-link"
              size="small"
              type="primary"
              @click="openShopUsers(row)"
            >
              账号管理
            </ElLink>
            <ElLink
              v-access:code="'platform:shopH5Domain:list'"
              class="shop-op-link"
              size="small"
              type="primary"
              @click="openShopH5Domains(row)"
            >
              域名管理
            </ElLink>
          </div>
          <div class="shop-op-row">
            <ElLink
              class="shop-op-link"
              size="small"
              type="primary"
              @click="openTrafficManage(row)"
            >
              流量管理
            </ElLink>
            <ElLink
              v-access:code="'platform:shop:edit'"
              class="shop-op-link"
              size="small"
              type="primary"
              @click="openEdit(row)"
            >
              编辑
            </ElLink>
            <ElLink
              v-access:code="'platform:shop:delete'"
              class="shop-op-link"
              size="small"
              type="primary"
              @click="deleteShop(row)"
            >
              删除
            </ElLink>
          </div>
        </div>
      </template>
    </Grid>

    <ShopAddModal v-model:open="addModalOpen" @success="reload" />
    <ShopEditModal
      v-model:open="editModalOpen"
      :shop="curModel"
      @success="reload"
    />
    <ShopUserManageModal
      v-if="manageShop"
      v-model:open="shopUserModalOpen"
      :app-id="manageShop.app_id"
      :app-name="manageShop.app_name"
    />
    <ShopH5DomainManageModal
      v-if="manageShop"
      v-model:open="h5DomainModalOpen"
      :app-id="manageShop.app_id"
      :app-name="manageShop.app_name"
    />
    <ShopTrafficManageModal
      v-if="manageShop"
      v-model:open="trafficModalOpen"
      :app-id="manageShop.app_id"
      :app-name="manageShop.app_name"
    />
  </Page>
</template>

<style scoped>
.shop-op-col {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 4px 0;
  line-height: 1.5;
}

.shop-op-row {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 10px;
}

.shop-op-link {
  margin-left: 0 !important;
  white-space: nowrap;
}

:deep(.shop-action-cell .vxe-cell) {
  align-items: flex-start !important;
  white-space: normal !important;
}

:deep(.shop-action-cell .vxe-cell--wrapper) {
  overflow: visible !important;
  white-space: normal !important;
}

:deep(.el-link.is-underline:hover) {
  opacity: 0.8;
}
:deep(.el-link.is-underline:hover::after) {
  display: none;
}
</style>
