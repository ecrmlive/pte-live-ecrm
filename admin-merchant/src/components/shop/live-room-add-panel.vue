<script setup lang="ts">
import type { ProductChooseItem } from '#/api/core/product';
import type { LiveVodVideoItem } from '#/api/core/live-vod';
import type { LiveRoomForm } from '#/api/core/live';

import { useVbenDrawer, useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { reactive, ref, watch } from 'vue';

import { createLiveRoomApi, getLiveAnchorListApi } from '#/api/core/live';
import { buildLiveRoomPayload, defaultLiveRoomForm } from '#/utils/live-room-payload';
import { dismissStaleModalOverlays } from '#/utils/modal-stack';

import LiveAnchorAddDialog from '#/components/shop/live-anchor-add-dialog.vue';
import LiveRoomFormFields from '#/components/shop/live-room-form-fields.vue';
import ProductMultiPickerDialog from '#/components/shop/product-multi-picker-dialog.vue';
import VodVideoLibrary from '#/components/shop/vod-video-library.vue';
import ProductAddModal from '#/views/native/product/product/product-add-modal.vue';

defineOptions({ name: 'LiveRoomAddPanel' });

const props = defineProps<{
  active?: boolean;
  /** 弹窗内嵌：不套 native-form-page 卡片壳 */
  embedded?: boolean;
  hideActions?: boolean;
}>();

const emit = defineEmits<{
  cancel: [];
  success: [];
}>();

const formRef = ref<InstanceType<typeof LiveRoomFormFields>>();
const submitting = ref(false);
const vodPickerOpen = ref(false);
const anchorAddOpen = ref(false);
const productPickerOpen = ref(false);
const anchorOptions = ref<Awaited<ReturnType<typeof getLiveAnchorListApi>>['list']['data']>([]);
const selectedProducts = ref<ProductChooseItem[]>([]);

const form = reactive<LiveRoomForm & ReturnType<typeof defaultLiveRoomForm>>(defaultLiveRoomForm());

async function loadAnchors() {
  try {
    const res = await getLiveAnchorListApi({ list_rows: 200, page: 1, status: 1 });
    anchorOptions.value = res.list.data ?? [];
  } catch {
    anchorOptions.value = [];
  }
}

function onAnchorChange(anchorId?: number) {
  const hit = anchorOptions.value.find((item) => item.anchor_id === anchorId);
  if (hit) {
    form.anchor_name = hit.nick_name;
    form.anchor_wechat = hit.wechat || '';
  }
}

function onAnchorAdded() {
  void loadAnchors();
}

const [ProductAddDialog, productAddModalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  connectedComponent: ProductAddModal,
});

function openProductAdd() {
  dismissStaleModalOverlays();
  productAddModalApi.open();
}

function onVodPicked(item: LiveVodVideoItem) {
  form.record_vod_file_id = item.file_id || '';
  form.record_vod_media_url = item.media_url || '';
  form.record_video_path = '';
  vodPickerOpen.value = false;
}

function onProductPickerClose(payload: {
  openDialog: boolean;
  params?: unknown;
  type: string;
}) {
  productPickerOpen.value = false;
  if (payload.type !== 'success' || !payload.params) return;
  const rows = (Array.isArray(payload.params) ? payload.params : [payload.params]) as ProductChooseItem[];
  const map = new Map(selectedProducts.value.map((item) => [item.product_id, item]));
  rows.forEach((row) => {
    if (row.product_id) {
      map.set(row.product_id, {
        ...row,
        product_image: row.image?.[0]?.file_path ?? row.product_image,
      });
    }
  });
  selectedProducts.value = [...map.values()];
}

function removeProduct(productId: number) {
  selectedProducts.value = selectedProducts.value.filter((item) => item.product_id !== productId);
}

function resetForm() {
  Object.assign(form, defaultLiveRoomForm());
  selectedProducts.value = [];
  formRef.value?.resetFields();
}

async function submit() {
  const valid = (await formRef.value?.validate()) ?? false;
  if (!valid) return;
  if (
    form.room_type === 2 &&
    !form.record_video_path?.trim() &&
    !form.record_vod_media_url?.trim() &&
    !form.record_vod_file_id?.trim()
  ) {
    ElMessage.error('录播房请选择录播视频');
    return;
  }
  submitting.value = true;
  try {
    const payload = buildLiveRoomPayload(form, {
      product_ids: selectedProducts.value.map((item) => item.product_id).join(','),
    });
    const res = await createLiveRoomApi(payload);
    const pushURL = res.data?.push_url;
    if (form.room_type === 1 && pushURL) {
      ElMessage.success({ duration: 8000, message: `创建成功，推流地址：${pushURL}`, showClose: true });
    } else {
      ElMessage.success(res.msg || '创建成功');
    }
    emit('success');
  } finally {
    submitting.value = false;
  }
}

function initForm() {
  resetForm();
  void loadAnchors();
}

const [VodPickerModal, vodPickerModalApi] = useVbenModal({
  onOpenChange(isOpen) {
    vodPickerOpen.value = isOpen;
  },
});

watch(vodPickerOpen, (visible) => {
  if (visible) {
    vodPickerModalApi.open();
    return;
  }
  vodPickerModalApi.close();
});

watch(
  () => props.active,
  (visible) => {
    if (visible) initForm();
  },
  { immediate: true },
);

defineExpose({ initForm, resetForm, submit, submitting });
</script>

<template>
  <div
    class="live-room-add-panel"
    :class="{ 'native-form-page native-form-shell': !embedded }"
  >
    <LiveRoomFormFields
      ref="formRef"
      :form="form"
      :anchor-options="anchorOptions"
      :selected-products="selectedProducts"
      @anchor-change="onAnchorChange"
      @open-anchor-add="anchorAddOpen = true"
      @open-product-add="openProductAdd"
      @open-product-pick="productPickerOpen = true"
      @pick-vod="vodPickerOpen = true"
      @remove-product="removeProduct"
    />

    <div v-if="!hideActions" class="live-room-add-panel__actions mt-6 flex justify-center gap-3">
      <el-button @click="emit('cancel')">取消</el-button>
      <el-button :loading="submitting" type="primary" @click="submit">保存</el-button>
    </div>

    <ProductAddDialog />
    <ProductMultiPickerDialog
      v-model:open="productPickerOpen"
      :exclude-ids="selectedProducts.map((item) => item.product_id)"
      @close-dialog="onProductPickerClose"
    />
    <LiveAnchorAddDialog v-model:open="anchorAddOpen" @success="onAnchorAdded" />
    <VodPickerModal
      :close-on-click-modal="false"
      :destroy-on-close="true"
      class="w-[960px]"
      title="选择录播视频"
    >
      <VodVideoLibrary picker @select="onVodPicked" />

      <template #footer>
        <ElButton @click="vodPickerOpen = false">取消</ElButton>
      </template>
    </VodPickerModal>
  </div>
</template>

<style scoped lang="scss">
.live-room-add-panel.native-form-page {
  min-height: auto;
  margin: 0;
  padding: 0;
  border: 0;
  border-radius: 0;
  box-shadow: none;
  background: transparent;
}

.live-room-add-panel :deep(.el-input),
.live-room-add-panel :deep(.el-select),
.live-room-add-panel :deep(.el-date-editor) {
  width: 100%;
  max-width: 460px;
}

.live-room-add-panel :deep(.el-radio-button__inner) {
  border-color: hsl(var(--border));
  background: hsl(var(--background));
  color: hsl(var(--foreground));
}

.live-room-add-panel :deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: hsl(var(--primary));
  border-color: hsl(var(--primary));
  color: hsl(var(--primary-foreground));
}
</style>
