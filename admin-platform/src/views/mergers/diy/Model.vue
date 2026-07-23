<template>
	<div class="pr">
		<div class="operate-pop" :class="{ disabled: form.selectedIndex == -1 }"
			v-if="diyData.items && diyData.items.length > 0">
			<div class="iconbox d-c-c" @click="DeleteFunc">
				<el-icon>
					<Delete />
				</el-icon>
			</div>
			<!-- 复制 -->
			<div class="iconbox d-c-c" :class="{ disabled: form.selectedIndex < 0 || isSpecialItem }"
				@click=" DocumentCopyFunc">
				<el-icon>
					<DocumentCopy />
				</el-icon>
			</div>

			<!-- 上移 -->
			<div class="iconbox d-c-c" :class="{ disabled: form.selectedIndex <= 0 || isSpecialItem }"
				@click=" ArrowUpBoldFunc">
				<el-icon>
					<ArrowUpBold />
				</el-icon>
			</div>

			<!-- 下移 -->
			<div class="iconbox d-c-c"
				:class="{ disabled: form.selectedIndex < 0 || form.selectedIndex === diyData.items.length-1 || isSpecialItem }"
				@click=" ArrowDownBoldFunc">
				<el-icon>
					<ArrowDownBold />
				</el-icon>
			</div>
		</div>
		<div class="scroll-ybox" style="padding: 10px 0; box-sizing: border-box;">
			<div class="diy-phone-container">
				<!--顶部设置栏-->
				<div class="diy-phone-item" :class="{ active: form.selectedIndex < 0 }">
					<div @click="onEditer(-1)" class="draggable-title" style="left:auto;right: -90px;top: 20px;">页面设置
					</div>
					<Setpages v-if="diyType != 'center'" :diyData="diyData" :isDiy="isDiy"></Setpages>
					<Setcenter v-if="diyType == 'center'" :diyData="diyData"></Setcenter>
				</div>
				<draggable class="wrapper" v-model="diyData.items" item-key="_diyUid" :move="checkMove"
					:options="{ animation: 120, filter: '.drag__nomove' }" draggable=".draggable">
					<template #item="{ element, index }">
						<div class="diy-phone-item" :class="{ active: form.selectedIndex == index,pstatic:element.type == 'surface'||element.type == 'service' || element.type == 'videoLive',
						 drag__nomove: element.type === 'option'
                   || element.type === 'search'
                   || element.type === 'topMerge',
				   draggable: index>1
						}">
							<div @click="onEditer(index)"
								v-if="element.type != 'surface'||element.type != 'service'&&element.type != 'videoLive'"
								class="draggable-title"
								:class="{right:element.type == 'guide'||element.type == 'blank'}">{{ element.name }}
							</div>
							<!-- 搜索栏 -->
							<template v-if="element.type == 'search'">
								<Search :item="element" :index="index" :selectedIndex="form.selectedIndex"></Search>
							</template>
							<!-- 图片轮播 -->
							<template v-else-if="element.type == 'banner'">
								<Banner :item="element" :index="index" :selectedIndex="form.selectedIndex"></Banner>
							</template>
							<!-- 图片-->
							<template v-else-if="element.type == 'imageSingle'">
								<ImageSingle :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</ImageSingle>
							</template>

							<!-- 橱窗-->
							<template v-else-if="element.type == 'window'">
								<Window :item="element" :index="index" :selectedIndex="form.selectedIndex"></Window>
							</template>
							<!-- 热区-->
							<template v-else-if="element.type == 'hotspot'">
								<Hotspot :item="element" :index="index" :selectedIndex="form.selectedIndex"></Hotspot>
							</template>
							<!-- 视频组-->
							<template v-else-if="element.type == 'video'">
								<Video :item="element" :index="index" :selectedIndex="form.selectedIndex"></Video>
							</template>
							<!--文章-->
							<template v-else-if="element.type == 'article'">
								<ArticleIndex :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</ArticleIndex>
							</template>
							<!--头条快报-->
							<template v-else-if="element.type == 'special'">
								<Special :item="element" :index="index" :selectedIndex="form.selectedIndex"></Special>
							</template>
							<!--公告组-->
							<template v-else-if="element.type == 'notice'">
								<Notice :item="element" :index="index" :selectedIndex="form.selectedIndex"></Notice>
							</template>
							<!--导航组-->
							<template v-else-if="element.type == 'navBar'">
								<NavBar :item="element" :index="index" :selectedIndex="form.selectedIndex"></NavBar>
							</template>
							<!--商品组-->
							<template v-else-if="element.type == 'product'">
								<ProductIndex :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</ProductIndex>
							</template>
							<!--优惠券-->
							<template v-else-if="element.type == 'coupon'">
								<Coupon :item="element" :index="index" :selectedIndex="form.selectedIndex"></Coupon>
							</template>
							<!--门店-->
							<template v-else-if="element.type == 'store'">
								<Store :item="element" :index="index" :selectedIndex="form.selectedIndex"></Store>
							</template>
							<!--客服-->
							<template v-else-if="element.type == 'service'">
								<Service :item="element" :index="index" :selectedIndex="form.selectedIndex"></Service>
							</template>
							<!--富文本-->
							<template v-else-if="element.type == 'richText'">
								<RichText :item="element" :index="index" :selectedIndex="form.selectedIndex"></RichText>
							</template>
							<!--辅助空白-->
							<template v-else-if="element.type == 'blank'">
								<Blank :item="element" :index="index" :selectedIndex="form.selectedIndex"></Blank>
							</template>
							<!--辅助线-->
							<template v-else-if="element.type == 'guide'">
								<Guide :item="element" :index="index" :selectedIndex="form.selectedIndex"></Guide>
							</template>
							<!--秒杀-->
							<template v-else-if="element.type == 'seckillProduct'">
								<Seckill :item="element" :index="index" :selectedIndex="form.selectedIndex"></Seckill>
							</template>
							<!--预告-->
							<template v-else-if="element.type == 'previewProduct'">
								<Preview :item="element" :index="index" :selectedIndex="form.selectedIndex"></Preview>
							</template>
							<!--拼团-->
							<template v-else-if="element.type == 'assembleProduct'">
								<assembleProduct :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</assembleProduct>
							</template>
							<!--砍价-->
							<template v-else-if="element.type == 'bargainProduct'">
								<BargainProduct :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</BargainProduct>
							</template>
							<!--新人-->
							<template v-else-if="element.type == 'newActivity'">
								<NewActivity :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</NewActivity>
							</template>
							<!--微信直播-->
							<template v-else-if="element.type == 'wxlive'">
								<Wxlive :item="element" :index="index" :selectedIndex="form.selectedIndex"></Wxlive>
							</template>
							<!--七禧直播-->
							<template v-else-if="element.type == 'qixiLive'">
								<QixiLive :item="element" :index="index" :selectedIndex="form.selectedIndex"></QixiLive>
							</template>
							<!--标题-->
							<template v-else-if="element.type == 'title'">
								<Title :item="element" :index="index" :selectedIndex="form.selectedIndex"></Title>
							</template>
							<!--个人信息-->
							<template v-else-if="element.type == 'base'">
								<Base :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</Base>
							</template>
							<!--我的订单-->
							<template v-else-if="element.type == 'order'">
								<Order :item="element" :index="index" :selectedIndex="form.selectedIndex"></Order>
							</template>
							<!--视频号直播-->
							<template v-else-if="element.type == 'videoLive'">
								<ShipinLiveindex :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</ShipinLiveindex>
							</template>
							<!--选项卡-->
							<template v-else-if="element.type == 'option'">
								<Option :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</Option>
							</template>
							<!--轮播搜索-->
							<template v-else-if="element.type == 'topMerge'">
								<TopMerge :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</TopMerge>
							</template>
							<!--悬浮按钮-->
							<template v-else-if="element.type == 'surface'">
								<Surface :item="element" :index="index" :selectedIndex="form.selectedIndex">
								</Surface>
							</template>
						</div>
					</template>
				</draggable>
				<div style="width: 100%;height: 10px;"></div>
			</div>
		</div>
	</div>
</template>

<script>
	import { ArrowDownBold, ArrowUpBold, Delete, DocumentCopy } from '@element-plus/icons-vue';
	import { ElMessageBox } from 'element-plus';

	import Setpages from './model/Setpages.vue';
	import Setcenter from './model/Setcenter.vue';
	import Search from './model/Search.vue';
	import Banner from './model/Banner.vue';
	import ImageSingle from './model/ImageSingle.vue';
	import Window from './model/Window.vue';
	import Hotspot from './model/Hotspot.vue';
	import Video from './model/Video.vue';
	import ArticleIndex from './model/Article.vue';
	import Special from './model/Special.vue';
	import Notice from './model/Notice.vue';
	import NavBar from './model/NavBar.vue';
	import ProductIndex from './model/Product.vue';
	import Coupon from './model/Coupon.vue';
	import Store from './model/Store.vue';
	import Service from './model/Service.vue';
	import RichText from './model/RichText.vue';
	import Blank from './model/Blank.vue';
	import Guide from './model/Guide.vue';
	import Seckill from './model/Seckill.vue';
	import Preview from './model/Preview.vue';
	import assembleProduct from './model/assembleProduct.vue';
	import BargainProduct from './model/BargainProduct.vue';
	import NewActivity from './model/NewActivity.vue';
	import Wxlive from './model/Wxlive.vue';
	import QixiLive from './model/QixiLive.vue';
	import Title from './model/Title.vue';
	import Base from './model/Base.vue';
	import Order from './model/Order.vue';
	import draggable from 'vuedraggable';
	import ShipinLiveindex from './model/ShipinLive.vue';
	import Option from './model/Option.vue';
	import TopMerge from './model/TopMerge.vue';
	import Surface from './model/Surface.vue';
	export default {
		components: {
			ArrowDownBold,
			ArrowUpBold,
			Delete,
			DocumentCopy,
			/*顶部状态栏*/
			Setpages,
			Setcenter,
			/*搜索组件*/
			Search,
			/*图片轮播组件*/
			Banner,
			/*图片组件*/
			ImageSingle,
			/*图片橱窗*/
			Window,
			/*热区*/
			Hotspot,
			/*视频*/
			Video,
			/*文章*/
			ArticleIndex,
			/*头条快报*/
			Special,
			/*公告组*/
			Notice,
			/*导航组*/
			NavBar,
			/*商品组*/
			ProductIndex,
			/*优惠券*/
			Coupon,
			/*门店*/
			Store,
			/*客服*/
			Service,
			/*富文本*/
			RichText,
			/*辅助空白*/
			Blank,
			/*辅助线*/
			Guide,
			/*拖动*/
			draggable,
			/*秒杀*/
			Seckill,
			/*拼团*/
			assembleProduct,
			/*砍价*/
			BargainProduct,
			/* 新人 */
			NewActivity,
			/*微信直播*/
			Wxlive,
			QixiLive,
			/*标题*/
			Title,
			Preview,
			Base,
			Order,
			Option,
			TopMerge,
			ShipinLiveindex,
			Surface
		},
		data() {
			return {};
		},
		props: {
			form: Object,
			defaultData: Object,
			diyData: Object,
			diyType: String,
			isDiy: Boolean
		},
		provide() {
			const self = this;
			return {
				diyModel: {
					onDeleleItem(index) {
						self.onDeleleItem(index);
					},
					onEditer(index) {
						self.onEditer(index);
					},
				},
			};
		},
		computed: {
			isSpecialItem() {
				const n = this.form.selectedIndex;
				if (n < 0) return false;
				console.log('isSpecialItem');
				console.log(this.diyData.items);
				console.log(n);
				console.log(this.diyData.items[n]);
				console.log(this.diyData.items[n].type);
				const t = this.diyData.items[n].type;
				return ['topMerge', 'option', 'search'].includes(t);
			}
		},
		created() {},
		methods: {
			checkMove(evt) {
				const draggedType = evt.draggedContext.element.type;
				const toIndex = evt.relatedContext.index; // 目标位置索引
				if (toIndex < 2 && draggedType !== 'topMerge' && draggedType !== 'search') {
					return false;
				}
				if (draggedType === 'option' && toIndex !== 1) {
					return false;
				}
				if (draggedType !== 'option' && toIndex === 1) {
					return false;
				}
				return true;
			},
			swapArray(arr, index1, index2) {
				arr[index1] = arr.splice(index2, 1, arr[index1])[0];
				return arr;
			},
			specialType(e) {
				let self = this;
				let n = e || self.form.selectedIndex;
				if (n == -1) {
					return true;
				}
				let type = this.diyData.items[n].type;
				if (type == 'option' || type == 'search' || type == 'topMerge') {
					return true;
				}
				return false;
			},
			DeleteFunc() {
				let self = this;
				let n = self.form.selectedIndex;
				console.log(self.form.selectedIndex)
				if (n < 0) {
					return;
				}
				self.diyData.items.splice(n, 1);
				self.form.selectedIndex = -1;
			},
			DocumentCopyFunc() {
				let self = this;
				let n = self.form.selectedIndex;
				if (n < 0) {
					return;
				}
				if (self.specialType()) {
					return
				}
				let item = self.diyData.items[n];
				self.diyData.items.splice(n, 0, item);
			},
			ArrowUpBoldFunc() {
				const n = this.form.selectedIndex;
				if (n < 0) return;

				if (this.specialType() || n === 0) return;

				const type = this.diyData.items[n].type;
				const targetIndex = n - 1;
				if (targetIndex < 2 && type !== 'topMerge' && type !== 'search') {
					return;
				}

				this.swapArray(this.diyData.items, n, targetIndex);
				this.form.selectedIndex = targetIndex;
			},
			ArrowDownBoldFunc() {
				let self = this;
				let n = self.form.selectedIndex;
				if (n < 0) {
					return;
				}
				if (self.specialType()) {
					return
				}
				if (n + 1 != self.diyData.items.length) {
					self.swapArray(self.diyData.items, n, n + 1);
					self.form.selectedIndex++;
				}
			},
			/*删除diy元素*/
			onDeleleItem: function(index) {
				let self = this;
				ElMessageBox.confirm('确定要删除吗?', '提示', {
					type: 'warning'
				}).then(() => {
					self.diyData.items.splice(index, 1);
					self.form.selectedIndex = -1;
				});
			},

			/*编辑当前选中的Diy元素*/
			onEditer: function(index) {
				let self = this;
				// 记录当前选中元素的索引
				self.form.selectedIndex = index;
				// 当前选中的元素数据
				self.form.curItem = self.form.selectedIndex < 0 ? self.diyData.page : self.diyData.items[self.form
					.selectedIndex];
				// 注册编辑器事件
				//self.initEditor();
			},

			/* 注册编辑器事件*/
			initEditor: function() {
				let self = this;
				// 注册dom事件
				self.$nextTick(function() {
					// 销毁 umeditor 组件
					if (self.form.umeditor.hasOwnProperty('key')) {
						self.form.umeditor.destroy();
					}
					// 注册html组件
					self.editorHtmlComponent();
					// 富文本事件
					if (self.form.curItem.type === 'richText') {
						//self.onRichText(self.form.curItem);
					}
				});
			},

			/*编辑器事件：html组件*/
			editorHtmlComponent: function() {
				let self = this;
				var editor = self.$refs['diy-editor'];
				// 单/多选框
				//editor.find('input[type=checkbox], input[type=radio]').uCheck();
				// select组件
				// $editor.find('select').selected();
			}
		}
	};
</script>

<style lang="scss" scoped>
/* Preview column chrome lives in diy-legacy-utils.scss (unscoped under .diy-editor-shell) */
</style>