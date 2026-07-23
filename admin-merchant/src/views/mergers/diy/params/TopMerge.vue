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

defineOptions({ name: 'DiyParamsTopMerge' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

import { getShopLinkProductCategoryApi } from '#/api/core/shop-link';
import draggable from 'vuedraggable';
import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

const styleType = ref<'content' | 'style'>('content');
const loading = ref(true);
const CategoryList = ref<Array<Record<string, unknown>>>([]);
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<Record<string, unknown> | null>(null);

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

function changeLink(index: number) {
  linkIndex.value = index;
  linkData.value = (props.curItem.images as Array<Record<string, unknown>>)[index] ?? null;
  isLinkset.value = true;
}

function closeLinkset(e: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  if (e && props.curItem.images) {
    const row = (props.curItem.images as Array<Record<string, unknown>>)[linkIndex.value];
    row.linkeType = e.type;
    row.linkUrl = e.url;
    row.name = e.name;
  }
}

void getData();


const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'params.data', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.images', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.searchText', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.showCategory', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.topUp', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.type', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor_color1', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.btnColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.btnOpColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.btnShape', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.categoryPadding', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.imgShape', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.themeType', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.topRadio', hideLabel: true, formItemClass: 'hidden' },
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, { fieldPaths: ["params.data","params.images","params.searchText","params.showCategory","params.topUp","params.type","style.bgcolor_color1","style.bottomRadio","style.btnColor","style.btnOpColor","style.btnShape","style.categoryPadding","style.imgShape","style.themeType","style.topRadio"] });
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
					</component>
				</div>
				<div class="form-vben-item" label="分类设置：">
					<component :is="RadioGroup" v-model="curItem.params.showCategory">
						<el-radio :label="true">显示</el-radio>
						<el-radio :label="false">隐藏</el-radio>
					</component>
				</div>
				<div class="form-vben-item" label="滑动置顶：">
					<component :is="RadioGroup" v-model="curItem.params.topUp">
						<el-radio label="1">开启</el-radio>
						<el-radio label="0">关闭</el-radio>
					</component>
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					搜索设置
				</div>
				<!--图片-->
				<div class="form-vben-item" label="默认logo图：">
					<div class="diy-special-cover">
						<img v-img-url="curItem.params.topLogo" alt=""
							@click="editor.onEditorSelectImage(curItem.params, 'topLogo')" />
						<div class="gray">建议尺寸78*64</div>
					</div>
				</div>
				<!--图片-->
				<div class="form-vben-item" label="顶部固定logo：" v-if="curItem.params.topUp == 1">
					<div class="diy-special-cover">
						<img v-img-url="curItem.params.fixedLogo" alt=""
							@click="editor.onEditorSelectImage(curItem.params, 'fixedLogo')" />
						<div class="gray">建议尺寸78*64</div>
					</div>
				</div>
				<div class="form-vben-item" label="提示文字：">
					<DiyInputField v-model="curItem.params.searchText" class="w-auto"></DiyInputField>
				</div>
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
				<div class="d-c-c pb16"><component :is="PrimaryButton" plain @click="editor.onEditorAddData">+添加</component>
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					图片设置
					<span class="gray f12">建议上传尺寸相同的图片，建议尺寸750px*340px</span>
				</div>
				<template v-if="curItem.images && curItem.images.length > 0">
					<draggable v-model="curItem.images" group="people" item-key="index" class="draggable-list">
						<template #item="{ element,index }">
							<div class="d-c-c param-img-item navbar">
								<el-icon class="el-icon-DeleteFilled" color="#fff"
									@click="editor.onEditorDeleleImg(index, selectedIndex)">
									<CloseBold />
								</el-icon>
								<div class="right pr">
									<div class="icon">
										<img :style="{ height: '170px' }" v-img-url="element.imgUrl" alt=""
											@click="editor.onEditorSelectImage(element, 'imgUrl')" />
									</div>
									<div class="d-s-c  ww100">
										<div class="url-box flex-1 d-s-c">
											<span class="key-name">链接</span>
											<DiyInputField v-model="element.linkUrl" style="padding-bottom: 10px;">
												<template #suffix>
													<el-icon @click="changeLink(index)" color="#333" size="16px">
														<ArrowRight />
													</el-icon>
												</template>
											</DiyInputField>
										</div>
									</div>
								</div>
							</div>
						</template>
					</draggable>
				</template>
				<div class="d-c-c pb16"><component :is="PrimaryButton" plain
						@click="editor.onEditorAddImg">+添加一个</component>
				</div>
			</div>
			<div v-if="styleType == 'style'">
				<div class="f16 gray3 form-subtitle">
					选项卡样式
				</div>
				<!--内容间距-->
				<div class="form-item">
					<div class="form-label">内容间距：</div>
					<DiySliderField v-model="curItem.style.categoryPadding" size="small" show-input
						:show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					指示器设置
				</div>
				<div class="form-vben-item" label="指示点形状：">
					<component :is="RadioGroup" v-model="curItem.style.imgShape" size="medium">
						<el-radio-button label="round">圆形</el-radio-button>
						<el-radio-button label="square">正方形</el-radio-button>
						<el-radio-button label="rectangle">长方形</el-radio-button>
					</component>
				</div>
				<div class="form-vben-item" label="指示点位置：">
					<component :is="RadioGroup" v-model="curItem.style.btnShape">
						<el-radio-button label="left">居左</el-radio-button>
						<el-radio-button label="center">居中</el-radio-button>
						<el-radio-button label="right">居右</el-radio-button>
					</component>
				</div>
				<!-- <div class="form-vben-item" label="色调：">
					<component :is="RadioGroup" v-model="curItem.style.themeType">
						<el-radio label="theme">跟随主题</el-radio>
						<el-radio label="custom">自定义</el-radio>
					</component>
				</div> -->
				<div class="form-item">
					<div class="form-label">选中样式：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.btnColor" default-color="#ffffff" />
						<DiyInputField class="ml10" v-model="curItem.style.btnColor" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
							 @click.stop="editor.onEditorResetColor(curItem.style, 'btnColor', '#ffffff')">重置</component>
					</div>
				</div>
				<div class="form-item">
					<div class="form-label">常规样式：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.btnOpColor" default-color="#ffffff" />
						<DiyInputField class="ml10" v-model="curItem.style.btnOpColor" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
							 @click.stop="editor.onEditorResetColor(curItem.style, 'btnOpColor', '#ffffff')">重置</component>
					</div>
				</div>
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					图片设置
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
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">
					组件设置
				</div>
				<div class="form-item">
					<div class="form-label">组件背景：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyInputField class="ml10" v-model="curItem.style.bgcolor_color1" placeholder="透明" />
						<DiyColorField v-model="curItem.style.bgcolor_color1" default-color="#ffffff" class="ml10" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"
							 @click.stop="editor.onEditorResetColor(curItem.style, 'bgcolor_color1', '#ffffff')">重置</component>
					</div>
				</div>
			</div>
		</div>
		<DiyLinkPickerDialog v-if="is_linkset" :is_linkset="is_linkset" :linkData='linkData' @close-dialog="closeLinkset">选择链接
		</DiyLinkPickerDialog>
	</div>
</template>

<style lang="scss" scoped>
	.diy-special-cover {
		img {
			width: 39px;
			height: 32px;
		}
	}

	.param-img-item.navbar {
		min-height: 132px;
		height: auto;
	}

	.param-img-item.navbar .icon img {
		width: 408px;
		height: 202px;
		background: #eeeeee;
		margin-bottom: 10px;
	}
</style>