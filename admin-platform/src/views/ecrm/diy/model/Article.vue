<template>

	<div
		:style="{
			background: item.style.bgcolor,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}"
		@click.stop="diyEditer(index)"
	>
		<div
			:style="{
				background: item.style.background,
				borderTopLeftRadius: item.style.topRadio + 'px',
				borderTopRightRadius: item.style.topRadio + 'px',
				borderBottomLeftRadius: item.style.bottomRadio + 'px',
				borderBottomRightRadius: item.style.bottomRadio + 'px'
			}"
			class="drag optional o-h"
			:class="{ selected: index === selectedIndex }"
		>
			<div class="diy-article">
				<div
					class="article-item"
					v-for="(article, index) in previewItems"
					:class="'show-type__' + article.show_type"
					:key="index"
				>
					<!-- 小图模式 -->
					<div class="article-item__image"><img v-img-url="article.image" alt="" /></div>
					<div class="article-item__left flex-1">
						<div class="article-item__title text-ellipsis-2">
							<span class="f18">{{ article.article_title }}</span>
						</div>
						<div class="article-item__footer d-b-c">
							<span class="gray9">{{ article.views_num }}次浏览</span>
							<span class="gray9">2022-02-22</span>
						</div>
					</div>
				</div>
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
				image: '',
				show_type: 10,
				views_num: 0,
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
.o-h {
	overflow: hidden;
}
.diy-article .show-type__10 {
	display: flex;
	padding: 10px;
	border-bottom: 1px solid #eeeeee;
}

.diy-article .show-type__10 .article-item__image {
	width: 123px;
	height: 70px;
	border-radius: 3px;
	overflow: hidden;
}

.diy-article .show-type__10 .article-item__image > img {
	width: 123px;
	height: 70px;
}

.diy-article .show-type__10 .article-item__title {
	height: 40px;
}
</style>
