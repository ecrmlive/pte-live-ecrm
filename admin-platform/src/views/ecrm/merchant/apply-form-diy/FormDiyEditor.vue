<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import draggable from 'vuedraggable';

import { DevicePreviewFrame } from '@pte-live/diy';

import ComponentPalette from './components/ComponentPalette.vue';
import FieldToolbar from './components/FieldToolbar.vue';
import FieldPreview from './preview/FieldPreview.vue';
import SystemFieldsPreview from './preview/SystemFieldsPreview.vue';
import {
  cloneField,
  createField,
} from './registry';
import FieldSettings from './settings/FieldSettings.vue';
import { SYSTEM_FIELDS } from './system-fields';
import type { ApplyFieldType, ApplyFormField, SystemField } from './types';

const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    /** 预览导航标题，默认入驻申请 */
    previewTitle?: string;
    /** 系统固定字段；不传则用商户入驻默认列表 */
    systemFields?: SystemField[];
  }>(),
  {
    previewTitle: '入驻申请',
  },
);

const previewSystemFields = computed(
  () => props.systemFields ?? SYSTEM_FIELDS,
);

const fields = defineModel<ApplyFormField[]>('fields', { required: true });
const selectedIndex = ref(-1);

watch(
  fields,
  (list) => {
    if (!list.length) {
      selectedIndex.value = -1;
      return;
    }
    if (selectedIndex.value < 0) {
      selectedIndex.value = 0;
      return;
    }
    if (selectedIndex.value >= list.length) {
      selectedIndex.value = list.length - 1;
    }
  },
  { deep: true, immediate: true },
);

const selectedField = computed(() => {
  if (selectedIndex.value < 0) return null;
  return fields.value[selectedIndex.value] ?? null;
});

function addField(type: ApplyFieldType) {
  if (props.disabled) return;
  fields.value.push(createField(type));
  selectedIndex.value = fields.value.length - 1;
}

function selectField(index: number) {
  selectedIndex.value = index;
}

function removeField(index: number) {
  if (props.disabled) return;
  fields.value.splice(index, 1);
}

function copyField(index: number) {
  if (props.disabled) return;
  const source = fields.value[index];
  if (!source) return;
  fields.value.splice(index + 1, 0, cloneField(source));
  selectedIndex.value = index + 1;
}

function moveField(index: number, direction: -1 | 1) {
  if (props.disabled) return;
  const target = index + direction;
  if (target < 0 || target >= fields.value.length) return;
  const list = fields.value;
  const current = list[index];
  const next = list[target];
  if (!current || !next) return;
  list[index] = next;
  list[target] = current;
  selectedIndex.value = target;
}

</script>

<template>
  <div class="form-diy-editor">
    <div class="diy-wrapper">
      <ComponentPalette :disabled="disabled" @add="addField" />

      <main class="diy-center">
        <DevicePreviewFrame
          :title="previewTitle"
          :show-back="true"
          :show-submit-bar="true"
          submit-text="提交"
          :side-gutter="52"
        >
          <SystemFieldsPreview :fields="previewSystemFields" />

          <draggable
            v-model="fields"
            item-key="id"
            handle=".field-wrap"
            :animation="180"
            :disabled="disabled"
            class="custom-list"
          >
            <template #item="{ element, index }">
              <div
                class="field-wrap"
                :class="{ on: selectedIndex === index }"
                @click="selectField(index)"
              >
                <FieldPreview :field="element" />
                <div v-if="selectedIndex === index" class="field-wrap__frame">
                  <FieldToolbar
                    :disabled="disabled"
                    :can-move-up="index > 0"
                    :can-move-down="index < fields.length - 1"
                    @up="moveField(index, -1)"
                    @down="moveField(index, 1)"
                    @copy="copyField(index)"
                    @remove="removeField(index)"
                  />
                </div>
              </div>
            </template>
          </draggable>

          <div v-if="!fields.length" class="phone__hint">
            点击左侧组件添加到表单
          </div>
        </DevicePreviewFrame>
      </main>

      <aside class="diy-right">
        <FieldSettings :field="selectedField" :disabled="disabled" />
      </aside>
    </div>
  </div>
</template>

<style scoped lang="scss">
.form-diy-editor {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  min-width: 1100px;
  background: #f0f2f5;
}

.diy-wrapper {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.diy-center {
  position: relative;
  display: flex;
  flex: 1;
  justify-content: center;
  align-items: flex-start;
  min-width: 0;
  padding: 20px 16px 16px;
  overflow: auto;
  /* 预览列外层滚动条（红框位置：机身/工具条右侧）隐藏，滚轮仍可用 */
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
  }
}

.phone__hint {
  padding: 20px 12px 28px;
  color: #bbb;
  text-align: center;
  font-size: 12px;
}

.custom-list {
  min-height: 8px;
  overflow: visible;
}

.field-wrap {
  position: relative;
  overflow: visible;
  background: #fff;
  cursor: move;
}

.field-wrap__frame {
  position: absolute;
  z-index: 6;
  top: 0;
  left: -2px;
  width: calc(100% + 4px);
  height: 100%;
  overflow: visible;
  border: 2px solid hsl(var(--primary));
  box-shadow: 0 0 10px 0 rgb(24 144 255 / 30%);
  pointer-events: none;
}

.field-wrap__frame :deep(.field-tools) {
  pointer-events: auto;
}

.diy-right {
  flex-shrink: 0;
  width: 400px;
  background: #fff;
  border-left: 1px solid #eee;
  overflow: auto;
}
</style>
