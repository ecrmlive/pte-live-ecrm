<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus';
import type { ProductChooseItem } from '#/api/core/product';
import type { LiveAnchorListItem, LiveRoomForm, LiveStreamInfo } from '#/api/core/live';

import { Plus } from '@element-plus/icons-vue';
import {
  ElButton,
  ElCol,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElIcon,
  ElInput,
  ElInputNumber,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElRow,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
  ElMessage,
} from 'element-plus';
import { computed, ref, watch } from 'vue';

import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

import ImageField from '#/components/shop/image-field.vue';
import NativeSectionTitle from '#/components/shop/native-section-title.vue';

defineOptions({ name: 'LiveRoomFormFields' });

const props = withDefaults(
  defineProps<{
    anchorOptions?: LiveAnchorListItem[];
    form: LiveRoomForm;
    isEdit?: boolean;
    selectedProducts?: ProductChooseItem[];
    showStream?: boolean;
    streamInfo?: LiveStreamInfo;
    streamLoading?: boolean;
  }>(),
  {
    anchorOptions: () => [],
    isEdit: false,
    selectedProducts: () => [],
    showStream: false,
    streamInfo: () => ({}),
    streamLoading: false,
  },
);

const emit = defineEmits<{
  anchorChange: [anchorId: number | undefined];
  openAnchorAdd: [];
  openProductAdd: [];
  openProductPick: [];
  pickVod: [];
  refreshStream: [];
  removeProduct: [productId: number];
}>();

const formRef = ref<FormInstance>();

const recordVideoPreviewUrl = computed(() => {
  const vod = String(props.form.record_vod_media_url || '').trim();
  if (vod && /^https?:\/\//i.test(vod)) return vod;
  if (vod) return resolveCosMediaUrl(vod);
  return resolveCosMediaUrl(props.form.record_video_path || '');
});
const recordVideoName = computed(() => {
  const vod = String(props.form.record_vod_media_url || '').trim();
  if (vod) {
    return props.form.record_vod_file_id
      ? `点播 ${props.form.record_vod_file_id}`
      : vod.split('/').pop() || vod;
  }
  const path = String(props.form.record_video_path || '').trim();
  if (!path) return '';
  return path.split('/').pop() || path;
});

const START_TIME_BUFFER_MS = 5 * 60 * 1000;

function getMinStartTime() {
  return new Date(Date.now() + START_TIME_BUFFER_MS);
}

function parseDateTime(str?: string) {
  if (!str) return null;
  const d = new Date(String(str).replace(/-/g, '/'));
  return Number.isNaN(d.getTime()) ? null : d;
}

function formatDateTimeValue(date: Date) {
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

const startTimeShortcuts = [
  {
    text: '此刻',
    value: () => getMinStartTime(),
  },
];

function isSameDay(a: Date, b: Date) {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

function range(end: number) {
  return Array.from({ length: end }, (_, i) => i);
}

function validateStartTime(_rule: unknown, value: string, callback: (err?: Error) => void) {
  if (props.isEdit) {
    callback();
    return;
  }
  if (!value) {
    callback();
    return;
  }
  const selected = parseDateTime(value);
  if (!selected) {
    callback(new Error('时间格式无效'));
    return;
  }
  if (selected.getTime() < getMinStartTime().getTime()) {
    callback(new Error('计划开播时间不能早于当前时间'));
    return;
  }
  callback();
}

const rules = computed<FormRules>(() => ({
  anchor_id: [{ message: '请选择主播', required: !props.isEdit, trigger: 'change' }],
  background_img: [{ message: '请上传直播背景图', required: true, trigger: 'change' }],
  cover_img: [{ message: '请上传直播封面图', required: true, trigger: 'change' }],
  name: [
    { message: '请输入直播间名称', required: true, trigger: 'blur' },
    { max: 30, message: '长度在 2 到 30 个字符', min: 2, trigger: 'blur' },
  ],
  start_time: [
    { message: '请选择计划开播时间', required: true, trigger: 'change' },
    { trigger: 'change', validator: validateStartTime },
  ],
}));

function streamStatusTag(v?: number) {
  const map: Record<number, { text: string; type: 'danger' | 'info' | 'success' | 'warning' }> = {
    0: { text: '未推流', type: 'info' },
    1: { text: '推流中', type: 'success' },
    2: { text: '断流', type: 'warning' },
    3: { text: '禁推', type: 'danger' },
    4: { text: '异常', type: 'danger' },
  };
  return map[v ?? 0] ?? map[0]!;
}

async function copyPushUrl() {
  const url = props.streamInfo?.push_url;
  if (!url) return;
  try {
    await navigator.clipboard.writeText(url);
    ElMessage.success('推流地址已复制');
  } catch {
    ElMessage.warning('复制失败，请手动复制');
  }
}

function onStartTimeChange(value: string) {
  if (!value || props.isEdit) return;
  const selected = parseDateTime(value);
  const min = getMinStartTime();
  if (selected && selected.getTime() < min.getTime()) {
    props.form.start_time = formatDateTimeValue(min);
  }
}

function pickerDay(comparingDate?: Date) {
  if (comparingDate) {
    return comparingDate instanceof Date ? comparingDate : new Date(comparingDate);
  }
  return parseDateTime(props.form.start_time) || new Date();
}

function disabledStartDate(time: Date) {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  return time.getTime() < today.getTime();
}

function disabledStartHours(_role: string, comparingDate?: Date) {
  const min = getMinStartTime();
  const day = pickerDay(comparingDate);
  if (!isSameDay(day, min)) return [];
  return range(min.getHours());
}

function disabledStartMinutes(hour: number, _role: string, comparingDate?: Date) {
  const min = getMinStartTime();
  const day = pickerDay(comparingDate);
  if (!isSameDay(day, min) || hour !== min.getHours()) return [];
  return range(min.getMinutes());
}

function disabledStartSeconds(hour: number, minute: number, _role: string, comparingDate?: Date) {
  const min = getMinStartTime();
  const day = pickerDay(comparingDate);
  if (!isSameDay(day, min) || hour !== min.getHours() || minute !== min.getMinutes()) return [];
  return range(min.getSeconds());
}

function onAnchorChange(anchorId?: number) {
  emit('anchorChange', anchorId);
}

function formatProductPrice(v?: number | string) {
  const n = Number(v);
  return Number.isFinite(n) ? n.toFixed(2) : '0.00';
}

function resolveOriginalPrice(row: ProductChooseItem) {
  const sku = (row as ProductChooseItem & { product_sku?: { line_price?: number | string } })
    .product_sku;
  if (sku?.line_price != null && Number(sku.line_price) > 0) {
    return sku.line_price;
  }
  if (row.line_price != null && Number(row.line_price) > 0) {
    return row.line_price;
  }
  return row.product_price ?? 0;
}

function resolveActivityPrice(row: ProductChooseItem) {
  const sku = (row as ProductChooseItem & { product_sku?: { product_price?: number | string } })
    .product_sku;
  if (sku?.product_price != null) {
    return sku.product_price;
  }
  return row.product_price ?? 0;
}

function resolveProductImage(row: ProductChooseItem) {
  if (row.product_image) return resolveCosMediaUrl(row.product_image);
  const path = row.image?.[0]?.file_path;
  return path ? resolveCosMediaUrl(path) : '';
}

function clearRecordVideo() {
  props.form.record_video_path = '';
  props.form.record_vod_file_id = '';
  props.form.record_vod_media_url = '';
}

function resetFields() {
  formRef.value?.resetFields();
  formRef.value?.clearValidate();
}

async function validate() {
  if (!formRef.value) return false;
  try {
    await formRef.value.validate();
    return true;
  } catch {
    return false;
  }
}

watch(
  () => props.form.cover_img,
  () => {
    formRef.value?.validateField('cover_img').catch(() => undefined);
  },
);

watch(
  () => props.form.background_img,
  () => {
    formRef.value?.validateField('background_img').catch(() => undefined);
  },
);

defineExpose({ resetFields, validate });
</script>

<template>
  <el-form
    ref="formRef"
    class="live-room-form-fields"
    label-width="120px"
    :model="form"
    :rules="rules"
  >
    <NativeSectionTitle title="基本信息" />
    <el-form-item label="直播间名称" prop="name">
      <el-input v-model="form.name" class="max-w460" placeholder="请输入直播间名称" />
    </el-form-item>
    <el-form-item label="计划开播时间" prop="start_time">
      <el-date-picker
        v-model="form.start_time"
        class="max-w460"
        popper-class="native-modal-popper"
        :disabled-date="disabledStartDate"
        :disabled-hours="disabledStartHours"
        :disabled-minutes="disabledStartMinutes"
        :disabled-seconds="disabledStartSeconds"
        format="YYYY-MM-DD HH:mm:ss"
        placeholder="选择时间"
        :shortcuts="startTimeShortcuts"
        type="datetime"
        value-format="YYYY-MM-DD HH:mm:ss"
        @change="onStartTimeChange"
      />
    </el-form-item>
    <el-form-item v-if="isEdit" label="结束时间">
      <el-date-picker
        v-model="form.end_time"
        class="max-w460"
        popper-class="native-modal-popper"
        format="YYYY-MM-DD HH:mm:ss"
        placeholder="选择结束时间"
        type="datetime"
        value-format="YYYY-MM-DD HH:mm:ss"
      />
    </el-form-item>
    <el-form-item label="类型" prop="room_type">
      <el-radio-group v-model="form.room_type" class="live-room-segment">
        <el-radio :value="1">直播</el-radio>
        <el-radio :value="2">录播</el-radio>
      </el-radio-group>
    </el-form-item>
    <el-form-item label="样式" prop="stream_orientation">
      <el-radio-group v-model="form.stream_orientation" class="live-room-segment">
        <el-radio :value="1">横屏</el-radio>
        <el-radio :value="2">竖屏</el-radio>
      </el-radio-group>
    </el-form-item>
    <el-form-item label="是否展示" prop="is_visible">
      <el-radio-group v-model="form.is_visible" class="live-room-segment">
        <el-radio :value="1">展示</el-radio>
        <el-radio :value="0">不展示</el-radio>
      </el-radio-group>
    </el-form-item>
    <el-form-item
      v-if="form.room_type === 2"
      label="录播视频"
      prop="record_vod_media_url"
    >
      <div class="record-video-field">
        <div class="record-video-actions">
          <el-button plain type="primary" @click="emit('pickVod')">
            {{ form.record_vod_media_url || form.record_video_path ? '重新选择' : '选择录播视频' }}
          </el-button>
          <el-button
            v-if="form.record_vod_media_url || form.record_video_path"
            link
            type="danger"
            @click="clearRecordVideo"
          >
            清除
          </el-button>
        </div>
        <div v-if="recordVideoPreviewUrl" class="record-video-preview">
          <video
            class="record-video-player"
            controls
            preload="metadata"
            :src="recordVideoPreviewUrl"
          />
          <div class="record-video-path">{{ recordVideoName }}</div>
        </div>
        <div v-else class="record-video-hint">从素材库「云点播视频」中选择录播文件</div>
      </div>
    </el-form-item>

    <NativeSectionTitle title="主播" />
    <el-form-item label="主播" prop="anchor_id">
      <div class="anchor-field">
        <el-select
          v-model="form.anchor_id"
          class="anchor-select"
          popper-class="native-modal-popper"
          :clearable="!isEdit"
          :disabled="isEdit"
          placeholder="请选择主播"
          @change="onAnchorChange"
        >
          <el-option
            v-for="item in anchorOptions"
            :key="item.anchor_id"
            :label="`${item.nick_name} (${item.phone || item.account})`"
            :value="item.anchor_id"
          />
        </el-select>
        <el-button v-if="!isEdit" type="primary" @click="emit('openAnchorAdd')">
          <el-icon class="mr-1"><Plus /></el-icon>
          新增
        </el-button>
      </div>
      <p v-if="isEdit" class="anchor-lock-tip">创建后不可更换主播</p>
    </el-form-item>

    <template v-if="!isEdit">
      <NativeSectionTitle title="商品（可选）" />
      <el-form-item label="关联商品">
        <div class="product-field">
          <div class="product-actions">
            <el-button type="primary" @click="emit('openProductPick')">选择商品</el-button>
            <el-button @click="emit('openProductAdd')">新增商品</el-button>
          </div>
          <div v-if="!selectedProducts.length" class="product-hint">
            可不选；创建后也可在列表「商品管理」中添加
          </div>
          <el-table
            v-else
            border
            class="selected-product-table"
            :data="selectedProducts"
            size="small"
          >
            <el-table-column label="图" width="56">
              <template #default="{ row }">
                <img
                  v-if="resolveProductImage(row)"
                  :alt="row.product_name"
                  height="32"
                  :src="resolveProductImage(row)"
                  width="32"
                />
              </template>
            </el-table-column>
            <el-table-column label="名称" min-width="140" prop="product_name" show-overflow-tooltip />
            <el-table-column label="库存" prop="product_stock" width="70" />
            <el-table-column align="right" label="划线价" width="84">
              <template #default="{ row }">¥{{ formatProductPrice(resolveOriginalPrice(row)) }}</template>
            </el-table-column>
            <el-table-column align="right" label="价格" width="84">
              <template #default="{ row }">¥{{ formatProductPrice(resolveActivityPrice(row)) }}</template>
            </el-table-column>
            <el-table-column align="center" label="操作" width="64">
              <template #default="{ row }">
                <el-button link size="small" type="primary" @click="emit('removeProduct', row.product_id)">
                  移除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-form-item>
    </template>

    <NativeSectionTitle title="图片素材" />
    <el-row :gutter="32" class="image-upload-row">
      <el-col :span="12">
        <el-form-item label="直播封面图" prop="cover_img">
          <ImageField
            v-model="form.cover_img"
            hint="建议正方形，不超过 512KB"
            :preview-size="120"
          />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="直播背景图" prop="background_img">
          <ImageField
            v-model="form.background_img"
            hint="建议适配手机竖屏，不超过 1MB"
            :preview-size="120"
          />
        </el-form-item>
      </el-col>
    </el-row>

    <NativeSectionTitle title="直播设置" />
    <el-form-item label="观看密码">
      <el-input v-model="form.watch_password" class="max-w460" placeholder="留空则无需密码" />
    </el-form-item>
    <el-form-item label="火力值">
      <el-input-number v-model="form.fire_value" :max="99999999" :min="0" />
    </el-form-item>
    <el-form-item label="分享介绍">
      <el-input
        v-model="form.share_intro"
        class="max-w460"
        placeholder="请输入分享介绍"
        :rows="3"
        type="textarea"
      />
    </el-form-item>
    <el-form-item label="公告">
      <el-input
        v-model="form.system_notice"
        class="max-w460"
        placeholder="请输入公告内容"
        :rows="3"
        type="textarea"
      />
    </el-form-item>

    <div v-if="isEdit && showStream && form.room_type === 1 && streamInfo?.push_url">
      <NativeSectionTitle title="推流信息" />
    </div>
    <el-form-item
      v-if="isEdit && showStream && form.room_type === 1 && streamInfo?.push_url"
      label="推流地址"
    >
      <div class="stream-field max-w460">
        <el-input :model-value="streamInfo.push_url" readonly />
        <div class="stream-field__meta">
          <span>流状态：</span>
          <el-tag :type="streamStatusTag(streamInfo.stream_status).type" size="small">
            {{ streamStatusTag(streamInfo.stream_status).text }}
          </el-tag>
          <el-button
            link
            :loading="streamLoading"
            size="small"
            type="primary"
            @click="emit('refreshStream')"
          >
            刷新地址
          </el-button>
          <el-button size="small" @click="copyPushUrl">复制</el-button>
        </div>
      </div>
    </el-form-item>

    <NativeSectionTitle title="功能开关" />
    <el-form-item label-width="0">
      <div class="switch-row">
        <div class="switch-item">
          <span class="switch-label">回放</span>
          <el-switch v-model="form.enable_replay" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">礼物</span>
          <el-switch v-model="form.enable_gift" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">在线人数</span>
          <el-switch v-model="form.show_online_count" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">累计人数</span>
          <el-switch v-model="form.show_total_count" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">火力值</span>
          <el-switch v-model="form.show_heat" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">连麦</span>
          <el-switch v-model="form.enable_linkmic" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">允许聊天</span>
          <el-switch v-model="form.allow_chat" :active-value="1" :inactive-value="0" />
        </div>
        <div class="switch-item">
          <span class="switch-label">显示首页</span>
          <el-switch v-model="form.show_home" :active-value="1" :inactive-value="0" />
        </div>
      </div>
      <div class="field-hint mt10">
        关闭后，观众端直播间内不显示「首页」入口（适用于分享扫码等场景）
      </div>
    </el-form-item>
  </el-form>
</template>

<style scoped lang="scss">
.live-room-form-fields {
  :deep(.el-form-item__label) {
    color: hsl(var(--foreground) / 88%);
  }

  .max-w460 {
    max-width: 460px;
    width: 100%;
  }

  .image-upload-row {
    margin-bottom: 4px;

    :deep(.el-form-item) {
      margin-bottom: 0;
    }
  }

  .field-hint {
    font-size: 12px;
    line-height: 1.5;
    color: hsl(var(--muted-foreground));
  }

  .mt10 {
    margin-top: 10px;
  }

  .switch-row {
    display: flex;
    flex-wrap: wrap;
    gap: 24px 32px;
    align-items: center;
    padding-left: 12px;
  }

  .switch-item {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .switch-label {
    font-size: 14px;
    color: hsl(var(--foreground) / 85%);
    white-space: nowrap;
  }

  .anchor-field {
    display: flex;
    gap: 8px;
    align-items: center;
    width: 100%;
    max-width: 460px;
  }

  .anchor-select {
    flex: 1;
    min-width: 240px;
  }

  .anchor-field :deep(.anchor-select.el-select) {
    width: 100%;
  }

  .anchor-lock-tip {
    margin: 6px 0 0;
    font-size: 12px;
    line-height: 1.5;
    color: hsl(var(--muted-foreground));
  }

  .product-field {
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 100%;
    max-width: 720px;
  }

  .product-actions {
    display: flex;
    gap: 8px;
  }

  .product-hint {
    font-size: 12px;
    line-height: 1.5;
    color: hsl(var(--muted-foreground));
  }

  .selected-product-table {
    width: 100%;
  }

  .record-video-field {
    display: flex;
    flex-direction: column;
    gap: 10px;
    max-width: 460px;
  }

  .record-video-actions {
    display: flex;
    gap: 8px;
    align-items: center;
  }

  .record-video-preview {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .record-video-player {
    width: 100%;
    max-width: 320px;
    max-height: 180px;
    background: #000;
    border-radius: 6px;
  }

  .record-video-path,
  .record-video-hint {
    font-size: 12px;
    line-height: 1.5;
    color: hsl(var(--muted-foreground));
    word-break: break-all;
  }

  .live-room-segment {
    display: flex;
    flex-wrap: wrap;
    gap: 4px 12px;
  }

  .live-room-segment :deep(.el-radio) {
    margin-right: 0;
  }

  .live-room-segment :deep(.el-radio__inner) {
    border-color: hsl(var(--border));
    background: hsl(var(--background));
  }

  .live-room-segment :deep(.el-radio__input.is-checked .el-radio__inner) {
    background: hsl(var(--primary));
    border-color: hsl(var(--primary));
  }

  .live-room-segment :deep(.el-radio__label) {
    color: hsl(var(--foreground) / 88%);
  }

  .stream-field__meta {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    margin-top: 8px;
    font-size: 12px;
    color: hsl(var(--muted-foreground));
  }
}
</style>
