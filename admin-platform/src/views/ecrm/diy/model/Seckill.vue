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
			<div class="diy-head d-b-c" :style="headerStyle">
				<div class="left d-s-c">
					<div v-if="item.params.titleType == 1 || !hasTitleImage" class="name" :style="{
							color: themeColor,
							fontSize: item.style.titleSize + 'px',
							fontWeight:item.style.titleWeight == 1?'bold':'',
							fontStyle:item.style.titleWeight == 2?'italic':''
						}">
						{{ item.params.title || '限时秒杀' }}
					</div>
					<img v-else class="titleImgt" :src="item.params.titleimage" alt="" />
					<div class="datetime d-s-c">
						<text class="text" :style="{ color: item.style.color }">距结束仅剩</text>
						<span class="hour" :style="{
								color: item.style.number_color,
								background:`linear-gradient(to right, ${item.style.title_color1 || '#fff'}, ${item.style.title_color2 || '#fff'})`
							}">
							30
						</span>
						<span class="text" :style="{ color: item.style.title_color1 }">:</span>
						<span class="hour" :style="{
								color: item.style.number_color,
								background:`linear-gradient(to right, ${item.style.title_color1 || '#fff'}, ${item.style.title_color2 || '#fff'})`
							}">
							00
						</span>
						<span class="text" :style="{ color: item.style.title_color1 }">:</span>
						<span class="hour" :style="{
								color: item.style.number_color,
								background:`linear-gradient(to right, ${item.style.title_color1 || '#fff'}, ${item.style.title_color2 || '#fff'})`
							}">
							00
						</span>
					</div>
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
							<!-- 商品图片 -->
							<img :style="{borderTopLeftRadius: item.style.productTopRadio + 'px',
						borderTopRightRadius: item.style.productTopRadio + 'px',
						borderBottomLeftRadius: item.style.productBottomRadio + 'px',
						borderBottomRightRadius: item.style.productBottomRadio + 'px'}" class="product-cover"
								v-img-url="product.image" />
							<div class="product-info">
								<!-- 商品名称 -->
								<div v-if="productFlags.productName" :style="{color: item.style.productName_color}"
									class="product-name "
									:class="productColumn == 1?'text-ellipsis-2':'text-ellipsis'">
									<span :class="{fb:item.style.nameWeight}">{{ product.product_name }}</span>
								</div>
								<div class="d-s-c ww100" v-if="productFlags.productSales">
									<div class="product-labels-boxs"
										:style="{backgroundImage: 'linear-gradient(to right, ' + (item.style.productSlider_color1 || '#fff') + ', ' + (item.style.productSlider_color2 || '#fff') + ')'}">
										<div class="product-labels">
											<div class="product-labels-width" style="width: 40%;"
												:style="{backgroundImage: 'linear-gradient(to right, ' + (item.style.productSlider_color1 || '#fff') + ', ' + (item.style.productSlider_color2 || '#fff') + ')'}">
											</div>
										</div>
									</div>
									<div class="f11" :style="{color: item.style.productSlider_color}">已抢40%</div>
								</div>
								<!-- 商品价格 -->
								<div class="price" :class="{hasbtn:productFlags.product_btn}">
									<div v-if="productFlags.productPrice" class="price-num"
										:style="{color: item.style.productPrice_color,backgroundImage: productColumn == 3?'linear-gradient(to right, ' + (item.style.productBtn_color1 || '#fff') + ', ' + (item.style.productBtn_color2 || '#fff') + ')':''}">
										<img v-if="productColumn == 3" class="price-icon"
											:src="seckillIconUrl" alt="" />
										<div v-if="productColumn == 4 && productFlags.product_btn"
											class="cart-btn"
											:style="{color:item.style.btn_text_color,backgroundImage: 'linear-gradient(to right, ' + (item.style.productBtn_color1 || '#fff') + ', ' + (item.style.productBtn_color2 || '#fff') + ')'}">
											<span class="cart-text">抢</span>
										</div>
										<span class="f11">¥</span>
										<span class="f16 fb">120</span>
									</div>
									<div class="f11 gray9 text-d-line" v-if="productFlags.linePrice"
										:style="{color: item.style.productLine_color,textAlign:productColumn == 3 || productColumn == 4?'center':''}">
										<span>¥</span>
										<span>233</span>
									</div>
								</div>
								<div v-if="productFlags.product_btn && (productColumn == 1 || productColumn == 2)"
									class="cart-btn"
									:style="{color:item.style.btn_text_color,backgroundImage: 'linear-gradient(to right, ' + (item.style.productBtn_color1 || '#fff') + ', ' + (item.style.productBtn_color2 || '#fff') + ')'}">
									<span class="cart-text">抢购</span>
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
	import { DIY_SECKILL_ICON_URL } from '#/utils/diy/order-icons';
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
			return {
				seckillIconUrl: DIY_SECKILL_ICON_URL,
			};
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
			hasTitleImage() {
				return Boolean(String((this.item.params || {}).titleimage || '').trim());
			},
			headerStyle() {
				const params = this.item.params || {};
				const style = this.item.style || {};
				const backgroundImage = String(params.bgimage || '').trim();
				const gradient = `linear-gradient(to right, ${style.titleBg_color1 || '#fff'}, ${style.titleBg_color2 || '#fff'})`;

				return {
					background:
						Number(params.titleBgType) === 1 && backgroundImage
							? `url(${backgroundImage})`
							: gradient,
					backgroundPosition: 'center',
					backgroundRepeat: 'no-repeat',
					backgroundSize: '100% 100%',
				};
			},
			productFlags() {
				const params = this.item.params || {};
				return {
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
		/* width: 100px; */
	}

	.diy-product .diy-head .datetime {
		margin-left: 13px;
		border: none;
		white-space: nowrap;
	}

	.diy-product .diy-head .datetime>span {
		display: inline-block;
	}

	.diy-product .diy-head .datetime .text {
		margin: 0 4px;
		line-height: 1;
	}

	.diy-product .diy-head .datetime .hour {
		display: inline-block;
		width: 16px;
		height: 16px;
		background: linear-gradient(to right, rgb(255, 255, 255), rgb(255, 255, 255));
		color: rgb(253, 59, 84);
		font-size: 12px;
		border-radius: 3px;
		padding: 0;
		line-height: 16px;
		text-align: center;
	}

	.diy-product .diy-head .datetime .box {
		padding: 2px;
		border-radius: 4px;
		background: #000000;
		color: #ffffff;
	}

	.product-list-box.column__2 {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 10px;
	}

	.product-list {
		.product-item {
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
					font-size: 13px;
					margin-bottom: 13px;
					flex: 1;

				}

				.price.hasbtn {
					min-height: 30px;
				}

				.product-labels-boxs {
					height: 8px;
					border-radius: 4px;
					// overflow: hidden;
					flex: 1;
					margin-right: 6px;

					.product-labels {
						background-color: rgba(#fff, 0.9);
						border-radius: 4px;

						.product-labels-width {
							height: 8px;
							border-radius: 4px;
							width: 40%;
						}
					}
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
		flex: 1;

		.product-item:last-child {
			margin-bottom: 0;
		}

		.product-item {
			margin-bottom: 11px;

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
					font-size: 13px;
					flex: 1;
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
					font-size: 12px;
					margin-bottom: 4px;
				}

				.price {
					margin-bottom: 0;

					.price-num {
						position: relative;
						padding: 0 12px;
						min-width: 68px;
						border-radius: 0 18px 18px 0;
						height: 18px;
						display: flex;
						justify-content: center;
						align-items: center;
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

				.price-num {
					display: flex;
					justify-content: center;
					align-items: center;
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

	.f11 {
		font-size: 11px;
	}
</style>
