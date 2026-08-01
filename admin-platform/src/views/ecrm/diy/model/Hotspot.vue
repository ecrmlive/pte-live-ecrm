<template>
	<div @click.stop="diyEditer(index)">
		<div class="drag optional" :class="{ selected: index === selectedIndex }">
			<div class="diy-hotspot" :style="{
					paddingTop: item.style.paddingTop + 'px',
					paddingBottom: item.style.paddingBottom + 'px',
					paddingLeft: item.style.paddingLeft + 'px',
					paddingRight: item.style.paddingLeft + 'px',
					background: item.style.background
				}">
				<div class="img-box" style="position: relative;" :style="{
					backgroundColor: item.style.bgcolor,
					borderTopLeftRadius: item.style.topRadio + 'px',
					borderTopRightRadius: item.style.topRadio + 'px',
					borderBottomLeftRadius: item.style.bottomRadio + 'px',
					borderBottomRightRadius: item.style.bottomRadio + 'px',
					overflow: 'hidden'
				}">
					<img v-img-url="item.params.image" style="width: 100%; display: block;" />
                   <!-- <div v-for="(hotspot, index) in item.data" :key="index"
                        class="hotspot-area"
                        :style="{
                            width: hotspot.width + '%',
                            height: hotspot.height + '%',
                            top: hotspot.top + '%',
                            left: hotspot.left + '%'
                        }">
                    </div> -->
				</div>
			</div>
			<div class="btn-edit-del">
				<div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
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
.diy-hotspot .img-box {
    width: 100%;
}
.hotspot-area {
    position: absolute;
    border: 1px dashed #ff0000;
    background: rgba(255, 0, 0, 0.2);
    box-sizing: border-box;
    z-index: 10;
}
</style>
