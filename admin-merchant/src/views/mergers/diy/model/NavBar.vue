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
		<div class="drag optional" :class="{ selected: index === selectedIndex }">
			<div
				class="diy-navBar"
				:style="{
					background: item.style.background,
					borderTopLeftRadius: item.style.topRadio + 'px',
					borderTopRightRadius: item.style.topRadio + 'px',
					borderBottomLeftRadius: item.style.bottomRadio + 'px',
					borderBottomRightRadius: item.style.bottomRadio + 'px'
				}"
			>
				<ul class="list" :class="'column-' + rowsNum">
					<template v-for="(navBar, index) in item.data" :key="index">
						<li v-if="!navBar.hide" class="item">
							<div class="item-image"><img v-img-url="navBar.imgUrl" /></div>
							<div class="item-text text-ellipsis" :style="{ color: navBar.color }">{{ navBar.text }}</div>
						</li>
					</template>
				</ul>
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
	computed: {
		rowsNum() {
			const raw = this.item?.style?.rowsNum;
			const num = Number(raw);
			return Number.isFinite(num) && num > 0 ? num : 5;
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
.diy-navBar .list {
	display: flex;
	flex-wrap: wrap;
	align-items: flex-start;
}
.diy-navBar .list .item {
	padding: 10px 0;
	display: flex;
	justify-content: center;
	align-items: center;
	flex-direction: column;
}
.diy-navBar .list.column-3 .item {
	width: 33.333333333%;
}
.diy-navBar .list.column-4 .item {
	width: 25%;
}
.diy-navBar .list.column-5 .item {
	width: 20%;
}
.diy-navBar .list .item-image {
	width: 60%;
}
.diy-navBar .list .item-image img {
	width: 100%;
}
.diy-navBar .list .item-text {
	width: 100%;
	padding: 4px 0;
	text-align: center;
}
</style>
