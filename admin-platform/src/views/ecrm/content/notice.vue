<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { confirm, Page, useVbenDrawer } from '@vben/common-ui';
import { VbenTiptap, VbenTiptapPreview } from '@vben/plugins/tiptap';

import { Plus } from '@element-plus/icons-vue';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { uploadAttachmentApi } from '#/api/core/attachment';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  fetchMerchantCategories,
  fetchMerchantTypes,
  type MerchantCategoryRow,
  type MerchantTypeRow,
} from '#/api/core/ecrm';
import {
  createPlatformNoticeApi,
  deletePlatformNoticeApi,
  getPlatformNoticeApi,
  listPlatformNoticesApi,
  type NoticeScopeType,
  type PlatformNotice,
  updatePlatformNoticeApi,
  updatePlatformNoticeStatusApi,
} from '#/api/core/platform-content';
import StorePickerModal, {
  type PickedStore,
} from '#/components/ecrm/store-picker-modal.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type NoticeForm = {
  content: string;
  is_show: number;
  scope_ids: number[];
  scope_items: PickedStore[];
  scope_type: NoticeScopeType;
  title: string;
};

const canManage = ref(false);
const editing = ref<PlatformNotice>();
const detail = ref<PlatformNotice>();
const storePickerOpen = ref(false);
const merchantTypes = ref<MerchantTypeRow[]>([]);
const merchantCategories = ref<MerchantCategoryRow[]>([]);
const form = reactive<NoticeForm>({
  content: '',
  is_show: 1,
  scope_ids: [],
  scope_items: [],
  scope_type: 'all',
  title: '',
});

const imageUpload: ImageUploadOptions = {
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => (await uploadAttachmentApi(file)).attachment_src,
  onUploadError: () => ElMessage.error('图片上传失败'),
};

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '时间选择',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'is_show',
    label: '启用状态',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入消息名称搜索' },
    fieldName: 'keyword',
    label: '消息名称',
  },
]);

const gridOptions: VxeGridProps<PlatformNotice> = {
  columns: [
    { field: 'title', minWidth: 240, title: '公告名称' },
    {
      field: 'scope_type',
      minWidth: 240,
      slots: { default: 'scope' },
      title: '店铺范围',
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '启用状态',
      width: 120,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '发送日期',
    },
    platformListActionColumn({ width: 210 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, values) => {
        const range = Array.isArray(values?.date_range)
          ? values.date_range
          : [];
        const isShow = Number(values?.is_show);
        const result = await listPlatformNoticesApi({
          date_from: range[0] ? `${range[0]} 00:00:00` : undefined,
          date_to: range[1] ? `${range[1]} 23:59:59` : undefined,
          is_show: isShow === 0 || isShow === 1 ? isShow : undefined,
          keyword: String(values?.keyword ?? '').trim() || undefined,
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'notice_id' },
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
  class: 'w-[min(96vw,1180px)]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: save,
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[min(96vw,1040px)]',
  footer: false,
  placement: 'right',
});

const scopeSelectionLabel = computed(() => {
  if (form.scope_type === 'store_name') {
    return form.scope_items.map((item) => item.mer_name).filter(Boolean);
  }
  if (form.scope_type === 'store_type') {
    return merchantTypes.value
      .filter((item) => form.scope_ids.includes(item.id))
      .map((item) => item.name);
  }
  if (form.scope_type === 'store_category') {
    return merchantCategories.value
      .filter((item) => form.scope_ids.includes(item.merchant_category_id))
      .map((item) => item.category_name);
  }
  return [];
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    content: '',
    is_show: 1,
    scope_ids: [],
    scope_items: [],
    scope_type: 'all',
    title: '',
  });
}

function applyNoticeToForm(row: PlatformNotice) {
  Object.assign(form, {
    content: row.content || '',
    is_show: row.is_show,
    scope_ids: [...(row.scope_ids || [])],
    scope_items: (row.scope_items || []).map((item) => ({
      mer_id: item.id,
      mer_name: item.name,
      mer_phone: '',
      real_name: '',
    })),
    scope_type: row.scope_type || 'all',
    title: row.title || '',
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增公告' }).open();
}

async function openEdit(row: PlatformNotice) {
  const latest = await getPlatformNoticeApi(row.notice_id);
  editing.value = latest;
  applyNoticeToForm(latest);
  formDrawerApi.setState({ title: '编辑公告' }).open();
}

async function openDetail(row: PlatformNotice) {
  detail.value = await getPlatformNoticeApi(row.notice_id);
  detailDrawerApi.setState({ title: '公告详情' }).open();
}

function changeScopeType() {
  form.scope_ids = [];
  form.scope_items = [];
}

function selectStores(stores: PickedStore[]) {
  form.scope_items = stores;
  form.scope_ids = stores.map((item) => item.mer_id);
}

function scopeLabel(row: PlatformNotice) {
  if (row.scope_type === 'all') return '全部';
  const names = (row.scope_items || [])
    .map((item) => item.name)
    .filter(Boolean);
  const prefix =
    row.scope_type === 'store_name'
      ? '指定店铺'
      : row.scope_type === 'store_type'
        ? '店铺类别'
        : '店铺分类';
  if (names.length === 0) return prefix;
  if (names.length <= 2) return `${prefix}：${names.join('、')}`;
  return `${prefix}：${names.slice(0, 2).join('、')} 等 ${names.length} 项`;
}

function validScope() {
  return form.scope_type === 'all' || form.scope_ids.length > 0;
}

async function save() {
  const content = form.content.trim();
  if (!form.title.trim() || !content || content === '<p></p>') {
    ElMessage.warning('请填写消息名称和公告内容');
    return;
  }
  if (!validScope()) {
    ElMessage.warning('请完成店铺范围关联选择');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      content,
      is_show: form.is_show,
      scope_ids: form.scope_ids,
      scope_type: form.scope_type,
      title: form.title.trim(),
    };
    if (editing.value) {
      await updatePlatformNoticeApi(editing.value.notice_id, payload);
    } else {
      await createPlatformNoticeApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('公告已保存');
    await gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function toggleStatus(row: PlatformNotice, value: boolean | number | string) {
  try {
    await updatePlatformNoticeStatusApi(row.notice_id, Number(value));
    row.is_show = Number(value);
    ElMessage.success('状态已保存');
  } catch {
    await gridApi.reload();
  }
}

async function remove(row: PlatformNotice) {
  try {
    await confirm({
      content: `删除公告“${row.title}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '删除公告',
    });
    await deletePlatformNoticeApi(row.notice_id);
    ElMessage.success('公告已删除');
    await gridApi.reload();
  } catch {
    // 用户取消或统一请求层处理。
  }
}

onMounted(async () => {
  const [codes, categoriesResult, typesResult] = await Promise.all([
    getAccessCodesApi(),
    fetchMerchantCategories().catch(() => ({ list: [] as MerchantCategoryRow[] })),
    fetchMerchantTypes().catch(() => ({ list: [] as MerchantTypeRow[] })),
  ]);
  canManage.value = codes.includes('content.notice.manage');
  merchantCategories.value = categoriesResult.list || [];
  merchantTypes.value = typesResult.list || [];
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          发布公告
        </ElButton>
      </template>
      <template #scope="{ row }">
        {{ scopeLabel(row) }}
      </template>
      <template #is_show="{ row }">
        <ElSwitch
          :disabled="!canManage"
          :model-value="row.is_show === 1"
          @change="(value) => toggleStatus(row, value ? 1 : 0)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="108px" class="notice-form">
        <ElFormItem label="消息名称" required>
          <ElInput
            v-model="form.title"
            maxlength="20"
            placeholder="请输入消息名称"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="选择店铺">
          <ElRadioGroup
            v-model="form.scope_type"
            @change="changeScopeType"
          >
            <ElRadio value="all">全部</ElRadio>
            <ElRadio value="store_name">店铺名称</ElRadio>
            <ElRadio value="store_type">店铺类别</ElRadio>
            <ElRadio value="store_category">店铺分类</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem
          v-if="form.scope_type === 'store_name'"
          label="关联店铺"
          required
        >
          <div class="notice-form__association">
            <ElButton type="primary" plain @click="storePickerOpen = true">
              选择店铺
            </ElButton>
            <div v-if="scopeSelectionLabel.length" class="notice-form__tags">
              <ElTag v-for="name in scopeSelectionLabel" :key="name">
                {{ name }}
              </ElTag>
            </div>
          </div>
        </ElFormItem>
        <ElFormItem
          v-else-if="form.scope_type === 'store_type'"
          label="关联类别"
          required
        >
          <ElSelect
            v-model="form.scope_ids"
            class="w-full"
            collapse-tags
            collapse-tags-tooltip
            multiple
            placeholder="请选择店铺类别"
          >
            <ElOption
              v-for="item in merchantTypes"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem
          v-else-if="form.scope_type === 'store_category'"
          label="关联分类"
          required
        >
          <ElSelect
            v-model="form.scope_ids"
            class="w-full"
            collapse-tags
            collapse-tags-tooltip
            multiple
            placeholder="请选择店铺分类"
          >
            <ElOption
              v-for="item in merchantCategories"
              :key="item.merchant_category_id"
              :label="item.category_name"
              :value="item.merchant_category_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="公告内容" required>
          <VbenTiptap
            v-model="form.content"
            class="w-full"
            :image-upload="imageUpload"
            :max-height="520"
            :min-height="360"
            :previewable="false"
            placeholder="请输入公告内容…"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <DetailDrawer>
      <div v-if="detail" class="notice-detail">
        <div class="notice-detail__item">
          <span>消息名称</span><strong>{{ detail.title }}</strong>
        </div>
        <div class="notice-detail__item">
          <span>店铺范围</span><strong>{{ scopeLabel(detail) }}</strong>
        </div>
        <div class="notice-detail__item">
          <span>发送日期</span><strong>{{ formatShanghaiDateTime(detail.create_time) }}</strong>
        </div>
        <VbenTiptapPreview :content="detail.content" :min-height="280" />
      </div>
    </DetailDrawer>

    <StorePickerModal
      v-model:open="storePickerOpen"
      :selected="form.scope_items"
      @confirm="selectStores"
    />
  </Page>
</template>

<style scoped>
.notice-form {
  margin: 8px auto 0;
  max-width: 980px;
}

.notice-form__association {
  display: flex;
  width: 100%;
  align-items: flex-start;
  gap: 12px;
}

.notice-form__tags {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: 8px;
  min-height: 32px;
  padding-top: 2px;
}

.notice-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.notice-detail__item {
  display: grid;
  grid-template-columns: 84px 1fr;
  gap: 12px;
  color: hsl(var(--muted-foreground));
}

.notice-detail__item strong {
  color: hsl(var(--foreground));
  font-weight: 500;
}
</style>
