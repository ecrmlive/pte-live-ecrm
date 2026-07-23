<template>
	<div class="phone-top" v-if="diyData.page?.params" @click.stop="diyEditer(-1)" :style="'background:#fff'">
		<div class="status-bar">
			<span class="icon iconfont icon-wifi" :style="'color:#333'">
			</span>
			<span class="time" :style="'color:#333'">19:00</span>
			<span class="icon iconfont icon-xinhao" :style="'color:#333'">
			</span>
			<span class="ml4 icon iconfont icon-iconset0250" :style="'color:#333'">
			</span>
		</div>
		<div class="page-title d-c-c">
			{{ diyData.page.params?.name || diyData.page.name || '页面标题' }}
		</div>
	</div>
</template>

<script>
	export default {
  inject: ['diyModel'],
		data() {
			return {};
		},
		props: {
			diyData: Object,
			isDiy: Boolean
		},
		computed: {
			bgColor() {
				return this.$props.isDiy ? '#fff' : this.$props.diyData.page.style.titleBackgroundColor
			},
			searchStyle() {
				const colorRgb = this.$props.diyData.page.style.titleBackgroundColor || '#999';
				const activeList = ['#ffffff', '#FFFFFF'];
				const flag = activeList.includes(colorRgb);
				return {
					color: flag ? '#ccc' : colorRgb,
					borderColor: flag ? '#ccc' : colorRgb,
				}
			},
		},
		methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },

			getImg(src) {
				var img_url = src
				var img = new Image()
				img.src = img_url;
				return Math.ceil(img.width * 60 / img.height) + 'px'
			},
		}
	};
</script>

<style scoped lang="scss">
	.page-title{
		height: 35px;
		font-size: 15px;
		color: #333;
	}
	.phone-top-search-box {
		flex: 1;
		height: 32px;
		line-height: 32px;
		background-color: #f7f7f7;
		border-radius: 3px;
		margin-left: 10px;
		color: #999999;
		font-size: 13px;
		text-align: left;
		padding-left: 10px;
		border-radius: 15px;
		font-weight: 400;
		border: 1px solid transparent;
	}

	.phone-top-search-box .Search {
		font-weight: 800;
		margin-right: 4px;
	}

	.navigation>img {
		width: 30px;
		height: 30px;
	}
</style>