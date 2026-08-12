<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { Page } from "@vben/common-ui";
import { ElButton, ElMessage, ElTabPane, ElTabs } from "element-plus";

import { getAccessCodesApi } from "#/api/core/auth";
import {
  getPlatformShopConfigApi,
  savePlatformShopConfigApi,
  type PlatformShopConfig,
} from "#/api/core/platform-mall-setting";
import SettingsTabLayout from "#/components/settings/SettingsTabLayout.vue";
import ImageField from "#/components/shop/image-field.vue";

type TabName = "dashboard" | "mall" | "merchant" | "product" | "trade";

const activeTab = ref<TabName>("product");
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformShopConfig>(defaultForm());
const savedForm = ref<PlatformShopConfig>(defaultForm());
const refundReasonsText = computed({
  get: () => form.value.refund_reasons.join("\n"),
  set: (value: string) => {
    form.value.refund_reasons = value.split("\n");
  },
});

function defaultForm(): PlatformShopConfig {
  return {
    auto_parse_clipboard: true,
    arrival_notice_enabled: true,
    product_comment_enabled: true,
    auto_positive_review_enabled: true,
    default_copy_times: 8,
    order_auto_cancel_minutes: 15,
    order_auto_receive_days: 7,
    after_sale_days: 1,
    merchant_refund_auto_days: 1,
    refund_reasons: ["商品质量问题", "不想要了", "未收到货"],
    platform_rights_enabled: true,
    platform_rights_days: 1,
    merge_payment_enabled: true,
    merchant_apply_enabled: true,
    merchant_qualification_required: true,
    merchant_margin_badge_enabled: false,
    merchant_margin_badge_image: "",
    merchant_category_limit: 5,
    mall_show_stores: true,
    mall_recommend_enabled: true,
    mall_recommend_distance_enabled: true,
    mall_recommend_sort: "star",
    live_stream_auto_approve: false,
    live_product_auto_approve: false,
    hot_ranking_enabled: true,
    hot_ranking_category_level: 2,
    hot_ranking_refresh_hours: 24,
    mall_search_mode: "fuzzy",
    product_ranking_period: "month",
    product_ranking_metric: "sales_amount",
    shop_ranking_period: "month",
    shop_ranking_metric: "product_count",
    dashboard_display_name: "数据大屏",
  };
}

function cloneConfig(config: PlatformShopConfig): PlatformShopConfig {
  return JSON.parse(JSON.stringify(config)) as PlatformShopConfig;
}

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformShopConfigApi();
    savedForm.value = cloneConfig(data.config);
    form.value = cloneConfig(data.config);
  } finally {
    loading.value = false;
  }
}

function reset() {
  form.value = cloneConfig(savedForm.value);
}

async function save() {
  const displayName = form.value.dashboard_display_name.trim();
  const reasons = form.value.refund_reasons.map((item) => item.trim()).filter(Boolean);
  if (!displayName) {
    ElMessage.warning("显示名称不能为空");
    return;
  }
  saving.value = true;
  try {
    form.value = await savePlatformShopConfigApi({
      ...form.value,
      dashboard_display_name: displayName,
      merchant_margin_badge_image: form.value.merchant_margin_badge_image.trim(),
      refund_reasons: reasons,
    });
    savedForm.value = cloneConfig(form.value);
    ElMessage.success("商城设置已保存");
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes("setting.shop.manage");
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <template #tabs>
        <ElTabs v-model="activeTab">
          <ElTabPane label="商品设置" name="product" />
          <ElTabPane label="交易设置" name="trade" />
          <ElTabPane label="商户设置" name="merchant" />
          <ElTabPane label="商城设置" name="mall" />
          <ElTabPane label="数据大屏" name="dashboard" />
        </ElTabs>
      </template>

      <el-form :disabled="!canManage" label-width="220px" class="mall-setting__form">
        <template v-if="activeTab === 'product'">
          <el-form-item label="自动解析复制口令">
            <div class="mall-setting__control">
              <el-switch v-model="form.auto_parse_clipboard" />
              <p>开启后小程序和 APP 自动读取粘贴板</p>
            </div>
          </el-form-item>
          <el-form-item label="到货通知">
            <div class="mall-setting__control">
              <el-switch v-model="form.arrival_notice_enabled" />
              <p>开启后售罄商品到货后将通知用户</p>
            </div>
          </el-form-item>
          <el-form-item label="商品评论开启">
            <div class="mall-setting__control">
              <el-switch v-model="form.product_comment_enabled" />
              <p>开启后，移动端展示商品评论</p>
            </div>
          </el-form-item>
          <el-form-item label="开启自动好评">
            <div class="mall-setting__control">
              <el-switch v-model="form.auto_positive_review_enabled" />
              <p>开启后，会自动对 7 天内未评价的订单默认五星好评</p>
            </div>
          </el-form-item>
          <el-form-item label="默认赠送复制次数">
            <div class="mall-setting__control">
              <el-input-number v-model="form.default_copy_times" :min="0" :max="1000000" />
              <p>默认给商户赠送的商品采集次数</p>
            </div>
          </el-form-item>
        </template>

        <template v-else-if="activeTab === 'trade'">
          <el-form-item label="订单自动关闭时间(分钟)">
            <div class="mall-setting__control">
              <el-input-number v-model="form.order_auto_cancel_minutes" :min="0" :max="10080" />
              <p>订单提交后待支付时长，0 为默认 15 分钟</p>
            </div>
          </el-form-item>
          <el-form-item label="订单自动收货时间(天)">
            <div class="mall-setting__control">
              <el-input-number v-model="form.order_auto_receive_days" :min="0" :max="365" />
              <p>订单自动收货时间自发货日起计算</p>
            </div>
          </el-form-item>
          <el-form-item label="售后时长(天)">
            <div class="mall-setting__control">
              <el-input-number v-model="form.after_sale_days" :min="0" :max="365" />
              <p>用户确认收货后可申请售后的时长</p>
            </div>
          </el-form-item>
          <el-form-item label="商户自动处理退款订单期限（天）" required>
            <div class="mall-setting__control">
              <el-input-number v-model="form.merchant_refund_auto_days" :min="0" :max="365" />
              <p>申请退款的订单超过设置天数后自动处理</p>
            </div>
          </el-form-item>
          <el-form-item label="退款理由">
            <div class="mall-setting__control mall-setting__wide-control">
              <el-input
                v-model="refundReasonsText"
                :rows="5"
                maxlength="1000"
                placeholder="每行一个退款理由"
                type="textarea"
              />
              <p>设置常用退款理由，用户提交退款申请时可快速选择</p>
            </div>
          </el-form-item>
          <el-form-item label="开启/关闭平台维权">
            <el-switch v-model="form.platform_rights_enabled" />
          </el-form-item>
          <el-form-item label="平台售后维权（天）">
            <div class="mall-setting__control">
              <el-input-number v-model="form.platform_rights_days" :min="0" :max="365" />
              <p>商家拒绝退款后，用户可在此期限内申请平台介入</p>
            </div>
          </el-form-item>
          <el-form-item label="合单支付">
            <div class="mall-setting__control">
              <el-switch v-model="form.merge_payment_enabled" />
              <p>开启后，用户可一次支付多个商户订单</p>
            </div>
          </el-form-item>
        </template>

        <template v-else-if="activeTab === 'merchant'">
          <el-form-item label="开启店铺入驻">
            <div class="mall-setting__control">
              <el-switch v-model="form.merchant_apply_enabled" />
              <p>开启后，商城移动端和 PC 端可提交商户入驻申请</p>
            </div>
          </el-form-item>
          <el-form-item label="店铺资质是否必传">
            <div class="mall-setting__control">
              <el-switch v-model="form.merchant_qualification_required" />
              <p>开启后，入驻时必须上传店铺经营资质</p>
            </div>
          </el-form-item>
          <el-form-item label="店铺保证金标识展示">
            <div class="mall-setting__control">
              <el-switch v-model="form.merchant_margin_badge_enabled" />
              <p>开启后，已缴纳保证金的店铺展示标识</p>
            </div>
          </el-form-item>
          <el-form-item label="保证金标识">
            <div class="mall-setting__control">
              <ImageField
                v-model="form.merchant_margin_badge_image"
                :disabled="!canManage"
                default-library="system"
                :preview-size="96"
              />
              <p>从素材库选择保证金标识</p>
            </div>
          </el-form-item>
          <el-form-item label="店铺商品分类限制">
            <div class="mall-setting__control">
              <el-input-number v-model="form.merchant_category_limit" :min="0" :max="10000" />
              <p>商户入驻时可选择的经营商品分类最大数量，0 代表不限制</p>
            </div>
          </el-form-item>
        </template>

        <template v-else-if="activeTab === 'mall'">
          <el-form-item label="是否展示店铺">
            <div class="mall-setting__control">
              <el-switch v-model="form.mall_show_stores" />
              <p>开启移动端商城正常展示店铺信息；关闭移动端则隐藏店铺信息</p>
            </div>
          </el-form-item>
          <el-form-item label="展示为你推荐">
            <div class="mall-setting__control">
              <el-switch v-model="form.mall_recommend_enabled" />
              <p>开启：商城首页展示“为你推荐”商品列表</p>
            </div>
          </el-form-item>
          <el-form-item label="为你推荐距离显示">
            <div class="mall-setting__control">
              <el-switch v-model="form.mall_recommend_distance_enabled" />
              <p>开启店铺距离展示时，移动端需配置地图 Key 才能正常显示</p>
            </div>
          </el-form-item>
          <el-form-item label="为你推荐方式">
            <div class="mall-setting__control">
              <el-radio-group v-model="form.mall_recommend_sort">
                <el-radio value="default">默认推荐</el-radio>
                <el-radio value="star">星级推荐</el-radio>
                <el-radio value="created_at">创建时间</el-radio>
              </el-radio-group>
              <p>设置移动端商城首页“为你推荐”中商品展示顺序</p>
            </div>
          </el-form-item>
          <el-form-item label="开启直播免审核">
            <div class="mall-setting__control">
              <el-switch v-model="form.live_stream_auto_approve" />
              <p>开启后，店铺创建直播间无需平台审核</p>
            </div>
          </el-form-item>
          <el-form-item label="开启直播商品免审核">
            <div class="mall-setting__control">
              <el-switch v-model="form.live_product_auto_approve" />
              <p>开启后，店铺创建直播商品无需平台审核</p>
            </div>
          </el-form-item>
          <el-form-item label="热卖排行开关">
            <div class="mall-setting__control">
              <el-switch v-model="form.hot_ranking_enabled" />
              <p>关闭后，商城不展示热卖排行</p>
            </div>
          </el-form-item>
          <el-form-item label="热卖排行分类等级">
            <el-radio-group v-model="form.hot_ranking_category_level">
              <el-radio :value="1">一级分类</el-radio>
              <el-radio :value="2">二级分类</el-radio>
              <el-radio :value="3">三级分类</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="热卖排行更新时长（小时）">
            <div class="mall-setting__control">
              <el-input-number v-model="form.hot_ranking_refresh_hours" :min="1" :max="720" />
              <p>按设置间隔自动根据当前累计销量计算一次排行榜单</p>
            </div>
          </el-form-item>
          <el-form-item label="商城搜索方式">
            <div class="mall-setting__control">
              <el-radio-group v-model="form.mall_search_mode">
                <el-radio value="fuzzy">模糊搜索</el-radio>
                <el-radio value="split">分词搜索</el-radio>
              </el-radio-group>
              <p>模糊搜索匹配范围更广；分词搜索结果更精准</p>
            </div>
          </el-form-item>
        </template>

        <template v-else>
          <el-form-item label="商品排行数据时间">
            <div class="mall-setting__control">
              <el-radio-group v-model="form.product_ranking_period">
                <el-radio value="today">今日</el-radio>
                <el-radio value="week">本周</el-radio>
                <el-radio value="month">本月</el-radio>
              </el-radio-group>
              <p>配置数据大屏商品排行展示周期</p>
            </div>
          </el-form-item>
          <el-form-item label="商品排行数据类型">
            <div class="mall-setting__control">
              <el-radio-group v-model="form.product_ranking_metric">
                <el-radio value="sales_quantity">销售数量</el-radio>
                <el-radio value="sales_amount">销售金额</el-radio>
              </el-radio-group>
              <p>配置数据大屏商品排行类型</p>
            </div>
          </el-form-item>
          <el-form-item label="店铺排行数据时间">
            <div class="mall-setting__control">
              <el-radio-group v-model="form.shop_ranking_period">
                <el-radio value="today">当日</el-radio>
                <el-radio value="week">本周</el-radio>
                <el-radio value="month">本月</el-radio>
              </el-radio-group>
              <p>配置数据大屏店铺排行展示周期</p>
            </div>
          </el-form-item>
          <el-form-item label="店铺排行数据类型">
            <div class="mall-setting__control">
              <el-radio-group v-model="form.shop_ranking_metric">
                <el-radio value="sales_amount">销售金额</el-radio>
                <el-radio value="product_count">商品数量</el-radio>
              </el-radio-group>
              <p>配置数据大屏店铺排行类型</p>
            </div>
          </el-form-item>
          <el-form-item label="显示名称">
            <div class="mall-setting__control mall-setting__wide-control">
              <el-input v-model="form.dashboard_display_name" maxlength="64" />
              <p>设置数据大屏标题</p>
            </div>
          </el-form-item>
        </template>
      </el-form>

      <template #actions>
        <ElButton :disabled="loading" @click="reset">重置</ElButton>
        <ElButton v-if="canManage" :loading="saving" type="primary" @click="save"> 保存 </ElButton>
      </template>
    </SettingsTabLayout>
  </Page>
</template>

<style scoped>
.mall-setting__form {
  max-width: 980px;
  padding: 8px 0 24px;
}

.mall-setting__form :deep(.el-form-item) {
  margin-bottom: 32px;
}

.mall-setting__form :deep(.el-form-item__label) {
  color: var(--el-text-color-primary);
  font-weight: 500;
}

.mall-setting__control {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
  width: 100%;
}

.mall-setting__control p {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.6;
}

.mall-setting__control :deep(.el-input-number) {
  width: 280px;
}

.mall-setting__wide-control {
  max-width: 680px;
}

.mall-setting__wide-control :deep(.el-textarea),
.mall-setting__wide-control :deep(.el-input) {
  width: 100%;
}
</style>
