<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { MarketingDecor } from '#/api/core/platform-marketing-decor';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElColorPicker,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  fetchProductLabels,
  type ProductLabelRow,
} from '#/api/core/ecrm';
import {
  createMarketingDecorApi,
  deleteMarketingDecorApi,
  listMarketingDecorApi,
  setMarketingDecorStatusApi,
  updateMarketingDecorApi,
} from '#/api/core/platform-marketing-decor';
import ImageField from '#/components/shop/image-field.vue';
import ImagesField from '#/components/shop/images-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

type LayoutType = 1 | 2 | 3;
type ExampleKind = 'list' | 'banner' | 'theme' | 'color';

const LAYOUT_LABEL: Record<LayoutType, string> = {
  1: '小图',
  2: '中图',
  3: '大图',
};

const EXAMPLE_URLS: Record<ExampleKind, string> = {
  list: 'https://picsum.photos/seed/qixi-topic-list-example/710/340',
  banner: 'https://picsum.photos/seed/qixi-topic-banner-example/750/750',
  theme: 'https://picsum.photos/seed/qixi-topic-theme-example/710/340',
  color: 'https://picsum.photos/seed/qixi-topic-color-example/710/340',
};

const EXAMPLE_TITLES: Record<ExampleKind, string> = {
  list: '活动列表图示例',
  banner: '活动轮播图示例',
  theme: '活动主题示例',
  color: '活动背景色示例',
};

const canRead = ref(false);
const canManage = ref(false);
const editing = ref<MarketingDecor>();
const labelOptions = ref<ProductLabelRow[]>([]);
const exampleKind = ref<ExampleKind>('list');

const form = reactive({
  label_id: undefined as number | undefined,
  name: '',
  cover_url: '',
  banners: [] as string[],
  theme_url: '',
  color: '',
  layout: 1 as LayoutType,
  status: 1,
  sort: 0,
});

const labelMap = computed(() => {
  const map = new Map<number, string>();
  for (const row of labelOptions.value) {
    map.set(Number(row.id), row.name);
  }
  return map;
});

const gridOptions: VxeGridProps<MarketingDecor> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'label_name',
      formatter: ({ row }) => labelNameOf(row),
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '关联标签',
    },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '活动名称',
    },
    {
      field: 'cover_url',
      slots: { default: 'listPic' },
      title: '活动列表图',
      width: 110,
    },
    {
      field: 'banners',
      slots: { default: 'banners' },
      title: '活动轮播图',
      width: 120,
    },
    {
      field: 'theme_url',
      slots: { default: 'theme' },
      title: '活动主题',
      width: 110,
    },
    {
      field: 'color',
      slots: { default: 'color' },
      title: '活动背景色',
      width: 110,
    },
    {
      field: 'layout',
      formatter: ({ row }) => layoutLabelOf(row),
      title: '商品布局',
      width: 100,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue, row }) =>
        formatShanghaiDateTime(cellValue || row.updated_at) || '—',
      minWidth: 170,
      showOverflow: false,
      title: '添加时间',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 110,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listMarketingDecorApi('topic', {
          page: page.currentPage,
          limit: page.pageSize,
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const [ExampleModal, exampleModalApi] = useVbenModal({
  title: '查看示例',
  class: 'w-[820px] max-w-[96vw]',
  footer: false,
});

function asNumber(raw: unknown): number {
  const n = Number(raw);
  return Number.isFinite(n) ? n : 0;
}

function asString(raw: unknown): string {
  return String(raw ?? '').trim();
}

function asStringList(raw: unknown): string[] {
  if (!Array.isArray(raw)) return [];
  return raw.map((x) => asString(x)).filter(Boolean);
}

function payloadOf(row: MarketingDecor): Record<string, unknown> {
  return (row.payload || {}) as Record<string, unknown>;
}

function labelIdOf(row: MarketingDecor): number {
  return asNumber(payloadOf(row).label_id);
}

function labelNameOf(row: MarketingDecor): string {
  const id = labelIdOf(row);
  if (!id) return '—';
  return labelMap.value.get(id) || `标签#${id}`;
}

function bannersOf(row: MarketingDecor): string[] {
  return asStringList(payloadOf(row).banner);
}

function themeOf(row: MarketingDecor): string {
  return asString(payloadOf(row).image);
}

function colorOf(row: MarketingDecor): string {
  return asString(payloadOf(row).color);
}

function layoutOf(row: MarketingDecor): LayoutType {
  const n = asNumber(payloadOf(row).type);
  return n === 2 || n === 3 ? n : 1;
}

function layoutLabelOf(row: MarketingDecor): string {
  return LAYOUT_LABEL[layoutOf(row)] || '—';
}

function mediaUrl(url: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    label_id: undefined,
    name: '',
    cover_url: '',
    banners: [],
    theme_url: '',
    color: '',
    layout: 1,
    status: 1,
    sort: 0,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增专场' }).open();
}

function openEdit(row: MarketingDecor) {
  editing.value = row;
  Object.assign(form, {
    label_id: labelIdOf(row) || undefined,
    name: row.name || '',
    cover_url: row.cover_url || '',
    banners: bannersOf(row),
    theme_url: themeOf(row),
    color: colorOf(row),
    layout: layoutOf(row),
    status: row.status === 1 ? 1 : 0,
    sort: Number(row.sort) || 0,
  });
  formDrawerApi.setState({ title: '编辑专场' }).open();
}

async function save() {
  if (!form.label_id) {
    ElMessage.warning('请选择关联标签');
    return;
  }
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请填写活动名称');
    return;
  }
  const cover = form.cover_url.trim();
  if (!cover) {
    ElMessage.warning('请上传活动列表图');
    return;
  }
  const layout = form.layout === 2 || form.layout === 3 ? form.layout : 1;
  const body = {
    name,
    cover_url: cover,
    status: form.status === 1 ? 1 : 0,
    sort: Number(form.sort) || 0,
    payload: {
      label_id: form.label_id,
      banner: form.banners.map((x) => x.trim()).filter(Boolean),
      image: form.theme_url.trim(),
      color: form.color.trim(),
      type: layout,
    },
  };

  formDrawerApi.lock();
  try {
    if (editing.value) {
      await updateMarketingDecorApi('topic', editing.value.id, body);
    } else {
      await createMarketingDecorApi('topic', body);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: MarketingDecor, enabled: boolean) {
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await setMarketingDecorStatusApi('topic', row.id, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function removeRow(row: MarketingDecor) {
  try {
    await confirm({
      content: `确定删除专场「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteMarketingDecorApi('topic', row.id);
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

onMounted(async () => {
  const [profile, codes, labelRes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    fetchProductLabels().catch(() => ({ list: [] as ProductLabelRow[] })),
  ]);
  labelOptions.value = (labelRes.list || []).filter((x) => Number(x.status) !== 0);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('marketing.topic.read') ||
      codes.includes('marketing.topic.manage'));
  canManage.value = roleOK && codes.includes('marketing.topic.manage');
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
          新增专场
        </ElButton>
      </template>

      <template #listPic="{ row }">
        <ElImage
          v-if="mediaUrl(row.cover_url)"
          :src="mediaUrl(row.cover_url)"
          fit="cover"
          class="topic-thumb"
          :preview-src-list="[mediaUrl(row.cover_url)]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>

      <template #banners="{ row }">
        <div v-if="bannersOf(row).length" class="topic-banner-cell">
          <ElImage
            v-for="(url, idx) in bannersOf(row).slice(0, 3)"
            :key="`${row.id}-b-${idx}`"
            :src="mediaUrl(url)"
            fit="cover"
            class="topic-thumb topic-thumb--sm"
            :preview-src-list="bannersOf(row).map((x) => mediaUrl(x))"
            :initial-index="idx"
          >
            <template #error>
              <span class="text-xs text-gray-400">—</span>
            </template>
          </ElImage>
        </div>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>

      <template #theme="{ row }">
        <ElImage
          v-if="mediaUrl(themeOf(row))"
          :src="mediaUrl(themeOf(row))"
          fit="cover"
          class="topic-thumb"
          :preview-src-list="[mediaUrl(themeOf(row))]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>

      <template #color="{ row }">
        <div v-if="colorOf(row)" class="topic-color-cell">
          <span
            class="topic-color-swatch"
            :style="{ backgroundColor: colorOf(row) }"
          />
          <span class="topic-color-text">{{ colorOf(row) }}</span>
        </div>
        <span v-else class="text-xs text-gray-400">默认</span>
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
      <ElForm label-width="120px" class="topic-form">
        <ElFormItem label="关联标签" required>
          <ElSelect
            v-model="form.label_id"
            class="w-full"
            clearable
            filterable
            placeholder="请选择"
          >
            <ElOption
              v-for="item in labelOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </ElSelect>
        </ElFormItem>

        <ElFormItem label="活动名称">
          <ElInput
            v-model="form.name"
            maxlength="20"
            show-word-limit
            placeholder="请输入活动名称"
          />
        </ElFormItem>

        <ElFormItem label="活动列表图" required>
          <div class="topic-pic">
            <ImageField v-model="form.cover_url" :preview-size="88" />
            <div class="topic-pic__side">
              <ElButton link type="primary" @click="openExample('list')">
                查看示例
              </ElButton>
              <div class="field-hint">建议尺寸：710 × 340</div>
            </div>
          </div>
        </ElFormItem>

        <ElFormItem label="活动轮播图">
          <div class="topic-pic topic-pic--block">
            <ImagesField v-model="form.banners" :limit="9" :preview-size="72" />
            <div class="topic-pic__side">
              <ElButton link type="primary" @click="openExample('banner')">
                查看示例
              </ElButton>
              <div class="field-hint">建议尺寸：750 × 750，可上传多张</div>
            </div>
          </div>
        </ElFormItem>

        <ElFormItem label="活动主题">
          <div class="topic-pic">
            <ImageField v-model="form.theme_url" :preview-size="88" />
            <div class="topic-pic__side">
              <ElButton link type="primary" @click="openExample('theme')">
                查看示例
              </ElButton>
              <div class="field-hint">建议尺寸：710 × 340</div>
            </div>
          </div>
        </ElFormItem>

        <ElFormItem label="活动背景色">
          <div class="topic-color-field">
            <ElColorPicker v-model="form.color" :show-alpha="false" />
            <ElInput
              v-model="form.color"
              class="topic-color-input"
              clearable
              placeholder="未设置则使用默认色"
            />
            <ElButton link type="primary" @click="openExample('color')">
              查看示例
            </ElButton>
          </div>
          <div class="field-hint">未设置背景色时，前端使用默认主题色</div>
        </ElFormItem>

        <ElFormItem label="商品布局" required>
          <ElRadioGroup v-model="form.layout">
            <ElRadio :label="1">小图</ElRadio>
            <ElRadio :label="2">中图</ElRadio>
            <ElRadio :label="3">大图</ElRadio>
          </ElRadioGroup>
        </ElFormItem>

        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="显示"
            inactive-text="隐藏"
          />
        </ElFormItem>

        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :step="1" class="w-full" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <ExampleModal>
      <div class="example-wrap">
        <ElImage
          :src="EXAMPLE_URLS[exampleKind]"
          fit="contain"
          class="example-img"
        />
        <div class="field-hint">
          <template v-if="exampleKind === 'list'">建议尺寸：710 × 340</template>
          <template v-else-if="exampleKind === 'banner'">建议尺寸：750 × 750</template>
          <template v-else-if="exampleKind === 'theme'">建议尺寸：710 × 340</template>
          <template v-else>背景色未设置时使用默认色；示意为商城专场页效果</template>
        </div>
      </div>
    </ExampleModal>
  </Page>
</template>

<style scoped>
.topic-thumb {
  display: block;
  width: 64px;
  height: 36px;
  border-radius: 4px;
}

.topic-thumb--sm {
  width: 36px;
  height: 36px;
}

.topic-banner-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.topic-color-cell {
  display: inline-flex;
  gap: 6px;
  align-items: center;
}

.topic-color-swatch {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 1px solid hsl(var(--border));
  border-radius: 4px;
}

.topic-color-text {
  font-size: 12px;
  color: hsl(var(--foreground) / 75%);
}

.topic-form {
  padding-top: 8px;
}

.field-hint {
  margin-top: 6px;
  color: hsl(var(--foreground) / 55%);
  font-size: 12px;
  line-height: 1.4;
}

.topic-pic {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-start;
}

.topic-pic--block {
  flex-direction: column;
}

.topic-pic__side {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 4px;
}

.topic-color-field {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  width: 100%;
}

.topic-color-input {
  width: 200px;
  max-width: 100%;
}

.example-wrap {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: center;
  padding: 8px 0 16px;
}

.example-img {
  width: 100%;
  max-width: 710px;
  min-height: 160px;
  border-radius: 4px;
  background: hsl(var(--accent));
}
</style>
