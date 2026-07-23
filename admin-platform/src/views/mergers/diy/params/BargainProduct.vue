<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';

import { diyInput, diySlider } from './shared/schema-helpers';
import { changeBargainLikeColumn, resetStyleColors } from './shared/marketing-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyColorField from './shared/diy-color-field.vue';
import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsBargainProduct' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const styleType = ref<'content' | 'style'>('content');
function changeColumn(column: number) { changeBargainLikeColumn(props.curItem, column); }
function ResetColor(name1: string, name2?: string) { resetStyleColors(editor, props.curItem, name1, name2); }


const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'params.column', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.linePrice', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.more', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.productName', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.productPrice', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.productSales', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.product_btn', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.showNum', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.title', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.titleBgType', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.titleType', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.background', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.btn_text_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.marginTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.moreColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.moreSize', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.nameWeight', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingBottom', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingLeft', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productBottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productBtn_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productBtn_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productLine_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productName_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productPrice_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productTopRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_sales_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleBg_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleBg_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleSize', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.titleWeight', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.topRadio', hideLabel: true, formItemClass: 'hidden' },
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, { fieldPaths: ["params.column","params.linePrice","params.more","params.productName","params.productPrice","params.productSales","params.product_btn","params.showNum","params.title","params.titleBgType","params.titleType","style.background","style.bgcolor_color1","style.bgcolor_color2","style.bottomRadio","style.btn_text_color","style.marginTop","style.moreColor","style.moreSize","style.nameWeight","style.paddingBottom","style.paddingLeft","style.paddingTop","style.productBottomRadio","style.productBtn_color1","style.productBtn_color2","style.productLine_color","style.productName_color","style.productPrice_color","style.productTopRadio","style.product_sales_color","style.titleBg_color1","style.titleBg_color2","style.titleColor","style.titleSize","style.titleWeight","style.topRadio"] });
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
				<div class="f16 gray3 form-subtitle">头部设置</div>
				<div class="form-item">
					<div class="form-label">背景风格：</div>
					<component :is="RadioGroup" v-model="curItem.params.titleBgType" size="small">
						<el-radio :value="2">背景色</el-radio>
						<el-radio :value="1">图片</el-radio>
					</component>
				</div>
				<div class="form-item" v-if="curItem.params.titleBgType == 1">
					<div class="form-label">背景图片：</div>
					<div class="flex-1">
						<div class="d-s-c f-w">
							<div class="item">
								<img width="80" v-img-url="curItem.params.bgimage" alt=""
									@click="editor.onEditorSelectImage(curItem.params, 'bgimage')" />
							</div>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
								 @click.stop="editor.onEditorResetColor(curItem.params, 'bgimage', '')">重置</component>
						</div>

						<div class="gray f12">建议图片尺寸750px(根据组件左右边距修改)*90px</div>
					</div>
				</div>
				<div class="form-item">
					<div class="form-label">标题风格：</div>
					<component :is="RadioGroup" v-model="curItem.params.titleType" size="small">
						<el-radio :value="1">文字</el-radio>
						<el-radio :value="2">图片</el-radio>
					</component>
				</div>
				<div class="form-item" v-if="curItem.params.titleType == 1">
					<div class="form-label">标题文字：</div>
					<DiyInputField maxlength="6" show-word-limit v-model="curItem.params.title"></DiyInputField>
				</div>
				<div class="form-item" v-if="curItem.params.titleType == 2">
					<div class="form-label">标题图片：</div>
					<div class="flex-1">
						<div class="d-s-c f-w">
							<div class="item">
								<img width="85" height="45" v-img-url="curItem.params.titleimage" alt=""
									@click="editor.onEditorSelectImage(curItem.params, 'titleimage')" />
							</div>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
								 @click.stop="editor.onEditorResetColor(curItem.params, 'titleimage', '')">重置</component>
						</div>

						<div class="gray f12">建议图片尺寸170px*90px</div>
					</div>
				</div>
				<div class="form-item">
					<div class="form-label">按钮文字：</div>
					<DiyInputField maxlength="6" show-word-limit v-model="curItem.params.more"></DiyInputField>
				</div>
				<div class="diy-field-hint gray f12 px-4">
					标题和按钮文字颜色跟随主题色
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">商品列表</div>
				<div class="form-item">
					<div class="form-label">选择风格：</div>
					<component :is="RadioGroup" v-model="curItem.params.column" size="small" @change="changeColumn">
						<el-radio :value="1">单列展示</el-radio>
						<el-radio :value="2">两列展示(纵向)</el-radio>
						<el-radio :value="3">三列展示</el-radio>
						<el-radio :value="4">左右滑动展示</el-radio>
					</component>
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">商品数据</div>
				<div class="form-item">
					<div class="form-label">商品数量：</div>
					<DiySliderField v-model="curItem.params.showNum" size="small" show-input :show-input-controls="false"
						input-size="small"></DiySliderField>
				</div>
				<!-- 商品名称 -->
				<div class="form-item">
					<div class="form-label">显示内容：</div>
					<div class="flex-1">
						<component :is="Checkbox" v-model="curItem.params.productName" label="商品名称" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox" v-model="curItem.params.productPrice" label="砍价价格" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox" v-model="curItem.params.linePrice" label="商品原价" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox" v-if="curItem.params.column==1" v-model="curItem.params.productSales" label="参与人数"
							:true-value="1" :false-value="0" size="small" />
					</div>
				</div>
				<div class="form-item" v-if="curItem.params.column==1 || curItem.params.column==2">
					<div class="form-label">砍价按钮：</div>
					<component :is="RadioGroup" v-model="curItem.params.product_btn" size="small">
						<el-radio :value="1">显示</el-radio>
						<el-radio :value="0">隐藏</el-radio>
					</component>
				</div>
			</div>
			<div v-if="styleType == 'style'">
				<div class="f16 gray3 form-subtitle">头部样式</div>
				<div class="form-item" v-if="curItem.params.titleBgType == 2">
					<div class="form-label">背景颜色：</div>
					<div class="flex-1">
						<div class="d-s-c" style="height: 36px;">
							<DiyInputField class="ml10" v-model="curItem.style.titleBg_color1" placeholder="透明" />
							<view class="ml10"><DiyColorField v-model="curItem.style.titleBg_color1" default-color="#FF4C01" class="diy-inline-color" /></view>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
								 @click.stop="editor.onEditorResetColor(curItem.style, 'titleBg_color1', '#FF4C01')">重置</component>
						</div>
						<div class="d-s-c" style="height: 36px;">
							<DiyInputField class="ml10" v-model="curItem.style.titleBg_color2" placeholder="透明" />
							<view class="ml10"><DiyColorField v-model="curItem.style.titleBg_color2" default-color="#FF4C01" class="diy-inline-color" /></view>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
								 @click.stop="editor.onEditorResetColor(curItem.style, 'titleBg_color2', '#FF4C01')">重置</component>
						</div>
					</div>
				</div>
				<template v-if="curItem.params.titleType == 1">
					<div class="form-item">
						<div class="form-label">标题文字：</div>
						<component :is="RadioGroup" v-model="curItem.style.titleWeight" size="small">
							<el-radio :value="0">正常</el-radio>
							<el-radio :value="1">加粗</el-radio>
							<el-radio :value="2">倾斜</el-radio>
						</component>
					</div>
					<div class="form-item">
						<div class="form-label">标题字号：</div>
						<DiySliderField v-model="curItem.style.titleSize" :min="12" :max="24" size="small" show-input
							:show-input-controls="false" input-size="small"></DiySliderField>
					</div>
				</template>
				<div class="form-item">
					<div class="form-label">按钮字号：</div>
					<DiySliderField v-model="curItem.style.moreSize" :min="12" :max="40" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<!--组件样式-->
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">商品样式</div>
				<div class="form-item">
					<div class="form-label">上圆角：</div>
					<DiySliderField v-model="curItem.style.productTopRadio" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<div class="form-item">
					<div class="form-label">下圆角：</div>
					<DiySliderField v-model="curItem.style.productBottomRadio" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<template v-if="curItem.params.productName">
					<div class="form-item">
						<div class="form-label">商品名称：</div>

						<div class="flex-1 d-s-c" style="height: 36px;">
							<DiyColorField v-model="curItem.style.productName_color" default-color="#333333" placeholder="透明" />
						</div>
					</div>
					<div class="form-item">
						<div class="form-label">商品名称：</div>
						<component :is="RadioGroup" v-model="curItem.style.nameWeight" size="small">
							<el-radio :value="0">正常</el-radio>
							<el-radio :value="1">加粗</el-radio>
						</component>
					</div>
				</template>
				<div class="form-item" v-if="curItem.params.linePrice">
					<div class="form-label">商品原价：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.productLine_color" default-color="#999999" placeholder="透明" />
					</div>
				</div>
				<div class="form-item" v-if="curItem.params.productPrice">
					<div class="form-label">砍价价格：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.productPrice_color" default-color="#FF4C01" placeholder="透明" />
					</div>
				</div>
				<div class="form-item" v-if="curItem.params.productSales">
					<div class="form-label">参与人数：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.product_sales_color" default-color="#999999" placeholder="透明" />
					</div>
				</div>
				<template v-if="curItem.params.product_btn == 1">
					<!--购买按钮-->
					<div class="form-item mult">
						<div class="form-label">按钮颜色：</div>
						<div class="flex-1">
							<div class="d-s-c" style="height: 36px;">
								<DiyInputField class="ml10" v-model="curItem.style.productBtn_color1" placeholder="透明" />
								<view class="ml10"><DiyColorField v-model="curItem.style.productBtn_color1" default-color="#FF4C01" class="diy-inline-color" /></view>
								<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
									 @click.stop="editor.onEditorResetColor(curItem.style, 'productBtn_color1', '#FF4C01')">重置</component>
							</div>
							<div class="d-s-c" style="height: 36px;">
								<DiyInputField class="ml10" v-model="curItem.style.productBtn_color2" placeholder="透明" />
								<view class="ml10"><DiyColorField v-model="curItem.style.productBtn_color2" default-color="#FF4C01" class="diy-inline-color" /></view>
								<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
									 @click.stop="editor.onEditorResetColor(curItem.style, 'productBtn_color2', '#FF4C01')">重置</component>
							</div>
						</div>
					</div>
					<div class="form-item" v-if="curItem.params.column != 3">
						<div class="form-label">内容颜色：</div>
						<div class="flex-1 d-s-c" style="height: 36px;">
							<DiyColorField v-model="curItem.style.btn_text_color" default-color="#fff" placeholder="透明" />
						</div>
					</div>
				</template>
				<!--组件样式-->
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">组件样式</div>
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
				<div class="form-item">
					<div class="form-label">底部背景：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.background" default-color="#ffffff" />
						<DiyInputField class="ml10" v-model="curItem.style.background" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
							 @click.stop="editor.onEditorResetColor(curItem.style, 'background', '#ffffff')">重置</component>
					</div>
				</div>
				<div class="form-item">
					<div class="form-label">组件背景：</div>
					<div class="flex-1">
						<div class="d-s-c" style="height: 36px;">
							<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color1" placeholder="透明" />
							<view class="ml10"><DiyColorField v-model="curItem.style.bgcolor_color1" default-color="#FF4C01" class="diy-inline-color" /></view>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
								 @click.stop="editor.onEditorResetColor(curItem.style, 'bgcolor_color1', '#FF4C01')">重置</component>
						</div>
						<div class="d-s-c" style="height: 36px;">
							<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color2" placeholder="透明" />
							<view class="ml10"><DiyColorField v-model="curItem.style.bgcolor_color2" default-color="#FF4C01" class="diy-inline-color" /></view>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
								 @click.stop="editor.onEditorResetColor(curItem.style, 'bgcolor_color2', '#FF4C01')">重置</component>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style lang="scss" scoped></style>