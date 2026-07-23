<template>
	<div class="drag optional" @click.stop="diyEditer(index)"
		:class="{ selected: index === selectedIndex }" :style="{
			background: item.style.background,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}">
		<div class="diy-banner" :style="'height:' + item.style.height * 0.5 + 'px;'">
			<div class="img-list pr">
				<img
					class="banner-img"
					style="display: block;"
					:style="{
						height: item.style.height * 0.5 + 'px',
						borderTopLeftRadius: item.style.topRadio + 'px',
						borderTopRightRadius: item.style.topRadio + 'px',
						borderBottomLeftRadius: item.style.bottomRadio + 'px',
						borderBottomRightRadius: item.style.bottomRadio + 'px'
					}"
					v-img-url="item.data[0] && item.data[0].imgUrl"
				/>
				<div class="dots center d-c-c">
					<div :key="index" :class="index == 0 ? 'active ' + item.style.imgShape : item.style.imgShape"
						v-for="(banner, index) in item.data"
						:style="index == 0 ? 'background:' + item.style.btnColor : 'background:' + item.style.btnColor + ';opacity: 0.4;'">
					</div>
				</div>
			</div>

		</div>
		<div class="btn-edit-del">
			<div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
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
	.p15 {
		padding: 15px;
	}

	.diy-banner.round {
		padding: 12px;
		box-sizing: content-box;
		overflow: hidden;
		text-align: center;
	}

	.diy-banner.round img {
		width: 100%;
		height: 100px;
		object-fit: fill;
		border-radius: 10px;
		margin-bottom: 12px;
		box-sizing: content-box;
	}

	.diy-banner.square {
		height: 100px;
		overflow: hidden;
		text-align: center;
	}

	.diy-banner.square img {
		width: 100%;
		height: 100px;
		object-fit: fill;
	}

	.diy-banner .dots {
		position: absolute;
		left: 0;
		right: 0;
		margin: 0 auto;
		bottom: 10px;
		width: 100%;
		z-index: 1;
	}

	.diy-banner.round .dots {
		position: absolute;
		left: 0;
		right: 0;
		margin: 0 auto;
		bottom: 20px;
	}

	.diy-banner .dots .square,
	.diy-banner .dots .round,
	.diy-banner .dots .rectangle {
		bottom: 40rpx;
		left: 0;
		right: 0;
		margin: auto;
	}

	.diy-banner .dots .square {
		width: 7px;
		height: 7px;
		margin: 0 2px;
		background: #ebedf0;
		opacity: 0.3;
	}

	.diy-banner .dots .round {
		width: 7px;
		height: 7px;
		margin: 0 2px;
		background: #ebedf0;
		opacity: 0.3;
		border-radius: 50%;
	}

	.diy-banner .dots .rectangle {
		width: 20px;
		height: 3px;
		margin: 0 2px;
		background: #ebedf0;
		opacity: 0.3;
		border-radius: 4rpx;
	}

	.diy-banner .dots .active {
		opacity: 1;
	}
</style>