<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformSeckillTimeApi,
  deletePlatformSeckillTimeApi,
  listPlatformSeckillTimesApi,
  setPlatformSeckillTimeStatusApi,
  updatePlatformSeckillTimeApi,
  type PlatformSeckillTime,
} from '#/api/core/platform-seckill';
import ImageField from '#/components/shop/image-field.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const form = reactive({
  title: '',
  start_time: 0,
  end_time: 6,
  status: 1,
  pic: '',
});

const startHourOptions = computed(() =>
  Array.from({ length: 24 }, (_, h) => ({
    label: `${h}:00`,
    value: h,
  })),
);

const endHourOptions = computed(() =>
  Array.from({ length: 24 }, (_, i) => {
    const h = i + 1;
    return { label: `${h}:00`, value: h };
  }),
);

function formatHour(value?: number) {
  if (value === undefined || value === null) return '—';
  return `${Number(value)}:00`;
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    title: '',
    start_time: 0,
    end_time: 6,
    status: 1,
    pic: '',
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '显示', value: 1 },
        { label: '不显示', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '是否显示',
  },
]);

const gridOptions: VxeGridProps<PlatformSeckillTime> = {
  columns: [
    { field: 'seckill_time_id', title: '编号', width: 80 },
    {
      field: 'title',
      minWidth: 140,
      showOverflow: false,
      title: '场次名称',
    },
    {
      field: 'start_time',
      title: '开始时间(整数小时)',
      width: 160,
      formatter: ({ cellValue }) => formatHour(cellValue),
    },
    {
      field: 'end_time',
      title: '结束时间(整点)',
      width: 140,
      formatter: ({ cellValue }) => formatHour(cellValue),
    },
    {
      field: 'pic',
      slots: { default: 'pic' },
      title: '图片',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 110,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const statusRaw = formValues?.status;
        const data = await listPlatformSeckillTimesApi({
          page: page.currentPage,
          limit: page.pageSize,
          status:
            statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'seckill_time_id' },
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
  placement: 'right',
  onConfirm: () => void save(),
});

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增秒杀配置' }).open();
}

function openEdit(row: PlatformSeckillTime) {
  editingID.value = row.seckill_time_id;
  Object.assign(form, {
    title: row.title || '',
    start_time: Number(row.start_time ?? 0),
    end_time: Number(row.end_time ?? 1),
    status: Number(row.status ?? 1),
    pic: row.pic || '',
  });
  formDrawerApi.setState({ title: '编辑配置' }).open();
}

async function save() {
  if (!form.title.trim()) {
    ElMessage.warning('请输入场次名称');
    return;
  }
  if (form.start_time >= form.end_time) {
    ElMessage.warning('结束时间必须大于开始时间');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const body = {
      title: form.title.trim(),
      start_time: form.start_time,
      end_time: form.end_time,
      status: form.status,
      pic: form.pic || '',
    };
    if (editingID.value) {
      await updatePlatformSeckillTimeApi(editingID.value, body);
      ElMessage.success('已更新秒杀配置');
    } else {
      await createPlatformSeckillTimeApi(body);
      ElMessage.success('已添加秒杀配置');
    }
    formDrawerApi.close();
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function toggleStatus(row: PlatformSeckillTime) {
  const next = row.status === 1 ? 0 : 1;
  await setPlatformSeckillTimeStatusApi(row.seckill_time_id, next);
  row.status = next;
  ElMessage.success(next === 1 ? '已显示' : '已隐藏');
}

async function remove(row: PlatformSeckillTime) {
  try {
    await ElMessageBox.confirm(
      `删除场次「${row.title}」后不可恢复，是否继续？`,
      '删除确认',
      { type: 'warning' },
    );
    await deletePlatformSeckillTimeApi(row.seckill_time_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增秒杀配置
        </ElButton>
      </template>
      <template #pic="{ row }">
        <ElImage
          v-if="row.pic"
          :src="resolveCosMediaUrl(row.pic)"
          fit="cover"
          class="seckill-pic"
        >
          <template #error>
            <div class="seckill-pic seckill-pic--empty">加载失败</div>
          </template>
        </ElImage>
        <div v-else class="seckill-pic seckill-pic--empty">—</div>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          @change="() => toggleStatus(row)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="120px" class="seckill-config-form" v-loading="saving">
        <ElFormItem label="场次名称" required>
          <ElInput
            v-model="form.title"
            maxlength="32"
            show-word-limit
            placeholder="请输入场次名称"
            class="max-w-md"
          />
        </ElFormItem>
        <ElFormItem label="开始时间" required>
          <ElSelect v-model="form.start_time" placeholder="请选择" class="max-w-xs">
            <ElOption
              v-for="opt in startHourOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="结束时间" required>
          <ElSelect v-model="form.end_time" placeholder="请选择" class="max-w-xs">
            <ElOption
              v-for="opt in endHourOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="是否启用">
          <ElSwitch
            :model-value="form.status === 1"
            @change="(v: string | number | boolean) => (form.status = v ? 1 : 0)"
          />
        </ElFormItem>
        <ElFormItem label="图片">
          <div>
            <ImageField v-model="form.pic" />
            <div class="form-tip">
              此图片将展示在移动端秒杀商品列表上方（建议尺寸：710×300px）
            </div>
          </div>
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.seckill-pic {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  overflow: hidden;
}

.seckill-pic--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: hsl(var(--muted) / 0.45);
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}

.form-tip {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.5;
  color: hsl(var(--muted-foreground));
}

.seckill-config-form :deep(.el-form-item) {
  margin-bottom: 18px;
}
</style>
