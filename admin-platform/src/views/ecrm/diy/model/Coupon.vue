<template>
	<div @click.stop="diyEditer(index)" class="drag optional" :class="{ selected: index === selectedIndex }">
		<div
			class="diy-coupon"
			:style="{
				background: item.style.bgcolor,
				paddingLeft: item.style.paddingLeft + 'px',
				paddingRight: item.style.paddingLeft + 'px',
				paddingTop: item.style.paddingTop + 'px',
				paddingBottom: item.style.paddingBottom + 'px'
			}"
		>
			<div
				class="pr d-b-c coupon-wrapper-box"
				:style="{
					background: item.style.bgtype == 1 ? item.style.background : 'none',
					borderTopLeftRadius: item.style.topRadio + 'px',
					borderTopRightRadius: item.style.topRadio + 'px',
					borderBottomLeftRadius: item.style.bottomRadio + 'px',
					borderBottomRightRadius: item.style.bottomRadio + 'px'
				}"
			>
				<img
					class="bg-couponbg"
					v-if="item.style.bgtype == 2"
					:src="item.style.bgimage"
					alt=""
					:style="{
						borderTopLeftRadius: item.style.topRadio + 'px',
						borderTopRightRadius: item.style.topRadio + 'px',
						borderBottomLeftRadius: item.style.bottomRadio + 'px',
						borderBottomRightRadius: item.style.bottomRadio + 'px'
					}"
				/>
				<div
					:key="'coupon-' + index"
					v-for="(_, index) in previewCount"
					class="coupon-wrapper"
				>
					<img class="bg-coupon-item" :src="couponBgUrl" alt="" />
					<div class="coupon-item d-c-c d-c">
						<div style="height: 78px;border-bottom: 1px dashed;" :style="{ borderColor: item.style.btncolor }">
							<div class="content-top" :style="{ color: item.style.pricecolor }">
								<span class="f12">￥</span>
								<span class="f23 fb">100</span>
							</div>
							<div class="content-bottom d-c f11" :style="{ color: item.style.cillcolor }"><span>满1000可用</span></div>
							<div class="d-c f11" :style="{ color: item.style.descolor }">指定商品可用</div>
						</div>
						<div style="height: 38px;" class="d-c-c">
							<div class="right-receive" :style="{ color: item.style.btnTxtcolor, borderRadius: item.style.btnRadio + 'px', backgroundColor: item.style.btncolor }">
								{{ item.params.btntext }}
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
	</div>
</template>

<script>
import { DIY_COUPON_BG_URL } from '#/utils/diy/order-icons';

import { resolveCouponPreviewCount } from '../params/shared/marketing-helpers';

export default {
  inject: ['diyModel'],
	data() {
		return {
			couponBgUrl: DIY_COUPON_BG_URL,
		};
	},
	props: ['item', 'index', 'selectedIndex'],
	computed: {
		previewCount() {
			return resolveCouponPreviewCount(this.item);
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
.diy-coupon {
	display: flex;
	justify-content: flex-start;
	align-items: center;
}

.diy-coupon .coupon-wrapper {
	position: relative;
	width: 108px;
	height: 116px;
}

.diy-coupon .coupon-wrapper .coupon-item {
	height: 116px;
	overflow: hidden;
	position: relative;
	z-index: 0;
}

.diy-coupon .coupon-wrapper .left-content {
	width: 80px;
	height: 56px;
	padding: 0;
	box-sizing: border-box;
	color: #ffffff;
	flex: 1;
	background: #ff4c01;
}

.diy-coupon .coupon-wrapper .left-content .price {
	font-size: 20px;
	font-weight: bold;
}

.diy-coupon .coupon-wrapper .content-top {
	font-size: 12px;
	text-align: center;
}
.diy-coupon .coupon-wrapper .content-bottom {
	font-size: 10px;
	text-align: center;
}
.diy-coupon .coupon-wrapper .right-receive {
	min-width: 71px;
	padding: 0 10px;
	box-sizing: border-box;
	height: 24px;
	flex-shrink: 0;
	display: flex;
	justify-content: center;
	align-items: center;
	text-align: center;
	font-size: 11px;
}
.lr-tb {
	width: 13px;
}
.coupon-wrapper-box {
	width: 100%;
	box-sizing: border-box;
	padding: 14px 11px;
	overflow-x: auto;
}
.bg-coupon-item {
	position: absolute;
	z-index: 0;
	left: 0;
	top: 0;
	width: 104px;
	height: 116px;
}
.f11 {
	font-size: 11px;
}
.f12 {
	font-size: 12px;
}
.f23 {
	font-size: 23px;
}
.bg-couponbg{
	position: absolute;
	z-index: 0;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
}
</style>
