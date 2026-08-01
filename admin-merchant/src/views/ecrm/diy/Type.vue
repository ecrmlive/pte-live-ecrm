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
						<p class="p-icon icon iconfont" :class="item.icon"></p>
						<p class="p-text">{{ item.name }}</p>
					</div>
				</div>
			</section>
		</div>
	</div>
</template>

<script>
const GROUP_ORDER = ['media', 'tools', 'shop', 'page'];
const HIDDEN_DIY_TYPES = new Set(['wxlive']);

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
		onAdd(type) {
			this.$emit('add-item', type);
		},
		/*初始化数据*/
		init() {
			const tempList = {};
			for (const key in this.defaultData) {
				const item = this.defaultData[key];
				if (!item || HIDDEN_DIY_TYPES.has(item.type) || HIDDEN_DIY_TYPES.has(key)) {
					continue;
				}
				if (!tempList[item.group]) {
					tempList[item.group] = { children: [] };
				}
				tempList[item.group].children.push(item);
			}

			const orderedList = {};
			for (const groupKey of GROUP_ORDER) {
				if (tempList[groupKey]) {
					orderedList[groupKey] = tempList[groupKey];
				}
			}
			for (const groupKey in tempList) {
				if (!orderedList[groupKey]) {
					orderedList[groupKey] = tempList[groupKey];
				}
			}
			this.typeList = orderedList;
		},
		typename(type) {
			let name = '';
			if (type == 'media') {
				name = '媒体组件';
			} else if (type == 'shop') {
				name = '商城组件';
			} else if (type == 'tools') {
				name = '工具组件';
			} else if (type == 'page') {
				name = '页面组件';
			}
			return name;
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
