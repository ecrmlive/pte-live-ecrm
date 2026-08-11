<script setup lang="ts">
import type { PlatformCoupon } from '#/api/core/platform-promotion';
import type {
  PlatformUserSetupConfig,
  PlatformUserSetupCoupon,
  PlatformUserSetupField,
  PlatformUserSetupFieldType,
} from '#/api/core/platform-setting-ext';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElCheckbox,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElLink,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElPopover,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTabPane,
  ElTabs,
} from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import { listPlatformCouponsApi } from '#/api/core/platform-promotion';
import {
  getPlatformUserSetupConfigApi,
  savePlatformUserSetupConfigApi,
} from '#/api/core/platform-setting-ext';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';
import ImageField from '#/components/shop/image-field.vue';

type TabKey = 'basic' | 'login';

const FIELD_TYPE_OPTIONS: Array<{
  label: string;
  value: PlatformUserSetupFieldType;
}> = [
  { label: '文本', value: 'input' },
  { label: '数字', value: 'int' },
  { label: '日期', value: 'date' },
  { label: '单选项', value: 'radio' },
  { label: '身份证', value: 'id_card' },
  { label: '邮箱', value: 'email' },
  { label: '手机号', value: 'phone' },
  { label: '地址', value: 'address' },
];

const FIELD_TYPE_LABEL: Record<PlatformUserSetupFieldType, string> = {
  input: '文本',
  int: '数字',
  phone: '手机号',
  date: '日期',
  radio: '单选项',
  address: '地址',
  id_card: '身份证',
  email: '邮箱',
};

const activeTab = ref<TabKey>('basic');
const loading = ref(false);
const saving = ref(false);
const canRead = ref(false);
const canManage = ref(false);
const dragFromIndex = ref<number | null>(null);

const form = reactive<PlatformUserSetupConfig>(emptyConfig());
const fieldForm = reactive({
  field: '',
  title: '',
  type: 'input' as PlatformUserSetupFieldType,
  msg: '',
  contentText: '选项一\n选项二',
});
const couponKeyword = ref('');
const couponLoading = ref(false);
const couponRows = ref<PlatformCoupon[]>([]);
const selectedCouponIds = ref<number[]>([]);

const [FieldDrawer, fieldDrawerApi] = useVbenDrawer({
  class: 'w-[640px] max-w-[96vw]',
  placement: 'right',
  title: '新增信息',
  onConfirm: async () => {
    const err = validateFieldForm();
    if (err) {
      ElMessage.warning(err);
      throw new Error(err);
    }
    const nextId =
      form.fields.reduce((max, row) => Math.max(max, Number(row.id || 0)), 0) +
      1;
    const content =
      fieldForm.type === 'radio'
        ? fieldForm.contentText
            .split(/\n|,|，/)
            .map((s) => s.trim())
            .filter(Boolean)
        : undefined;
    form.fields.push({
      id: nextId,
      field: fieldForm.field.trim().toLowerCase(),
      title: fieldForm.title.trim(),
      type: fieldForm.type,
      msg: fieldForm.msg.trim(),
      content,
      is_used: 1,
      is_require: 0,
      is_show: 1,
      is_default: 0,
      sort: form.fields.length,
    });
  },
});

const [CouponModal, couponModalApi] = useVbenModal({
  title: '新增优惠券',
  class: 'w-[860px] max-w-[96vw]',
  onConfirm: () => {
    const picked = couponRows.value.filter((row) =>
      selectedCouponIds.value.includes(Number(row.coupon_id)),
    );
    if (!picked.length) {
      ElMessage.warning('请选择优惠券');
      throw new Error('请选择优惠券');
    }
    const exists = new Set(
      form.register_give_coupon.map((c) => Number(c.coupon_id)),
    );
    for (const row of picked) {
      const id = Number(row.coupon_id);
      if (exists.has(id)) continue;
      form.register_give_coupon.push(mapCoupon(row));
      exists.add(id);
    }
  },
});

const footerPrimaryText = computed(() =>
  activeTab.value === 'basic' ? '保存' : '提交',
);

function emptyConfig(): PlatformUserSetupConfig {
  return {
    user_default_avatar: '',
    fields: [],
    is_phone_login: 0,
    first_avatar_switch: 1,
    open_update_info: 1,
    wechat_phone_switch: 0,
    newcomer_status: 0,
    register_popup_pic: '',
    register_money_status: 0,
    register_give_money: 0,
    register_integral_status: 0,
    register_give_integral: 0,
    register_coupon_status: 0,
    register_give_coupon: [],
  };
}

function applyConfig(data: PlatformUserSetupConfig) {
  Object.assign(form, emptyConfig(), data, {
    fields: Array.isArray(data.fields)
      ? data.fields.map((row, idx) => ({
          ...row,
          is_used: row.is_used === 1 ? 1 : 0,
          is_require: row.is_require === 1 ? 1 : 0,
          is_show: row.is_show === 1 ? 1 : 0,
          is_default: row.is_default === 1 ? 1 : 0,
          sort: idx,
          content: Array.isArray(row.content) ? [...row.content] : undefined,
        }))
      : [],
    register_give_coupon: Array.isArray(data.register_give_coupon)
      ? data.register_give_coupon.map((c) => ({ ...c }))
      : [],
  });
}

function typeLabel(type: string) {
  return FIELD_TYPE_LABEL[type as PlatformUserSetupFieldType] || type;
}

function couponTypeLabel(row: PlatformUserSetupCoupon) {
  // 对齐 CRMEB / 平台券展示
  switch (Number(row.coupon_type)) {
    case 0:
      return '平台通用券';
    case 1:
      return '品类券';
    case 2:
      return '商品券';
    default:
      return '优惠券';
  }
}

function couponMinPriceText(row: PlatformUserSetupCoupon) {
  const n = Number(row.use_min_price || 0);
  return n > 0 ? n.toFixed(2) : '不限制';
}

function couponValidText(row: PlatformUserSetupCoupon) {
  if (Number(row.is_timeout) === 1 || Number(row.coupon_type) === 1) {
    const start = row.use_start_time || '-';
    const end = row.use_end_time || '-';
    return `${start}-${end}`;
  }
  return `${Number(row.coupon_time || 0)}天`;
}

function mapCoupon(row: PlatformCoupon): PlatformUserSetupCoupon {
  return {
    coupon_id: Number(row.coupon_id),
    title: row.title || `优惠券#${row.coupon_id}`,
    coupon_type: Number(row.coupon_type ?? 0),
    coupon_price: Number(row.coupon_price ?? 0),
    use_min_price: Number(row.use_min_price ?? 0),
    coupon_time: Number(row.coupon_time ?? 0),
    is_timeout: Number(row.is_timeout ?? 0),
    use_start_time: row.use_start_time || undefined,
    use_end_time: row.use_end_time || undefined,
  };
}

function onUsedChange(row: PlatformUserSetupField) {
  if (row.is_used !== 1) {
    row.is_require = 0;
    row.is_show = 0;
  }
}

function validateFieldForm(): string | null {
  const field = fieldForm.field.trim().toLowerCase();
  const title = fieldForm.title.trim();
  const msg = fieldForm.msg.trim();
  if (!/^[a-z][a-z0-9_]{0,63}$/.test(field)) {
    return '字段名须以小写字母开头，仅含小写字母/数字/下划线';
  }
  if (form.fields.some((r) => r.field === field)) {
    return '字段名已存在';
  }
  if (!title) return '请填写信息名称';
  if (!fieldForm.type) return '请选择信息格式';
  if (!msg) return '请填写提示文案';
  if (fieldForm.type === 'radio') {
    const opts = fieldForm.contentText
      .split(/\n|,|，/)
      .map((s) => s.trim())
      .filter(Boolean);
    if (opts.length < 2) return '单选至少两个选项（换行或逗号分隔）';
  }
  return null;
}

function openAddField() {
  if (!canManage.value) return;
  fieldForm.field = '';
  fieldForm.title = '';
  fieldForm.type = 'input';
  fieldForm.msg = '';
  fieldForm.contentText = '选项一\n选项二';
  fieldDrawerApi.open();
}

async function removeField(row: PlatformUserSetupField) {
  if (!canManage.value) return;
  if (row.is_default === 1) {
    ElMessage.warning('系统默认字段不能删除');
    return;
  }
  try {
    await ElMessageBox.confirm(`确认删除信息「${row.title}」？`, '删除确认', {
      type: 'warning',
    });
  } catch {
    return;
  }
  form.fields = form.fields.filter((item) => item.id !== row.id);
}

function removeCoupon(couponId: number) {
  if (!canManage.value) return;
  form.register_give_coupon = form.register_give_coupon.filter(
    (c) => Number(c.coupon_id) !== Number(couponId),
  );
}

async function openCouponPicker() {
  if (!canManage.value) return;
  selectedCouponIds.value = [];
  couponKeyword.value = '';
  couponModalApi.open();
  await loadCoupons();
}

async function loadCoupons() {
  couponLoading.value = true;
  try {
    const data = await listPlatformCouponsApi({
      page: 1,
      limit: 50,
      keyword: couponKeyword.value.trim() || undefined,
      status: 1,
    });
    couponRows.value = Array.isArray(data.list) ? data.list : [];
  } finally {
    couponLoading.value = false;
  }
}

function onCouponSelectionChange(rows: PlatformCoupon[]) {
  selectedCouponIds.value = rows.map((r) => Number(r.coupon_id));
}

function moveField(from: number, to: number) {
  if (from < 0 || to < 0 || from === to) return;
  if (from >= form.fields.length || to >= form.fields.length) return;
  const list = [...form.fields];
  const [moved] = list.splice(from, 1);
  if (!moved) return;
  list.splice(to, 0, moved);
  form.fields = list.map((row, idx) => ({ ...row, sort: idx }));
}

function onFieldDragStart(index: number, evt: DragEvent) {
  if (!canManage.value) {
    evt.preventDefault();
    return;
  }
  dragFromIndex.value = index;
  if (evt.dataTransfer) {
    evt.dataTransfer.effectAllowed = 'move';
    evt.dataTransfer.setData('text/plain', String(index));
  }
}

function onFieldDragOver(evt: DragEvent) {
  if (!canManage.value || dragFromIndex.value === null) return;
  evt.preventDefault();
  if (evt.dataTransfer) evt.dataTransfer.dropEffect = 'move';
}

function resolveFieldRowIndex(target: EventTarget | null): number {
  const el = target instanceof Element ? target : null;
  const tr = el?.closest?.('tbody tr') as HTMLTableRowElement | null;
  if (!tr?.parentElement) return -1;
  return Array.from(tr.parentElement.children).indexOf(tr);
}

function onFieldsTableDrop(evt: DragEvent) {
  if (!canManage.value) return;
  evt.preventDefault();
  const from = dragFromIndex.value;
  dragFromIndex.value = null;
  if (from === null) return;
  const to = resolveFieldRowIndex(evt.target);
  if (to < 0) return;
  moveField(from, to);
}

function onFieldDragEnd() {
  dragFromIndex.value = null;
}

function validate(): string | null {
  if (!form.fields.length) return '请至少保留一条用户信息字段';
  if (form.newcomer_status !== 1) return null;
  if (form.register_money_status === 1 && !(form.register_give_money > 0)) {
    return '请填写余额赠送金额';
  }
  if (
    form.register_integral_status === 1 &&
    !(form.register_give_integral > 0)
  ) {
    return '请填写积分赠送数值';
  }
  if (
    form.register_coupon_status === 1 &&
    form.register_give_coupon.length === 0
  ) {
    return '请添加赠送优惠券';
  }
  return null;
}

async function load() {
  if (!canRead.value) {
    ElMessage.warning('无权查看用户设置');
    return;
  }
  loading.value = true;
  try {
    const data = await getPlatformUserSetupConfigApi();
    applyConfig(data.config);
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!canManage.value) {
    ElMessage.warning('无权保存用户设置');
    return;
  }
  const err = validate();
  if (err) {
    ElMessage.warning(err);
    return;
  }
  saving.value = true;
  try {
    const saved = await savePlatformUserSetupConfigApi({
      ...form,
      user_default_avatar: form.user_default_avatar.trim(),
      register_popup_pic: form.register_popup_pic.trim(),
      fields: form.fields.map((row, idx) => ({
        ...row,
        field: row.field.trim().toLowerCase(),
        title: row.title.trim(),
        msg: row.msg.trim(),
        sort: idx,
        content:
          row.type === 'radio' && Array.isArray(row.content)
            ? row.content
            : undefined,
      })),
      register_give_coupon: form.register_give_coupon.map((c) => ({ ...c })),
    });
    applyConfig(saved);
    ElMessage.success(activeTab.value === 'basic' ? '保存成功' : '提交成功');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canRead.value =
    codes.includes('user.setup.read') || codes.includes('user.setup.manage');
  canManage.value = codes.includes('user.setup.manage');
  await load();
});

</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <template #tabs>
        <ElTabs v-model="activeTab">
          <ElTabPane label="基础信息" name="basic" />
          <ElTabPane label="登录注册" name="login" />
        </ElTabs>
      </template>

      <ElForm
        :disabled="!canManage"
        class="user-setup-form"
        label-width="160px"
      >
        <template v-if="activeTab === 'basic'">
          <ElFormItem label="用户默认头像：">
            <div>
              <ImageField
                v-model="form.user_default_avatar"
                :disabled="!canManage"
                default-library="system"
                shape="circle"
                :preview-size="120"
                hint="建议尺寸：120*120px"
              />
            </div>
          </ElFormItem>

          <ElFormItem label="用户信息设置：">
            <div
              class="user-setup-form__fields"
              @dragover="onFieldDragOver($event)"
              @drop="onFieldsTableDrop($event)"
            >
              <ElTable
                :data="form.fields"
                row-key="id"
                border
                class="user-setup-form__table"
              >
                <ElTableColumn label="#" width="56" align="center">
                  <template #default="{ $index }">
                    <span
                      class="drag-handle"
                      :class="{ 'is-disabled': !canManage }"
                      title="拖拽排序"
                      :draggable="canManage"
                      @dragstart="onFieldDragStart($index, $event)"
                      @dragend="onFieldDragEnd"
                    >
                      ⋮⋮
                    </span>
                  </template>
                </ElTableColumn>
                <ElTableColumn
                  prop="title"
                  label="信息名称"
                  min-width="120"
                  show-overflow-tooltip
                />
                <ElTableColumn label="使用" width="70" align="center">
                  <template #default="{ row }">
                    <ElCheckbox
                      v-model="row.is_used"
                      :true-value="1"
                      :false-value="0"
                      :disabled="!canManage"
                      @change="onUsedChange(row)"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="必填" width="70" align="center">
                  <template #default="{ row }">
                    <ElCheckbox
                      v-model="row.is_require"
                      :true-value="1"
                      :false-value="0"
                      :disabled="!canManage || row.is_used !== 1"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="用户端显示" width="100" align="center">
                  <template #default="{ row }">
                    <ElCheckbox
                      v-model="row.is_show"
                      :true-value="1"
                      :false-value="0"
                      :disabled="!canManage || row.is_used !== 1"
                    />
                  </template>
                </ElTableColumn>
                <ElTableColumn label="信息格式" min-width="100">
                  <template #default="{ row }">
                    {{ typeLabel(row.type) }}
                  </template>
                </ElTableColumn>
                <ElTableColumn
                  prop="msg"
                  label="提示信息"
                  min-width="140"
                  show-overflow-tooltip
                />
                <ElTableColumn label="操作" width="80" fixed="right">
                  <template #default="{ row }">
                    <ElButton
                      link
                      type="primary"
                      :disabled="!canManage || row.is_default === 1"
                      @click="removeField(row)"
                    >
                      删除
                    </ElButton>
                  </template>
                </ElTableColumn>
              </ElTable>
              <div class="user-setup-form__add-field">
                <ElButton
                  :disabled="!canManage"
                  @click="openAddField"
                >
                  新增信息
                </ElButton>
              </div>
              <div class="user-setup-form__notes">
                <p>
                  1.开启使用后，后台添加用户时可填写此信息；开启必填后，后台添加用户时此信息必须填写；开启用户端展示后，在商城用户个人信息中展示
                </p>
                <p>
                  2.自定义添加日期和单选格式的字段，暂不支持用户列表搜索，如业务需要建议进一步开发；其它字段均支持用户列表搜索
                </p>
              </div>
            </div>
          </ElFormItem>
        </template>

        <template v-else>
          <div class="user-setup-form__section-bar">登录设置</div>
          <ElFormItem label="强制手机号绑定：">
            <div>
              <ElSwitch
                v-model="form.is_phone_login"
                :active-value="1"
                :inactive-value="0"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="form-tip">开启，商城登录时强制手机号登录绑定</div>
            </div>
          </ElFormItem>
          <ElFormItem label="获取头像昵称：">
            <div>
              <ElSwitch
                v-model="form.first_avatar_switch"
                :active-value="1"
                :inactive-value="0"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="form-tip">
                开启，小程序首次登录弹出获取头像昵称弹窗；关闭，则不弹出
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="修改头像和昵称：">
            <div>
              <ElSwitch
                v-model="form.open_update_info"
                :active-value="1"
                :inactive-value="0"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="form-tip">
                开启，用户可自主修改头像或昵称；关闭，则不能修改
              </div>
            </div>
          </ElFormItem>
          <ElFormItem label="手机号快速验证组件：">
            <div>
              <ElSwitch
                v-model="form.wechat_phone_switch"
                :active-value="1"
                :inactive-value="0"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="form-tip">
                开启，使用微信手机号快速验证组件；关闭，则不使用
              </div>
            </div>
          </ElFormItem>

          <div class="user-setup-form__section-bar">注册有礼</div>
          <ElFormItem label="注册有礼启用：">
            <div>
              <ElSwitch
                v-model="form.newcomer_status"
                :active-value="1"
                :inactive-value="0"
                inline-prompt
                active-text="开启"
                inactive-text="关闭"
                :width="56"
              />
              <div class="form-tip">新用户注册后，会给用户赠送礼品</div>
            </div>
          </ElFormItem>
          <template v-if="form.newcomer_status === 1">
            <ElFormItem label="注册有礼弹窗：">
              <div>
                <ImageField
                  v-model="form.register_popup_pic"
                  :disabled="!canManage"
                  default-library="system"
                  :preview-size="135"
                />
                <div class="form-tip">
                  建议尺寸：270px*426px，如不上传此图，则默认使用系统样式
                  <ElPopover
                    placement="right"
                    trigger="click"
                    :width="460"
                    popper-class="register-gift-example-popper"
                  >
                    <template #reference>
                      <ElLink
                        type="primary"
                        :underline="false"
                        class="example-link-inline"
                      >
                        查看示例
                      </ElLink>
                    </template>
                    <div class="example-phones">
                      <div class="example-phone">
                        <div class="example-phone__label">默认样式</div>
                        <div class="example-phone__frame example-phone__frame--gold">
                          <div class="example-phone__notch" />
                          <div class="example-phone__body">
                            <div class="example-phone__gift example-phone__gift--gold">
                              <div class="example-phone__badge">礼</div>
                              <div class="example-phone__rewards">
                                <div>100.00 余额</div>
                                <div>50 积分</div>
                              </div>
                              <div class="example-phone__btn">立即领取</div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="example-phone">
                        <div class="example-phone__label">上传素材样式</div>
                        <div class="example-phone__frame example-phone__frame--purple">
                          <div class="example-phone__notch" />
                          <div class="example-phone__body">
                            <div class="example-phone__gift example-phone__gift--purple">
                              <div class="example-phone__art" />
                              <div class="example-phone__rewards example-phone__rewards--light">
                                <div>100.00 余额</div>
                                <div>50 积分</div>
                              </div>
                              <div class="example-phone__btn">立即领取</div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </ElPopover>
                </div>
              </div>
            </ElFormItem>
            <ElFormItem label="赠送余额：">
              <div>
                <ElSwitch
                  v-model="form.register_money_status"
                  :active-value="1"
                  :inactive-value="0"
                  inline-prompt
                  active-text="开启"
                  inactive-text="关闭"
                  :width="56"
                />
                <div
                  v-if="form.register_money_status === 1"
                  class="user-setup-form__amount"
                >
                  <ElInputNumber
                    v-model="form.register_give_money"
                    :min="0"
                    :max="99999"
                    :precision="2"
                    controls-position="right"
                    class="user-setup-form__amount-input"
                  />
                </div>
                <div class="form-tip">
                  新用户注册即奖励充值余额，大于或等于0（0为不赠送）
                </div>
              </div>
            </ElFormItem>
            <ElFormItem label="赠送积分：">
              <div>
                <ElSwitch
                  v-model="form.register_integral_status"
                  :active-value="1"
                  :inactive-value="0"
                  inline-prompt
                  active-text="开启"
                  inactive-text="关闭"
                  :width="56"
                />
                <div
                  v-if="form.register_integral_status === 1"
                  class="user-setup-form__amount"
                >
                  <ElInputNumber
                    v-model="form.register_give_integral"
                    :min="0"
                    :max="999999"
                    :precision="0"
                    controls-position="right"
                    class="user-setup-form__amount-input"
                  />
                </div>
                <div class="form-tip">
                  新用户注册即奖励积分，大于或等于0（0为不赠送）
                </div>
              </div>
            </ElFormItem>
            <ElFormItem label="赠送优惠券：">
              <div class="user-setup-form__coupon">
                <ElSwitch
                  v-model="form.register_coupon_status"
                  :active-value="1"
                  :inactive-value="0"
                  inline-prompt
                  active-text="开启"
                  inactive-text="关闭"
                  :width="56"
                />
                <div class="form-tip">新用户注册后即赠送优惠券</div>
                <template v-if="form.register_coupon_status === 1">
                  <ElTable
                    :data="form.register_give_coupon"
                    border
                    class="user-setup-form__table mt-3"
                  >
                    <ElTableColumn
                      prop="title"
                      label="优惠券名称"
                      min-width="140"
                      show-overflow-tooltip
                    />
                    <ElTableColumn label="优惠券类型" min-width="110">
                      <template #default="{ row }">
                        {{ couponTypeLabel(row) }}
                      </template>
                    </ElTableColumn>
                    <ElTableColumn label="优惠券面值" width="100">
                      <template #default="{ row }">
                        {{ Number(row.coupon_price || 0).toFixed(2) }}
                      </template>
                    </ElTableColumn>
                    <ElTableColumn label="最低消费额" width="110">
                      <template #default="{ row }">
                        {{ couponMinPriceText(row) }}
                      </template>
                    </ElTableColumn>
                    <ElTableColumn label="有效期限" min-width="160">
                      <template #default="{ row }">
                        {{ couponValidText(row) }}
                      </template>
                    </ElTableColumn>
                    <ElTableColumn label="操作" width="80" fixed="right">
                      <template #default="{ row }">
                        <ElButton
                          link
                          type="primary"
                          :disabled="!canManage"
                          @click="removeCoupon(row.coupon_id)"
                        >
                          删除
                        </ElButton>
                      </template>
                    </ElTableColumn>
                  </ElTable>
                  <ElButton
                    class="mt-3"
                    link
                    type="primary"
                    :disabled="!canManage"
                    @click="openCouponPicker"
                  >
                    +新增优惠券
                  </ElButton>
                </template>
              </div>
            </ElFormItem>
          </template>
        </template>
      </ElForm>

      <template #actions>
        <ElButton
          type="primary"
          :disabled="!canManage"
          :loading="saving"
          @click="save"
        >
          {{ footerPrimaryText }}
        </ElButton>
      </template>
    </SettingsTabLayout>

    <FieldDrawer>
      <ElForm label-width="108px" class="field-drawer-form">
        <ElFormItem label="字段名" required>
          <ElInput
            v-model="fieldForm.field"
            maxlength="64"
            placeholder="以英文开头的字母、数字、下划线组合，用于代码中筛选信息名称，在后台前端不展示"
          />
        </ElFormItem>
        <ElFormItem label="信息名称" required>
          <ElInput
            v-model="fieldForm.title"
            maxlength="64"
            placeholder="请输入信息名称"
          />
        </ElFormItem>
        <ElFormItem label="信息格式" required>
          <ElSelect v-model="fieldForm.type" class="w-full" placeholder="请选择">
            <ElOption
              v-for="opt in FIELD_TYPE_OPTIONS"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem
          v-if="fieldForm.type === 'radio'"
          label="配置内容"
          required
        >
          <ElInput
            v-model="fieldForm.contentText"
            type="textarea"
            :rows="4"
            placeholder="请输入配置内容（每行一个选项，至少两项）"
          />
        </ElFormItem>
        <ElFormItem label="提示文案" required>
          <ElInput
            v-model="fieldForm.msg"
            maxlength="255"
            placeholder="请输入提示文案"
          />
        </ElFormItem>
      </ElForm>
    </FieldDrawer>

    <CouponModal>
      <div class="coupon-picker">
        <div class="coupon-picker__toolbar">
          <ElInput
            v-model="couponKeyword"
            clearable
            placeholder="搜索优惠券名称"
            class="coupon-picker__input"
            @keyup.enter="loadCoupons"
          />
          <ElButton type="primary" @click="loadCoupons">搜索</ElButton>
        </div>
        <ElTable
          v-loading="couponLoading"
          :data="couponRows"
          row-key="coupon_id"
          border
          max-height="420"
          @selection-change="onCouponSelectionChange"
        >
          <ElTableColumn type="selection" width="48" />
          <ElTableColumn prop="title" label="名称" min-width="160" />
          <ElTableColumn label="面值" width="90">
            <template #default="{ row }">
              {{ Number(row.coupon_price || 0).toFixed(2) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="最低消费" width="100">
            <template #default="{ row }">
              {{ Number(row.use_min_price || 0).toFixed(2) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="有效期" width="120">
            <template #default="{ row }">
              {{
                Number(row.coupon_type || 0) === 1
                  ? '时间段'
                  : `${row.coupon_time || 0}天`
              }}
            </template>
          </ElTableColumn>
        </ElTable>
      </div>
    </CouponModal>
  </Page>
</template>

<style scoped>
:deep(.settings-tab-layout__card) {
  max-width: 1200px;
}

.user-setup-form {
  width: 100%;
  max-width: none;
  padding: 8px 0 0;
}

.user-setup-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.user-setup-form :deep(.el-form-item__label) {
  align-items: flex-start;
  padding-top: 6px;
  color: hsl(var(--foreground));
  font-weight: 400;
}

.user-setup-form__fields {
  width: 100%;
}

.user-setup-form__table {
  width: 100%;
}

.user-setup-form__add-field {
  margin-top: 12px;
}

.user-setup-form__notes {
  margin-top: 12px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.7;
}

.user-setup-form__notes p {
  margin: 0;
}

.user-setup-form__section-bar {
  padding: 8px 12px;
  margin: 4px 0 18px;
  color: #1d4ed8;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.4;
  background: #eff6ff;
  border-left: 3px solid #3b82f6;
}

.user-setup-form__amount {
  margin-top: 10px;
}

.user-setup-form__amount-input {
  width: 200px;
}

.user-setup-form__coupon {
  width: 100%;
}

.form-tip {
  margin-top: 6px;
  max-width: 720px;
  color: #909399;
  font-size: 12px;
  line-height: 1.6;
}

.drag-handle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  color: hsl(var(--muted-foreground));
  cursor: grab;
  user-select: none;
  letter-spacing: -2px;
}

.drag-handle.is-disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.example-link-inline {
  margin-left: 4px;
  font-size: 12px;
  vertical-align: baseline;
}

.field-drawer-form {
  padding: 8px 4px 0;
}

.coupon-picker__toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 12px;
}

.coupon-picker__input {
  max-width: 280px;
}
</style>

<!-- Popover 内容挂到 body，需非 scoped 才能命中示例手机框 -->
<style>
.register-gift-example-popper .example-phones {
  display: flex;
  gap: 20px;
  justify-content: center;
  padding: 4px 2px;
}

.register-gift-example-popper .example-phone {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.register-gift-example-popper .example-phone__label {
  color: #303133;
  font-size: 13px;
  font-weight: 500;
}

.register-gift-example-popper .example-phone__frame {
  position: relative;
  width: 168px;
  height: 300px;
  padding: 10px 8px 12px;
  overflow: hidden;
  background: #111;
  border: 2px solid #2a2a2a;
  border-radius: 22px;
  box-shadow: 0 8px 20px rgb(0 0 0 / 12%);
}

.register-gift-example-popper .example-phone__notch {
  width: 54px;
  height: 8px;
  margin: 0 auto 8px;
  background: #222;
  border-radius: 8px;
}

.register-gift-example-popper .example-phone__body {
  display: flex;
  align-items: center;
  justify-content: center;
  height: calc(100% - 16px);
  overflow: hidden;
  border-radius: 12px;
}

.register-gift-example-popper .example-phone__frame--gold .example-phone__body {
  background: linear-gradient(165deg, #3a1a12 0%, #1a0c08 100%);
}

.register-gift-example-popper
  .example-phone__frame--purple
  .example-phone__body {
  background: linear-gradient(165deg, #2b1550 0%, #120820 100%);
}

.register-gift-example-popper .example-phone__gift {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 118px;
  padding: 14px 10px 12px;
  border-radius: 10px;
}

.register-gift-example-popper .example-phone__gift--gold {
  background: linear-gradient(180deg, #f7d48a 0%, #e8a23a 45%, #c47820 100%);
  box-shadow: inset 0 0 0 1px rgb(255 255 255 / 35%);
}

.register-gift-example-popper .example-phone__gift--purple {
  background: linear-gradient(180deg, #8b5cf6 0%, #6d28d9 55%, #4c1d95 100%);
  box-shadow: inset 0 0 0 1px rgb(255 255 255 / 25%);
}

.register-gift-example-popper .example-phone__badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  margin-bottom: 10px;
  color: #fff7e6;
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(160deg, #ef4444 0%, #b91c1c 100%);
  border-radius: 10px;
  box-shadow: 0 4px 10px rgb(185 28 28 / 35%);
}

.register-gift-example-popper .example-phone__art {
  width: 56px;
  height: 56px;
  margin-bottom: 10px;
  background:
    radial-gradient(circle at 30% 30%, #f9a8d4 0%, transparent 45%),
    radial-gradient(circle at 70% 60%, #c4b5fd 0%, transparent 50%),
    linear-gradient(140deg, #f472b6 0%, #7c3aed 100%);
  border-radius: 50%;
  box-shadow: 0 4px 12px rgb(124 58 237 / 40%);
}

.register-gift-example-popper .example-phone__rewards {
  margin-bottom: 12px;
  color: #7c2d12;
  font-size: 12px;
  font-weight: 600;
  line-height: 1.55;
  text-align: center;
}

.register-gift-example-popper .example-phone__rewards--light {
  color: #faf5ff;
}

.register-gift-example-popper .example-phone__btn {
  width: 100%;
  padding: 6px 0;
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  text-align: center;
  background: linear-gradient(90deg, #f97316 0%, #ef4444 100%);
  border-radius: 999px;
}

.register-gift-example-popper
  .example-phone__gift--purple
  .example-phone__btn {
  background: linear-gradient(90deg, #ec4899 0%, #a855f7 100%);
}
</style>
