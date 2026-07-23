<template>

	<div
		@click.stop="diyEditer(index)"
		:style="{
			background: item.style.bgcolor,
			padding: item.style.paddingTop + 'px ' + item.style.paddingLeft + 'px ' + item.style.paddingBottom + 'px ' + item.style.paddingLeft + 'px'
		}"
	>
		<div class="drag optional" :class="{ selected: index === selectedIndex }">
			<div
				class="diy-Order"
				:style="{
					background: item.style.background,
					borderRadius: item.style.topRadio + 'px ' + item.style.topRadio + 'px ' + item.style.bottomRadio + 'px ' + item.style.bottomRadio + 'px'
				}"
			>
				<ul class="list column-5">
					<li class="item" v-for="(name, iconIndex) in nameList" :key="styleType + '-' + iconIndex">
						<div class="item-image"><img :src="orderIcons[iconIndex]" /></div>
						<div class="item-text text-ellipsis">{{ name }}</div>
					</li>
				</ul>
			</div>
			<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
		</div>
	</div>
</template>

<script>
import { DIY_ORDER_ICON_LIST } from '#/utils/diy/order-icons';

export default {
  inject: ['diyModel'],
	data() {
		return {
			nameList: ['待付款', '待发货', '待收货', '待评价', '退款/售后'],
			imgList: DIY_ORDER_ICON_LIST,
		};
	},
	props: ['item', 'index', 'selectedIndex'],
	computed: {
		styleType() {
			const raw = this.item?.style?.type;
			const type = Number(raw);
			return Number.isFinite(type) && type > 0 ? type : 1;
		},
		orderIcons() {
			return this.imgList[this.styleType - 1] || this.imgList[0];
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
.diy-Order .list {
	display: flex;
	flex-wrap: wrap;
	align-items: flex-start;
}
.diy-Order .list .item {
	padding: 10px 0;
	display: flex;
	justify-content: center;
	align-items: center;
	flex-direction: column;
}
.diy-Order .list.column-3 .item {
	width: 33.333333333%;
}
.diy-Order .list.column-4 .item {
	width: 25%;
}
.diy-Order .list.column-5 .item {
	width: 20%;
}
.diy-Order .list .item-image {
	width: 60%;
}
.diy-Order .list .item-image img {
	width: 80%;
	margin: 0 auto;
}
.diy-Order .list .item-text {
	width: 100%;
	padding: 4px 0;
	text-align: center;
}
</style>
