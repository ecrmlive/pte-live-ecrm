<template>

	<div
		:style="{
			background: item.style.bgcolor,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}"
		class="drag optional"
		:class="{ selected: index === selectedIndex }"
		@click.stop="diyEditer(index)"
	>
		<div
			class="diy-qixilive"
			:style="{
				background: item.style.background,
				borderTopLeftRadius: item.style.topRadio + 'px',
				borderTopRightRadius: item.style.topRadio + 'px',
				borderBottomLeftRadius: item.style.bottomRadio + 'px',
				borderBottomRightRadius: item.style.bottomRadio + 'px'
			}"
		>
			<div
				class="qixilive-head d-b-c"
				:style="{
					backgroundImage: item.style.background_image ? 'url(' + headBg + ')' : 'none',
					backgroundSize: '100% 100%',
					backgroundRepeat: 'no-repeat'
				}"
			>
				<div class="head-title" :style="{ color: themeColor }">{{ sectionTitle }}</div>
				<div v-if="showMore" class="head-more d-c-c" :style="{ color: themeColor }">
					<span>{{ moreTitle }}</span>
					<el-icon size="14px"><ArrowRight /></el-icon>
				</div>
			</div>
			<ul class="qixilive-list d-s-c f-w">
				<li class="item" v-for="(live, idx) in previewItems" :key="idx">
					<div class="box">
						<div class="pic"><img v-img-url="live.image" /></div>
						<div>{{ live.name }}</div>
					</div>
				</li>
			</ul>
		</div>
		<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
	</div>
</template>

<script>
import { getPageThemeApi } from '#/api/core/page';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { getThemeColorByIndex } from '#/utils/shop-theme';

import { resolveLivePreviewItems } from '../params/shared/marketing-helpers';

export default {
  inject: ['diyModel'],
	props: ['item', 'index', 'selectedIndex'],
	data() {
		return {
			themeColor: getThemeColorByIndex('0')
		};
	},
	computed: {
		previewItems() {
			return resolveLivePreviewItems(this.item);
		},
		headBg() {
			return resolveCosMediaUrl(this.item.style.background_image || '');
		},
		sectionTitle() {
			return (this.item.params && this.item.params.title) || '热门直播';
		},
		moreTitle() {
			return (this.item.params && this.item.params.moreTitle) || '更多';
		},
		showMore() {
			const v = this.item.params && this.item.params.showMore;
			return v === 1 || v === '1' || v === true;
		}
	},
	created() {
		this.loadThemeColor();
	},
	methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },

		loadThemeColor() {
			getPageThemeApi()
				.then((res) => {
					const theme = res.vars?.values?.theme ?? '0';
					this.themeColor = getThemeColorByIndex(theme);
				})
				.catch(() => {});
		}
	}
};
</script>

<style lang="scss" scoped>
.diy-qixilive {
	overflow: visible;
}

.diy-qixilive .qixilive-head {
	min-height: 40px;
	height: 40px;
	padding: 0 12px;
	border-top-left-radius: 6px;
	border-top-right-radius: 6px;
	overflow: visible;
	box-sizing: border-box;
	flex-shrink: 0;
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.diy-qixilive .head-title {
	font-size: 16px;
	font-weight: 700;
	line-height: 1;
}

.diy-qixilive .head-more {
	font-size: 12px;
	line-height: 1;
	gap: 2px;
	flex-shrink: 0;
	white-space: nowrap;
	display: flex;
	align-items: center;
}

.diy-qixilive .qixilive-list {
	display: flex;
	flex-wrap: wrap;
	justify-content: flex-start;
	align-items: flex-start;
	margin: 0;
	padding: 0;
	list-style: none;
}

.diy-qixilive .qixilive-list .item {
	width: 50%;
	box-sizing: border-box;
}

.diy-qixilive .qixilive-list .item .box {
	padding: 10px;
}

.diy-qixilive .qixilive-list .item .pic {
	position: relative;
	width: 100%;
	height: 120px;
	overflow: hidden;
}

.diy-qixilive .qixilive-list .item .pic img {
	display: block;
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.diy-qixilive .qixilive-list .item .box > div:last-child {
	color: #333;
	font-size: 13px;
	line-height: 1.4;
}
</style>
