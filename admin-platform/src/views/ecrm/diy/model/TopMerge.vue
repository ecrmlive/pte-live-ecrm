<template>
	<div class="drag optional" @click.stop="diyEditer(index)"
		:class="{ selected: index === selectedIndex }">
		<div class="bg-topMerge">
			<img class="bgimg" v-img-url="item.images[0] && item.images[0].imgUrl" />
		</div>
		<div class="bg-topMerge-color"
			:style="{backgroundImage:`linear-gradient(rgba(245, 245, 245, 0) 0%, rgba(245, 245, 245, 0) 50%,${ item.style.bgcolor_color1 || '#fff'} 100%)`}">
		</div>
		<div class="diy-TopMerge">
			<div class="navigation d-s-c">
				<div v-if="item.params.showLocation" class="location-chip">
					<span>⌖</span>{{ item.params.locationText || '定位中' }}<b>›</b>
				</div>
				<img v-else class="logo-img" v-img-url="item.params.topLogo" alt="" />
				<div class="phone-top-search-box d-s-c">
					<el-icon class="mr10" color="#999">
						<Search />
					</el-icon>{{item.params.searchText}}
				</div>
			</div>
			<div class="option-list optionType" v-if="item.params.showCategory">
				<div class="item" :style="{marginRight:item.style.categoryPadding + 'px'}"
					v-for="(citem,cindex) in item.data" :key="cindex">
					{{citem.text}}
				</div>
			</div>
			<div class="img-list pr" :class="`imageType${item.params.type}`">
				<img class="banner-img" v-if="item.params.type == 1" :style="{
						borderTopLeftRadius: item.style.topRadio + 'px',
						borderTopRightRadius: item.style.topRadio + 'px',
						borderBottomLeftRadius: item.style.bottomRadio + 'px',
						borderBottomRightRadius: item.style.bottomRadio + 'px'
					}" v-img-url="item.images[0] && item.images[0].imgUrl" />
				<div v-if="item.params.type != 1 " class="ww100 d-b-c">
					<div class="smallImg"
						:style="{borderTopRightRadius: item.style.topRadio + 'px',borderBottomRightRadius: item.style.bottomRadio + 'px'}">
						<img class="banner-img" v-img-url="item.images[1] && item.images[1].imgUrl" />
					</div>
					<img class="banner-img " :style="{borderTopLeftRadius: item.style.topRadio + 'px',
						borderTopRightRadius: item.style.topRadio + 'px',
						borderBottomLeftRadius: item.style.bottomRadio + 'px',
						borderBottomRightRadius: item.style.bottomRadio + 'px'}" v-img-url="item.images[0] && item.images[0].imgUrl" />
					<div class="smallImg"
						:style="{borderTopLeftRadius: item.style.topRadio + 'px',borderBottomLeftRadius: item.style.bottomRadio + 'px'}">
						<img class="banner-img " v-img-url="item.images[2] && item.images[2].imgUrl" />
					</div>
				</div>
				<div class="dots center d-c-c"
					:class="{'d-s-c':item.style.btnShape =='left','d-c-c':item.style.btnShape =='center','d-e-c':item.style.btnShape =='right'}">
					<div :key="index" :class="index == 0 ? 'active ' + item.style.imgShape : item.style.imgShape"
						v-for="(banner, index) in item.images"
						:style="index == 0 ? 'background:' + item.style.btnColor : 'background:' + item.style.btnOpColor">
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
	.optional {
		position: relative;
		padding-bottom: 11px;
	}

	.navigation {
		display: flex;
		align-items: center;
	}

	.img-list {
		position: relative;
		width: 100%;
	}

	.bg-topMerge-color {
		position: absolute;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		overflow: hidden;
		z-index: 1;
	}

	.bg-topMerge {
		position: absolute;
		width: 100%;
		height: 100%;
		top: 0;
		z-index: 0;
		filter: blur(0);
		overflow: hidden;

		.bgimg {
			width: 100%;
			height: 100%;
			filter: blur(15px);
			transform: scale(1.5);
		}
	}

	.navigation {
		padding: 9px 12px 0 12px;

		.logo-img {
			display: block;
			width: 39px;
			height: 32px;
		}

		.location-chip {
			display: flex;
			max-width: 100px;
			align-items: center;
			gap: 3px;
			overflow: hidden;
			color: #4a5261;
			font-size: 12px;
			text-overflow: ellipsis;
			white-space: nowrap;

			span { color: #1677ff; font-size: 16px; }
			b { color: #8b95a7; font-size: 16px; font-weight: 400; }
		}

		.phone-top-search-box {
			margin-left: 15px;
			flex: 1;
			min-width: 0;
			height: 30px;
			border-radius: 30px;
			font-size: 13px;
			line-height: 30px;
			color: #999999;
			padding: 0 10px;
			background-color: #fff;
			display: flex;
			align-items: center;
		}
	}

	.p15 {
		padding: 15px;
	}

	.option-list.optionType {
		display: flex;
		justify-content: flex-start;
		align-items: center;
		overflow-x: auto;
		padding-left: 20px;
		height: 42px;

		.item {
			flex-shrink: 0;
			white-space: nowrap;
			font-size: 13px;
			color: #fff;
			position: relative;
		}

		.item.acitve {
			content: '';
			position: absolute;
			width: 19px;
			height: 2px;
			background: #FFFFFF;
			border-radius: 2px;
			left: 0;
			right: 0;
			margin: auto;
			bottom: -8px;
		}
	}

	.diy-TopMerge {
		position: relative;
		z-index: 1;
	}

	.diy-TopMerge .imageType1 {
		overflow: hidden;
		text-align: center;
		width: calc(100% - 20px);
		max-width: 355px;
		height: 159px;
		margin: 0 auto;

		.banner-img {
			display: block;
			object-fit: cover;
			object-position: center;
			width: 100%;
			height: 159px;
		}
	}

	.diy-TopMerge .imageType2 {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: calc(100% - 20px);
		max-width: 355px;
		margin: 0 auto;

		.banner-img {
			display: block;
			object-fit: fill;
			width: 335px;
			height: 159px;

		}

		.smallImg {
			width: 10px;
			height: 139px;
			background: #fff;
			overflow: hidden;

			.banner-img {
				width: 10px;
				height: 139px;
			}
		}

	}

	.diy-TopMerge .dots {
		padding: 0 20px;
		box-sizing: border-box;
		position: absolute;
		left: 0;
		right: 0;
		margin: 0 auto;
		bottom: 10px;
		width: 100%;
		z-index: 1;
		display: flex;
		align-items: center;
	}

	.diy-TopMerge.round .dots {
		position: absolute;
		left: 0;
		right: 0;
		margin: 0 auto;
		bottom: 20px;
	}

	.diy-TopMerge .dots .square,
	.diy-TopMerge .dots .round,
	.diy-TopMerge .dots .rectangle {
		bottom: 20px;
		left: 0;
		right: 0;
		margin: auto;
	}

	.diy-TopMerge .dots .square {
		width: 7px;
		height: 7px;
		margin: 0 2px;
		background: #ebedf0;
		opacity: 0.3;
	}

	.diy-TopMerge .dots .round {
		width: 7px;
		height: 7px;
		margin: 0 2px;
		background: #ebedf0;
		opacity: 0.3;
		border-radius: 50%;
	}

	.diy-TopMerge .dots .rectangle {
		width: 20px;
		height: 3px;
		margin: 0 2px;
		background: #ebedf0;
		opacity: 0.3;
		border-radius: 2px;
	}

	.diy-TopMerge .dots .active {
		opacity: 1;
	}
</style>
