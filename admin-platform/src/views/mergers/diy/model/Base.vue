<template>

	<div @click.stop="diyEditer(index)" :style="{
			background: item.style.bgcolor,
			padding: item.style.paddingTop + 'px ' + item.style.paddingLeft + 'px ' + item.style.paddingBottom + 'px ' + item.style.paddingLeft + 'px'
		}">
		<div class="drag optional" :class="{ selected: index === selectedIndex }">
			<div class="bg-box">
				<div class="d-s-c" style="z-index: 1;position: relative;margin-left: 10px;margin-top: 25px;">
					<div class="item-image"><img v-img-url="loginNameUrl" /></div>
					<div class="d-c d-c-s white" style="height: 51px;">
						<div class="f12">用户昵称</div>
						<div>ID:123456</div>
					</div>
				</div>
				<div class="bg-base" :class="'bg-base-' + styleType">
					<img class="centerbg" :src="centerbg" alt="" />
				</div>
				<div class="diy-Base" :style="{ background: item.style.background }">
					<div class="list column-3">
						<div class="item">
							<div class="item-text fb f16">100<span class="f12">.00</span></div>
							<div class="item-text text-ellipsis gray9">账户余额</div>
						</div>
						<div class="item" :style="'margin:0 ' + item.style.padding + 'px;'">
							<div class="item-text fb f16">20<span class="f12">.00</span></div>
							<div class="item-text text-ellipsis gray9">积分</div>
						</div>
						<div class="item">
							<div class="item-text fb f16">10</div>
							<div class="item-text text-ellipsis gray9">优惠券</div>
						</div>
					</div>
				</div>
			</div>

			<div class="btn-edit-del">
				<div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
			</div>
		</div>
	</div>
</template>

<script>
	import { DIY_CENTER_BG_URL, DIY_LOGIN_NAME_URL } from '#/utils/diy/order-icons';
	export default {
  inject: ['diyModel'],
		data() {
			return {
				centerbg: DIY_CENTER_BG_URL,
				loginNameUrl: DIY_LOGIN_NAME_URL,
			};
		},
	props: ['item', 'index', 'selectedIndex'],
	computed: {
		styleType() {
			const raw = this.item?.style?.type;
			const type = Number(raw);
			return Number.isFinite(type) && type > 0 ? type : 1;
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
	.diy-Base {
		position: absolute;
		bottom: 0;
		width: 355px;
		left: 8px;
		padding: 0 10px;
		box-sizing: border-box;
		border-radius: 8px;
	}

	.diy-Base .list {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
	}

	.diy-Base .list .item {
		padding: 10px 0;
		display: flex;
		justify-content: center;
		align-items: center;
		flex-direction: column;
	}

	.diy-Base .list.column-3 .item {
		// width: 33.333333333%;
		width: 20%;
	}

	.diy-Base .list.column-4 .item {
		width: 25%;
	}

	.diy-Base .list.column-5 .item {
		width: 20%;
	}

	.diy-Base .list .item-image {
		width: 60%;
	}

	.diy-Base .list .item-image img {
		width: 100%;
	}

	.diy-Base .list .item-text {
		width: 100%;
		padding: 4px 0;
		text-align: center;
	}

	.bg-box {
		overflow: hidden;
		height: 170px;
		position: relative;
	}

	.bg-base {
		width: 100%;
		height: 151px;
		position: absolute;
		right: 0;
		left: 0;
		top: 0;
		bottom: 0;
		margin: 0 auto;
		.centerbg{
			position: absolute;
			left: 0;
			top: 0;
			width: 100%;
		}
	}

	.bg-base-1 {
		background: linear-gradient(-180deg, #fd6303 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-2 {
		background: linear-gradient(-180deg, #42ca4d 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-3 {
		background: linear-gradient(-180deg, #ff448f 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-4 {
		background: linear-gradient(-180deg, #409eff 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-5 {
		background: linear-gradient(-180deg, #101012 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-6 {
		background: linear-gradient(-180deg, #e2aa62 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-7 {
		background: linear-gradient(-180deg, #a253ff 0%, rgba(240, 40, 17, 0) 100%);
	}

	.bg-base-8 {
		background: linear-gradient(-180deg, #fe0024 0%, rgba(240, 40, 17, 0) 100%);
	}

	.item-image {
		width: 51px;
		height: 51px;
		border-radius: 50%;
		background-color: #fff;
		overflow: hidden;
		margin-right: 10px;

		img {
			width: 51px;
			height: 51px;
			border-radius: 50%;
			background-color: #fff;
		}
	}

	.d-c-s {
		display: flex;
		justify-content: center;
		align-items: flex-start;
	}
</style>