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
		<div class="diy-product" :style="borderStyle">
			<div class="product-list-box" :key="'column-' + productColumn" :class="`column__${productColumn}`">
				<div class="product-list" :class="`column__${productColumn}`"
					v-for="(listitem,listindex) in productColumn == 2 || productColumn == 4 ? 2 : 1">
					<template v-for="(product, index) in previewProducts" :key="index">
						<div class="product-item" :style="productStyle"
							v-if=" ![2,4].includes(productColumn) || index % 2 === listindex">
							<!-- 商品图片 -->
							<img :style="{borderTopLeftRadius: item.style.productTopRadio + 'px',
						borderTopRightRadius: item.style.productTopRadio + 'px',
						borderBottomLeftRadius: item.style.productBottomRadio + 'px',
						borderBottomRightRadius: item.style.productBottomRadio + 'px'}" class="product-cover"
								v-img-url="product.image" />
							<div class="product-info">
								<!-- 商品名称 -->
								<div :style="{color: item.style.product_name_color}"
									class="product-name text-ellipsis-2">
									<span v-if="productFlags.productName"
										:class="{fb:item.style.nameWeight}">{{ product.product_name }}</span>
								</div>
								<!-- 商品价格 -->
								<div class="price d-s-c">
									<div v-if="productFlags.productPrice"
										:style="{color: productPriceColor}">
										<span class="f11">¥</span>
										<span class="f16 fb">{{ product.product_price }}</span>
									</div>
									<div class="f11 ml5 gray9 text-d-line"
										v-if="productFlags.linePrice && product.line_price > 0"
										:style="{color: linePriceColor}">
										<span>¥</span>
										<span>{{ product.line_price }}</span>
									</div>
								</div>
								<div class="d-s-c">
									<!-- 商品销量 -->
									<div v-if="productFlags.productSales" class="product-sale">
										<span
											:style="{color: item.style.product_sales_color}">已售{{ product.product_sales }}件</span>
									</div>
									<!-- 好评率 -->
									<div v-if="productFlags.comment" class="product-comment">
										<span
											:style="{color: item.style.product_comment_color}">好评率{{ product.product_sales }}</span>
									</div>
								</div>
								<template v-if="productFlags.showCart == 1">
									<div v-if="productFlags.cartType == 0" class="cart-btn"
										:style="{color:item.style.cart_text_color,backgroundImage: 'linear-gradient(to right, ' + (item.style.cart_color1 || '#fff') + ', ' + (item.style.cart_color2 || '#fff') + ')'}">
										<span class="cart-text">{{productFlags.cartText || '购买'}}</span>
									</div>
									<div v-if="productFlags.cartType == 1" class="cart-btn icon"
										:style="{backgroundImage: 'linear-gradient(to right, ' + (item.style.cart_color1 || '#fff') + ', ' + (item.style.cart_color2 || '#fff') + ')'}">
										<span class=" icon iconfont icon-icozhuanhuan"
											:style="{color:item.style.cart_text_color}"></span>
									</div>
									<div v-if="productFlags.cartType == 2" class="cart-btn icon"
										:style="{backgroundImage: 'linear-gradient(to right, ' + (item.style.cart_color1 || '#fff') + ', ' + (item.style.cart_color2 || '#fff') + ')'}">
										<span class=" icon iconfont icon-jia"
											:style="{color:item.style.cart_text_color}"></span>
									</div>
								</template>
							</div>
						</div>
					</template>

				</div>

			</div>
		</div>
		<div class="btn-edit-del">
			<div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
		</div>
	</div>
</template>

<script>
	import { resolveProductPreviewProducts } from '../params/shared/marketing-helpers';

	export default {
  inject: ['diyModel'],
		data() {
			return {};
		},
		created() {},
		props: ['item', 'index', 'selectedIndex'],
		computed: {
			previewProducts() {
				return resolveProductPreviewProducts(this.item);
			},
			productColumn() {
				const params = this.item.params || {};
				const style = this.item.style || {};
				const column = params.column ?? style.column ?? 2;
				return Number(column);
			},
			productFlags() {
				const params = this.item.params || {};
				const show = (this.item.style && this.item.style.show) || {};
				return {
					productName: params.productName ?? show.productName,
					productPrice: params.productPrice ?? show.productPrice,
					linePrice: params.linePrice ?? show.linePrice,
					productSales: params.productSales ?? show.productSales,
					comment: params.comment ?? show.comment,
					showCart: params.showCart,
					cartType: params.cartType,
					cartText: params.cartText,
				};
			},
			productPriceColor() {
				return this.item.style?.product_price_color || '#FF4C01';
			},
			linePriceColor() {
				return this.item.style?.line_price_color || '#999999';
			},
			borderStyle() {
				const isSpecialColumn = [3, 4, 6].includes(this.productColumn);
				const topRadio = isSpecialColumn ? `${this.item.style.topRadio}px` : '';
				const bottomRadio = isSpecialColumn ? `${this.item.style.bottomRadio}px` : '';
				const bgGradient = isSpecialColumn ?
					`linear-gradient(to right, ${this.item.style.bgcolor_color1 || '#fff'}, ${this.item.style.bgcolor_color2 || '#fff'})` :
					'';

				return {
					borderTopLeftRadius: topRadio,
					borderTopRightRadius: topRadio,
					borderBottomLeftRadius: bottomRadio,
					borderBottomRightRadius: bottomRadio,
					backgroundImage: bgGradient,
				};
			},
			productStyle() {
				const isSpecialColumn = ![3, 4, 6].includes(this.productColumn);
				const topRadio = isSpecialColumn ? `${this.item.style.topRadio}px` : '';
				const bottomRadio = isSpecialColumn ? `${this.item.style.bottomRadio}px` : '';
				const bgGradient = isSpecialColumn ?
					`linear-gradient(to right, ${this.item.style.bgcolor_color1 || '#fff'}, ${this.item.style.bgcolor_color2 || '#fff'})` :
					'';

				return {
					borderTopLeftRadius: topRadio,
					borderTopRightRadius: topRadio,
					borderBottomLeftRadius: bottomRadio,
					borderBottomRightRadius: bottomRadio,
					backgroundImage: bgGradient,
				};
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
	.product-list-box.column__2,
	.product-list-box.column__4 {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
	}

	.product-list {
		.product-item {
			overflow: hidden;
		}
	}

	.product-list.column__1,
	.product-list.column__4 {
		.product-item {
			display: flex;
			justify-content: space-between;
			align-items: center;
		}
	}

	.product-list.column__2,
	.product-list.column__3,
	.product-list.column__5,
	.product-list.column__6 {
		.product-item {
			display: flex;
			justify-content: space-between;
			align-items: center;
			flex-direction: column;
		}
	}

	/* 单列展示 */
	.product-list.column__1 {
		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			padding: 10px 10px 15px 10px;
			margin-bottom: 11px;

			.product-cover {
				width: 110px;
				height: 110px;
				margin-right: 10px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				min-height: 110px;
				display: flex;
				justify-content: flex-start;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 14px;
					margin-bottom: 16px;
					flex: 1;

				}

				.product-labels {
					margin-bottom: 8px;
				}

				.price {
					margin-bottom: 5px;
				}

				.product-sale {
					font-size: 12px;
					margin-right: 6px;
				}

				.product-comment {
					font-size: 12px;
				}

				.cart-btn {
					position: absolute;
					right: 0;
					bottom: 0;
				}
			}
		}
	}

	/* 两列展示(纵向) */
	.product-list.column__2:first-child {
		margin-right: 9px;
	}

	.product-list.column__2 {
		flex: 1;

		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			margin-bottom: 11px;

			.product-cover {
				width: 100%;
				height: 171px;
				margin-bottom: 7px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				width: 100%;
				box-sizing: border-box;
				padding: 0 9px 10px 9px;
				display: flex;
				justify-content: space-between;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;
				min-height: 33px;

				.product-name {
					font-size: 12px;
					flex: 1;
				}

				.price {
					margin-bottom: 5px;
				}

				.product-sale {
					font-size: 12px;
					margin-right: 6px;
				}

				.cart-btn {
					position: absolute;
					right: 10px;
					bottom: 10px;
				}
			}
		}
	}

	/* 三列展示 */
	.product-list.column__3 {
		display: flex;
		justify-content: flex-start;
		align-items: flex-start;
		flex-wrap: wrap;
		padding: 10px;
		padding-bottom: 0;

		.product-item:nth-child(3n) {
			margin-right: 0;
		}

		.product-item {
			margin-bottom: 11px;
			width: 31.3%;
			margin-right: 10px;

			.product-cover {
				width: 100%;
				height: 105px;
				margin-bottom: 7px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				width: 100%;
				box-sizing: border-box;
				padding: 0 9px 0 9px;
				display: flex;
				justify-content: space-between;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 12px;
					line-height: 18px;
					height: 36px;
				}

				.price {
					margin-bottom: 0;
				}

				.cart-btn {
					position: absolute;
					right: 0;
					bottom: 0;
				}
			}
		}
	}

	/* 两列展示(横向) */
	.product-list.column__4 {
		flex: 1;

		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			padding: 10px 10px 15px 10px;
			margin-bottom: 11px;

			.product-cover {
				width: 71px;
				height: 71px;
				margin-right: 7px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				min-height: 71px;
				display: flex;
				justify-content: space-between;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 12px;
					flex: 1;
				}

				.price {
					margin-bottom: 5px;
				}
			}
		}
	}

	/* 大图 */
	.product-list.column__5 {
		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			margin-bottom: 11px;

			.product-cover {
				width: 100%;
				height: 201px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				width: 100%;
				padding: 9px 10px 15px 10px;
				box-sizing: border-box;
				min-height: 110px;
				display: flex;
				justify-content: flex-start;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 14px;
					margin-bottom: 16px;

				}

				.product-labels {
					margin-bottom: 8px;
				}

				.price {
					margin-bottom: 5px;
				}

				.product-sale {
					font-size: 12px;
					margin-right: 6px;
				}

				.product-comment {
					font-size: 12px;
				}

				.cart-btn {
					position: absolute;
					right: 10px;
					bottom: 15px;
				}
			}
		}
	}

	/* 横向滚动 */
	.product-list.column__6 {
		display: flex;
		justify-content: flex-start;
		align-items: flex-start;
		padding: 10px;
		padding-bottom: 0;
		overflow-x: auto;

		.product-item {
			width: 96px;
			margin-right: 11px;
			padding-bottom: 11px;
			flex-shrink: 0;

			.product-cover {
				width: 96px;
				height: 96px;
				margin-bottom: 8px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				width: 100%;
				box-sizing: border-box;
				display: flex;
				justify-content: space-between;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 12px;
					line-height: 18px;
					height: 36px;
				}

				.price {
					min-height: 24px;
					margin-bottom: 0;
				}

				.cart-btn {
					position: absolute;
					right: 0;
					bottom: 0;
				}
			}
		}
	}

	.ml5 {
		margin-left: 5px;
	}

	.cart-btn {
		min-width: 53px;
		height: 23px;
		line-height: 1;
		padding: 0 10px;
		box-sizing: border-box;
		background: #409eff;
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
</style>