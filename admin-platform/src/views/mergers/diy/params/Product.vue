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

defineOptions({ name: 'DiyParamsProduct' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

import { getShopLinkProductCategoryApi } from '#/api/core/shop-link';
import draggable from 'vuedraggable';
import { changeProductColumn, resetStyleColors } from './shared/marketing-helpers';

const styleType = ref<'content' | 'style'>('content');
const loading = ref(true);
const CategoryList = ref<Array<Record<string, unknown>>>([]);
const currCategory = ref<Array<number | string>>([]);

function changeColumn(column: number) {
  changeProductColumn(props.curItem, column);
}

function ResetColor(titleColor1: string, titleColor2?: string) {
  editor.onEditorResetColor(props.curItem.style as Record<string, unknown>, titleColor1, '#ffffff');
  if (titleColor2) {
    editor.onEditorResetColor(props.curItem.style as Record<string, unknown>, titleColor2, '#ffffff');
  }
}

function currCategoryAuto(list: Array<Record<string, unknown>>) {
  const arr: Array<number | string> = [];
  const target = (props.curItem.params as Record<string, unknown>)?.auto as Record<string, unknown>;
  const categoryId = target?.category;
  for (const item of list) {
    if (item.category_id === categoryId) {
      arr.push(item.category_id as number | string);
      break;
    }
    if (Array.isArray(item.child) && item.child.length > 0) {
      for (const child of item.child as Array<Record<string, unknown>>) {
        if (child.category_id === categoryId) {
          arr.push(item.category_id as number | string);
          arr.push(child.category_id as number | string);
          break;
        }
      }
    }
  }
  return arr;
}

async function getData() {
  loading.value = true;
  try {
    const res = await getShopLinkProductCategoryApi();
    CategoryList.value = res.list ?? [];
    currCategory.value = currCategoryAuto(CategoryList.value);
  } finally {
    loading.value = false;
  }
}

function changeCategory(_e: unknown, cascaderRef: { getCheckedNodes: () => Array<{ data: { category_id: number } }> }) {
  const item = cascaderRef.getCheckedNodes();
  const auto = (props.curItem.params as Record<string, unknown>).auto as Record<string, unknown>;
  auto.category = item[0]?.data.category_id;
}

void getData();


const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'params.data', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.auto.productSort', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.auto.showNum', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.cartText', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.cartType', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.column', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.comment', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.linePrice', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.productName', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.productPrice', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.productSales', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.showCart', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.source', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.background', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.cart_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.cart_color2', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.cart_text_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.line_price_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.marginTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.nameWeight', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingBottom', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingLeft', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productBottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.productTopRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_comment_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_name_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_price_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.product_sales_color', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.topRadio', hideLabel: true, formItemClass: 'hidden' },
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, { fieldPaths: ["params.data","params.auto.productSort","params.auto.showNum","params.cartText","params.cartType","params.column","params.comment","params.linePrice","params.productName","params.productPrice","params.productSales","params.showCart","params.source","style.background","style.bgcolor_color1","style.bgcolor_color2","style.bottomRadio","style.cart_color1","style.cart_color2","style.cart_text_color","style.line_price_color","style.marginTop","style.nameWeight","style.paddingBottom","style.paddingLeft","style.paddingTop","style.productBottomRadio","style.productTopRadio","style.product_comment_color","style.product_name_color","style.product_price_color","style.product_sales_color","style.topRadio"] });
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
				<!--列表设置-->
				<div class="f16 gray3 form-subtitle">列表设置</div>
				<div class="form-item">
					<div class="form-label">选择风格：</div>
					<component :is="RadioGroup" v-model="curItem.params.column" size="small" @change="changeColumn">
						<el-radio :value="1">单列展示</el-radio>
						<el-radio :value="2">两列展示(纵向)</el-radio>
						<el-radio :value="3">三列展示</el-radio>
						<el-radio :value="4">两列展示(横向)</el-radio>
						<el-radio :value="5">大图展示</el-radio>
						<el-radio :value="6">左右滑动展示</el-radio>
					</component>
				</div>
				<!--商品设置-->
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">商品设置</div>
				<div class="form-item">
					<div class="form-label">商品来源：</div>
					<component :is="RadioGroup" v-model="curItem.params.source" size="small">
						<el-radio :value="'auto'">自动获取</el-radio>
						<el-radio :value="'choice'">手动选择</el-radio>
					</component>
				</div>
				<!-- 自动获取 -->
				<template v-if="curItem.params.source == 'auto'">
					<!-- 商品分类 -->
					<div class="form-item" v-if="CategoryList && CategoryList.length > 0">
						<div class="form-label">商品分类：</div>
						<el-cascader class="ww100" v-model="currCategory" ref="cascader" :options="CategoryList"
							:props="{ checkStrictly: true, children: 'child', value: 'category_id', label: 'name' }"
							@change="(e) => changeCategory(e, $refs.cascader as any)"></el-cascader>
					</div>
					<!-- 显示数量 -->
					<div class="form-item">
						<div class="form-label">显示数量：</div>
						<DiySliderField v-model="curItem.params.auto.showNum" size="small" show-input :max="50"
							:show-input-controls="false" input-size="small"></DiySliderField>
					</div>
					<div class="form-item">
						<div class="form-label">商品排序：</div>
						<component :is="RadioGroup" v-model="curItem.params.auto.productSort" size="small">
							<el-radio :value="'all'">综合</el-radio>
							<el-radio :value="'sales'">销量</el-radio>
							<el-radio :value="'price'">价格</el-radio>
						</component>
					</div>
				</template>
				<!-- 手动选择 -->
				<template v-if="curItem.params.source == 'choice'">
					<div class="form-item">
						<div class="form-label">商品列表：</div>
						<div class="flex-1">
							<draggable v-model="curItem.data" item-key="index"
								:options="{ draggable: '.item', animation: 500 }" class="choice-product-list">
								<template #item="{ element, index }">
									<div class="d-s-c f-w">
										<div class="item">
											<div class="delete-box">
												<el-icon :size="20"
													@click="editor.onEditorDeleleData(index, selectedIndex)">
													<CircleCloseFilled />
												</el-icon>
											</div>
											<img v-img-url="element.image" alt="" />
										</div>
									</div>
								</template>
							</draggable>
							<div><component :is="PrimaryButton" icon="Plus"
									@click.stop="editor.openProduct(curItem.data, true)">选择商品</component>
							</div>
						</div>
					</div>
				</template>
				<!--组件样式-->
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">显示内容</div>
				<!-- 商品名称 -->
				<div class="form-item">
					<div class="form-label">显示内容：</div>
					<div class="flex-1">
						<component :is="Checkbox" v-model="curItem.params.productName" label="商品名称" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox" v-model="curItem.params.productPrice" label="商品价格" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox"
							v-if="curItem.params.column!=3 &&curItem.params.column!=4&&curItem.params.column!=6"
							v-model="curItem.params.productSales" label="商品销量" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox"
							v-if="curItem.params.column!=3 &&curItem.params.column!=4&&curItem.params.column!=6"
							v-model="curItem.params.linePrice" label="划线价格" :true-value="1" :false-value="0"
							size="small" />
						<component :is="Checkbox"
							v-if="curItem.params.column!=2 && curItem.params.column!=3 && curItem.params.column!=4 &&curItem.params.column!=6"
							v-model="curItem.params.comment" label="好评率" :true-value="1" :false-value="0"
							size="small" />
					</div>
				</div>
				<template v-if="curItem.params.column!=4">
					<div class="form-chink"></div>
					<div class="f16 gray3 form-subtitle">购物车按钮</div>
					<div class="form-item">
						<div class="form-label">是否显示：</div>
						<component :is="RadioGroup" v-model="curItem.params.showCart" size="small">
							<el-radio :value="1">显示</el-radio>
							<el-radio :value="0">隐藏</el-radio>
						</component>
					</div>
					<template v-if="curItem.params.showCart == 1">
						<div class="form-item">
							<div class="form-label">按钮样式：</div>
							<component :is="RadioGroup" v-model="curItem.params.cartType" size="small">
								<el-radio :value="0">
									<div class="cart-btn" :style="{color:'#fff'}">
										{{curItem.params.cartText || '购买'}}
									</div>
								</el-radio>
								<el-radio :value="1">
									<div class="cart-btn icon">
										<span class=" icon iconfont icon-icozhuanhuan" :style="{color:'#fff'}"></span>
									</div>
								</el-radio>
								<el-radio :value="2">
									<div class="cart-btn icon">
										<span class=" icon iconfont icon-jia" :style="{color:'#fff'}"></span>
									</div>
								</el-radio>
							</component>
						</div>
						<div class="form-item" v-if="curItem.params.cartType == 0">
							<div class="form-label">按钮文字：</div>
							<DiyInputField maxlength="6" show-word-limit v-model="curItem.params.cartText"></DiyInputField>
						</div>
					</template>


				</template>

			</div>

			<div v-if="styleType == 'style'">
				<!--商品样式-->
				<div class="f16 gray3 form-subtitle">商品样式</div>
				<!--上圆角-->
				<div class="form-item">
					<div class="form-label">上圆角：</div>
					<DiySliderField v-model="curItem.style.productTopRadio" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<!--下圆角-->
				<div class="form-item">
					<div class="form-label">下圆角：</div>
					<DiySliderField v-model="curItem.style.productBottomRadio" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<div class="form-item" v-if="curItem.params.productName">
					<div class="form-label">商品名称：</div>
					<component :is="RadioGroup" v-model="curItem.style.nameWeight" size="small">
						<el-radio :value="0">正常</el-radio>
						<el-radio :value="1">加粗</el-radio>
					</component>
				</div>
				<div class="form-item" v-if="curItem.params.productName">
					<div class="form-label">商品名称：</div>
					<DiyColorField v-model="curItem.style.product_name_color" default-color="#ffffff" placeholder="透明" />
				</div>
				<div class="form-item" v-if="curItem.params.productPrice">
					<div class="form-label">商品价格：</div>
					<DiyColorField v-model="curItem.style.product_price_color" default-color="#ffffff" placeholder="透明" />
				</div>
				<div class="form-item" v-if="curItem.params.linePrice">
					<div class="form-label">划线价格：</div>
					<DiyColorField v-model="curItem.style.line_price_color" default-color="#ffffff" placeholder="透明" />
				</div>
				<div class="form-item" v-if="curItem.params.comment">
					<div class="form-label">好评率：</div>
					<DiyColorField v-model="curItem.style.product_comment_color" default-color="#ffffff" placeholder="透明" />
				</div>
				<div class="form-item" v-if="curItem.params.productSales">
					<div class="form-label">已售数量：</div>
					<DiyColorField v-model="curItem.style.product_sales_color" default-color="#ffffff" placeholder="透明" />
				</div>
				<!--组件样式-->
				<template v-if="curItem.params.showCart == 1">
					<div class="form-chink"></div>
					<div class="f16 gray3 form-subtitle">购物车按钮</div>
					<div class="form-item">
						<div class="form-label">按钮颜色：</div>
						<div class="flex-1 d-s-c" style="height: 36px;">
							<DiyInputField class="ml10" v-model="curItem.style.cart_color1" placeholder="透明" />
							<DiyInputField class="ml10" v-model="curItem.style.cart_color2" placeholder="透明" />
							<view class="ml10"><DiyColorField v-model="curItem.style.cart_color1" default-color="#FF4C01" class="diy-inline-color" /></view>
							<view class="ml10"><DiyColorField v-model="curItem.style.cart_color2" default-color="#FF4C01" class="diy-inline-color" /></view>
							<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="ResetColor('cart_color1', 'cart_color2')">重置</component>
						</div>
					</div>
					<div class="form-item">
						<div class="form-label">内容颜色：</div>
						<DiyColorField v-model="curItem.style.cart_text_color" default-color="#ffffff" placeholder="透明" />
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
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color1" placeholder="透明" />
						<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color2" placeholder="透明" />
						<view class="ml10"><DiyColorField v-model="curItem.style.bgcolor_color1" default-color="#FF4C01" class="diy-inline-color" /></view>
						<view class="ml10"><DiyColorField v-model="curItem.style.bgcolor_color2" default-color="#FF4C01" class="diy-inline-color" /></view>
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
							 @click.stop="ResetColor('bgcolor_color1', 'bgcolor_color2')">重置</component>
					</div>
				</div>
			</div>


		</div>
	</div>
</template>

<style lang="scss" scoped>
	.choice-product-list {
		display: flex;
		justify-content: flex-start;
		flex-wrap: wrap;
		padding: 20px 0;
	}

	.choice-product-list .item {
		position: relative;
		width: 80px;
		height: 80px;
		margin-right: 10px;
		margin-bottom: 10px;
		border: 1px solid #dddddd;
	}

	.choice-product-list .item .delete-box {
		position: absolute;
		width: 20px;
		height: 20px;
		top: -10px;
		right: -10px;
		font-size: 20px;
		cursor: pointer;
		color: #999999;
	}

	.choice-product-list .item .delete-box:hover {
		color: rgb(255, 51, 0);
	}

	.choice-product-list .item.plus-btn {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	.choice-product-list .item.plus-btn>i {
		font-size: 30px;
		color: #cccccc;
	}

	.choice-product-list img {
		width: 100%;
		height: 100%;
	}

	.cart-btn {
		min-width: 53px;
		height: 23px;
		line-height: 23px;
		padding: 0 10px;
		box-sizing: border-box;
		background: #409eff;
		color: #fff;
		font-size: 12px;
		color: #fff;
		display: flex;
		justify-content: center;
		align-items: center;
		border-radius: 100px;
	}

	.cart-btn.icon {
		padding: 0;
		min-width: 23px;
		width: 23px;
		height: 23px;
		line-height: 23px;
		border-radius: 50%;

		.iconfont.icon {
			font-size: 12px;
			line-height: 1;
		}
	}
</style>