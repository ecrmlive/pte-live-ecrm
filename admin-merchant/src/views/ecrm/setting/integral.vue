<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { getMerchantIntegralPolicyApi, saveMerchantIntegralPolicyApi } from '#/api/core/merchant-integral';

const loading = ref(false);
const saving = ref(false);
const form = ref({ enabled: false, points_per_yuan: 100, max_deduction_bps: 2000 });

const maxDeductionPercent = computed(() => (form.value.max_deduction_bps / 100).toFixed(1));

async function load() {
  loading.value = true;
  try {
    const data = await getMerchantIntegralPolicyApi();
    form.value = {
      enabled: data.enabled,
      points_per_yuan: data.points_per_yuan,
      max_deduction_bps: data.max_deduction_bps,
    };
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    const data = await saveMerchantIntegralPolicyApi(form.value);
    form.value = {
      enabled: data.enabled,
      points_per_yuan: data.points_per_yuan,
      max_deduction_bps: data.max_deduction_bps,
    };
    ElMessage.success('积分抵扣设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void load();
});
</script>

<template>
  <Page title="积分抵扣设置" description="配置本店下单时是否允许积分抵扣，以及抵扣比例上限。">
    <el-card v-loading="loading" shadow="never">
      <el-alert
        class="mb-4"
        type="info"
        :closable="false"
        title="max_deduction_bps 为万分比：2000 表示最多抵扣订单金额的 20%。"
      />
      <el-form label-width="180px" class="max-w-3xl">
        <el-form-item label="启用积分抵扣">
          <el-switch v-model="form.enabled" />
        </el-form-item>
        <el-form-item label="每元所需积分">
          <el-input-number v-model="form.points_per_yuan" :min="1" :max="100000" />
          <span class="ml-2 text-sm text-muted-foreground">100 表示 1 元 = 100 积分</span>
        </el-form-item>
        <el-form-item label="最大抵扣比例（bps）">
          <el-input-number v-model="form.max_deduction_bps" :min="1" :max="10000" :step="100" />
          <span class="ml-2 text-sm text-muted-foreground">当前约 {{ maxDeductionPercent }}%</span>
        </el-form-item>
      </el-form>
      <div class="mt-4 flex justify-end">
        <el-button @click="load">重置</el-button>
        <el-button :loading="saving" type="primary" @click="save">保存设置</el-button>
      </div>
    </el-card>
  </Page>
</template>
