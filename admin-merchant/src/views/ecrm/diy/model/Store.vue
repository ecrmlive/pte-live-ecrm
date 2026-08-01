<template>

	<div
		:style="{
			background: item.style.bgcolor,
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px'
		}"
		@click.stop="diyEditer(index)"
	>
		<div
			class="drag optional"
			:class="{ selected: index === selectedIndex }"
			:style="{
				background: item.style.background,
				borderTopLeftRadius: item.style.topRadio + 'px',
				borderTopRightRadius: item.style.topRadio + 'px',
				borderBottomLeftRadius: item.style.bottomRadio + 'px',
				borderBottomRightRadius: item.style.bottomRadio + 'px'
			}"
		>
			<div class="diy-shop">
				<div class="diy-head d-b-c">
					<div class="left d-s-c"><div class="name">线下门店</div></div>
				</div>
				<div class="shop-item" :key="index" v-for="(shop, index) in previewItems">
					<div class="shop-item-logo">
						<img v-if="shop.logo" :src="shop.logo.file_path" alt="门店logo" />
						<img v-else :src="shop.logo_image" alt="门店logo" />
					</div>
					<div class="shop-item-content">
						<div class="shop-item-title">
							<span>{{ shop.store_name }}</span>
						</div>
						<div class="shop-item-address">
							<span>门店地址：{{ shop.address }}</span>
						</div>
						<div class="shop-item-phone">
							<span>联系电话：{{ shop.phone }}</span>
						</div>
					</div>
				</div>
			</div>
			<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
		</div>
	</div>
</template>

<script>
import { resolveAutoSourcePreviewItems } from '../params/shared/marketing-helpers';

export default {
  inject: ['diyModel'],
	data() {
		return {};
	},
	props: ['item', 'index', 'selectedIndex'],
	computed: {
		previewItems() {
			return resolveAutoSourcePreviewItems(this.item, {
				address: '门店地址',
				logo_image: '',
				phone: '联系电话',
				store_name: '此处是门店',
			});
		},
	},
	mounted() {},
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
.diy-shop .diy-head {
	height: 40px;
	padding: 0 10px;
}

.diy-shop .diy-head .name {
	font-size: 18px;
	font-weight: bold;
}

.diy-shop .shop-item {
	display: flex;
	justify-content: flex-start;
	align-items: center;
}

.diy-shop .shop-item-logo,
.diy-shop .shop-item-logo img {
	width: 60px;
	height: 60px;
	border-radius: 50%;
}

.diy-shop .shop-item {
	padding: 10px;
}

.diy-shop .shop-item .shop-item-content {
	margin-left: 10px;
}

.diy-shop .shop-item .shop-item-title {
	font-size: 16px;
	color: rgb(255, 51, 0);
}

.diy-shop .shop-item .shop-item-address {
	padding-top: 6px;
	color: #999999;
}

.diy-shop .shop-item .shop-item-phone {
	padding-top: 6px;
	color: #999999;
}
</style>
