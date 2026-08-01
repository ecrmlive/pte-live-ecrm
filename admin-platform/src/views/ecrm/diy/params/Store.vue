<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { CircleCloseFilled, Plus } from '@element-plus/icons-vue';
import { ElIcon } from 'element-plus';
import { computed } from 'vue';

import {
  diyBgColors,
  diyInput,
  diyRadioGroup,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsStore' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    data?: Array<Record<string, unknown>>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const schema = computed((): VbenFormSchema[] => [
  diySection('内容设置'),
  diyRadioGroup(
    'params.source',
    '商品来源：',
    [
      { label: '自动获取', value: 'auto' },
      { label: '手动选择', value: 'choice' },
    ],
  ),
  ...(props.curItem.params?.source === 'auto'
    ? [diyInput('params.auto.showNum', '展示数量：')]
    : []),
  diySection('样式设置'),
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema);
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <template v-if="curItem.params?.source == 'choice'">
      <div class="form-item ml-[100px]">
        <div class="form-label mb-2">门店列表：</div>
        <div class="choice-shop-list">
          <div
            v-for="(shop, index) in curItem.data"
            :key="index"
            class="item"
          >
            <div class="delete-box">
              <ElIcon @click="editor.onEditorDeleleData(index, selectedIndex ?? 0)">
                <CircleCloseFilled />
              </ElIcon>
            </div>
            <img v-if="shop.logo" :src="(shop.logo as { file_path?: string }).file_path" alt="" />
            <img v-else :src="shop.logo_image as string" alt="" />
          </div>
          <div class="item plus-btn" @click.stop="editor.openStore(true)">
            <ElIcon><Plus /></ElIcon>
            <p>选择门店</p>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.choice-shop-list {
  display: flex;
  justify-content: flex-start;
  flex-wrap: wrap;
}
.choice-shop-list .item {
  position: relative;
  width: 80px;
  height: 80px;
  margin-right: 10px;
  border: 1px solid #dddddd;
}
.choice-shop-list .item .delete-box {
  position: absolute;
  width: 20px;
  height: 20px;
  top: -10px;
  right: -10px;
  font-size: 20px;
  cursor: pointer;
  color: #999999;
}
.choice-shop-list .item .delete-box:hover {
  color: rgb(255, 51, 0);
}
.choice-shop-list .item.plus-btn {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.choice-shop-list .item.plus-btn > i {
  font-size: 30px;
  color: #cccccc;
}
.choice-shop-list img {
  width: 100%;
  height: 100%;
}
</style>
