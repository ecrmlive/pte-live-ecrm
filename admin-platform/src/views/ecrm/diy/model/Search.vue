<template>

	<div class="drag optional" :class="{ selected: index === selectedIndex }" :style="{
			background: item.style.background,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}" @click.stop="diyEditer(index)">
		<div class="diy-search navigation d-s-c" :style="{
					background: item.style.bgcolor,
					borderTopLeftRadius: item.style.topRadio + 'px',
					borderTopRightRadius: item.style.topRadio + 'px',
					borderBottomLeftRadius: item.style.bottomRadio + 'px',
					borderBottomRightRadius: item.style.bottomRadio + 'px'
				}">
			<img v-if="item.params.title_type == 'image'" class="logo-img" v-img-url="item.params.toplogo" alt="" />
			<div v-if="item.params.title_type == 'text'" class="f15" :style="{color: item.style.titleTextColor}">
				{{item.params.title}}</div>
			<div class="phone-top-search-box d-s-c"
				:style="{background: item.style.searchBackGround,color: item.style.searchColor}">
				<el-icon class="mr10" :color="item.style.searchColor|| '#999'">
					<Search />
				</el-icon>{{item.params.searchText}}
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
		created() {},
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
	.diy-search {
		padding: 9px 12px 9px 12px;

		.logo-img {
			display: block;
			width: 39px;
			height: 32px;
		}

		.phone-top-search-box {
			margin-left: 15px;
			flex: 1;
			height: 30px;
			border-radius: 30px;
			font-size: 13px;
			line-height: 30px;
			color: #999999;
			padding: 0 10px;
		}
	}
</style>