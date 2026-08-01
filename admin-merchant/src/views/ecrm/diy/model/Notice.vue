<template>
	<div>
		<div @click.stop="diyEditer(index)">
			<div
				class="drag optional"
				:class="{ selected: index === selectedIndex }"
				:style="{
					background: item.style.bgcolor,
					paddingLeft: item.style.paddingLeft + 'px',
					paddingRight: item.style.paddingLeft + 'px',
					paddingTop: item.style.paddingTop + 'px',
					paddingBottom: item.style.paddingBottom + 'px'
				}"
			>
				<div
					class="diy-notice"
					:style="{
						background: item.style.background,
						borderTopLeftRadius: item.style.topRadio + 'px',
						borderTopRightRadius: item.style.topRadio + 'px',
						borderBottomLeftRadius: item.style.bottomRadio + 'px',
						borderBottomRightRadius: item.style.bottomRadio + 'px',
						padding: '0 ' + item.style.padding + 'px'
					}"
				>
					<div class="notice-icon"><img v-img-url="item.params.icon" /></div>
					<div class="notice-text flex-1 text-ellipsis">
						<span :style="{ color: item.style.textColor }">{{ item.params.text }}</span>
					</div>
				</div>
				<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
			</div>
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
.diy-notice {
	display: flex;
	justify-content: flex-start;
	align-items: center;
}
.diy-notice .notice-icon {
	width: 32px;
	height: 32px;
}
.diy-notice .notice-text {
	margin-left: 10px;
	overflow: hidden;
	white-space: nowrap;
}
</style>
