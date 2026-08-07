<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';
import type {
  AgentPosterData,
  AgentPosterFormPayload,
} from '#/api/core/plus-agent';

import { useVbenDrawer } from '@vben/common-ui';
import { User } from '@element-plus/icons-vue';
import { ElButton, ElMessage } from 'element-plus';
import { computed,  markRaw, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import {
  addAgentPosterApi,
  createDefaultAgentPosterData,
  editAgentPosterApi,
  getAgentPosterAddMetaApi,
  getAgentPosterEditMetaApi,
} from '#/api/core/plus-agent';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

import qrcodePlaceholder from '#/assets/img/qrcode.png';
import ImageField from '#/components/shop/image-field.vue';
import PosterLayoutField from '#/views/native/plus/agent/poster-layout-field.vue';
import { nativeSectionTitle } from '#/utils/native-form-schema';

defineOptions({ name: 'AgentPosterFormDialog' });

type DragType = 'code' | 'name' | 'photo';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  posterId?: number;
}>();

const emit = defineEmits<{
  success: [];
}>();

const loading = ref(false);
const pageLoading = ref(false);
const previewData = ref<AgentPosterData>(createDefaultAgentPosterData());
const previewImage = ref('');

const posterId = computed(() => Number(props.posterId ?? 0));
const isEdit = computed(() => posterId.value > 0);
const modalTitle = computed(() => (isEdit.value ? '编辑分销海报' : '添加分销海报'));

const backdropUrl = computed(() =>
  resolveCosMediaUrl(previewImage.value || previewData.value.backdrop.src || ''),
);

const schema = computed((): VbenFormSchema[] => [
  nativeSectionTitle('_basic', '基础设置', { formItemClass: 'col-span-full native-section-title-item' }),
  {
    component: 'Input',
    componentProps: { class: 'max-w-[460px]', placeholder: '请输入海报名称' },
    fieldName: 'poster_name',
    label: '海报名称',
    rules: 'required',
  },
  {
    component: 'Input',
    componentProps: { class: 'w-[180px]', type: 'number' },
    fieldName: 'sort',
    label: '海报排序',
    rules: 'required',
  },
  {
    component: markRaw(ImageField),
    componentProps: {
      hint: '尺寸：750px × 1334px',
    },
    fieldName: 'poster_image',
    label: '海报背景图',
  },
  nativeSectionTitle('_poster', '海报设置', { formItemClass: 'col-span-full native-section-title-item' }),
  {
    component: markRaw(PosterLayoutField),
    defaultValue: createDefaultAgentPosterData(),
    fieldName: 'poster_data',
    formItemClass: 'col-span-full poster-layout-form-item',
    hideLabel: true,
  },
]);

const [Form, formApi] = useVbenForm(
  reactive({
    commonConfig: {
      labelWidth: 120,
    },
    handleSubmit: async (values) => {
      loading.value = true;
      try {
        const posterData = values.poster_data as AgentPosterData;
        const posterImage = String(values.poster_image ?? posterData.backdrop.src ?? '');
        const payload = {
          poster_data: {
            ...posterData,
            backdrop: { src: posterImage },
          },
          poster_image: posterImage,
          poster_name: String(values.poster_name ?? ''),
          sort: Number(values.sort ?? 100),
        };
        if (isEdit.value) {
          await editAgentPosterApi({ ...payload, poster_id: posterId.value });
          ElMessage.success('编辑成功');
        } else {
          await addAgentPosterApi(payload);
          ElMessage.success('添加成功');
        }
        open.value = false;
        emit('success');
      } finally {
        loading.value = false;
      }
    },
    handleValuesChange: (values) => {
      const posterData = values.poster_data as AgentPosterData | undefined;
      const posterImage = values.poster_image as string | undefined;
      if (posterData) {
        previewData.value = posterData;
      }
      if (posterImage) {
        previewImage.value = posterImage;
        previewData.value = {
          ...previewData.value,
          backdrop: { src: posterImage },
        };
      }
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema,
    showDefaultActions: false,
  }),
);

function normalizePosterData(raw: AgentPosterData) {
  const data = { ...createDefaultAgentPosterData(), ...raw };
  data.avatar.left = Number(data.avatar.left ?? 0);
  data.avatar.top = Number(data.avatar.top ?? 0);
  data.nickName.left = Number(data.nickName.left ?? 0);
  data.nickName.top = Number(data.nickName.top ?? 0);
  data.qrcode.left = Number(data.qrcode.left ?? 0);
  data.qrcode.top = Number(data.qrcode.top ?? 0);
  data.avatar.width = Number(data.avatar.width ?? 30);
  data.qrcode.width = Number(data.qrcode.width ?? 50);
  data.nickName.fontSize = Number(data.nickName.fontSize ?? 20);
  return data;
}

function extractEditForm(res: Record<string, unknown>) {
  const nested = res.data;
  if (nested && typeof nested === 'object') {
    return nested as AgentPosterFormPayload;
  }
  return res as unknown as AgentPosterFormPayload;
}

async function loadDetail() {
  pageLoading.value = true;
  try {
    if (isEdit.value) {
      const res = await getAgentPosterEditMetaApi(posterId.value);
      const payload = extractEditForm(res);
      const posterData = normalizePosterData(payload.poster_data);
      const posterImage = String(payload.poster_image ?? posterData.backdrop.src);
      previewData.value = posterData;
      previewImage.value = posterImage;
      await formApi.setValues({
        poster_data: posterData,
        poster_image: posterImage,
        poster_name: String(payload.poster_name ?? ''),
        sort: Number(payload.sort ?? 100),
      });
      return;
    }
    const res = await getAgentPosterAddMetaApi();
    const posterData = normalizePosterData(res.data ?? createDefaultAgentPosterData());
    const posterImage = posterData.backdrop.src;
    previewData.value = posterData;
    previewImage.value = posterImage;
    await formApi.setValues({
      poster_data: posterData,
      poster_image: posterImage,
      poster_name: '',
      sort: 100,
    });
  } finally {
    pageLoading.value = false;
  }
}

function assignDragPosition(type: DragType, x: number, y: number) {
  const left = Math.max(0, Math.round(x));
  const top = Math.max(0, Math.round(y));
  const next = { ...previewData.value };
  if (type === 'photo') {
    next.avatar = { ...next.avatar, left, top };
  } else if (type === 'code') {
    next.qrcode = { ...next.qrcode, left, top };
  } else {
    next.nickName = { ...next.nickName, left, top };
  }
  previewData.value = next;
  void formApi.setFieldValue('poster_data', next);
}

function startDrag(event: MouseEvent, type: DragType) {
  event.preventDefault();
  const startX = event.clientX;
  const startY = event.clientY;
  const originLeft =
    type === 'photo'
      ? previewData.value.avatar.left
      : type === 'code'
        ? previewData.value.qrcode.left
        : previewData.value.nickName.left;
  const originTop =
    type === 'photo'
      ? previewData.value.avatar.top
      : type === 'code'
        ? previewData.value.qrcode.top
        : previewData.value.nickName.top;

  const onMove = (moveEvent: MouseEvent) => {
    assignDragPosition(
      type,
      originLeft + moveEvent.clientX - startX,
      originTop + moveEvent.clientY - startY,
    );
  };
  const onUp = () => {
    document.removeEventListener('mousemove', onMove);
    document.removeEventListener('mouseup', onUp);
  };
  document.addEventListener('mousemove', onMove);
  document.addEventListener('mouseup', onUp);
}

async function submit() {
  await formApi.validateAndSubmitForm();
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void loadDetail();
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :destroy-on-close="true"
    :title="modalTitle"
    class="native-form-dialog w-[min(1080px,96vw)]"
  >
    <div v-loading="pageLoading" class="poster-editor native-form-scroll">
      <div class="poster-editor__preview">
        <div class="poster-canvas">
          <img v-if="backdropUrl" :src="backdropUrl" alt="" class="poster-canvas__bg" />
          <div class="poster-canvas__layer">
            <div
              class="poster-canvas__avatar"
              :class="{ 'is-circle': previewData.avatar.style === 'circle' }"
              :style="{
                width: `${previewData.avatar.width}px`,
                height: `${previewData.avatar.width}px`,
                top: `${previewData.avatar.top}px`,
                left: `${previewData.avatar.left}px`,
              }"
              @mousedown="startDrag($event, 'photo')"
            >
              <el-icon :size="28"><User /></el-icon>
            </div>
            <div
              class="poster-canvas__name"
              :style="{
                fontSize: `${previewData.nickName.fontSize}px`,
                color: previewData.nickName.color,
                top: `${previewData.nickName.top}px`,
                left: `${previewData.nickName.left}px`,
              }"
              @mousedown="startDrag($event, 'name')"
            >
              昵称
            </div>
            <div
              class="poster-canvas__qrcode"
              :class="{ 'is-circle': previewData.qrcode.style === 'circle' }"
              :style="{
                width: `${previewData.qrcode.width}px`,
                height: `${previewData.qrcode.width}px`,
                top: `${previewData.qrcode.top}px`,
                left: `${previewData.qrcode.left}px`,
              }"
              @mousedown="startDrag($event, 'code')"
            >
              <img :src="qrcodePlaceholder" alt="" />
            </div>
          </div>
        </div>
      </div>

      <div class="poster-editor__form">
        <Form />
      </div>
    </div>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="submit">提交</ElButton>
    </template>
  </Modal>
</template>

<style scoped lang="scss">
.poster-editor {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.poster-editor__form {
  :deep(.native-section-title-item) {
    margin-bottom: 0;

    .el-form-item__label {
      display: none;
    }

    .el-form-item__content {
      margin-left: 0 !important;
      max-width: none;
    }
  }

  :deep(.poster-layout-form-item .el-form-item__content) {
    margin-left: 0 !important;
    max-width: none;
  }
}

.poster-editor__preview {
  position: sticky;
  top: 0;
  flex-shrink: 0;
}

.poster-editor__form {
  flex: 1;
  min-width: 0;
}

.poster-canvas {
  position: relative;
  width: 300px;
  height: 534px;
  overflow: hidden;
  border: 1px solid hsl(var(--border));
  border-radius: 16px;
  background: hsl(var(--muted));
}

.poster-canvas__bg {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.poster-canvas__layer {
  position: absolute;
  inset: 0;
}

.poster-canvas__avatar,
.poster-canvas__name,
.poster-canvas__qrcode {
  position: absolute;
  cursor: move;
  user-select: none;
}

.poster-canvas__avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: #eee;
}

.poster-canvas__avatar.is-circle,
.poster-canvas__qrcode.is-circle {
  border-radius: 50%;
}

.poster-canvas__name {
  min-height: 24px;
  font-weight: 700;
  line-height: 1.2;
  white-space: nowrap;
}

.poster-canvas__qrcode {
  overflow: hidden;
  background: #fff;
}

.poster-canvas__qrcode img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
