<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  getReservationConfigApi,
  listReservationProductsApi,
  saveReservationConfigApi,
  type ReservationProduct,
  type ReservationSlot,
} from '#/api/core/merchant-reservation';

const loading = ref(false);
const saving = ref(false);
const rows = ref<ReservationProduct[]>([]);
const total = ref(0);
const configOpen = ref(false);
const current = ref<ReservationProduct>();
const slots = ref<ReservationSlot[]>([]);
const query = reactive({ limit: 20, page: 1 });
const form = reactive({
  reservation_type: 1,
  show_reservation_days: 7,
  slotText: '',
});

async function load() {
  loading.value = true;
  try {
    const result = await listReservationProductsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function slotsToText(items: ReservationSlot[]) {
  return items
    .map((item) => `${item.start_time}-${item.end_time},${item.stock}`)
    .join('\n');
}

function parseSlots(text: string) {
  return text
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      const [timePart, stockPart] = line.split(',');
      const [start_time = '', end_time = ''] = (timePart || '').split('-');
      return {
        start_time: start_time.trim(),
        end_time: end_time.trim(),
        stock: Number(stockPart) || 0,
      };
    })
    .filter((item) => item.start_time && item.end_time);
}

async function openConfig(row: ReservationProduct) {
  current.value = row;
  const result = await getReservationConfigApi(row.product_id);
  form.reservation_type = result.config?.reservation_type ?? row.reservation_type ?? 1;
  form.show_reservation_days = result.config?.show_reservation_days ?? row.show_reservation_days ?? 7;
  slots.value = result.slots || [];
  form.slotText = slotsToText(slots.value);
  configOpen.value = true;
}

async function saveConfig() {
  if (!current.value) return;
  saving.value = true;
  try {
    await saveReservationConfigApi(current.value.product_id, {
      reservation_type: form.reservation_type,
      show_reservation_days: form.show_reservation_days,
      slots: parseSlots(form.slotText),
    });
    ElMessage.success('预约配置已保存');
    configOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="预约设置" description="管理本店预约服务商品的可预约天数与时段库存；仅 type=预约 的商品会出现在列表中。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="product_id">
        <el-table-column label="商品 ID" prop="product_id" width="100" />
        <el-table-column label="商品名称" min-width="200" prop="store_name" show-overflow-tooltip />
        <el-table-column label="售价" width="100">
          <template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="库存" prop="stock" width="90" />
        <el-table-column label="可预约天数" prop="show_reservation_days" width="110" />
        <el-table-column label="预约类型" width="100">
          <template #default="{ row }">{{ row.reservation_type === 2 ? '按时段' : '按日期' }}</template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="100">
          <template #default="{ row }"><el-button link type="primary" @click="openConfig(row)">配置</el-button></template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :page-sizes="[10, 20, 50]"
          :total="total"
          background
          layout="total, sizes, prev, pager, next"
          @current-change="(page: number) => { query.page = page; load(); }"
          @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }"
        />
      </div>
    </el-card>

    <el-dialog v-model="configOpen" destroy-on-close title="预约配置" width="640px">
      <template v-if="current">
        <div class="mb-4 text-base font-medium">{{ current.store_name }}</div>
        <el-form label-width="112px">
          <el-form-item label="预约类型">
            <el-radio-group v-model="form.reservation_type">
              <el-radio :value="1">按日期</el-radio>
              <el-radio :value="2">按时段</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="可预约天数"><el-input-number v-model="form.show_reservation_days" :min="1" :max="90" /></el-form-item>
          <el-form-item label="时段与库存">
            <el-input
              v-model="form.slotText"
              :rows="6"
              placeholder="每行一条：09:00-10:00,20&#10;10:00-11:00,15"
              type="textarea"
            />
          </el-form-item>
        </el-form>
      </template>
      <template #footer>
        <el-button @click="configOpen = false">取消</el-button>
        <el-button :loading="saving" type="primary" @click="saveConfig">保存</el-button>
      </template>
    </el-dialog>
  </Page>
</template>
