<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref, shallowRef } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
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
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createMerchantProductApi,
  deleteMerchantProductApi,
  getMerchantProductCategoriesApi,
  listMerchantProductsApi,
  listMerchantProductRecycleBinApi,
  restoreMerchantProductApi,
  setMerchantProductShowApi,
  setMerchantProductStockApi,
  updateMerchantProductApi,
  type MerchantCategoryNode,
  type MerchantProduct,
  type MerchantProductSaveInput,
  type MerchantRecycleProduct,
} from '#/api/core/merchant-catalog';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const canCreate = ref(false);
const canUpdate = ref(false);
const canDelete = ref(false);
const canShow = ref(false);
const canStock = ref(false);
const canRestore = ref(false);
const editingID = ref<number>();
const stockRow = ref<MerchantProduct>();
const stockValue = ref(0);
const recycleRows = ref<MerchantRecycleProduct[]>([]);
const recycleTotal = ref(0);
const recycleQuery = reactive({ page: 1, limit: 20 });
const categoryOptions = shallowRef<Array<{ label: string; value: number }>>([]);

const form = reactive<MerchantProductSaveInput>({
  cate_id: 0,
  cost: 0,
  image: '',
  is_show: 1,
  keyword: '',
  ot_price: 0,
  price: 0,
  slider_image: '',
  spec_type: 0,
  stock: 0,
  store_info: '',
  store_name: '',
  type: 0,
  unit_name: '件',
});

const statusText = (status: number) =>
  ({ '-2': '已下架', '-1': '审核驳回', 0: '待审核', 1: '已通过' })[status] ||
  '未知';

const statusType = (status: number) =>
  ({ '-2': 'info', '-1': 'danger', 0: 'warning', 1: 'success' })[status] ||
  'info';

function flattenCategories(
  nodes: MerchantCategoryNode[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.store_category_id },
    ...flattenCategories(node.children || [], `${prefix}— `),
  ]);
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    cate_id: 0,
    cost: 0,
    image: '',
    is_show: 1,
    keyword: '',
    ot_price: 0,
    price: 0,
    slider_image: '',
    spec_type: 0,
    stock: 0,
    store_info: '',
    store_name: '',
    type: 0,
    unit_name: '件',
  });
}

const formOptions = computed((): VbenFormProps =>
  listFormOptionsDefaults([
    LIST_DATE_RANGE_FIELD,
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '商品名称 / 关键词',
      },
      fieldName: 'keyword',
      label: '商品搜索',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '待审核', value: 0 },
          { label: '已通过', value: 1 },
          { label: '已驳回', value: -1 },
          { label: '已下架', value: -2 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'status',
      label: '审核状态',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: categoryOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'cate_id',
      label: '商品分类',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '上架中', value: 1 },
          { label: '已下架', value: 0 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'is_show',
      label: '上架状态',
    },
  ]),
);

const gridOptions: VxeGridProps<MerchantProduct> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'product_id', title: 'ID', width: 80 },
    {
      field: 'store_name',
      minWidth: 190,
      showOverflow: true,
      title: '商品名称',
    },
    { field: 'cate_name', minWidth: 120, title: '分类' },
    {
      field: 'price',
      title: '售价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'stock', title: '库存', width: 90 },
    {
      field: 'status',
      slots: { default: 'audit' },
      title: '审核',
      width: 96,
    },
    {
      field: 'is_show',
      slots: { default: 'show' },
      title: '上架',
      width: 88,
    },
    {
      field: 'refusal',
      minWidth: 160,
      showOverflow: true,
      title: '驳回原因',
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 248 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const isShow = formValues?.is_show;
        const cateId = formValues?.cate_id;
        const data = await listMerchantProductsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            typeof status === 'number' && [-2, -1, 0, 1].includes(Number(status))
              ? Number(status)
              : undefined,
          cate_id:
            typeof cateId === 'number' && Number(cateId) > 0
              ? Number(cateId)
              : undefined,
          is_show: isShow === 0 || isShow === 1 ? Number(isShow) : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});

const [ProductModal, productModalApi] = useVbenModal({
  onConfirm: async () => {
    if (!form.store_name.trim()) {
      ElMessage.warning('请填写商品名称');
      return;
    }
    if (!form.cate_id) {
      ElMessage.warning('请选择平台分类');
      return;
    }
    if (form.price < 0 || form.stock < 0) {
      ElMessage.warning('价格和库存不能小于 0');
      return;
    }
    saving.value = true;
    productModalApi.lock();
    try {
      if (editingID.value) {
        await updateMerchantProductApi(editingID.value, form);
      } else {
        await createMerchantProductApi(form);
      }
      ElMessage.success(
        editingID.value ? '商品已更新，可能需重新审核' : '商品已创建',
      );
      productModalApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      productModalApi.unlock();
    }
  },
});

const [StockModal, stockModalApi] = useVbenModal({
  title: '调整库存',
  confirmText: '保存',
  onConfirm: async () => {
    if (!stockRow.value || stockValue.value < 0) return;
    stockModalApi.lock();
    try {
      await setMerchantProductStockApi(stockRow.value.product_id, stockValue.value);
      stockRow.value.stock = stockValue.value;
      ElMessage.success('库存已更新');
      stockModalApi.close();
      gridApi.reload();
    } finally {
      stockModalApi.unlock();
    }
  },
});

const [RecycleModal, recycleModalApi] = useVbenModal({
  title: '商品回收站',
  showConfirmButton: false,
  cancelText: '关闭',
});

function openCreate() {
  resetForm();
  productModalApi.setState({ title: '新增商品' });
  productModalApi.open();
}

function openEdit(row: MerchantProduct) {
  editingID.value = row.product_id;
  Object.assign(form, {
    cate_id: row.cate_id,
    cost: -1,
    image: row.image,
    is_show: row.is_show,
    keyword: row.keyword,
    ot_price: row.ot_price,
    price: row.price,
    slider_image: row.image,
    spec_type: 0,
    stock: row.stock,
    store_info: row.store_info,
    store_name: row.store_name,
    type: 0,
    unit_name: row.unit_name || '件',
  });
  productModalApi.setState({ title: '编辑商品' });
  productModalApi.open();
}

async function loadCategories() {
  const result = await getMerchantProductCategoriesApi();
  categoryOptions.value = flattenCategories(result.list || []);
}

async function loadRecycleBin() {
  const result = await listMerchantProductRecycleBinApi(recycleQuery);
  recycleRows.value = result.list || [];
  recycleTotal.value = result.total || 0;
}

async function openRecycleBin() {
  recycleQuery.page = 1;
  await loadRecycleBin();
  recycleModalApi.open();
}

async function restore(row: MerchantRecycleProduct) {
  await restoreMerchantProductApi(row.product_id);
  ElMessage.success('商品已恢复，等待平台重新审核');
  await Promise.all([loadRecycleBin(), gridApi.reload()]);
}

async function changeShow(row: MerchantProduct) {
  const next = row.is_show === 1 ? 0 : 1;
  try {
    await setMerchantProductShowApi(row.product_id, next);
    row.is_show = next;
    ElMessage.success(next ? '商品已上架' : '商品已下架');
    gridApi.reload();
  } catch {
    ElMessage.error('商品状态更新失败');
  }
}

function openStock(row: MerchantProduct) {
  stockRow.value = row;
  stockValue.value = row.stock;
  stockModalApi.open();
}

async function remove(row: MerchantProduct) {
  try {
    await ElMessageBox.confirm(
      `将商品「${row.store_name}」移入回收站，30 天内可恢复，是否继续？`,
      '移入回收站',
      { type: 'warning' },
    );
    await deleteMerchantProductApi(row.product_id);
    ElMessage.success('商品已移入回收站');
    gridApi.reload();
  } catch {
    // 用户取消
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([
    getAccessCodesApi(),
    loadCategories(),
  ]);
  canCreate.value = permissions.includes('product.create');
  canUpdate.value = permissions.includes('product.update');
  canDelete.value = permissions.includes('product.delete');
  canShow.value = permissions.includes('product.show');
  canStock.value = permissions.includes('product.stock');
  canRestore.value = permissions.includes('product.restore');
});
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton v-if="canRestore" @click="openRecycleBin">回收站</ElButton>
      <ElButton v-if="canCreate" type="primary" @click="openCreate">
        新增商品
      </ElButton>
    </template>

    <Grid>
      <template #audit="{ row }">
        <ElTag :type="statusType(row.status)">
          {{ statusText(row.status) }}
        </ElTag>
      </template>
      <template #show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '上架中' : '已下架' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton v-if="canUpdate" link type="primary" @click="openEdit(row)">
          编辑
        </ElButton>
        <ElButton v-if="canStock" link type="primary" @click="openStock(row)">
          库存
        </ElButton>
        <ElButton
          v-if="canShow"
          link
          type="warning"
          @click="changeShow(row)"
        >
          {{ row.is_show === 1 ? '下架' : '上架' }}
        </ElButton>
        <ElButton v-if="canDelete" link type="danger" @click="remove(row)">
          删除
        </ElButton>
      </template>
    </Grid>

    <ProductModal class="w-[760px] max-w-[96vw]">
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="88px">
        <ElFormItem label="商品名称" required>
          <ElInput v-model="form.store_name" />
        </ElFormItem>
        <ElFormItem label="平台分类" required>
          <ElSelect v-model="form.cate_id" filterable class="w-full">
            <ElOption
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="销售价">
          <ElInputNumber
            v-model="form.price"
            :min="0"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="划线价">
          <ElInputNumber
            v-model="form.ot_price"
            :min="0"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="库存">
          <ElInputNumber v-model="form.stock" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="单位">
          <ElInput v-model="form.unit_name" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="关键词">
          <ElInput v-model="form.keyword" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="主图">
          <ElInput
            v-model="form.image"
            placeholder="图片 URL（素材库接入后可直接选择）"
          />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="商品简介">
          <ElInput v-model="form.store_info" :rows="4" type="textarea" />
        </ElFormItem>
        <ElFormItem label="初始上架">
          <ElSwitch
            v-model="form.is_show"
            :active-value="1"
            :inactive-value="0"
          />
        </ElFormItem>
      </ElForm>
    </ProductModal>

    <StockModal class="w-[380px] max-w-[96vw]">
      <ElForm label-width="72px">
        <ElFormItem label="商品">
          <span>{{ stockRow?.store_name }}</span>
        </ElFormItem>
        <ElFormItem label="库存">
          <ElInputNumber v-model="stockValue" :min="0" class="w-full" />
        </ElFormItem>
      </ElForm>
    </StockModal>

    <RecycleModal class="w-[760px] max-w-[96vw]">
      <ElAlert
        class="mb-4"
        :closable="false"
        show-icon
        title="恢复后的商品会重新进入平台审核，超过恢复期限后不可恢复。"
        type="info"
      />
      <ElTable :data="recycleRows" row-key="product_id">
        <ElTableColumn label="商品" min-width="240" prop="store_name" />
        <ElTableColumn label="移入时间" min-width="170" prop="deleted_at">
          <template #default="{ row }">
            {{ formatShanghaiDateTime(row.deleted_at) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="可恢复至" min-width="170" prop="restore_until">
          <template #default="{ row }">
            {{ formatShanghaiDateTime(row.restore_until) }}
          </template>
        </ElTableColumn>
        <ElTableColumn fixed="right" label="操作" width="100">
          <template #default="{ row }">
            <ElButton link type="primary" @click="restore(row)">恢复</ElButton>
          </template>
        </ElTableColumn>
      </ElTable>
      <div class="mt-4 flex justify-end">
        <el-pagination
          :current-page="recycleQuery.page"
          :page-size="recycleQuery.limit"
          :total="recycleTotal"
          background
          layout="total, prev, pager, next"
          @current-change="(page: number) => { recycleQuery.page = page; loadRecycleBin(); }"
        />
      </div>
    </RecycleModal>
  </Page>
</template>
