<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformAgentZoneConfigApi,
  savePlatformAgentZoneConfigApi,
  type MerchantApplyFormField,
  type PlatformAgentZoneConfig,
} from '#/api/core/platform-mall-setting';

import FormDiyEditor from '../merchant/apply-form-diy/FormDiyEditor.vue';
import type { ApplyFormField } from '../merchant/apply-form-diy/types';
import { AGENT_SYSTEM_FIELDS } from './agent-system-fields';

type TabName = 'commission' | 'form';

/** 对齐 CRMEB ZoneDefaultCommission 说明表 */
const COMMISSION_RULE_ROWS = [
  {
    orderType: '1. 三级区域关联店铺的订单',
    calcLogicTips: '三级代理拿 “自身层级提成”，一二级代理拿 “上下级差价提成”',
    calcLogicList: [
      '省代提成 = 平台抽成 × (省代总比例 - 市代比例)',
      '市代提成 = 平台抽成 × (市代比例 - 区代比例)',
      '区代提成 = 平台抽成 × 区代比例',
    ],
    exampleList: [
      '省代提成 = 10万 × (8% - 5%)=3000元',
      '市代提成 = 10万 × (5% - 3%)=2000元',
      '区代提成 = 10万 × 3%=3000元',
    ],
  },
  {
    orderType: '2. 二级区域关联店铺的订单',
    calcLogicTips: '二级代理拿 “自身层级提成”，一级代理拿 “上下级差价提成”',
    calcLogicList: [
      '省代提成 = 平台抽成 × (省代总比例 - 市代比例)',
      '市代提成 = 平台抽成 × 市代比例',
    ],
    exampleList: [
      '省代提成 = 10万 × (8% - 5%)=3000元',
      '市代提成 = 10万 × 5%=5000元',
    ],
  },
  {
    orderType: '3. 一级区域关联店铺的订单',
    calcLogicTips: '一级代理拿区域代理全级提成',
    calcLogicList: ['省代提成 = 平台抽成 × 省代总比例'],
    exampleList: ['省代提成 = 10万 × 8%=8000元'],
  },
];

const activeTab = ref<TabName>('commission');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformAgentZoneConfig>({
  one_agent_commission: 0,
  two_agent_commission: 0,
  three_agent_commission: 0,
  form_fields: [],
});

function normalizeLoadedFields(
  list: MerchantApplyFormField[],
): ApplyFormField[] {
  return list.map((field) => ({
    ...field,
    content_type: field.content_type || 'text',
    max_upload:
      field.type === 'image' ? field.max_upload || 8 : field.max_upload,
    city_level:
      field.type === 'city'
        ? field.city_level || 'province_city_district'
        : field.city_level,
    default_visible: field.default_visible,
    default_mode: field.default_mode,
    specify_value: field.specify_value || '',
  }));
}

function serializeFields(list: ApplyFormField[]): MerchantApplyFormField[] {
  return list.map((field) => ({
    ...field,
    title: field.title.trim(),
    placeholder: (field.placeholder || '').trim(),
    default_value: (field.default_value || '').trim(),
    content_type: field.content_type || 'text',
    options:
      field.type === 'radio' ||
      field.type === 'checkbox' ||
      field.type === 'select'
        ? field.options?.filter(Boolean)
        : undefined,
    max_upload: field.type === 'image' ? field.max_upload || 8 : undefined,
    city_level: field.type === 'city' ? field.city_level : undefined,
    default_visible:
      field.type === 'date' ||
      field.type === 'daterange' ||
      field.type === 'time' ||
      field.type === 'timerange'
        ? field.default_visible || 'show'
        : undefined,
    default_mode:
      field.type === 'date' ||
      field.type === 'daterange' ||
      field.type === 'time' ||
      field.type === 'timerange'
        ? field.default_mode || 'current'
        : undefined,
    specify_value:
      field.type === 'date' ||
      field.type === 'daterange' ||
      field.type === 'time' ||
      field.type === 'timerange'
        ? (field.specify_value || '').trim()
        : undefined,
  }));
}

function validateCommission(): string | null {
  const one = Number(form.value.one_agent_commission);
  const two = Number(form.value.two_agent_commission);
  const three = Number(form.value.three_agent_commission);
  const rates = [one, two, three];
  if (rates.some((v) => Number.isNaN(v) || v < 0 || v > 100)) {
    return '提成比例必须为 0 ~ 100 之间的数值';
  }
  if (one < two) {
    return '一级代理提成比例必须大于等于二级代理提成比例';
  }
  if (two < three) {
    return '二级代理提成比例必须大于等于三级代理提成比例';
  }
  return null;
}

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformAgentZoneConfigApi();
    form.value = {
      one_agent_commission: Number(data.config.one_agent_commission) || 0,
      two_agent_commission: Number(data.config.two_agent_commission) || 0,
      three_agent_commission: Number(data.config.three_agent_commission) || 0,
      form_fields: normalizeLoadedFields(
        Array.isArray(data.config.form_fields) ? data.config.form_fields : [],
      ),
    };
  } catch {
    ElMessage.error('加载代理设置失败');
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!canManage.value) {
    ElMessage.warning('无权保存代理设置');
    return;
  }
  const commissionError = validateCommission();
  if (commissionError) {
    ElMessage.warning(commissionError);
    activeTab.value = 'commission';
    return;
  }
  const titles = new Set<string>();
  for (const field of form.value.form_fields) {
    const title = field.title.trim();
    if (!title) {
      ElMessage.warning('自定义字段标题不能为空');
      activeTab.value = 'form';
      return;
    }
    if (titles.has(title)) {
      ElMessage.error(`存在重复标题：${title}`);
      activeTab.value = 'form';
      return;
    }
    titles.add(title);
  }
  saving.value = true;
  try {
    const saved = await savePlatformAgentZoneConfigApi({
      one_agent_commission: Number(form.value.one_agent_commission) || 0,
      two_agent_commission: Number(form.value.two_agent_commission) || 0,
      three_agent_commission: Number(form.value.three_agent_commission) || 0,
      form_fields: serializeFields(form.value.form_fields),
    });
    form.value = {
      one_agent_commission: Number(saved.one_agent_commission) || 0,
      two_agent_commission: Number(saved.two_agent_commission) || 0,
      three_agent_commission: Number(saved.three_agent_commission) || 0,
      form_fields: normalizeLoadedFields(
        Array.isArray(saved.form_fields) ? saved.form_fields : [],
      ),
    };
    ElMessage.success('代理设置已保存');
  } catch {
    ElMessage.error('保存代理设置失败');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  // 权限与配置加载解耦：GET 404/失败不得阻断 canManage，否则表单与 DIY 会永久禁用
  const codesPromise = getAccessCodesApi()
    .then((codes) => {
      canManage.value = codes.includes('region.agent_settings.manage');
    })
    .catch(() => {
      canManage.value = false;
    });
  await Promise.all([codesPromise, load()]);
});</script>

<template>
  <Page auto-content-height content-class="!bg-card !p-0">
    <div v-loading="loading" class="agent-setting">
      <div class="agent-setting__tabs">
        <ElTabs v-model="activeTab">
          <ElTabPane label="默认提成" name="commission" />
          <ElTabPane label="代理申请表单" name="form" />
        </ElTabs>
      </div>

      <div v-show="activeTab === 'commission'" class="agent-setting__commission">
        <div class="commission-tip">
          <p class="commission-tip__intro">
            提示：关联店铺的订单提成需在 3 级代理间按 “层级比例”
            <span class="red">递减</span>
            分配，分 3 类订单情况：
          </p>
          <ElTable
            :data="COMMISSION_RULE_ROWS"
            border
            size="small"
            class="commission-tip__table"
          >
            <ElTableColumn prop="orderType" label="订单类型" width="220" />
            <ElTableColumn label="提成计算逻辑" min-width="320">
              <template #default="{ row }">
                <p class="rule-tips">{{ row.calcLogicTips }}</p>
                <p
                  v-for="(line, idx) in row.calcLogicList"
                  :key="`logic-${idx}`"
                  class="rule-item"
                >
                  {{ line }}
                </p>
              </template>
            </ElTableColumn>
            <ElTableColumn label="示例" min-width="280">
              <template #default="{ row }">
                <p class="rule-tips">
                  示例（假设：省代总提成比例 8%，市代提成比例 5%，区代提成
                  3%，平台抽成 10 万元）
                </p>
                <p
                  v-for="(line, idx) in row.exampleList"
                  :key="`ex-${idx}`"
                  class="rule-item"
                >
                  {{ line }}
                </p>
              </template>
            </ElTableColumn>
          </ElTable>
        </div>

        <ElForm label-width="120px" class="commission-form">
          <ElFormItem label="一级代理提成:">
            <ElInputNumber
              v-model="form.one_agent_commission"
              :min="0"
              :max="100"
              :precision="2"
              :step="0.01"
              :disabled="!canManage"
              controls-position="right"
            />
            <span class="unit">%</span>
          </ElFormItem>
          <ElFormItem label="二级代理提成:">
            <ElInputNumber
              v-model="form.two_agent_commission"
              :min="0"
              :max="100"
              :precision="2"
              :step="0.01"
              :disabled="!canManage"
              controls-position="right"
            />
            <span class="unit">%</span>
          </ElFormItem>
          <ElFormItem label="三级代理提成:">
            <ElInputNumber
              v-model="form.three_agent_commission"
              :min="0"
              :max="100"
              :precision="2"
              :step="0.01"
              :disabled="!canManage"
              controls-position="right"
            />
            <span class="unit">%</span>
          </ElFormItem>
        </ElForm>

        <div class="agent-setting__footer">
          <ElButton
            v-if="canManage"
            :loading="saving"
            type="primary"
            @click="save"
          >
            保存
          </ElButton>
        </div>
      </div>

      <div v-show="activeTab === 'form'" class="agent-setting__form">
        <FormDiyEditor
          v-model:fields="form.form_fields"
          :disabled="!canManage"
          preview-title="代理申请"
          :system-fields="AGENT_SYSTEM_FIELDS"
        />
        <div class="agent-setting__footer">
          <ElButton
            v-if="canManage"
            :loading="saving"
            type="primary"
            @click="save"
          >
            保存
          </ElButton>
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped lang="scss">
.agent-setting {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: hsl(var(--card));
}

.agent-setting__tabs {
  flex-shrink: 0;
  padding: 0 16px;
  border-bottom: 1px solid hsl(var(--border));
}

.agent-setting__tabs :deep(.el-tabs__header) {
  margin: 0;
}

.agent-setting__tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.agent-setting__tabs :deep(.el-tabs__item) {
  height: 44px;
  line-height: 44px;
}

.agent-setting__commission {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  overflow: auto;
}

.commission-tip {
  margin: 16px 24px 8px;
  padding: 14px 16px;
  border: 1px solid #ffe58f;
  border-radius: 4px;
  background: #fffbe6;
}

.commission-tip__intro {
  margin: 0 0 12px;
  color: #595959;
  font-size: 13px;
  line-height: 1.6;
}

.commission-tip__intro .red {
  color: #f5222d;
  font-weight: 600;
}

.commission-tip__table {
  width: 100%;
  background: #fff;
}

.rule-tips {
  margin: 0 0 6px;
  color: #8c8c8c;
  font-size: 12px;
  line-height: 1.5;
}

.rule-item {
  margin: 0;
  color: #262626;
  font-size: 12px;
  line-height: 1.6;
}

.commission-form {
  flex: 1;
  padding: 24px 32px 16px;
}

.commission-form :deep(.el-form-item__label) {
  font-weight: 400;
}

.unit {
  margin-left: 8px;
  color: #595959;
}

.agent-setting__form {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
}

.agent-setting__footer {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  height: 60px;
  background: #fff;
  box-shadow: 0 -1px 4px rgb(0 0 0 / 10%);
}
</style>
