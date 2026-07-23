<script setup lang="ts">
import type { DeliveryAreaNode } from '#/api/core/shop-setting';

import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { ref, watch } from 'vue';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  areaIndex: number;
  areas: DeliveryAreaNode[];
}>();

const emit = defineEmits<{
  cancel: [];
  confirm: [DeliveryAreaNode[]];
}>();

const provinceIndex = ref<number | null>(null);
const options = ref<DeliveryAreaNode[]>([]);
const allActive = ref(false);
const indeterminate = ref(false);

function cloneAreas(list: DeliveryAreaNode[]) {
  return structuredClone(list);
}

function setList(list: DeliveryAreaNode[]) {
  for (const item of list) {
    if (item.checked !== true) item.checked = false;
    if (item.indeterminate !== true) item.indeterminate = false;
    if (item.index == null) {
      item.index = [];
      item.disabled = false;
    }
    if (Array.isArray(item.children)) {
      let thisNum = 0;
      let otherNum = 0;
      const count = item.children.length;
      for (const child of item.children) {
        if (child.checked !== true) child.checked = false;
        if (child.index === props.areaIndex && item.index!.includes(props.areaIndex)) {
          child.checked = true;
          child.disabled = false;
          thisNum++;
        } else if (child.index != null && item.index!.includes(child.index as number)) {
          child.checked = true;
          child.disabled = true;
          otherNum++;
        } else {
          child.checked = false;
          child.disabled = false;
        }
        if (thisNum === count || otherNum === count) {
          item.checked = true;
          item.indeterminate = false;
          item.disabled = otherNum === count;
        } else if (thisNum === 0 || otherNum === 0) {
          item.checked = false;
          item.indeterminate = false;
          item.disabled = false;
        } else {
          item.disabled = false;
          item.checked = false;
          item.indeterminate = thisNum > 0;
        }
      }
    }
  }
}

function isChooseAll(list: DeliveryAreaNode[]) {
  let allcount = 0;
  for (const item of list) {
    if (Array.isArray(item.index)) {
      if (
        (item.checked === true && item.index.includes(props.areaIndex)) ||
        item.indeterminate === true
      ) {
        allcount++;
      }
    } else if (item.checked === true || item.indeterminate === true) {
      allcount++;
    }
  }
  if (allcount === list.length) return 2;
  return allcount > 0 ? 1 : 0;
}

function syncAllState() {
  const n = isChooseAll(options.value);
  allActive.value = n === 2;
  indeterminate.value = n === 1;
}

function initDialog() {
  setList(props.areas);
  options.value = cloneAreas(props.areas);
  provinceIndex.value = null;
  syncAllState();
}

function allProvinceFunc() {
  for (const item of options.value) {
    if (item.disabled) continue;
    item.checked = allActive.value;
    if (allActive.value) {
      if (!item.index!.includes(props.areaIndex)) item.index!.push(props.areaIndex);
    } else {
      const idx = item.index!.indexOf(props.areaIndex);
      if (idx !== -1) item.index!.splice(idx, 1);
    }
    item.indeterminate = indeterminate.value;
    item.children?.forEach((child) => {
      if (child.disabled) return;
      child.checked = allActive.value;
      child.index = allActive.value ? props.areaIndex : null;
    });
  }
}

function handleCheckedProvinceChange(i: number) {
  provinceIndex.value = i;
  const province = options.value[i];
  if (!province) return;
  province.indeterminate = false;
  if (province.checked && !province.disabled) {
    if (!province.index!.includes(props.areaIndex)) province.index!.push(props.areaIndex);
  } else {
    const idx = province.index!.indexOf(props.areaIndex);
    if (idx !== -1) province.index!.splice(idx, 1);
    if (!province.disabled) province.checked = false;
  }
  province.children?.forEach((child) => {
    if (child.disabled) return;
    child.checked = province.checked;
    child.index = props.areaIndex;
  });
  syncAllState();
}

function allCityFunc() {
  const province = provinceIndex.value == null ? null : options.value[provinceIndex.value];
  if (!province) return;
  const flag = province.checked;
  if (flag) {
    if (!province.index!.includes(props.areaIndex)) province.index!.push(props.areaIndex);
  } else {
    const idx = province.index!.indexOf(props.areaIndex);
    if (idx !== -1) province.index!.splice(idx, 1);
  }
  province.indeterminate = false;
  province.children?.forEach((child) => {
    if (child.disabled) return;
    child.checked = flag;
    child.index = flag ? props.areaIndex : null;
  });
  syncAllState();
}

function handleCheckedCityChange(i: number) {
  const province = provinceIndex.value == null ? null : options.value[provinceIndex.value];
  if (!province?.children) return;
  const child = province.children[i];
  const flag = child?.checked;
  const provIdx = province.index!.indexOf(props.areaIndex);
  const n = isChooseAll(province.children);
  if (flag) {
    if (provIdx === -1) province.index!.push(props.areaIndex);
    child!.index = props.areaIndex;
  } else {
    if (provIdx !== -1) province.index!.splice(provIdx, 1);
    child!.index = null;
  }
  if (n === 0) {
    province.checked = false;
    province.indeterminate = false;
  } else if (n === 2) {
    province.checked = true;
    province.indeterminate = false;
  } else {
    province.checked = false;
    province.indeterminate = true;
  }
  syncAllState();
}

function cityShow(i: number) {
  provinceIndex.value = i;
}

function closeArea() {
  open.value = false;
  emit('cancel');
}

function confirmArea() {
  open.value = false;
  emit('confirm', options.value);
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    initDialog();
    modalApi.open();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[700px]"
    title="添加可配送区域"
  >
    <div class="flex gap-3">
      <div class="h-[400px] flex-1 overflow-y-auto border p-2">
        <el-checkbox v-model="allActive" :indeterminate="indeterminate" @change="allProvinceFunc">
          全选
        </el-checkbox>
        <div
          v-for="(item, index) in options"
          :key="String(item.value)"
          class="cursor-pointer py-1"
          :class="{ 'bg-muted': index === provinceIndex }"
          @click="cityShow(index)"
        >
          <el-checkbox
            v-model="item.checked"
            :disabled="item.disabled"
            :indeterminate="item.indeterminate"
            :label="item.value"
            @change="handleCheckedProvinceChange(index)"
          >
            {{ item.label }}
          </el-checkbox>
        </div>
      </div>
      <div class="h-[400px] flex-1 overflow-y-auto border p-2">
        <template v-if="provinceIndex != null && options[provinceIndex]">
          <el-checkbox
            v-model="options[provinceIndex]!.checked"
            :disabled="options[provinceIndex]!.disabled"
            :indeterminate="options[provinceIndex]!.indeterminate"
            @change="allCityFunc"
          >
            全选
          </el-checkbox>
          <div v-for="(city, index) in options[provinceIndex]!.children" :key="String(city.value)" class="py-1">
            <el-checkbox
              v-model="city.checked"
              :disabled="city.disabled"
              :label="city.value"
              @change="handleCheckedCityChange(index)"
            >
              {{ city.label }}
            </el-checkbox>
          </div>
        </template>
      </div>
    </div>
    <template #footer>
      <ElButton @click="closeArea">取消</ElButton>
      <ElButton type="primary" @click="confirmArea">确定</ElButton>
    </template>
  </Modal>
</template>
