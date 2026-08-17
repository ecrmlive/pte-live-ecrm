<template>
	<div v-if="typeList != null" class="diy-type-list diy-component-list">
		<div class="min-group">
			<section
				v-for="(group, key) in typeList"
				:key="key"
				:name="key"
				class="diy-type-group"
			>
				<div class="native-section-title diy-type-group__title">{{ typename(key) }}</div>
				<div class="bd">
					<div
						v-for="(item, index) in group.children"
						:key="index"
						class="item"
						@click="onAdd(item.type)"
					>
						<p class="p-icon icon iconfont" :class="componentIcon(item)"></p>
						<p class="p-text">{{ item.name }}</p>
					</div>
				</div>
			</section>
		</div>
	</div>
</template>

<script>
const GROUPS = [
	{ key: 'media', name: '媒体组件' },
	{ key: 'shop', name: '商城组件' },
	{ key: 'marketing', name: '营销组件' },
	{ key: 'tools', name: '工具组件' },
	{ key: 'page', name: '页面组件' },
];
const HIDDEN_DIY_TYPES = new Set(['wxlive']);

/**
 * 仅为历史数据或新增组件缺失 icon/group 时兜底。
 * 正常情况下始终使用接口下发的真实 name、icon 和 group，避免前端映射
 * 覆盖组件本身的语义或把实际组件从组件库中漏掉。
 */
const COMPONENT_FALLBACKS = {
	discountGroup: { group: 'marketing', icon: 'icon-zhekou' },
	community: { group: 'marketing', icon: 'icon-huodongtuiguang' },
	bottomNav: { group: 'tools', icon: 'icon-daohang' },
	ranking: { group: 'shop', icon: 'icon-paihangbang' },
};
const DEFAULT_DIY_COMPONENT_ICON = 'icon-yingyong';

export default {
	data() {
		return {
			/*类别列表*/
			typeList: null,
			activeName: 0,
		};
	},
	emits: ['add-item'],
	props: {
		defaultData: Object,
	},
	created() {
		this.init();
	},
	filters: {},
	methods: {
		componentIcon(item) {
			return item?.icon || COMPONENT_FALLBACKS[item?.type]?.icon || DEFAULT_DIY_COMPONENT_ICON;
		},
		componentGroup(item) {
			return item?.group || COMPONENT_FALLBACKS[item?.type]?.group || 'tools';
		},
		onAdd(type) {
			this.$emit('add-item', type);
		},
		/*初始化数据*/
		init() {
			const groupedItems = {};
			for (const key in this.defaultData) {
				const item = this.defaultData[key];
				if (item?.type && !HIDDEN_DIY_TYPES.has(item.type) && !HIDDEN_DIY_TYPES.has(key)) {
					const group = this.componentGroup(item);
					if (!groupedItems[group]) {
						groupedItems[group] = { children: [] };
					}
					groupedItems[group].children.push(item);
				}
			}

			const orderedList = {};
			for (const group of GROUPS) {
				if (groupedItems[group.key]) {
					orderedList[group.key] = groupedItems[group.key];
				}
			}
			for (const groupKey in groupedItems) {
				if (!orderedList[groupKey]) {
					orderedList[groupKey] = groupedItems[groupKey];
				}
			}
			this.typeList = orderedList;
		},
		typename(type) {
			return GROUPS.find((group) => group.key === type)?.name ?? '其他组件';
		},
	},
};
</script>

<style scoped>
.diy-type-group + .diy-type-group {
	margin-top: 4px;
}

.diy-type-group__title {
	margin-top: 16px;
	margin-bottom: 10px;
}

.diy-type-group:first-child .diy-type-group__title {
	margin-top: 0;
}
</style>
