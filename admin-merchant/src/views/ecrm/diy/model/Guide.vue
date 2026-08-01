<template>

	<div>
		<div @click.stop="diyEditer(index)" class="drag optional" :class="{ selected: index === selectedIndex }">
			<div
				class="diy-guide"
				:style="{
					padding: item.style.paddingTop + 'px ' + item.style.paddingLeft + 'px ' + item.style.paddingBottom + 'px ' + item.style.paddingLeft + 'px',
					background: item.style.background
				}"
			>
				<p class="line" :style="{ borderTopWidth: item.style.lineHeight + 'px', borderTopColor: item.style.lineColor, borderTopStyle: item.style.lineStyle }"></p>
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

<style lang="scss" scoped></style>
