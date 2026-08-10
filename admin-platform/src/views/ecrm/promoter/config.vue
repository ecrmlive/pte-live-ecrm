<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElCheckbox,
  ElCheckboxGroup,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSwitch,
  ElTabs,
  ElTabPane,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformDistributionConfigApi,
  savePlatformDistributionConfigApi,
  type PlatformDistributionConfig,
} from '#/api/core/platform-mall-setting';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

type TabKey = 'distribution' | 'commission' | 'withdraw' | 'gift';

const activeTab = ref<TabKey>('distribution');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const note = ref('');

const form = ref<PlatformDistributionConfig>(emptyForm());

const showLowMoney = computed(() => form.value.promoter_type === 3);

function emptyForm(): PlatformDistributionConfig {
  return {
    extension_status: false,
    extension_self: false,
    extension_limit: false,
    extension_limit_day: 15,
    promoter_type: 0,
    promoter_low_money: 0,
    extension_pop: 0,
    extension_one_rate: 0,
    extension_two_rate: 0,
    user_extract_min: 1,
    lock_brokerage_timer: 0,
    sys_extension_type: 0,
    withdraw_type: ['1'],
    extract_switch: 1,
    transfer_scene_id: 0,
    max_bag_number: 10,
  };
}

function round4(v: number) {
  return Math.round(Number(v || 0) * 10_000) / 10_000;
}

function validate(): string | null {
  const f = form.value;
  if (!Number.isInteger(f.extension_limit_day) || f.extension_limit_day <= 0) {
    return '分销绑定时间必须大于 0 天';
  }
  if (f.promoter_type === 3 && (!(f.promoter_low_money > 0))) {
    return '满额分销最低金额必须大于 0';
  }
  if (f.extension_one_rate < 0 || f.extension_two_rate < 0) {
    return '返佣比例不能小于 0';
  }
  if (round4(f.extension_one_rate) < round4(f.extension_two_rate)) {
    return '一级比例不能小于二级比例';
  }
  if (round4(f.extension_one_rate + f.extension_two_rate) > 1) {
    return '一二级比例之和不能超过 1（即 100%）';
  }
  if (f.user_extract_min < 0) {
    return '最低提现金额不能小于 0';
  }
  if (!f.withdraw_type?.length) {
    return '请至少选择一种提现方式';
  }
  return null;
}

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformDistributionConfigApi();
    note.value = data.note || '';
    form.value = {
      ...emptyForm(),
      ...data.config,
      withdraw_type: Array.isArray(data.config.withdraw_type)
        ? [...data.config.withdraw_type]
        : ['1'],
    };
  } finally {
    loading.value = false;
  }
}

async function save() {
  const err = validate();
  if (err) {
    ElMessage.warning(err);
    return;
  }
  saving.value = true;
  try {
    form.value = await savePlatformDistributionConfigApi({
      ...form.value,
      extension_one_rate: round4(form.value.extension_one_rate),
      extension_two_rate: round4(form.value.extension_two_rate),
      withdraw_type: [...form.value.withdraw_type],
    });
    ElMessage.success('分销配置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [codes] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = codes.includes('promoter.config.manage');
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <template #tabs>
        <ElTabs v-model="activeTab">
          <ElTabPane label="分销设置" name="distribution" />
          <ElTabPane label="返佣设置" name="commission" />
          <ElTabPane label="提现设置" name="withdraw" />
          <ElTabPane label="分销礼包设置" name="gift" />
        </ElTabs>
      </template>

      <ElForm
        :disabled="!canManage"
        label-width="180px"
        class="distribution-config-form"
      >
        <template v-if="activeTab === 'distribution'">
          <ElFormItem label="分销启用">
            <div>
              <ElSwitch
                v-model="form.extension_status"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="tips">开启后商城分销功能可用；关闭后不可用</div>
            </div>
          </ElFormItem>
          <ElFormItem label="分销模式">
            <div>
              <ElRadioGroup v-model="form.promoter_type">
                <ElRadio :value="0">礼包分销</ElRadio>
                <ElRadio :value="1">手动分销</ElRadio>
                <ElRadio :value="2">人人分销</ElRadio>
                <ElRadio :value="3">满额分销</ElRadio>
              </ElRadioGroup>
              <div class="tips">
                礼包：购买分销礼包后成为分销员；手动：仅后台设置；人人：注册即可；满额：消费达门槛自动开通
              </div>
            </div>
          </ElFormItem>
          <ElFormItem v-if="showLowMoney" label="满额分销最低金额（元）">
            <div>
              <ElInputNumber
                v-model="form.promoter_low_money"
                :min="0.01"
                :max="1000000"
                :precision="2"
                :step="1"
              />
              <div class="tips">用户消费金额达到设置金额（含）即可自动开通分销权限</div>
            </div>
          </ElFormItem>
          <ElFormItem label="分销内购">
            <div>
              <ElSwitch
                v-model="form.extension_self"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="tips">
                开启：分销员自购享受一级返佣，上一级享受二级；关闭：自购无返佣
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="分销限时开关">
            <div>
              <ElSwitch
                v-model="form.extension_limit"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="tips">
                开启后按绑定时段返佣；关闭为永久绑定（不建议频繁修改）
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="分销绑定时间（天）">
            <div>
              <ElInputNumber
                v-model="form.extension_limit_day"
                :min="1"
                :max="3650"
                :precision="0"
                :step="1"
              />
              <div class="tips">绑定成功至自动解绑的天数；解绑后按新关系结算</div>
            </div>
          </ElFormItem>
          <ElFormItem label="佣金悬浮窗">
            <div>
              <ElRadioGroup v-model="form.extension_pop">
                <ElRadio :value="0">全部可见</ElRadio>
                <ElRadio :value="1">推广员可见</ElRadio>
                <ElRadio :value="2">非推广员可见</ElRadio>
                <ElRadio :value="3">关闭</ElRadio>
              </ElRadioGroup>
              <div class="tips">商品详情页最高佣金悬浮框展示范围</div>
            </div>
          </ElFormItem>
        </template>

        <template v-else-if="activeTab === 'commission'">
          <ElFormItem label="一级分销比例">
            <div>
              <ElInputNumber
                v-model="form.extension_one_rate"
                :min="0"
                :max="1"
                :precision="4"
                :step="0.01"
              />
              <div class="tips">例：0.15 = 订单金额的 15%；须 ≥ 二级比例</div>
            </div>
          </ElFormItem>
          <ElFormItem label="二级分销比例">
            <div>
              <ElInputNumber
                v-model="form.extension_two_rate"
                :min="0"
                :max="1"
                :precision="4"
                :step="0.01"
              />
              <div class="tips">一二级比例之和不能超过 1（100%）</div>
            </div>
          </ElFormItem>
          <ElFormItem label="佣金冻结时间（天）">
            <div>
              <ElInputNumber
                v-model="form.lock_brokerage_timer"
                :min="0"
                :max="3650"
                :precision="0"
                :step="1"
              />
              <div class="tips">
                自确认收货起算；0 表示无冻结期，解冻后方可提现
              </div>
            </div>
          </ElFormItem>
        </template>

        <template v-else-if="activeTab === 'withdraw'">
          <ElFormItem label="佣金最低提现金额（元）">
            <div>
              <ElInputNumber
                v-model="form.user_extract_min"
                :min="0"
                :max="1000000"
                :precision="2"
                :step="1"
              />
              <div class="tips">佣金达到该金额后可申请提现</div>
            </div>
          </ElFormItem>
          <ElFormItem label="提现方式">
            <div>
              <ElCheckboxGroup v-model="form.withdraw_type">
                <ElCheckbox value="0">银行卡</ElCheckbox>
                <ElCheckbox value="1">微信</ElCheckbox>
                <ElCheckbox value="2">支付宝</ElCheckbox>
                <ElCheckbox value="4">余额</ElCheckbox>
              </ElCheckboxGroup>
              <div class="tips">未选时默认微信；多数方式需后台手动打款</div>
            </div>
          </ElFormItem>
          <ElFormItem label="微信到账方式">
            <div>
              <ElRadioGroup v-model="form.sys_extension_type">
                <ElRadio :value="0">线下转账</ElRadio>
                <ElRadio :value="1">企业付款到零钱</ElRadio>
                <ElRadio :value="2">商家转账到零钱</ElRadio>
              </ElRadioGroup>
              <div class="tips">
                自动到账需开通对应微信能力并配置证书；本页不保存密钥
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="微信提现方式">
            <div>
              <ElRadioGroup v-model="form.extract_switch">
                <ElRadio :value="1">线下转账</ElRadio>
                <ElRadio :value="2">自动转账</ElRadio>
              </ElRadioGroup>
            </div>
          </ElFormItem>
          <ElFormItem label="转账场景 ID">
            <div>
              <ElInputNumber
                v-model="form.transfer_scene_id"
                :min="0"
                :max="999999"
                :precision="0"
                :step="1"
              />
              <div class="tips">
                商家转账场景，可在微信商户平台「产品中心 - 商家转账」申请，如 1001
              </div>
            </div>
          </ElFormItem>
        </template>

        <template v-else>
          <ElFormItem label="店铺礼包最大数量">
            <div>
              <ElInputNumber
                v-model="form.max_bag_number"
                :min="0"
                :max="9999"
                :precision="0"
                :step="1"
              />
              <div class="tips">每个商户可设置的分销礼包数量上限</div>
            </div>
          </ElFormItem>
        </template>
      </ElForm>

      <p v-if="note" class="distribution-config-note">{{ note }}</p>

      <template #actions>
        <ElButton @click="load">重置</ElButton>
        <ElButton
          type="primary"
          :disabled="!canManage"
          :loading="saving"
          @click="save"
        >
          保存
        </ElButton>
      </template>
    </SettingsTabLayout>
  </Page>
</template>

<style scoped>
.distribution-config-form {
  width: 100%;
  max-width: none;
}

.distribution-config-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.distribution-config-form :deep(.el-form-item__label) {
  align-items: center;
  color: hsl(var(--foreground));
  font-weight: 400;
}

.distribution-config-form :deep(.el-switch) {
  --el-switch-on-color: hsl(var(--primary));
}

.tips {
  margin-top: 6px;
  max-width: 720px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
}

.distribution-config-note {
  margin: 8px 0 0;
  max-width: 720px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
}
</style>
