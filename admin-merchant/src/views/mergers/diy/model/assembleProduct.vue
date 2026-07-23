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
			<div class="product-list-box" :key="'column-' + productColumn + '-num-' + previewProducts.length"
				:class="`column__${productColumn}`">
				<div class="product-list" :class="`column__${productColumn}`"
					v-for="(listitem,listindex) in productColumn == 2  ? 2 : 1">
					<template v-for="(product, index) in previewProducts" :key="index">
						<div class="product-item" v-if=" productColumn !== 2 || index % 2 === listindex">
							<span class="assembel-num-pop"
								:style="{color: item.style.assemble_numtext_color,background: item.style.assemble_numbtn_color}"
								v-if="productColumn != 1 && productFlags.assembleNum">2人团</span>
							<!-- 商品图片 -->
							<img :style="{borderTopLeftRadius: item.style.productTopRadio + 'px',
						borderTopRightRadius: item.style.productTopRadio + 'px',
						borderBottomLeftRadius: item.style.productBottomRadio + 'px',
						borderBottomRightRadius: item.style.productBottomRadio + 'px'}" class="product-cover"
								v-img-url="product.image" />
							<div class="product-info">
								<!-- 商品名称 -->
								<div class="flex-1 ww100">
									<div v-if="productFlags.productName" :style="{color: item.style.productName_color}"
										class="product-name text-ellipsis">
										<span :class="{fb:item.style.nameWeight}">{{ product.product_name }}</span>
									</div>
									<div v-if="productColumn == 1" class="d-s-c">
										<span class="assembel-num" v-if="productFlags.assembleNum"
											:style="{color: item.style.assemble_numtext_color,background: item.style.assemble_numbtn_color}">2人成团</span>
										<span class="assembel-sales" v-if="productFlags.productSales"
											:style="{color: item.style.product_sales_color}">已拼1026人</span>
									</div>
								</div>
								<!-- 商品价格 -->
								<div class="price" :class="{hasbtn:productFlags.product_btn}">
									<div v-if="productFlags.productPrice" class="price-num"
										:style="{color: item.style.productPrice_color}">
										<span class="f11">拼团价</span>
										<span class="f11">¥</span>
										<span class="f18 fb">120</span>
									</div>
									<div class="f11 gray9 text-d-line" v-if="productFlags.linePrice"
										:style="{color: item.style.productLine_color}">
										<span class="">¥</span>
										<span class="">233</span>
									</div>
								</div>
								<div v-if="productFlags.product_btn && (productColumn == 1 || productColumn == 2)"
									class="cart-btn"
									:style="{color:item.style.btn_text_color,backgroundImage: 'linear-gradient(to right, ' + (item.style.productBtn_color1 || '#fff') + ', ' + (item.style.productBtn_color2 || '#fff') + ')'}">
									<span class="cart-text">去开团</span>
								</div>
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
	import {
		marketingDisplayFlag,
		resolveMarketingColumn,
		resolveMarketingPreviewProducts,
	} from '../params/shared/marketing-helpers';
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
			productColumn() {
				return resolveMarketingColumn(this.item);
			},
			previewProducts() {
				return resolveMarketingPreviewProducts(this.item);
			},
			productFlags() {
				const params = this.item.params || {};
				return {
					assembleNum: marketingDisplayFlag(params.assembleNum),
					linePrice: marketingDisplayFlag(params.linePrice),
					productName: marketingDisplayFlag(params.productName),
					productPrice: marketingDisplayFlag(params.productPrice),
					productSales: marketingDisplayFlag(params.productSales),
					product_btn: marketingDisplayFlag(params.product_btn),
				};
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

	.product-list-box.column__2 {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 10px;
	}

	.product-list {
		.product-item {
			position: relative;
			overflow: hidden;

		}
	}

	.product-list.column__1 {
		.product-item {
			display: flex;
			justify-content: space-between;
			align-items: center;
		}
	}

	.product-list.column__2,
	.product-list.column__3,
	.product-list.column__5,
	.product-list.column__4 {
		.product-item {
			display: flex;
			justify-content: space-between;
			align-items: center;
			flex-direction: column;
		}
	}

	/* 单列展示 */
	.product-list.column__1 {
		padding: 10px;

		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			margin-bottom: 11px;

			.product-cover {
				width: 96px;
				height: 96px;
				margin-right: 10px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				min-height: 96px;
				display: flex;
				justify-content: flex-start;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 15px;
					margin-bottom: 4px;
				}

				.assembel-num {
					display: flex;
					justify-content: center;
					align-items: center;
					min-width: 51px;
					padding: 0 6px;
					box-sizing: border-box;
					height: 18px;
					line-height: 1;
					border-radius: 4px;
					font-size: 11px;
					margin-right: 10px;
				}

				.assembel-sales {
					font-size: 12px;
				}

				.price.hasbtn {
					min-height: 30px;
				}

				.product-sale {
					font-size: 12px;
					margin-right: 6px;
				}

				.cart-btn {
					width: 63px;
					height: 23px;
					position: absolute;
					right: 0;
					bottom: 0;
				}
			}
		}
	}

	/* 两列展示(纵向) */
	.product-list.column__2:first-child {
		margin-right: 10px;
	}

	.product-list.column__2 {
		flex: 1 1 0;
		min-width: 0;

		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			margin-bottom: 11px;

			.assembel-num-pop {
				position: absolute;
				left: 0;
				top: 15px;
				min-width: 48px;
				box-sizing: border-box;
				padding: 0 10px;
				height: 19px;
				line-height: 1;
				background: #000000;
				border-radius: 0px 10px 10px 0px;
				color: #FFFFFF;
				display: flex;
				justify-content: center;
				align-items: center;
			}

			.product-cover {
				width: 100%;
				height: 162px;
				margin-bottom: 7px;
				display: block;
				flex-shrink: 0;
			}

			.product-info {
				width: 100%;
				box-sizing: border-box;
				padding: 0px;
				display: flex;
				justify-content: space-between;
				align-items: flex-start;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 15px;
				}

				.price {
					min-height: 30px;
				}

				.product-sale {
					font-size: 13px;
					margin-top: 6px;
				}

				.cart-btn {
					width: 63px;
					height: 23px;
					position: absolute;
					right: 0;
					bottom: 0;
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

			.assembel-num-pop {
				position: absolute;
				left: 0;
				top: 0;
				min-width: 48px;
				box-sizing: border-box;
				padding: 0 10px;
				height: 19px;
				line-height: 1;
				background: #000000;
				border-radius: 10px 0px 10px 0px;
				color: #FFFFFF;
				display: flex;
				justify-content: center;
				align-items: center;
			}

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
				display: flex;
				justify-content: space-between;
				align-items: center;
				flex-direction: column;
				flex: 1;
				position: relative;

				.product-name {
					font-size: 14px;
					margin-bottom: 4px;
				}

				.price {
					margin-bottom: 0;
					width: 100%;
					.price-num {
						position: relative;
						border-radius: 0 18px 18px 0;
						margin-bottom: 4px;

						.price-icon {
							position: absolute;
							left: -6px;
							top: 0;
							bottom: 0;
							margin: auto;
							z-index: 1;
							width: 12px;
							height: 18px;
						}
					}


				}

				.cart-btn {
					position: absolute;
					right: 0;
					bottom: 0;
				}
			}
		}
	}

	/* 横向滚动 */
	.product-list.column__4 {
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

			.assembel-num-pop {
				position: absolute;
				left: 6px;
				top: 6px;
				min-width: 48px;
				box-sizing: border-box;
				padding: 0 10px;
				height: 19px;
				line-height: 1;
				background: #000000;
				border-radius: 10px;
				color: #FFFFFF;
				display: flex;
				justify-content: center;
				align-items: center;
			}

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
					font-size: 15px;
				}

				.price-num {
					font-size: 16px;
					line-height: 1;
				}

				.price {
					width: 100%;
					min-height: 24px;
					margin-bottom: 0;
				}

				.cart-btn {
					position: relative;
					width: 20px;
					min-width: 20px;
					height: 20px;
					font-size: 11px;
					padding: 0;

					margin-right: 3px;
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