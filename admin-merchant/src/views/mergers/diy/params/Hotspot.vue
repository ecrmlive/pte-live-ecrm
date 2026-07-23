<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ArrowRight } from '@element-plus/icons-vue';
import { ElIcon } from 'element-plus';
import { computed, onUnmounted, ref } from 'vue';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import {
  diyBgColors,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyColorField from './shared/diy-color-field.vue';
import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsHotspot' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown> & {
    data?: Array<Record<string, unknown>>;
    params?: Record<string, unknown>;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const editorRef = ref<HTMLElement | null>(null);
const currentHotspotIndex = ref(-1);
const isLinkset = ref(false);
const linkData = ref<Record<string, unknown> | null>(null);

const dragState = ref({
  containerHeight: 0,
  containerWidth: 0,
  index: -1,
  isDragging: false,
  startHeight: 0,
  startLeft: 0,
  startTop: 0,
  startWidth: 0,
  startX: 0,
  startY: 0,
  type: null as 'move' | 'resize' | null,
});

const schema = computed((): VbenFormSchema[] => [
  diySection('边距设置'),
  ...diyBgColors('#ffffff', '#F2F2F2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema);

function addHotspot() {
  props.curItem.data = props.curItem.data ?? [];
  props.curItem.data.push({
    height: 20,
    left: 0,
    linkUrl: '',
    top: 0,
    width: 20,
  });
  currentHotspotIndex.value = props.curItem.data.length - 1;
}

function deleteHotspot(index: number) {
  props.curItem.data?.splice(index, 1);
  if (currentHotspotIndex.value === index) {
    currentHotspotIndex.value = -1;
  } else if (currentHotspotIndex.value > index) {
    currentHotspotIndex.value--;
  }
}

function onMouseDown(e: MouseEvent, index: number, type: 'move' | 'resize') {
  const container = editorRef.value;
  if (!container || !props.curItem.data?.[index]) return;
  const rect = container.getBoundingClientRect();
  const hotspot = props.curItem.data[index]!;

  dragState.value = {
    containerHeight: rect.height,
    containerWidth: rect.width,
    index,
    isDragging: true,
    startHeight: Number.parseFloat(String(hotspot.height)),
    startLeft: Number.parseFloat(String(hotspot.left)),
    startTop: Number.parseFloat(String(hotspot.top)),
    startWidth: Number.parseFloat(String(hotspot.width)),
    startX: e.clientX,
    startY: e.clientY,
    type,
  };
  currentHotspotIndex.value = index;
  document.addEventListener('mousemove', onMouseMove);
  document.addEventListener('mouseup', onMouseUp);
}

function onMouseMove(e: MouseEvent) {
  if (!dragState.value.isDragging || !props.curItem.data) return;
  const {
    containerHeight,
    containerWidth,
    index,
    startHeight,
    startLeft,
    startTop,
    startWidth,
    startX,
    startY,
    type,
  } = dragState.value;
  const hotspot = props.curItem.data[index]!;
  const deltaXPercent = ((e.clientX - startX) / containerWidth) * 100;
  const deltaYPercent = ((e.clientY - startY) / containerHeight) * 100;

  if (type === 'move') {
    hotspot.left = Math.max(
      0,
      Math.min(100 - Number(hotspot.width), startLeft + deltaXPercent),
    );
    hotspot.top = Math.max(
      0,
      Math.min(100 - Number(hotspot.height), startTop + deltaYPercent),
    );
  } else if (type === 'resize') {
    hotspot.width = Math.max(
      5,
      Math.min(100 - Number(hotspot.left), startWidth + deltaXPercent),
    );
    hotspot.height = Math.max(
      5,
      Math.min(100 - Number(hotspot.top), startHeight + deltaYPercent),
    );
  }
}

function onMouseUp() {
  dragState.value.isDragging = false;
  document.removeEventListener('mousemove', onMouseMove);
  document.removeEventListener('mouseup', onMouseUp);
}

function changeLink(index: number) {
  isLinkset.value = true;
  currentHotspotIndex.value = index;
  linkData.value = props.curItem.data?.[index] ?? null;
}

function closeLinkset(e: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  if (e && props.curItem.data?.[currentHotspotIndex.value]) {
    const hotspot = props.curItem.data[currentHotspotIndex.value]!;
    hotspot.linkeType = e.type;
    hotspot.linkUrl = e.url;
    hotspot.name = e.name;
  }
}

onUnmounted(() => {
  document.removeEventListener('mousemove', onMouseMove);
  document.removeEventListener('mouseup', onMouseUp);
});
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <div class="form-chink" />
    <div class="f16 gray3 form-subtitle mb-2 px-4">热区设置</div>
    <div class="form-item ml-[100px]">
      <div class="form-label mb-2">背景图片：</div>
      <div
        class="img-box flex cursor-pointer items-center justify-center"
        style="width: 100px; height: 100px; border: 1px dashed #ccc"
        @click="editor.onEditorSelectImage(curItem.params!, 'image')"
      >
        <img
          v-if="curItem.params?.image"
          v-img-url="curItem.params.image"
          style="max-width: 100%; max-height: 100%"
        />
        <div v-else class="upload-btn" style="font-size: 30px; color: #ccc">+</div>
      </div>
    </div>
    <div class="form-item ml-5 block">
      <div class="form-label mb-2">热区编辑：</div>
      <div
        ref="editorRef"
        class="hotspot-editor relative w-full select-none border border-[#ccc]"
      >
        <img
          v-img-url="curItem.params?.image || ''"
          draggable="false"
          style="width: 100%; display: block"
        />
        <div
          v-for="(hotspot, index) in curItem.data"
          :key="index"
          class="hotspot-box"
          :class="{ active: currentHotspotIndex === index }"
          :style="{
            width: `${hotspot.width}%`,
            height: `${hotspot.height}%`,
            top: `${hotspot.top}%`,
            left: `${hotspot.left}%`,
          }"
          @click.stop="currentHotspotIndex = index"
          @mousedown.stop="onMouseDown($event, index, 'move')"
        >
          <span class="index-label">{{ index + 1 }}</span>
          <div
            class="resize-handle"
            @mousedown.stop="onMouseDown($event, index, 'resize')"
          />
          <div class="delete-btn" @click.stop="deleteHotspot(index)">×</div>
        </div>
      </div>
      <div class="d-c-c mt-2">
        <component :is="PrimaryButton" plain size="small" @click="addHotspot">添加热区</component>
      </div>
    </div>
    <div v-if="curItem.data && curItem.data.length > 0" class="mt-2">
      <div class="f16 gray3 form-subtitle px-4">链接设置</div>
      <div
        v-for="(hotspot, index) in curItem.data"
        :key="index"
        class="form-item mb-2 ml-5 p-2"
        :style="{
          background: '#f5f5f5',
          border:
            currentHotspotIndex === index
              ? '1px solid #409eff'
              : '1px solid transparent',
        }"
        @click="currentHotspotIndex = index"
      >
        <div class="d-s-c">
          <div class="form-label" style="width: auto">热区 {{ index + 1 }}：</div>
        </div>
        <div class="d-s-c">
          <DiyLinkInputField
            v-model="hotspot.linkUrl"
            placeholder="选择链接"
            readonly
            @click="changeLink(index)"
          >
<template #suffix>
<ElIcon color="#333" size="16px"><ArrowRight /></ElIcon>
</template>
</DiyLinkInputField>
        </div>
        <div class="d-s-c">
          <div class="flex-1 tr">
            <component :is="PrimaryButton" link size="small" class="!text-destructive" @click.stop="deleteHotspot(index)">删除</component>
          </div>
        </div>
      </div>
    </div>
    <DiyLinkPickerDialog
      v-if="isLinkset"
      :is_linkset="isLinkset"
      :link-data="linkData"
      @close-dialog="closeLinkset"
    >
      选择链接
    </DiyLinkPickerDialog>
  </div>
</template>

<style lang="scss" scoped>
.hotspot-box {
  position: absolute;
  border: 1px dashed #409eff;
  background: rgba(64, 158, 255, 0.3);
  cursor: move;
}
.hotspot-box.active {
  border: 2px solid #409eff;
  background: rgba(64, 158, 255, 0.5);
  z-index: 10;
}
.index-label {
  position: absolute;
  top: 0;
  left: 0;
  background: #409eff;
  color: #fff;
  font-size: 12px;
  padding: 0 4px;
}
.resize-handle {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 10px;
  height: 10px;
  background: #fff;
  border: 1px solid #409eff;
  cursor: se-resize;
}
.delete-btn {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 16px;
  height: 16px;
  background: #f56c6c;
  color: #fff;
  border-radius: 50%;
  text-align: center;
  line-height: 14px;
  font-size: 12px;
  cursor: pointer;
  display: none;
}
.hotspot-box:hover .delete-btn {
  display: block;
}
</style>
