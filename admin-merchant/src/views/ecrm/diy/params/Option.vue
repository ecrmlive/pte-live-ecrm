<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';

import { diyInput, diySlider } from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyColorField from './shared/diy-color-field.vue';
import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsOption' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

import { getShopLinkProductCategoryApi } from '#/api/core/shop-link';
import draggable from 'vuedraggable';
import { resetStyleColors } from './shared/marketing-helpers';

const styleType = ref<'content' | 'style'>('content');
const loading = ref(true);
const CategoryList = ref<Array<Record<string, unknown>>>([]);

async function getData() {
  loading.value = true;
  try {
    const res = await getShopLinkProductCategoryApi();
    CategoryList.value = res.list ?? [];
  } finally {
    loading.value = false;
  }
}

function changeCategory(_e: unknown, index: number, cascaderRef: { getCheckedNodes: () => Array<{ data: { category_id: number; name: string } }> }) {
  const item = cascaderRef.getCheckedNodes();
  const row = (props.curItem.data as Array<Record<string, unknown>>)[index];
  row.name = item[0]?.data.name;
  row.category_id = item[0]?.data.category_id;
}

function ResetColor(titleColor1: string, titleColor2?: string) {
  resetStyleColors(editor, props.curItem, titleColor1, titleColor2);
}

void getData();


const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'params.data', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.topUp', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.type', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.activeText', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.active_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.active_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.background', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.marginTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingBottom', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingLeft', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.themeType', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.topRadio', hideLabel: true, formItemClass: 'hidden' },
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, { fieldPaths: ["params.data","params.topUp","params.type","style.activeText","style.active_color1","style.active_color2","style.background","style.bgcolor_color1","style.bgcolor_color2","style.bottomRadio","style.marginTop","style.paddingBottom","style.paddingLeft","style.paddingTop","style.themeType","style.topRadio"] });
</script>

<template>
	<div>
		<div class="diy-vben-params">
		<Form class="hidden" />
			<div class="common-form common-form-new">
				<span>{{ curItem.name }}</span>
				<div class="diy-changes">
					<div class="diy-change" :class="{active:styleType == 'content'}" @click="styleType = 'content'">内容
					</div>
					<div class="diy-change" :class="{active:styleType == 'style'}" @click="styleType = 'style'">样式</div>
				</div>
			</div>
			<div v-if="styleType == 'content'">
				<div class="f16 gray3 form-subtitle">展示设置</div>
				<!--选择风格-->
				<div class="form-vben-item" label="选择风格：">
					<component :is="RadioGroup" v-model="curItem.params.type">
						<el-radio label="1">风格一</el-radio>
						<el-radio label="2"> 风格二</el-radio>
						<el-radio label="3">风格三</el-radio>
					</component>
				</div>
				<!-- <div class="form-vben-item" label="滑动置顶：">
					<component :is="RadioGroup" v-model="curItem.params.topUp">
						<el-radio label="1">开启</el-radio>
						<el-radio label="0">关闭</el-radio>
					</component>
				</div> -->
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					选项卡设置
					<span class="gray f12">鼠标拖拽版块可调整选项卡顺序</span>
				</div>
				<template v-if="curItem.data && curItem.data.length > 0">
					<draggable v-model="curItem.data" group="people" item-key="index" class="draggable-list">
						<template #item="{ element,index }">
							<div class="d-c-c param-img-item navbar" v-if="index!=0">
								<div class="right">
									<el-icon class="el-icon-DeleteFilled"
										@click="editor.onEditorDeleleData(index, selectedIndex)">
										<CloseBold />
									</el-icon>
									<div class="url-box mb16 flex-1 d-s-c ww100">
										<span class="key-name">显示文字</span>
										<DiyInputField maxlength="6" show-word-limit v-model="element.text"></DiyInputField>
									</div>
									<div class="d-s-c  ww100">
										<div class="url-box flex-1 d-s-c">
											<span class="key-name">商品分类</span>
											<el-cascader class="ww100" v-model="element.currCategory"
												:ref="'cascader'+index" :options="CategoryList"
												:props="{ checkStrictly: true, children: 'child', value: 'category_id', label: 'name' }"
												@change="changeCategory($event,index)"></el-cascader>
										</div>
									</div>
								</div>
							</div>
						</template>
					</draggable>
				</template>
				<div class="d-c-c pb16"><component :is="PrimaryButton" plain @click="editor.onEditorAddData">+新增</component>
				</div>
			</div>
			<div v-if="styleType == 'style'">
				<div class="f16 gray3 form-subtitle">
					选项卡样式
				</div>
				<!--色调-->
				<!-- <div class="form-vben-item" label="色调：">
					<component :is="RadioGroup" v-model="curItem.style.themeType">
						<el-radio label="theme">跟随主题</el-radio>
						<el-radio label="custom">自定义</el-radio>
					</component>
				</div> -->
				<div class="form-item">
					<div class="form-label">装饰元素：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyInputField class="ml10" v-model="curItem.style.active_color1" placeholder="透明" />
						<DiyInputField v-if="curItem.params.type!=3" class="ml10" v-model="curItem.style.active_color2"
							placeholder="透明" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="ResetColor('active_color1', 'active_color2')">重置</component>
					</div>
				</div>
				<div class="form-item">
					<div class="form-label">选中文字：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.activeText" default-color="#ffffff" />
					</div>
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					通用样式
				</div>
				<div class="form-item">
					<div class="form-label">组件背景：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color1" placeholder="透明" />
						<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color2" placeholder="透明" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
							 @click.stop="ResetColor('bgcolor_color1', 'bgcolor_color2')">重置</component>
					</div>
				</div>
				<div class="form-item">
					<div class="form-label">底部背景：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.background" default-color="#ffffff" />
					</div>
				</div>
				<!--上下边距-->
				<div class="form-item">
					<div class="form-label">上边距：</div>
					<DiySliderField v-model="curItem.style.paddingTop" size="small" show-input :show-input-controls="false"
						input-size="small"></DiySliderField>
				</div>
				<!--上下边距-->
				<div class="form-item">
					<div class="form-label">下边距：</div>
					<DiySliderField v-model="curItem.style.paddingBottom" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<!--左右边距-->
				<div class="form-item">
					<div class="form-label">左右边距：</div>
					<DiySliderField v-model="curItem.style.paddingLeft" size="small" show-input :show-input-controls="false"
						input-size="small"></DiySliderField>
				</div>
				<!--页面上间距:-->
				<div class="form-item">
					<div class="form-label">页面上间距：</div>
					<DiySliderField v-model="curItem.style.marginTop" size="small" show-input :show-input-controls="false"
						input-size="small"></DiySliderField>
				</div>
				<!--上圆角-->
				<div class="form-item">
					<div class="form-label">上圆角：</div>
					<DiySliderField v-model="curItem.style.topRadio" size="small" show-input :show-input-controls="false"
						input-size="small"></DiySliderField>
				</div>
				<!--下圆角-->
				<div class="form-item">
					<div class="form-label">下圆角：</div>
					<DiySliderField v-model="curItem.style.bottomRadio" size="small" show-input :show-input-controls="false"
						input-size="small"></DiySliderField>
				</div>
			</div>

		</div>
	</div>
</template>

<style lang="scss" scoped></style>