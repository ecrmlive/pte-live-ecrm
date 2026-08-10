<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { VbenTiptap } from '@vben/plugins/tiptap';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
} from 'element-plus';

import { uploadAttachmentApi } from '#/api/core/attachment';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformIntegralConfigApi,
  savePlatformIntegralConfigApi,
  type PlatformIntegralConfig,
} from '#/api/core/platform-mall-setting';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);

const form = reactive<PlatformIntegralConfig>({
  integral_status: 1,
  integral_money: 0.1,
  integral_order_rate: 1,
  integral_freeze: 0,
  integral_clear_time: 24,
  integral_user_give: 50,
  integral_community_give: 10,
  integral_community_give_limit: 10,
  rule: '',
});

const imageUpload: ImageUploadOptions = {
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const row = await uploadAttachmentApi(file);
    return row.attachment_src;
  },
  onUploadError: () => {
    ElMessage.error('图片上传失败');
  },
};

function applyConfig(data: PlatformIntegralConfig) {
  form.integral_status = data.integral_status === 1 ? 1 : 0;
  form.integral_money = Number(data.integral_money ?? 0.1);
  form.integral_order_rate = Number(data.integral_order_rate ?? 1);
  form.integral_freeze = Math.trunc(Number(data.integral_freeze ?? 0));
  form.integral_clear_time = Math.trunc(Number(data.integral_clear_time ?? 24));
  form.integral_user_give = Math.trunc(Number(data.integral_user_give ?? 50));
  form.integral_community_give = Math.trunc(
    Number(data.integral_community_give ?? 10),
  );
  form.integral_community_give_limit = Math.trunc(
    Number(data.integral_community_give_limit ?? 10),
  );
  form.rule = data.rule || '';
}

function validate(): string | null {
  if (form.integral_status !== 0 && form.integral_status !== 1) {
    return '请选择积分开关';
  }
  // 开启态才校验完整积分规则字段；关闭态仅开关 + 种草两项
  if (form.integral_status === 1) {
    if (
      !Number.isFinite(form.integral_money) ||
      form.integral_money < 0
    ) {
      return '积分抵用金额请填写不小于 0 的数字';
    }
    if (
      !Number.isFinite(form.integral_order_rate) ||
      form.integral_order_rate < 0
    ) {
      return '下单赠送积分比例请填写不小于 0 的数字';
    }
    if (!Number.isInteger(form.integral_freeze) || form.integral_freeze < 0) {
      return '下单赠送积分冻结期请填写不小于 0 的整数';
    }
    if (
      !Number.isInteger(form.integral_clear_time) ||
      form.integral_clear_time < 0
    ) {
      return '积分清除时间设置请填写不小于 0 的整数';
    }
    if (
      !Number.isInteger(form.integral_user_give) ||
      form.integral_user_give < 0
    ) {
      return '邀请好友赠送积分请填写不小于 0 的整数';
    }
  }
  if (
    !Number.isInteger(form.integral_community_give) ||
    form.integral_community_give < 0 ||
    form.integral_community_give > 9999
  ) {
    return '发布种草可获得积分请填写 0～9999 的整数';
  }
  if (
    !Number.isInteger(form.integral_community_give_limit) ||
    form.integral_community_give_limit < 0 ||
    form.integral_community_give_limit > 9999
  ) {
    return '发布种草篇数限量请填写 0～9999 的整数';
  }
  return null;
}

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformIntegralConfigApi();
    applyConfig(data.config);
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
    const saved = await savePlatformIntegralConfigApi({
      integral_status: form.integral_status === 1 ? 1 : 0,
      integral_money: Number(form.integral_money),
      integral_order_rate: Number(form.integral_order_rate),
      integral_freeze: Math.trunc(Number(form.integral_freeze)),
      integral_clear_time: Math.trunc(Number(form.integral_clear_time)),
      integral_user_give: Math.trunc(Number(form.integral_user_give)),
      integral_community_give: Math.trunc(Number(form.integral_community_give)),
      integral_community_give_limit: Math.trunc(
        Number(form.integral_community_give_limit),
      ),
      rule: form.rule || '',
    });
    applyConfig(saved);
    ElMessage.success('保存成功');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value =
    codes.includes('marketing.integral.config') ||
    codes.includes('marketing.points.manage');
  await load();
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <ElForm
        :disabled="!canManage"
        class="integral-config-form"
        label-width="180px"
      >
        <ElFormItem label="积分：" required>
          <div>
            <ElRadioGroup v-model="form.integral_status">
              <ElRadio :label="0">关闭</ElRadio>
              <ElRadio :label="1">开启</ElRadio>
            </ElRadioGroup>
            <div class="form-tip">
              商城积分功能开启关闭，关闭后前端不展示积分模块
            </div>
          </div>
        </ElFormItem>

        <template v-if="form.integral_status === 1">
          <ElFormItem label="积分抵用金额：" required>
            <div>
              <div class="integral-config-form__inline">
                <ElInputNumber
                  v-model="form.integral_money"
                  :min="0"
                  :precision="2"
                  :step="0.01"
                  controls-position="right"
                />
                <span class="integral-config-form__unit">元</span>
              </div>
              <div class="form-tip">积分抵用比例(1积分抵多少金额)</div>
            </div>
          </ElFormItem>

          <ElFormItem label="下单赠送积分比例：" required>
            <div>
              <div class="integral-config-form__inline">
                <ElInputNumber
                  v-model="form.integral_order_rate"
                  :min="0"
                  :precision="2"
                  :step="1"
                  controls-position="right"
                />
                <span class="integral-config-form__unit">分</span>
              </div>
              <div class="form-tip">消费1元赠送多少积分</div>
            </div>
          </ElFormItem>

          <ElFormItem label="下单赠送积分冻结期：" required>
            <div>
              <div class="integral-config-form__inline">
                <ElInputNumber
                  v-model="form.integral_freeze"
                  :min="0"
                  :precision="0"
                  :step="1"
                  controls-position="right"
                />
                <span class="integral-config-form__unit">天</span>
              </div>
              <div class="form-tip">
                获得的积分多少天之后才能使用；0 无冻结；约20分钟检查一次
              </div>
            </div>
          </ElFormItem>

          <ElFormItem label="积分清除时间设置：" required>
            <div>
              <div class="integral-config-form__inline">
                <ElInputNumber
                  v-model="form.integral_clear_time"
                  :min="0"
                  :precision="0"
                  :step="1"
                  controls-position="right"
                />
                <span class="integral-config-form__unit">月</span>
              </div>
              <div class="form-tip">每满 N 个月清除一次</div>
            </div>
          </ElFormItem>

          <ElFormItem label="邀请好友赠送积分：" required>
            <div>
              <div class="integral-config-form__inline">
                <ElInputNumber
                  v-model="form.integral_user_give"
                  :min="0"
                  :precision="0"
                  :step="1"
                  controls-position="right"
                />
                <span class="integral-config-form__unit">分</span>
              </div>
              <div class="form-tip">邀请好友注册登录商城赠送积分</div>
            </div>
          </ElFormItem>
        </template>

        <ElFormItem label="发布种草可获得积分：" required>
          <div>
            <div class="integral-config-form__inline">
              <ElInputNumber
                v-model="form.integral_community_give"
                :max="9999"
                :min="0"
                :precision="0"
                :step="1"
                controls-position="right"
              />
              <span class="integral-config-form__unit">分</span>
            </div>
            <div class="form-tip">发布一篇种草获得多少积分</div>
          </div>
        </ElFormItem>

        <ElFormItem label="发布种草篇数限量：" required>
          <div>
            <ElInputNumber
              v-model="form.integral_community_give_limit"
              :max="9999"
              :min="0"
              :precision="0"
              :step="1"
              controls-position="right"
            />
            <div class="form-tip">每日最多几篇种草可获积分</div>
          </div>
        </ElFormItem>

        <ElFormItem v-if="form.integral_status === 1" label="积分说明：">
          <VbenTiptap
            v-model="form.rule"
            :editable="canManage"
            :image-upload="imageUpload"
            :max-height="420"
            :min-height="280"
            :previewable="false"
            placeholder="请输入积分说明…"
          />
        </ElFormItem>
      </ElForm>

      <template #actions>
        <ElButton
          type="primary"
          :disabled="!canManage"
          :loading="saving"
          @click="save"
        >
          提交
        </ElButton>
      </template>
    </SettingsTabLayout>
  </Page>
</template>

<style scoped>
.integral-config-form {
  width: 100%;
  max-width: none;
  padding: 8px 0 0;
}

.integral-config-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.integral-config-form :deep(.el-form-item__label) {
  align-items: flex-start;
  padding-top: 6px;
  color: hsl(var(--foreground));
  font-weight: 400;
}

.integral-config-form__inline {
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.integral-config-form__unit {
  color: hsl(var(--foreground));
  font-size: 14px;
  line-height: 32px;
}

.form-tip {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.5;
  color: hsl(var(--muted-foreground));
}
</style>
