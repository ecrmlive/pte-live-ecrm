<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import { getPageThemeApi } from '#/api/core/page';
import { getThemeColorByIndex } from '#/utils/shop-theme';

import { applyTitleStylePreset } from './shared/marketing-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyColorField from './shared/diy-color-field.vue';
import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsTitle' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const isLinkset = ref(false);
const linkName = ref('');

async function changeType(type: number | string) {
  const numericType = Number(type);
  const style = (props.curItem.style ??= {}) as Record<string, unknown>;
  style.type = numericType;
  applyTitleStylePreset(props.curItem, numericType);
  try {
    const res = await getPageThemeApi();
    const theme = res.vars?.values?.theme ?? '0';
    const themeColor = getThemeColorByIndex(theme);
    style.textColor = themeColor;
    style.lineColor = themeColor;
    style.moretextColor = themeColor;
    if (numericType === 7 || numericType === 8) {
      style.subtextColor = themeColor;
    }
  } catch {
    // keep preset colors when theme API unavailable
  }
}

function changeLink(name: string) {
  isLinkset.value = true;
  linkName.value = name;
}

function closeLinkset(e: { url?: string } | null) {
  isLinkset.value = false;
  if (e?.url) {
    (props.curItem.params as Record<string, unknown>)[linkName.value] = e.url;
  }
}

parseIntFields(props.curItem, ['style.paddingTop']);
const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'params.morelinkUrl', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.moretitle', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.show_icon', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.sublinkUrl', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.subtitle', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'params.title', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.background', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bgcolor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.bottomRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.isLine', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.isMore', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.isSub', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.lineColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.moretextColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingBottom', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingLeft', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.paddingTop', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.subbackground', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.subtextColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.subtextSize', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.textColor', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.textSize', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.topRadio', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.type', hideLabel: true, formItemClass: 'hidden' },
  { component: 'Input', fieldName: 'style.weight', hideLabel: true, formItemClass: 'hidden' },
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, { fieldPaths: ["params.morelinkUrl","params.moretitle","params.show_icon","params.sublinkUrl","params.subtitle","params.title","style.background","style.bgcolor","style.bottomRadio","style.isLine","style.isMore","style.isSub","style.lineColor","style.moretextColor","style.paddingBottom","style.paddingLeft","style.paddingTop","style.subbackground","style.subtextColor","style.subtextSize","style.textColor","style.textSize","style.topRadio","style.type","style.weight"] });
</script>

<template>
	<div>
		<div class="common-form">
			<span>{{ curItem.name }}</span>
		</div>
		<div class="diy-vben-params">
		<Form class="hidden" />
			<div class="f16 gray3 form-subtitle">风格设置</div>
			<div class="form-vben-item" label="风格：">
				<component :is="RadioGroup" @change="changeType" v-model="curItem.style.type">
					<el-radio-button label="1">风格一</el-radio-button>
					<el-radio-button label="2">风格二</el-radio-button>
					<el-radio-button label="3">风格三</el-radio-button>
					<el-radio-button label="4">风格四</el-radio-button>
					<el-radio-button label="5">风格五</el-radio-button>
					<el-radio-button label="6">风格六</el-radio-button>
					<el-radio-button label="7">风格七</el-radio-button>
					<el-radio-button label="8">风格八</el-radio-button>
				</component>
			</div>
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">标题内容</div>
			<div class="form-vben-item" label="标题名称："><DiyInputField v-model="curItem.params.title" class="w-auto"></DiyInputField></div>
			<div class="form-vben-item" label="链接：">
				<DiyInputField @click="changeLink('sublinkUrl')" v-model="curItem.params.sublinkUrl">
					<template #suffix>
						<el-icon color="#333" size="16px"><ArrowRight /></el-icon>
					</template>
				</DiyInputField>
			</div>
			<div class="form-vben-item"
				label="副标题名称："
				v-if="curItem.style.type == 4 || curItem.style.type == 5 || curItem.style.type == 6 || curItem.style.type == 7 || curItem.style.type == 8"
			>
				<DiyInputField v-model="curItem.params.subtitle" class="w-auto"></DiyInputField>
			</div>
			<div class="form-chink"></div>
			<template v-if="curItem.style.type == 8">
				<div class="f16 gray3 form-subtitle">"更多"按钮内容</div>
				<div class="form-vben-item" label="按钮文字："><DiyInputField v-model="curItem.params.moretitle" class="w-auto"></DiyInputField></div>
				<div class="form-vben-item" label="链接：">
					<DiyInputField @click="changeLink('morelinkUrl')" v-model="curItem.params.morelinkUrl">
						<template #suffix>
							<el-icon color="#333" size="16px"><ArrowRight /></el-icon>
						</template>
					</DiyInputField>
				</div>
				<div class="form-chink"></div>
			</template>

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
			<div class="form-chink"></div>
			<div class="f16 gray3 form-subtitle">标题样式</div>
			<!--下圆角-->
			<div class="form-item">
				<div class="form-label">文字大小：</div>
				<DiySliderField v-model="curItem.style.textSize" :min="12" :max="20" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
			</div>
			<div class="form-vben-item" label="文字加粗：">
				<component :is="RadioGroup" v-model="curItem.style.weight">
					<el-radio-button :label="400">标准</el-radio-button>
					<el-radio-button :label="800">加粗</el-radio-button>
				</component>
			</div>
			<div class="form-item">
				<div class="form-label">文字颜色：</div>
				<div class="flex-1 d-s-c" style="height: 36px;">
					<DiyColorField v-model="curItem.style.textColor" default-color="#ffffff" />
					<DiyInputField class="ml10" v-model="curItem.style.textColor" />
					<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="editor.onEditorResetColor(curItem.style, 'textColor', '#333333')">重置</component>
				</div>
			</div>
			<div class="form-vben-item" label="是否显示辅助图：" v-if="curItem.style.type == 8">
				<component :is="RadioGroup" v-model="curItem.style.isLine">
					<el-radio-button :label="1">是</el-radio-button>
					<el-radio-button :label="0">否</el-radio-button>
				</component>
			</div>
			<div class="form-item" v-if="curItem.style.type == 1 || curItem.style.type == 2 || curItem.style.type == 8">
				<div class="form-label">辅助图颜色：</div>
				<div class="flex-1 d-s-c" style="height: 36px;">
					<DiyColorField v-model="curItem.style.lineColor" default-color="#ffffff" />
					<DiyInputField class="ml10" v-model="curItem.style.lineColor" />
					<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="editor.onEditorResetColor(curItem.style, 'lineColor', '#eeeeee')">重置</component>
				</div>
			</div>
			<template v-if="curItem.style.type == 4 || curItem.style.type == 5 || curItem.style.type == 6 || curItem.style.type == 7 || curItem.style.type == 8">
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">副标题样式</div>
				<div class="form-vben-item" label="是否显示：" v-if="curItem.style.type == 8">
					<component :is="RadioGroup" v-model="curItem.style.isSub">
						<el-radio-button :label="1">是</el-radio-button>
						<el-radio-button :label="0">否</el-radio-button>
					</component>
				</div>
				<!--下圆角-->
				<div class="form-item">
					<div class="form-label">文字大小：</div>
					<DiySliderField v-model="curItem.style.subtextSize" :min="12" :max="40" size="small" show-input :show-input-controls="false" input-size="small"></DiySliderField>
				</div>
				<div class="form-item">
					<div class="form-label">文字颜色：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.subtextColor" default-color="#ffffff" />
						<DiyInputField class="ml10" v-model="curItem.style.subtextColor" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="editor.onEditorResetColor(curItem.style, 'subtextColor', '#DDDDDD')">重置</component>
					</div>
				</div>
				<div class="form-item" v-if="curItem.style.type == 7 || curItem.style.type == 8">
					<div class="form-label">背景颜色：</div>
					<DiyColorField v-model="curItem.style.subbackground" default-color="#ffffff" placeholder="透明" />
				</div>
			</template>
			<template v-if="curItem.style.type == 8">
				<div class="form-chink"></div>
				<div class="f16 gray3 form-subtitle">"更多"按钮样式</div>
				<div class="form-vben-item" label="是否显示：">
					<component :is="RadioGroup" v-model="curItem.style.isMore">
						<el-radio-button :label="1">是</el-radio-button>
						<el-radio-button :label="0">否</el-radio-button>
					</component>
				</div>
				<div class="form-item">
					<div class="form-label">文字颜色：</div>
					<div class="flex-1 d-s-c" style="height: 36px;">
						<DiyColorField v-model="curItem.style.moretextColor" default-color="#ffffff" />
						<DiyInputField class="ml10" v-model="curItem.style.moretextColor" />
						<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"  @click.stop="editor.onEditorResetColor(curItem.style, 'moretextColor', '#999999')">重置</component>
					</div>
				</div>
			</template>
			<!-- <div class="form-vben-item" label="图标显示：">
				<component :is="RadioGroup" v-model="curItem.params.show_icon">
					<el-radio label="yes">显示</el-radio>
					<el-radio label="no">不显示</el-radio>
				</component>
			</div>
			<div class="form-vben-item" label="标题图标：">
				<div class="diy-notice-icon">
					<img v-img-url="curItem.params.icon" alt="" style="width: 100%;height: auto;" @click="editor.onEditorSelectImage(curItem.params, 'icon')" />
				</div>
				<div class="ww100">建议尺寸32×32</div>
			</div> -->
			<!-- 公告内容 -->
		</div>
		<DiyLinkPickerDialog v-if="is_linkset" :is_linkset="is_linkset" @close-dialog="closeLinkset">选择链接</DiyLinkPickerDialog>
	</div>
</template>

<style lang="scss" scoped>
.diy-notice-icon,
.diy-notice-icon img {
	width: 32px;
	height: 32px;
}
</style>
