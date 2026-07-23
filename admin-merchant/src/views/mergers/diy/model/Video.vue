<template>
	<div
		class="pr"
		:style="{
			background: item.style.bgcolor,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}"
	>
		<div class="cover-box" @click.stop="diyEditer(index)"></div>
		<div class="drag-optional" :class="{ selected: index === selectedIndex }">
			<div class="diy-video">
				<video
					:style="{
						height: item.style.height + 'px',
						borderTopLeftRadius: item.style.topRadio + 'px',
						borderTopRightRadius: item.style.topRadio + 'px',
						borderBottomLeftRadius: item.style.bottomRadio + 'px',
						borderBottomRightRadius: item.style.bottomRadio + 'px'
					}"
					:src="item.params.videoUrl"
					:poster="item.params.poster"
					:autoplay="item.params.autoplay == 1"
					controls
				>
					您的浏览器不支持 video 标签
				</video>
			</div>
			<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
		</div>
	</div>
</template>

<script>
export default {
  inject: ['diyModel'],
	data() {
		return {};
	},
	props: ['item', 'index', 'selectedIndex'],
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
.diy-video {
	overflow: hidden;
}
.diy-video video {
	display: block;

	width: 100%;
}
.cover-box {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	z-index: 10;
}
</style>
