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

defineOptions({ name: 'DiyParamsPreview' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const styleType = ref<'content' | 'style'>('content');


const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'params.more', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.showNum', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.title', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.background', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgTag', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.moreColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.moreSize', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingBottom', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingLeft', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productBg_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productLine_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productName_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productPrice_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_imgRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_lineprice', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_name', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_price', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_tag', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_topRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.tagColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleSize', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleType', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.topRadio', hideLabel: true, formItemClass: 'hidden' },
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, { fieldPaths: ["params.more","params.showNum","params.title","style.background","style.bgTag","style.bgcolor","style.bottomRadio","style.moreColor","style.moreSize","style.paddingBottom","style.paddingLeft","style.paddingTop","style.productBg_color","style.productLine_color","style.productName_color","style.productPrice_color","style.product_bottomRadio","style.product_imgRadio","style.product_lineprice","style.product_name","style.product_price","style.product_tag","style.product_topRadio","style.tagColor","style.titleColor","style.titleSize","style.titleType","style.topRadio"] });
</script>

<template>

	<div>
		<div class="common-form">
			<span>{{ curItem.name }}</span>
		</div>
		<div class="diy-vben-params">
		<Form class="hidden" />
			<div class="f16 gray3 form-subtitle">商品数据</div>
			<!--商品数量-->
			<div class="form-item">
				<div class="form-label">商品数量：</div>
				<DiySliderField v-model="curItem.params.showNum" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">头部风格</div>
			<div class="form-vben-item" label="主标题类型：">
				<component :is="RadioGroup" v-model="curItem.style.titleType" size="medium">
					<el-radio-button :label="1">文字</el-radio-button>
					<el-radio-button :label="2">图片</el-radio-button>
				</component>
			</div>
			<!--图片-->
			<div class="form-vben-item" label="标题图片：" v-if="curItem.style.titleType == 2">
				<div class="diy-special-cover">
					<img  style="width: 220px;" v-img-url="curItem.style.title_image" alt="" @click="editor.onEditorSelectImage(curItem.style, 'title_image')" />
					<!-- <div class="gray9 f12">建议图片高度88px</div> -->
				</div>
			</div>
			<div class="form-vben-item" label="标题文字：" v-if="curItem.style.titleType == 1"><DiyInputField v-model="curItem.params.title" class="w-auto"></DiyInputField></div>
			<!--文字大小-->
			<div class="form-item" v-if="curItem.style.titleType == 1">
				<div class="form-label">标题文字大小：</div>
				<DiySliderField v-model="curItem.style.titleSize" :min="12" :max="24" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-vben-item" label="右侧文字："><DiyInputField v-model="curItem.params.more" class="w-auto"></DiyInputField></div>
			<div class="diy-field-hint gray f12 px-4">
				标题和按钮文字颜色跟随主题色
			</div>
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">右侧文字样式</div>
			<!--文字大小-->
			<div class="form-item">
				<div class="form-label">文字大小：</div>
				<DiySliderField v-model="curItem.style.moreSize" :min="12" :max="40" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">显示内容</div>
			<!--商品名称-->
			<div class="form-vben-item" label="商品名称：">
				<component :is="RadioGroup" v-model="curItem.style.product_name" size="medium">
					<el-radio-button :label="1">显示</el-radio-button>
					<el-radio-button :label="2">不显示</el-radio-button>
				</component>
			</div>
			<!--销售价-->
			<div class="form-vben-item" label="销售价：">
				<component :is="RadioGroup" v-model="curItem.style.product_price" size="medium">
					<el-radio-button :label="1">显示</el-radio-button>
					<el-radio-button :label="2">不显示</el-radio-button>
				</component>
			</div>
			<!--划线价-->
			<div class="form-vben-item" label="划线价：">
				<component :is="RadioGroup" v-model="curItem.style.product_lineprice" size="medium">
					<el-radio-button :label="1">显示</el-radio-button>
					<el-radio-button :label="2">不显示</el-radio-button>
				</component>
			</div>
			<!--划线价-->
			<div class="form-vben-item" label="标签：">
				<component :is="RadioGroup" v-model="curItem.style.product_tag" size="medium">
					<el-radio-button :label="1">显示</el-radio-button>
					<el-radio-button :label="2">不显示</el-radio-button>
				</component>
			</div>
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">头部样式</div>
			<div class="form-vben-item" label="背景图片：">
				<div class="diy-special-cover">
					<img   style="width: 220px;" v-img-url="curItem.style.bgimage" alt="" @click="editor.onEditorSelectImage(curItem.style, 'bgimage')" />
					<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="editor.onEditorResetColor(curItem.style, 'bgimage', '')">重置</component>
					<div class="gray f12">建议图片宽度750px(根据组件左右边距修改)高度88px</div>
				</div>
			</div>
			<!--组件样式-->
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">商品样式</div>
			<!--图片圆角-->
			<div class="form-item">
				<div class="form-label">图片圆角：</div>
				<DiySliderField v-model="curItem.style.product_imgRadio" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-item">
				<div class="form-label">商品背景：</div>
				<DiyColorField v-model="curItem.style.productBg_color" default-color="#ffffff" placeholder="透明" />
			</div>
			<!--商品上圆角-->
			<div class="form-item">
				<div class="form-label">上圆角：</div>
				<DiySliderField v-model="curItem.style.product_topRadio" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<!--商品下圆角-->
			<div class="form-item">
				<div class="form-label">下圆角：</div>
				<DiySliderField v-model="curItem.style.product_bottomRadio" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-item" v-if="curItem.style.product_name">
				<div class="form-label">商品名称：</div>
				<DiyColorField v-model="curItem.style.productName_color" default-color="#ffffff" placeholder="透明" />
			</div>
			<div class="form-item" v-if="curItem.style.product_price">
				<div class="form-label">销售价：</div>
				<DiyColorField v-model="curItem.style.productPrice_color" default-color="#ffffff" placeholder="透明" />
			</div>
			<div class="form-item" v-if="curItem.style.product_lineprice">
				<div class="form-label">划线价：</div>
				<DiyColorField v-model="curItem.style.productLine_color" default-color="#ffffff" placeholder="透明" />
			</div>
			<div class="form-item" v-if="curItem.style.product_tag">
				<div class="form-label">标签颜色：</div>
				<DiyColorField v-model="curItem.style.tagColor" default-color="#ffffff" placeholder="透明" />
			</div>
			<div class="form-item" v-if="curItem.style.product_tag">
				<div class="form-label">标签背景：</div>
				<DiyColorField v-model="curItem.style.bgTag" default-color="#ffffff" placeholder="透明" />
			</div>
			<!--组件样式-->
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">组件样式</div>
			<!--上下边距-->
			<div class="form-item">
				<div class="form-label">上边距：</div>
				<DiySliderField v-model="curItem.style.paddingTop" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<!--上下边距-->
			<div class="form-item">
				<div class="form-label">下边距：</div>
				<DiySliderField v-model="curItem.style.paddingBottom" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<!--左右边距-->
			<div class="form-item">
				<div class="form-label">左右边距：</div>
				<DiySliderField v-model="curItem.style.paddingLeft" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<!--上圆角-->
			<div class="form-item">
				<div class="form-label">上圆角：</div>
				<DiySliderField v-model="curItem.style.topRadio" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<!--下圆角-->
			<div class="form-item">
				<div class="form-label">下圆角：</div>
				<DiySliderField v-model="curItem.style.bottomRadio" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-item">
				<div class="form-label">底部背景：</div>
				<DiyColorField v-model="curItem.style.bgcolor" default-color="#ffffff" placeholder="透明" />
			</div>
			<div class="form-item">
				<div class="form-label">组件背景：</div>
				<DiyColorField v-model="curItem.style.background" default-color="#ffffff" placeholder="透明" />
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped></style>
