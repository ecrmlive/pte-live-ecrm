<template>

	<div class="drag optional" :class="{ selected: index === selectedIndex }"
		@click.stop="diyEditer(index)" :style="{
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px',
			marginTop:item.style.marginTop + 'px',
			background: item.style.background
			
		}">
		<div class="diy-product o-h" :style="borderStyle">
			<div class="diy-head d-b-c" :style="{
					background:item.params.titleBgType == 2?`linear-gradient(to right, ${item.style.titleBg_color1 || '#fff'}, ${item.style.titleBg_color2 || '#fff'})` :'url(' + item.params.bgimage + ')'
				}">
				<div class="left d-s-c">
					<div v-if="item.params.titleType == 1" class="name" :style="{
							color: themeColor,
							fontSize: item.style.titleSize + 'px',
							fontWeight:item.style.titleWeight == 1?'bold':'',
							fontStyle:item.style.titleWeight == 2?'italic':''
						}">
						{{ item.params.title }}
					</div>
					<img v-if="item.params.titleType == 2" class="titleImgt" :src="item.params.titleimage" alt="" />
				</div>
				<div class="right white d-c-c" style="line-height: 1;" :style="{
						color: themeColor,
						fontSize: item.style.moreSize + 'px'
					}">
					{{ item.params.more }}
					<el-icon size="14px">
						<ArrowRight />
					</el-icon>
				</div>
			</div>
			<div class="product-list-box column__4" :key="'num-' + previewProducts.length">
				<div class="product-list column__4">
					<div v-for="(product, index) in previewProducts" :key="index" class="product-item">
						<!-- 商品图片 -->
						<img :style="{borderTopLeftRadius: item.style.productTopRadio + 'px',
					borderTopRightRadius: item.style.productTopRadio + 'px',
					borderBottomLeftRadius: item.style.productBottomRadio + 'px',
					borderBottomRightRadius: item.style.productBottomRadio + 'px'}" class="product-cover" v-img-url="product.image" />
						<div class="product-info">
							<!-- 商品名称 -->
							<div class="flex-1 ww100">
								<div :style="{color: item.style.productName_color}"
									class="product-name text-ellipsis-2">
									<span class="product-tag"
										:style="{color:item.style.tag_text_color,backgroundImage: 'linear-gradient(to right, ' + (item.style.tag_bg1 || '#fff') + ', ' + (item.style.tag_bg2 || '#fff') + ')'}">新人专享</span>
									<span :class="{fb:item.style.nameWeight}">{{ product.product_name }}</span>
								</div>
							</div>
							<!-- 商品价格 -->
							<div class="price hasbtn ">
								<div class="price-num" :style="{color: item.style.productPrice_color}">
									<span class="f11">¥</span>
									<span class="f18 fb">120</span>
								</div>
								<div class="cart-btn icon"
									:style="{backgroundImage: 'linear-gradient(to right, ' + (item.style.productBtn_color1 || '#fff') + ', ' + (item.style.productBtn_color2 || '#fff') + ')'}">
									<span class=" icon iconfont icon-jia"
										:style="{color:item.style.btn_text_color}"></span>
								</div>
							</div>
						</div>
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
	import { resolveMarketingPreviewProducts } from '../params/shared/marketing-helpers';
	import diyPageThemeMixin from './shared/diy-page-theme-mixin';

	export default {
  mixins: [diyPageThemeMixin],
  inject: ['diyModel'],
		data() {
			return {};
		},
		created() {},
		props: ['item', 'index', 'selectedIndex'],
		computed: {
			previewProducts() {
				return resolveMarketingPreviewProducts(this.item);
			},
			borderStyle() {
				const topRadio = `${this.item.style.topRadio}px`;
				const bottomRadio = `${this.item.style.bottomRadio}px`;
				const bgGradient =
					`linear-gradient(to right, ${this.item.style.bgcolor_color1 || '#fff'}, ${this.item.style.bgcolor_color2 || '#fff'})`;

				return {
					borderTopLeftRadius: topRadio,
					borderTopRightRadius: topRadio,
					borderBottomLeftRadius: bottomRadio,
					borderBottomRightRadius: bottomRadio,
					backgroundImage: bgGradient,
				};
			}
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
	.titleImgt {
		width: 85px;
		height: 45px;
	}

	.diy-product {
		overflow: hidden;
	}

	.diy-product .diy-head {
		padding: 0 10px;
		height: 45px;
		background-size: 100% 100% !important;
	}

	.diy-product .diy-head .name {
		font-size: 16px;
	}

	.product-list {
		.product-item {
			position: relative;
			overflow: hidden;

		}
	}

	.product-list.column__4 {
		.product-item {
			display: flex;
			justify-content: space-between;
			align-items: center;
			flex-direction: column;
		}
	}

	.product-list-box {
		padding: 0 7px 9px 7px;
	}

	/* 横向滚动 */
	.product-list.column__4 {
		display: flex;
		justify-content: flex-start;
		align-items: flex-start;
		overflow-x: auto;
		border-radius: 8px;
		background-color: #fff;

		.product-item {
			width: 89px;
			margin-right: 0;
			padding: 14px 0;
			flex-shrink: 0;
			border-right: 1px solid #eee;

			.product-cover {
				width: 74px;
				height: 74px;
				margin-bottom: 15px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				width: 74px;
				box-sizing: border-box;
				display: flex;
				justify-content: space-between;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 15px;

					.product-tag {
						color: #fff;
						background-color: #fff;
						text-align: center;
						line-height: 1;
						border-radius: 3px;
						font-size: 10px;
						padding: 3px 5px;
						margin-right: 5px;
					}
				}

				.price-num {
					font-size: 16px;
					line-height: 1;
				}

				.price {
					margin-top: 10px;
					width: 100%;
					min-height: 24px;
					margin-bottom: 0;
					display: flex;
					justify-content: space-between;
					align-items: flex-end;
				}

				.cart-btn {
					position: relative;
					width: 25px;
					min-width: 25px;
					height: 25px;
					font-size: 11px;
					padding: 0;
					margin-left: 3px;
				}
			}
		}
	}

	.cart-btn {
		min-width: 63px;
		height: 23px;
		line-height: 1;
		padding: 0 10px;
		box-sizing: border-box;
		font-size: 12px;
		color: #fff;
		display: flex;
		justify-content: center;
		align-items: center;
		border-radius: 100px;

		.cart-text {
			font-size: 12px;
		}
	}

	.cart-btn.icon {
		padding: 0;
		min-width: 23px;
		width: 23px;
		height: 23px;
		line-height: 23px;
		border-radius: 50%;

		.iconfont.icon {
			font-size: 12px;
			line-height: 1;
		}
	}

	.f11 {
		font-size: 11px;
	}
</style>