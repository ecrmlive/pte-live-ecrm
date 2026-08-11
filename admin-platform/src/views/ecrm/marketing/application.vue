<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type {
  SignupActivity,
  SignupFormField,
  SignupFormOption,
  SignupRecord,
} from '#/api/core/platform-application';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElColorPicker,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElPagination,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTabs,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createSignupActivityApi,
  deleteSignupActivityApi,
  exportSignupRecordsApi,
  getSignupActivityApi,
  listSignupActivitiesApi,
  listSignupFormOptionsApi,
  listSignupRecordsApi,
  setSignupActivityStatusApi,
  updateSignupActivityApi,
} from '#/api/core/platform-application';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit' | 'view';
type ExampleKind = 'cover' | 'poster' | 'color';

const EXAMPLE_URLS: Record<ExampleKind, string> = {
  cover: 'https://picsum.photos/seed/qixi-signup-cover-example/750/350',
  poster: 'https://picsum.photos/seed/qixi-signup-poster-example/750/1250',
  color: 'https://picsum.photos/seed/qixi-signup-color-example/750/1334',
};

const EXAMPLE_TITLES: Record<ExampleKind, string> = {
  cover: '封面图示例',
  poster: '分享海报示例',
  color: '活动背景色示例',
};

const canRead = ref(false);
const canManage = ref(false);
const drawerMode = ref<DrawerMode>('create');
const formTab = ref('basic');
const viewTab = ref('info');
const editingId = ref(0);
const viewing = ref<SignupActivity>();
const formOptions = ref<SignupFormOption[]>([]);
const exampleKind = ref<ExampleKind>('cover');
const exporting = ref(false);
const statsKeyword = ref('');
const statsLoading = ref(false);
const statsRows = ref<SignupRecord[]>([]);
const statsFields = ref<SignupFormField[]>([]);
const statsTotal = ref(0);
const statsPage = ref(1);
const statsLimit = ref(10);

const form = reactive({
  name: '',
  info: '',
  cover_url: '',
  poster_url: '',
  color: '#E8F5E9',
  form_id: undefined as number | undefined,
  quota: 0,
  status: 1,
  sort: 0,
  date_range: [] as string[],
});

const formNameMap = computed(() => {
  const map = new Map<number, string>();
  for (const row of formOptions.value) {
    map.set(Number(row.id), row.name);
  }
  return map;
});

const previewFields = computed(() => {
  if (drawerMode.value === 'view') {
    return viewing.value?.form_fields || [];
  }
  const id = Number(form.form_id || 0);
  return formOptions.value.find((x) => x.id === id)?.fields || [];
});

const formOptionsSchema = computed<VbenFormProps['schema']>(() => [
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '活动名称 / 关键字',
    },
    fieldName: 'keyword',
    label: '活动搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '未开始', value: 0 },
        { label: '进行中', value: 1 },
        { label: '已结束', value: -1 },
      ],
      placeholder: '全部',
    },
    fieldName: 'activity_status',
    label: '报名状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: formOptions.value.map((x) => ({
        label: x.name,
        value: x.id,
      })),
      placeholder: '全部',
    },
    fieldName: 'form_id',
    label: '关联表单',
  },
]);

const gridOptions: VxeGridProps<SignupActivity> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '活动名称',
    },
    {
      field: 'cover_url',
      slots: { default: 'cover' },
      title: '封面图',
      width: 110,
    },
    {
      field: 'poster_url',
      slots: { default: 'poster' },
      title: '分享海报图',
      width: 110,
    },
    {
      field: 'activity_status_text',
      slots: { default: 'activityStatus' },
      title: '报名状态',
      width: 100,
    },
    {
      field: 'signup_count_text',
      minWidth: 120,
      title: '已报人数/总人数',
    },
    {
      field: 'form_name',
      formatter: ({ row }) => row.form_name || formNameMap.value.get(row.form_id) || '—',
      minWidth: 130,
      showOverflow: 'tooltip',
      title: '关联表单',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 170,
      showOverflow: false,
      title: '创建时间',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 110,
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const values = (formValues || {}) as Record<string, unknown>;
        const statusRaw = values.activity_status;
        const formIdRaw = values.form_id;
        const data = await listSignupActivitiesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(values.keyword ?? '').trim() || undefined,
          activity_status:
            statusRaw === 0 || statusRaw === 1 || statusRaw === -1 || statusRaw === '0' || statusRaw === '1' || statusRaw === '-1'
              ? statusRaw
              : undefined,
          form_id: formIdRaw ? Number(formIdRaw) : undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
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

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: listFormOptionsDefaults(formOptionsSchema.value),
  gridOptions,
});

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const [ViewDrawer, viewDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  cancelText: '关闭',
});

const [ExampleModal, exampleModalApi] = useVbenModal({
  title: '查看示例',
  class: 'w-[820px] max-w-[96vw]',
  footer: false,
});

function mediaUrl(url?: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function statusTagType(status: number) {
  if (status === 1) return 'success';
  if (status === 0) return 'warning';
  if (status === -1) return 'info';
  return 'danger';
}

function resetForm() {
  editingId.value = 0;
  Object.assign(form, {
    name: '',
    info: '',
    cover_url: '',
    poster_url: '',
    color: '#E8F5E9',
    form_id: undefined,
    quota: 0,
    status: 1,
    sort: 0,
    date_range: [],
  });
  formTab.value = 'basic';
}

function fillForm(row: SignupActivity) {
  editingId.value = row.id;
  Object.assign(form, {
    name: row.name || '',
    info: row.info || '',
    cover_url: row.cover_url || '',
    poster_url: row.poster_url || '',
    color: row.color || '#E8F5E9',
    form_id: row.form_id || undefined,
    quota: Number(row.quota) || 0,
    status: row.status === 1 ? 1 : 0,
    sort: Number(row.sort) || 0,
    date_range:
      row.starts_at && row.ends_at
        ? [row.starts_at.slice(0, 10), row.ends_at.slice(0, 10)]
        : [],
  });
}

function openCreate() {
  drawerMode.value = 'create';
  resetForm();
  formDrawerApi.setState({ title: '创建报名活动', showConfirmButton: true }).open();
}

function openEdit(row: SignupActivity) {
  drawerMode.value = 'edit';
  fillForm(row);
  formTab.value = 'basic';
  formDrawerApi.setState({ title: '编辑报名活动', showConfirmButton: true }).open();
}

async function openView(row: SignupActivity) {
  drawerMode.value = 'view';
  viewTab.value = 'info';
  statsKeyword.value = '';
  statsPage.value = 1;
  try {
    viewing.value = await getSignupActivityApi(row.id);
  } catch {
    viewing.value = row;
  }
  viewDrawerApi.setState({ title: '查看报名活动' }).open();
  void loadStats();
}

function openEditFromView() {
  if (!viewing.value) return;
  viewDrawerApi.close();
  openEdit(viewing.value);
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请输入活动名称');
    return;
  }
  if (!form.cover_url.trim()) {
    ElMessage.warning('请上传封面图');
    return;
  }
  if (!form.poster_url.trim()) {
    ElMessage.warning('请上传活动分享海报');
    return;
  }
  if (!form.date_range?.[0] || !form.date_range?.[1]) {
    ElMessage.warning('请选择活动起止日期');
    return;
  }
  if (!form.form_id) {
    ElMessage.warning('请关联系统表单');
    formTab.value = 'form';
    return;
  }

  const body = {
    name,
    info: form.info.trim(),
    cover_url: form.cover_url.trim(),
    poster_url: form.poster_url.trim(),
    color: form.color.trim(),
    form_id: Number(form.form_id),
    quota: Number(form.quota) || 0,
    status: form.status === 1 ? 1 : 0,
    sort: Number(form.sort) || 0,
    starts_at: `${form.date_range[0]} 00:00:00`,
    ends_at: `${form.date_range[1]} 23:59:59`,
  };

  formDrawerApi.lock();
  try {
    if (editingId.value) {
      await updateSignupActivityApi(editingId.value, body);
    } else {
      await createSignupActivityApi(body);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: SignupActivity, enabled: boolean) {
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await setSignupActivityStatusApi(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function removeRow(row: SignupActivity) {
  try {
    await confirm({
      content: `确定删除报名活动「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteSignupActivityApi(row.id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

function openExample(kind: ExampleKind) {
  exampleKind.value = kind;
  exampleModalApi.setState({ title: EXAMPLE_TITLES[kind] }).open();
}

async function loadStats() {
  if (!viewing.value?.id || !canRead.value) return;
  statsLoading.value = true;
  try {
    const data = await listSignupRecordsApi(viewing.value.id, {
      page: statsPage.value,
      limit: statsLimit.value,
      keyword: statsKeyword.value.trim() || undefined,
    });
    statsRows.value = data.list || [];
    statsFields.value = data.fields || viewing.value.form_fields || [];
    statsTotal.value = data.total || 0;
  } finally {
    statsLoading.value = false;
  }
}

function searchStats() {
  statsPage.value = 1;
  void loadStats();
}

async function exportStats() {
  if (!viewing.value?.id || exporting.value) return;
  exporting.value = true;
  try {
    const result = await exportSignupRecordsApi(viewing.value.id, {
      keyword: statsKeyword.value.trim() || undefined,
    });
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8;',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || '报名用户.csv';
    link.click();
    URL.revokeObjectURL(url);
  } finally {
    exporting.value = false;
  }
}

function formCell(row: SignupRecord, key: string) {
  return row.form_cols?.[key] || String(row.form_value?.[key] ?? '') || '—';
}

onMounted(async () => {
  const [profile, codes, forms] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    listSignupFormOptionsApi().catch(() => ({ list: [] as SignupFormOption[] })),
  ]);
  formOptions.value = forms.list || [];
  gridApi.setState({
    formOptions: listFormOptionsDefaults(formOptionsSchema.value),
  });
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('marketing.application.read') ||
      codes.includes('marketing.application.manage'));
  canManage.value = roleOK && codes.includes('marketing.application.manage');
  if (canRead.value) gridApi.reload();
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
          创建报名活动
        </ElButton>
      </template>

      <template #cover="{ row }">
        <ElImage
          v-if="mediaUrl(row.cover_url)"
          :src="mediaUrl(row.cover_url)"
          fit="cover"
          class="app-thumb"
          :preview-src-list="[mediaUrl(row.cover_url)]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>

      <template #poster="{ row }">
        <ElImage
          v-if="mediaUrl(row.poster_url)"
          :src="mediaUrl(row.poster_url)"
          fit="cover"
          class="app-thumb app-thumb--poster"
          :preview-src-list="[mediaUrl(row.poster_url)]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>

      <template #activityStatus="{ row }">
        <ElTag :type="statusTagType(row.activity_status)" size="small">
          {{ row.activity_status_text || '—' }}
        </ElTag>
      </template>

      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openView(row)">查看</ElButton>
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
          type="primary"
          @click="removeRow(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElTabs v-model="formTab" class="app-tabs">
        <ElTabPane label="基础设置" name="basic">
          <ElForm label-width="120px" class="app-form">
            <ElFormItem label="活动名称" required>
              <ElInput
                v-model="form.name"
                maxlength="128"
                show-word-limit
                placeholder="请输入活动名称"
              />
            </ElFormItem>
            <ElFormItem label="活动简介">
              <ElInput
                v-model="form.info"
                type="textarea"
                :rows="4"
                maxlength="500"
                show-word-limit
                placeholder="请输入活动简介"
              />
            </ElFormItem>
            <ElFormItem label="封面图" required>
              <div class="app-pic">
                <ImageField v-model="form.cover_url" :preview-size="88" />
                <div class="app-pic__side">
                  <ElButton link type="primary" @click="openExample('cover')">
                    查看示例
                  </ElButton>
                  <div class="field-hint">建议尺寸：750 × 350</div>
                </div>
              </div>
            </ElFormItem>
            <ElFormItem label="活动分享海报" required>
              <div class="app-pic">
                <ImageField v-model="form.poster_url" :preview-size="88" />
                <div class="app-pic__side">
                  <ElButton link type="primary" @click="openExample('poster')">
                    查看示例
                  </ElButton>
                  <div class="field-hint">建议尺寸：750 × 1250</div>
                </div>
              </div>
            </ElFormItem>
            <ElFormItem label="活动背景色">
              <div class="app-color">
                <ElColorPicker v-model="form.color" />
                <ElInput
                  v-model="form.color"
                  class="app-color__input"
                  placeholder="#E8F5E9"
                />
                <ElButton link type="primary" @click="openExample('color')">
                  查看示例
                </ElButton>
              </div>
            </ElFormItem>
            <ElFormItem label="活动起止日期" required>
              <ElDatePicker
                v-model="form.date_range"
                type="daterange"
                value-format="YYYY-MM-DD"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                class="w-full"
              />
            </ElFormItem>
            <ElFormItem label="活动人数上限">
              <ElInputNumber
                v-model="form.quota"
                :min="0"
                :step="1"
                controls-position="right"
              />
              <div class="field-hint ml-3">0 = 不限制</div>
            </ElFormItem>
            <ElFormItem label="排序">
              <ElInputNumber
                v-model="form.sort"
                :min="0"
                :step="1"
                controls-position="right"
              />
            </ElFormItem>
            <ElFormItem label="是否开启">
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
        </ElTabPane>

        <ElTabPane label="关联表单" name="form">
          <ElForm label-width="120px" class="app-form">
            <ElFormItem label="关联系统表单" required>
              <ElSelect
                v-model="form.form_id"
                class="w-full"
                clearable
                filterable
                placeholder="请选择系统表单"
              >
                <ElOption
                  v-for="item in formOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem v-if="previewFields.length" label="表单字段">
              <div class="phone-preview">
                <div class="phone-preview__title">手机预览</div>
                <div
                  v-for="field in previewFields"
                  :key="field.key"
                  class="phone-preview__field"
                >
                  <div class="phone-preview__label">{{ field.label }}</div>
                  <div
                    v-if="field.type === 'image'"
                    class="phone-preview__image"
                  >
                    上传图片
                  </div>
                  <div v-else class="phone-preview__input">请输入</div>
                </div>
              </div>
            </ElFormItem>
          </ElForm>
        </ElTabPane>
      </ElTabs>
    </FormDrawer>

    <ViewDrawer>
      <div v-if="viewing" class="view-head">
        <div>
          <div class="view-head__title">{{ viewing.name }}</div>
          <div class="view-head__meta">活动ID：{{ viewing.id }}</div>
        </div>
        <ElButton v-if="canManage" type="primary" @click="openEditFromView">
          编辑
        </ElButton>
      </div>

      <div v-if="viewing" class="view-summary">
        <div class="view-summary__item">
          <span class="label">报名状态</span>
          <ElTag :type="statusTagType(viewing.activity_status)" size="small">
            {{ viewing.activity_status_text }}
          </ElTag>
        </div>
        <div class="view-summary__item">
          <span class="label">已报名人数</span>
          <span>{{ viewing.total }}</span>
        </div>
        <div class="view-summary__item">
          <span class="label">人数上限</span>
          <span>{{ viewing.quota === 0 ? '不限制' : viewing.quota }}</span>
        </div>
        <div class="view-summary__item">
          <span class="label">创建时间</span>
          <span>{{ formatShanghaiDateTime(viewing.created_at) || '—' }}</span>
        </div>
      </div>

      <ElTabs v-if="viewing" v-model="viewTab" class="app-tabs">
        <ElTabPane label="活动信息" name="info">
          <div class="view-info">
            <div class="view-info__row">
              <span class="label">活动名称</span>
              <span>{{ viewing.name }}</span>
            </div>
            <div class="view-info__row">
              <span class="label">活动简介</span>
              <span>{{ viewing.info || '—' }}</span>
            </div>
            <div class="view-info__row">
              <span class="label">封面图</span>
              <ElImage
                v-if="mediaUrl(viewing.cover_url)"
                :src="mediaUrl(viewing.cover_url)"
                fit="cover"
                class="app-thumb app-thumb--lg"
                :preview-src-list="[mediaUrl(viewing.cover_url)]"
              />
              <span v-else>—</span>
            </div>
            <div class="view-info__row">
              <span class="label">分享海报</span>
              <ElImage
                v-if="mediaUrl(viewing.poster_url)"
                :src="mediaUrl(viewing.poster_url)"
                fit="cover"
                class="app-thumb app-thumb--poster app-thumb--lg"
                :preview-src-list="[mediaUrl(viewing.poster_url)]"
              />
              <span v-else>—</span>
            </div>
            <div class="view-info__row">
              <span class="label">背景色</span>
              <span class="topic-color-cell">
                <span
                  v-if="viewing.color"
                  class="topic-color-swatch"
                  :style="{ backgroundColor: viewing.color }"
                />
                {{ viewing.color || '默认' }}
              </span>
            </div>
            <div class="view-info__row">
              <span class="label">活动时间</span>
              <span>
                {{ formatShanghaiDateTime(viewing.starts_at) || '—' }}
                ~
                {{ formatShanghaiDateTime(viewing.ends_at) || '—' }}
              </span>
            </div>
            <div class="view-info__row">
              <span class="label">关联表单</span>
              <span>{{ viewing.form_name || '—' }}</span>
            </div>
            <div class="view-info__row view-info__row--block">
              <span class="label">手机预览</span>
              <div class="phone-preview">
                <div
                  v-for="field in previewFields"
                  :key="field.key"
                  class="phone-preview__field"
                >
                  <div class="phone-preview__label">{{ field.label }}</div>
                  <div
                    v-if="field.type === 'image'"
                    class="phone-preview__image"
                  >
                    上传图片
                  </div>
                  <div v-else class="phone-preview__input">请输入</div>
                </div>
                <div v-if="!previewFields.length" class="text-xs text-gray-400">
                  暂无表单字段
                </div>
              </div>
            </div>
          </div>
        </ElTabPane>

        <ElTabPane label="活动统计" name="stats">
          <div class="stats-toolbar">
            <ElInput
              v-model="statsKeyword"
              clearable
              class="stats-toolbar__input"
              placeholder="昵称 / ID / 手机号"
              @keyup.enter="searchStats"
            />
            <ElButton type="primary" @click="searchStats">查询</ElButton>
            <ElButton :loading="exporting" @click="exportStats">导出</ElButton>
          </div>

          <div v-loading="statsLoading" class="stats-table-wrap">
            <table class="stats-table">
              <thead>
                <tr>
                  <th>序号</th>
                  <th>用户名称/ID</th>
                  <th>手机号</th>
                  <th v-for="field in statsFields" :key="field.key">
                    {{ field.label }}
                  </th>
                  <th>提交时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in statsRows" :key="row.id">
                  <td>{{ row.index }}</td>
                  <td>
                    <div>{{ row.nickname || '—' }}</div>
                    <div class="text-xs text-gray-400">ID: {{ row.user_id }}</div>
                  </td>
                  <td>{{ row.mobile || '—' }}</td>
                  <td v-for="field in statsFields" :key="`${row.id}-${field.key}`">
                    <ElImage
                      v-if="
                        field.type === 'image' && mediaUrl(formCell(row, field.key))
                      "
                      :src="mediaUrl(formCell(row, field.key))"
                      fit="cover"
                      class="app-thumb app-thumb--sm"
                      :preview-src-list="[mediaUrl(formCell(row, field.key))]"
                    />
                    <span v-else>{{ formCell(row, field.key) }}</span>
                  </td>
                  <td>
                    {{ formatShanghaiDateTime(row.created_at) || '—' }}
                  </td>
                </tr>
                <tr v-if="!statsRows.length">
                  <td :colspan="4 + statsFields.length" class="stats-empty">
                    暂无报名记录
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="stats-pager">
            <ElPagination
              v-model:current-page="statsPage"
              v-model:page-size="statsLimit"
              :total="statsTotal"
              layout="total, prev, pager, next"
              @current-change="loadStats"
              @size-change="
                () => {
                  statsPage = 1;
                  loadStats();
                }
              "
            />
          </div>
        </ElTabPane>
      </ElTabs>
    </ViewDrawer>

    <ExampleModal>
      <div class="example-wrap">
        <ElImage
          :src="EXAMPLE_URLS[exampleKind]"
          fit="contain"
          class="example-img"
        />
      </div>
    </ExampleModal>
  </Page>
</template>

<style scoped>
.app-thumb {
  width: 56px;
  height: 36px;
  border-radius: 4px;
  overflow: hidden;
}

.app-thumb--poster {
  width: 36px;
  height: 56px;
}

.app-thumb--sm {
  width: 40px;
  height: 40px;
}

.app-thumb--lg {
  width: 120px;
  height: 70px;
}

.app-thumb--poster.app-thumb--lg {
  width: 70px;
  height: 120px;
}

.app-tabs {
  min-height: 420px;
}

.app-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.app-pic {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.app-pic__side {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 4px;
}

.field-hint {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.4;
}

.app-color {
  display: flex;
  gap: 10px;
  align-items: center;
  width: 100%;
}

.app-color__input {
  max-width: 160px;
}

.phone-preview {
  width: 280px;
  padding: 16px;
  border: 1px solid var(--el-border-color);
  border-radius: 16px;
  background: #fafafa;
}

.phone-preview__title {
  margin-bottom: 12px;
  font-weight: 600;
}

.phone-preview__field + .phone-preview__field {
  margin-top: 12px;
}

.phone-preview__label {
  margin-bottom: 6px;
  font-size: 13px;
}

.phone-preview__input,
.phone-preview__image {
  height: 36px;
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  background: #fff;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 36px;
  text-align: center;
}

.phone-preview__image {
  height: 72px;
  line-height: 72px;
}

.view-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.view-head__title {
  font-size: 18px;
  font-weight: 600;
}

.view-head__meta {
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.view-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
  padding: 14px 16px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: var(--el-fill-color-blank);
}

.view-summary__item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.view-summary__item .label,
.view-info__row .label {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.view-info__row {
  display: grid;
  grid-template-columns: 96px 1fr;
  gap: 12px;
  align-items: start;
  padding: 10px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.view-info__row--block {
  align-items: start;
}

.topic-color-cell {
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.topic-color-swatch {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 1px solid var(--el-border-color);
}

.stats-toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 12px;
}

.stats-toolbar__input {
  width: 260px;
}

.stats-table-wrap {
  overflow: auto;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.stats-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.stats-table th,
.stats-table td {
  padding: 10px 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  text-align: left;
  vertical-align: middle;
  white-space: nowrap;
}

.stats-table th {
  background: var(--el-fill-color-light);
  font-weight: 600;
}

.stats-empty {
  text-align: center !important;
  color: var(--el-text-color-secondary);
}

.stats-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

.example-wrap {
  display: flex;
  justify-content: center;
  padding: 8px;
}

.example-img {
  max-width: 100%;
  max-height: 70vh;
}
</style>
