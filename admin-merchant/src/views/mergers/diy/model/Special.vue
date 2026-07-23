<template>

	<div
		@click.stop="diyEditer(index)"
		:style="{
			background: item.style.bgcolor,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}"
	>
		<div
			class="drag optional"
			:style="{
				background: item.style.background,
				borderTopLeftRadius: item.style.topRadio + 'px',
				borderTopRightRadius: item.style.topRadio + 'px',
				borderBottomLeftRadius: item.style.bottomRadio + 'px',
				borderBottomRightRadius: item.style.bottomRadio + 'px'
			}"
			:class="{ selected: index === selectedIndex }"
		>
			<div class="diy-special">
				<div class="special-left"><img v-img-url="item.style.image" alt="" /></div>
				<div class="special-content" :class="'display_' + item.style.display">
					<ul class="special-content-list">
						<li class="content-item text-ellipsis" v-for="(article, index) in previewItems" :key="index">
							<span>{{ article.article_title }}</span>
						</li>
					</ul>
				</div>
				<div class="special-more"><i class="el-icon-arrow-right"></i></div>
			</div>
			<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
		</div>
	</div>
</template>

<script>
import { resolveAutoSourcePreviewItems } from '../params/shared/marketing-helpers';

export default {
  inject: ['diyModel'],
	data() {
		return {};
	},
	props: ['item', 'index', 'selectedIndex'],
	computed: {
		previewItems() {
			return resolveAutoSourcePreviewItems(this.item, {
				article_title: '此处是文章标题',
			});
		},
	},
	methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },
}
};
</script>

<style lang="scss" scoped>
.diy-special {
	padding: 4px;
	display: flex;
	justify-content: space-between;
	align-items: center;
}
.diy-special .special-left {
	width: 105px;
	height: 27px;
}
.diy-special .special-left img {
	width: 100%;
	height: 100%;
}
.diy-special .special-more {
	font-size: 18px;
}
.diy-special .special-content {
	flex: 1;
	margin: 0 10px;
	overflow: hidden;
}
.diy-special .content-item {
	padding: 4px 0;
}
.diy-special .special-content.display_1 {
	height: 22px;
}
.diy-special .special-content.display_2 {
	height: 44px;
}
</style>
